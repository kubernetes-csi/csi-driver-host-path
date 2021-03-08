/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hostpath

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/golang/glog"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	fs "k8s.io/kubernetes/pkg/volume/util/fs"
	"k8s.io/kubernetes/pkg/volume/util/volumepathhandler"
	utilexec "k8s.io/utils/exec"
)

const (
	kib    int64 = 1024
	mib    int64 = kib * 1024
	gib    int64 = mib * 1024
	gib100 int64 = gib * 100
	tib    int64 = gib * 1024
	tib100 int64 = tib * 100

	// storageKind is the special parameter which requests
	// storage of a certain kind (only affects capacity checks).
	storageKind = "kind"
)

type hostPath struct {
	config Config

	// gRPC calls involving any of the fields below must be serialized
	// by locking this mutex before starting. Internal helper
	// functions assume that the mutex has been locked.
	mutex     sync.Mutex
	volumes   map[string]hostPathVolume
	snapshots map[string]hostPathSnapshot
}

type hostPathVolume struct {
	VolName        string     `json:"volName"`
	VolID          string     `json:"volID"`
	VolSize        int64      `json:"volSize"`
	VolPath        string     `json:"volPath"`
	VolAccessType  accessType `json:"volAccessType"`
	ParentVolID    string     `json:"parentVolID,omitempty"`
	ParentSnapID   string     `json:"parentSnapID,omitempty"`
	Ephemeral      bool       `json:"ephemeral"`
	NodeID         string     `json:"nodeID"`
	Kind           string     `json:"kind"`
	ReadOnlyAttach bool       `json:"readOnlyAttach"`
	IsAttached     bool       `json:"isAttached"`
	IsStaged       bool       `json:"isStaged"`
	IsPublished    bool       `json:"isPublished"`
}

type hostPathSnapshot struct {
	Name         string               `json:"name"`
	Id           string               `json:"id"`
	VolID        string               `json:"volID"`
	Path         string               `json:"path"`
	CreationTime *timestamp.Timestamp `json:"creationTime"`
	SizeBytes    int64                `json:"sizeBytes"`
	ReadyToUse   bool                 `json:"readyToUse"`
}

type Config struct {
	DriverName        string
	Endpoint          string
	ProxyEndpoint     string
	NodeID            string
	VendorVersion     string
	MaxVolumesPerNode int64
	MaxVolumeSize     int64
	Capacity          Capacity
	Ephemeral         bool
	ShowVersion       bool
	EnableAttach      bool
}

var (
	vendorVersion = "dev"
)

const (
	// Directory where data for volumes and snapshots are persisted.
	// This can be ephemeral within the container or persisted if
	// backed by a Pod volume.
	dataRoot = "/csi-data-dir"

	// Extension with which snapshot files will be saved.
	snapshotExt = ".snap"
)

func NewHostPathDriver(cfg Config) (*hostPath, error) {
	if cfg.DriverName == "" {
		return nil, errors.New("no driver name provided")
	}

	if cfg.NodeID == "" {
		return nil, errors.New("no node id provided")
	}

	if cfg.Endpoint == "" {
		return nil, errors.New("no driver endpoint provided")
	}

	if err := os.MkdirAll(dataRoot, 0750); err != nil {
		return nil, fmt.Errorf("failed to create dataRoot: %v", err)
	}

	glog.Infof("Driver: %v ", cfg.DriverName)
	glog.Infof("Version: %s", cfg.VendorVersion)

	hp := &hostPath{
		config:    cfg,
		volumes:   map[string]hostPathVolume{},
		snapshots: map[string]hostPathSnapshot{},
	}
	if err := hp.discoveryExistingVolumes(); err != nil {
		return nil, err
	}
	hp.discoverExistingSnapshots()
	return hp, nil
}

func getSnapshotID(file string) (bool, string) {
	glog.V(4).Infof("file: %s", file)
	// Files with .snap extension are volumesnapshot files.
	// e.g. foo.snap, foo.bar.snap
	if filepath.Ext(file) == snapshotExt {
		return true, strings.TrimSuffix(file, snapshotExt)
	}
	return false, ""
}

func (h *hostPath) discoverExistingSnapshots() {
	glog.V(4).Infof("discovering existing snapshots in %s", dataRoot)
	files, err := ioutil.ReadDir(dataRoot)
	if err != nil {
		glog.Errorf("failed to discover snapshots under %s: %v", dataRoot, err)
	}
	for _, file := range files {
		isSnapshot, snapshotID := getSnapshotID(file.Name())
		if isSnapshot {
			glog.V(4).Infof("adding snapshot %s from file %s", snapshotID, getSnapshotPath(snapshotID))
			h.snapshots[snapshotID] = hostPathSnapshot{
				Id:         snapshotID,
				Path:       getSnapshotPath(snapshotID),
				ReadyToUse: true,
			}
		}
	}
}

func (hp *hostPath) discoveryExistingVolumes() error {
	cmdPath, err := exec.LookPath("findmnt")
	if err != nil {
		return fmt.Errorf("findmnt not found: %w", err)
	}

	out, err := exec.Command(cmdPath, "--json").CombinedOutput()
	if err != nil {
		glog.V(3).Infof("failed to execute command: %+v", cmdPath)
		return err
	}

	if len(out) < 1 {
		return fmt.Errorf("mount point info is nil")
	}

	mountInfos, err := parseMountInfo([]byte(out))
	if err != nil {
		return fmt.Errorf("failed to parse the mount infos: %+v", err)
	}

	mountInfosOfPod := MountPointInfo{}
	for _, mountInfo := range mountInfos {
		if mountInfo.Target == podVolumeTargetPath {
			mountInfosOfPod = mountInfo
			break
		}
	}

	// getting existing volumes based on the mount point infos.
	// It's a temporary solution to recall volumes.
	// TODO: discover what kind of storage was used and the nominal size.
	for _, pv := range mountInfosOfPod.ContainerFileSystem {
		if !strings.Contains(pv.Target, csiSignOfVolumeTargetPath) {
			continue
		}

		hpv, err := parseVolumeInfo(pv)
		if err != nil {
			return err
		}

		if hpv.Kind != "" && hp.config.Capacity.Enabled() {
			if _, err := hp.config.Capacity.Alloc(hpv.Kind, hpv.VolSize); err != nil {
				return fmt.Errorf("existing volume(s) do not match new capacity configuration: %v", err)
			}
		}
		hp.volumes[hpv.VolID] = *hpv
	}

	glog.V(4).Infof("Existing Volumes: %+v", hp.volumes)
	return nil
}

func (hp *hostPath) Run() error {
	s := NewNonBlockingGRPCServer()
	// hp itself implements ControllerServer, NodeServer, and IdentityServer.
	s.Start(hp.config.Endpoint, hp, hp, hp)
	s.Wait()

	return nil
}

func (hp *hostPath) getVolumeByID(volumeID string) (hostPathVolume, error) {
	if hostPathVol, ok := hp.volumes[volumeID]; ok {
		return hostPathVol, nil
	}
	return hostPathVolume{}, status.Errorf(codes.NotFound, "volume id %s does not exist in the volumes list", volumeID)
}

func (hp *hostPath) getVolumeByName(volName string) (hostPathVolume, error) {
	for _, hostPathVol := range hp.volumes {
		if hostPathVol.VolName == volName {
			return hostPathVol, nil
		}
	}
	return hostPathVolume{}, status.Errorf(codes.NotFound, "volume name %s does not exist in the volumes list", volName)
}

func (hp *hostPath) getSnapshotByName(name string) (hostPathSnapshot, error) {
	for _, snapshot := range hp.snapshots {
		if snapshot.Name == name {
			return snapshot, nil
		}
	}
	return hostPathSnapshot{}, status.Errorf(codes.NotFound, "snapshot name %s does not exist in the snapshots list", name)
}

// getVolumePath returns the canonical path for hostpath volume
func getVolumePath(volID string) string {
	return filepath.Join(dataRoot, volID)
}

// createVolume allocates capacity, creates the directory for the hostpath volume, and
// adds the volume to the list.
//
// It returns the volume path or err if one occurs. That error is suitable as result of a gRPC call.
func (hp *hostPath) createVolume(volID, name string, cap int64, volAccessType accessType, ephemeral bool, kind string) (hpv *hostPathVolume, finalErr error) {
	// Check for maximum available capacity
	if cap > hp.config.MaxVolumeSize {
		return nil, status.Errorf(codes.OutOfRange, "Requested capacity %d exceeds maximum allowed %d", cap, hp.config.MaxVolumeSize)
	}
	if hp.config.Capacity.Enabled() {
		actualKind, err := hp.config.Capacity.Alloc(kind, cap)
		if err != nil {
			return nil, err
		}
		// Free the capacity in case of any error - either a volume gets created or it doesn't.
		defer func() {
			if finalErr != nil {
				hp.config.Capacity.Free(actualKind, cap)
			}
		}()
		kind = actualKind
	} else if kind != "" {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("capacity tracking disabled, specifying kind %q is invalid", kind))
	}

	path := getVolumePath(volID)

	switch volAccessType {
	case mountAccess:
		err := os.MkdirAll(path, 0777)
		if err != nil {
			return nil, err
		}
	case blockAccess:
		executor := utilexec.New()
		size := fmt.Sprintf("%dM", cap/mib)
		// Create a block file.
		_, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				out, err := executor.Command("fallocate", "-l", size, path).CombinedOutput()
				if err != nil {
					return nil, fmt.Errorf("failed to create block device: %v, %v", err, string(out))
				}
			} else {
				return nil, fmt.Errorf("failed to stat block device: %v, %v", path, err)
			}
		}

		// Associate block file with the loop device.
		volPathHandler := volumepathhandler.VolumePathHandler{}
		_, err = volPathHandler.AttachFileDevice(path)
		if err != nil {
			// Remove the block file because it'll no longer be used again.
			if err2 := os.Remove(path); err2 != nil {
				glog.Errorf("failed to cleanup block file %s: %v", path, err2)
			}
			return nil, fmt.Errorf("failed to attach device %v: %v", path, err)
		}
	default:
		return nil, fmt.Errorf("unsupported access type %v", volAccessType)
	}

	hostpathVol := hostPathVolume{
		VolID:         volID,
		VolName:       name,
		VolSize:       cap,
		VolPath:       path,
		VolAccessType: volAccessType,
		Ephemeral:     ephemeral,
		Kind:          kind,
	}
	glog.V(4).Infof("adding hostpath volume: %s = %+v", volID, hostpathVol)
	hp.volumes[volID] = hostpathVol
	return &hostpathVol, nil
}

// updateVolume updates the existing hostpath volume.
func (hp *hostPath) updateVolume(volID string, volume hostPathVolume) error {
	glog.V(4).Infof("updating hostpath volume: %s", volID)

	if _, err := hp.getVolumeByID(volID); err != nil {
		return err
	}

	hp.volumes[volID] = volume
	return nil
}

// deleteVolume deletes the directory for the hostpath volume.
func (hp *hostPath) deleteVolume(volID string) error {
	glog.V(4).Infof("starting to delete hostpath volume: %s", volID)

	vol, err := hp.getVolumeByID(volID)
	if err != nil {
		// Return OK if the volume is not found.
		return nil
	}

	if vol.VolAccessType == blockAccess {
		volPathHandler := volumepathhandler.VolumePathHandler{}
		path := getVolumePath(volID)
		glog.V(4).Infof("deleting loop device for file %s if it exists", path)
		if err := volPathHandler.DetachFileDevice(path); err != nil {
			return fmt.Errorf("failed to remove loop device for file %s: %v", path, err)
		}
	}

	path := getVolumePath(volID)
	if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	if hp.config.Capacity.Enabled() {
		hp.config.Capacity.Free(vol.Kind, vol.VolSize)
	}
	delete(hp.volumes, volID)
	glog.V(4).Infof("deleted hostpath volume: %s = %+v", volID, vol)
	return nil
}

// hostPathIsEmpty is a simple check to determine if the specified hostpath directory
// is empty or not.
func hostPathIsEmpty(p string) (bool, error) {
	f, err := os.Open(p)
	if err != nil {
		return true, fmt.Errorf("unable to open hostpath volume, error: %v", err)
	}
	defer f.Close()

	_, err = f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// loadFromSnapshot populates the given destPath with data from the snapshotID
func (hp *hostPath) loadFromSnapshot(size int64, snapshotId, destPath string, mode accessType) error {
	snapshot, ok := hp.snapshots[snapshotId]
	if !ok {
		return status.Errorf(codes.NotFound, "cannot find snapshot %v", snapshotId)
	}
	if snapshot.ReadyToUse != true {
		return fmt.Errorf("snapshot %v is not yet ready to use", snapshotId)
	}
	if snapshot.SizeBytes > size {
		return status.Errorf(codes.InvalidArgument, "snapshot %v size %v is greater than requested volume size %v", snapshotId, snapshot.SizeBytes, size)
	}
	snapshotPath := snapshot.Path

	var cmd []string
	switch mode {
	case mountAccess:
		cmd = []string{"tar", "zxvf", snapshotPath, "-C", destPath}
	case blockAccess:
		cmd = []string{"dd", "if=" + snapshotPath, "of=" + destPath}
	default:
		return status.Errorf(codes.InvalidArgument, "unknown accessType: %d", mode)
	}

	executor := utilexec.New()
	glog.V(4).Infof("Command Start: %v", cmd)
	out, err := executor.Command(cmd[0], cmd[1:]...).CombinedOutput()
	glog.V(4).Infof("Command Finish: %v", string(out))
	if err != nil {
		return fmt.Errorf("failed pre-populate data from snapshot %v: %w: %s", snapshotId, err, out)
	}
	return nil
}

// loadFromVolume populates the given destPath with data from the srcVolumeID
func (hp *hostPath) loadFromVolume(size int64, srcVolumeId, destPath string, mode accessType) error {
	hostPathVolume, ok := hp.volumes[srcVolumeId]
	if !ok {
		return status.Error(codes.NotFound, "source volumeId does not exist, are source/destination in the same storage class?")
	}
	if hostPathVolume.VolSize > size {
		return status.Errorf(codes.InvalidArgument, "volume %v size %v is greater than requested volume size %v", srcVolumeId, hostPathVolume.VolSize, size)
	}
	if mode != hostPathVolume.VolAccessType {
		return status.Errorf(codes.InvalidArgument, "volume %v mode is not compatible with requested mode", srcVolumeId)
	}

	switch mode {
	case mountAccess:
		return loadFromFilesystemVolume(hostPathVolume, destPath)
	case blockAccess:
		return loadFromBlockVolume(hostPathVolume, destPath)
	default:
		return status.Errorf(codes.InvalidArgument, "unknown accessType: %d", mode)
	}
}

func loadFromFilesystemVolume(hostPathVolume hostPathVolume, destPath string) error {
	srcPath := hostPathVolume.VolPath
	isEmpty, err := hostPathIsEmpty(srcPath)
	if err != nil {
		return fmt.Errorf("failed verification check of source hostpath volume %v: %w", hostPathVolume.VolID, err)
	}

	// If the source hostpath volume is empty it's a noop and we just move along, otherwise the cp call will fail with a a file stat error DNE
	if !isEmpty {
		args := []string{"-a", srcPath + "/.", destPath + "/"}
		executor := utilexec.New()
		out, err := executor.Command("cp", args...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed pre-populate data from volume %v: %s: %w", hostPathVolume.VolID, out, err)
		}
	}
	return nil
}

func loadFromBlockVolume(hostPathVolume hostPathVolume, destPath string) error {
	srcPath := hostPathVolume.VolPath
	args := []string{"if=" + srcPath, "of=" + destPath}
	executor := utilexec.New()
	out, err := executor.Command("dd", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed pre-populate data from volume %v: %w: %s", hostPathVolume.VolID, err, out)
	}
	return nil
}

func (hp *hostPath) getSortedVolumeIDs() []string {
	ids := make([]string, len(hp.volumes))
	index := 0
	for volId := range hp.volumes {
		ids[index] = volId
		index += 1
	}

	sort.Strings(ids)
	return ids
}

func filterVolumeName(targetPath string) string {
	pathItems := strings.Split(targetPath, "kubernetes.io~csi/")
	if len(pathItems) < 2 {
		return ""
	}

	return strings.TrimSuffix(pathItems[1], "/mount")
}

func filterVolumeID(sourcePath string) string {
	volumeSourcePathRegex := regexp.MustCompile(`\[(.*)\]`)
	volumeSP := string(volumeSourcePathRegex.Find([]byte(sourcePath)))
	if volumeSP == "" {
		return ""
	}

	return strings.TrimSuffix(strings.TrimPrefix(volumeSP, "[/var/lib/csi-hostpath-data/"), "]")
}

func parseVolumeInfo(volume MountPointInfo) (*hostPathVolume, error) {
	volumeName := filterVolumeName(volume.Target)
	volumeID := filterVolumeID(volume.Source)
	sourcePath := getSourcePath(volumeID)
	_, fscapacity, _, _, _, _, err := fs.FsInfo(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get capacity info: %+v", err)
	}

	hp := hostPathVolume{
		VolName:       volumeName,
		VolID:         volumeID,
		VolSize:       fscapacity,
		VolPath:       getVolumePath(volumeID),
		VolAccessType: mountAccess,
	}

	return &hp, nil
}
