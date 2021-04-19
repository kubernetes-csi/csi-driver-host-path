package state

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/kubernetes/pkg/volume/util/volumepathhandler"
	utilexec "k8s.io/utils/exec"
)

func (hps *HostPathDriverState) ListVolumes() (HostPathVolumes, error) {
	return hps.HostPathVolumes, nil
}

func (hps *HostPathDriverState) GetVolumeByID(volID string) (HostPathVolume, error) {
	hpv, ok := hps.HostPathVolumes[volID]
	if !ok {
		return HostPathVolume{}, status.Errorf(codes.NotFound, "volume id %s does not exist in the volumes list", volID)

	}

	return hpv, nil
}

func (hps *HostPathDriverState) GetVolumeByName(volName string) (HostPathVolume, error) {
	for _, volume := range hps.HostPathVolumes {
		if volume.VolName == volName {
			return volume, nil
		}
	}

	return HostPathVolume{}, status.Errorf(codes.NotFound, "volume name %s does not exist in the volumes list", volName)
}

// CreateVolume allocates capacity, creates the directory for the hostpath volume, and
// adds the volume to the list.
//
// It returns the volume path or err if one occurs. That error is suitable as result of a gRPC call.
func (hps *HostPathDriverState) CreateVolume(volID, name string, cap int64, volAccessType AccessType, ephemeral bool, kind string, maxVolumeSize int64, capacity Capacity) (hpv *HostPathVolume, finalErr error) {
	// Check for maximum available capacity
	if cap > maxVolumeSize {
		return nil, status.Errorf(codes.OutOfRange, "Requested capacity %d exceeds maximum allowed %d", cap, maxVolumeSize)
	}
	if capacity.Enabled() {
		if kind == "" {
			// Pick some kind with sufficient remaining capacity.
			for k, c := range capacity {
				if hps.SumVolumeSizes(k)+cap <= c.Value() {
					kind = k
					break
				}
			}
		}
		if kind == "" {
			// Still nothing?!
			return nil, status.Errorf(codes.OutOfRange, "requested capacity %d of arbitrary storage exceeds all remaining capacity", cap)
		}
		used := hps.SumVolumeSizes(kind)
		available := capacity[kind]
		if used+cap > available.Value() {

			return nil, status.Errorf(codes.OutOfRange, "requested capacity %d exceeds remaining capacity for %q, %s out of %s already used",
				cap, kind, resource.NewQuantity(used, resource.BinarySI).String(), available.String())
		}
	} else if kind != "" {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("capacity tracking disabled, specifying kind %q is invalid", kind))
	}

	path := GetVolumePath(volID, hps.DataRoot)

	switch volAccessType {
	case MountAccess:
		err := os.MkdirAll(path, 0777)
		if err != nil {
			return nil, err
		}
	case BlockAccess:
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

	hostpathVol := HostPathVolume{
		VolID:         volID,
		VolName:       name,
		VolSize:       cap,
		VolPath:       path,
		VolAccessType: volAccessType,
		Ephemeral:     ephemeral,
		Kind:          kind,
	}

	glog.V(4).Infof("adding hostpath volume: %s = %+v", volID, hostpathVol)
	hps.HostPathVolumes[volID] = hostpathVol
	if err := hps.flushVolumesToFile(); err != nil {
		return nil, fmt.Errorf("failed to store volume data into local file: %s. Error: %s", hps.VolumeDataFilePath, err)

	}
	return &hostpathVol, nil
}

func (hps *HostPathDriverState) UpdateVolume(volID string, volume HostPathVolume) error {
	glog.V(4).Infof("updating hostpath volume: %s", volID)

	if _, err := hps.GetVolumeByID(volID); err != nil {
		return err
	}

	hps.HostPathVolumes[volID] = volume
	if err := hps.flushVolumesToFile(); err != nil {
		return fmt.Errorf("failed to update volume data into local file: %s. Error: %s", hps.VolumeDataFilePath, err)

	}

	return nil
}

// deleteVolume deletes the directory for the hostpath volume.
func (hps *HostPathDriverState) DeleteVolume(volID string, capacity Capacity) error {
	glog.V(4).Infof("starting to delete hostpath volume: %s", volID)

	vol, err := hps.GetVolumeByID(volID)
	if err != nil {
		// Return OK if the volume is not found.
		return nil
	}

	if vol.VolAccessType == BlockAccess {
		volPathHandler := volumepathhandler.VolumePathHandler{}
		path := GetVolumePath(volID, hps.DataRoot)
		glog.V(4).Infof("deleting loop device for file %s if it exists", path)
		if err := volPathHandler.DetachFileDevice(path); err != nil {
			return fmt.Errorf("failed to remove loop device for file %s: %v", path, err)
		}
	}

	path := GetVolumePath(volID, hps.DataRoot)
	if err := os.RemoveAll(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	delete(hps.HostPathVolumes, volID)
	if err := hps.flushVolumesToFile(); err != nil {
		return fmt.Errorf("failed to update volume data into local file: %s. Error: %s", hps.VolumeDataFilePath, err)

	}
	glog.V(4).Infof("deleted hostpath volume: %s = %+v", volID, vol)
	return nil
}

func (hps *HostPathDriverState) GetVolumeLocker() *sync.RWMutex {
	return hps.VolumesFileRWLock
}

func (hps *HostPathDriverState) SumVolumeSizes(kind string) (sum int64) {
	for _, volume := range hps.HostPathVolumes {
		if volume.Kind == kind {
			sum += volume.VolSize
		}
	}
	return
}

// GetVolumePath returns the canonical path for hostpath volume
func GetVolumePath(volID string, dataRoot string) string {
	return filepath.Join(dataRoot, volID)
}

func (hps *HostPathDriverState) GetSortedVolumeIDs() []string {
	ids := make([]string, len(hps.HostPathVolumes))
	index := 0
	for volId := range hps.HostPathVolumes {
		ids[index] = volId
		index += 1
	}

	sort.Strings(ids)
	return ids
}

func (hps *HostPathDriverState) flushVolumesToFile() error {
	err := os.MkdirAll(hps.DataRoot, 0777)
	if err != nil {
		return err
	}

	data, err := json.Marshal(hps.HostPathVolumes)
	if err != nil {
		glog.Errorf("failed to unmarshal existing volumes: %s", err)
		return err
	}

	_, err = os.Stat(hps.VolumeDataFilePath)
	if err != nil {
		_, err = os.Create(hps.VolumeDataFilePath)
		if err != nil {
			glog.Errorf("failed to create volume data file: %s", err)
			a := err.Error()
			fmt.Println(a)
			return err
		}
	}

	err = ioutil.WriteFile(hps.VolumeDataFilePath, data, 0644)
	if err != nil {
		glog.Errorf("failed to discover existing volumes under %s: %v", hps.VolumeDataFilePath, err)
		return err
	}

	glog.V(4).Info("discover existing volumes successfully")
	return nil
}

func (hps *HostPathDriverState) ListSnapshots() (HostPathSnapshots, error) {
	return hps.HostPathSnapshots, nil
}

func (hps *HostPathDriverState) GetSnapshotByName(name string) (HostPathSnapshot, error) {
	for _, snapshot := range hps.HostPathSnapshots {
		if snapshot.Name == name {
			return snapshot, nil
		}
	}
	return HostPathSnapshot{}, status.Errorf(codes.NotFound, "snapshot name %s does not exist in the snapshots list", name)
}

func (hps *HostPathDriverState) GetSnapshotByID(id string) (HostPathSnapshot, error) {
	snapshot, ok := hps.HostPathSnapshots[id]
	if !ok {
		return HostPathSnapshot{}, status.Errorf(codes.NotFound, "snapshot id %s does not exist in the snapshots list", id)
	}

	return snapshot, nil
}

func (hps *HostPathDriverState) CreateSnapshot(name, snapshotId, volumeId, snapshotFilePath string, creationTime *timestamp.Timestamp, size int64, readyToUse bool) (HostPathSnapshot, error) {
	snapshot := HostPathSnapshot{}
	snapshot.Name = name
	snapshot.Id = snapshotId
	snapshot.VolID = volumeId
	snapshot.Path = snapshotFilePath
	snapshot.CreationTime = creationTime
	snapshot.SizeBytes = size
	snapshot.ReadyToUse = true

	hps.HostPathSnapshots[snapshotId] = snapshot
	return snapshot, nil
}

func (hps *HostPathDriverState) DeleteSnapshot(snapshotId string) error {
	// Lock before acting on global state. A production-quality
	// driver might use more fine-grained locking.
	hps.SnapshotsFileRWLock.Lock()
	defer hps.SnapshotsFileRWLock.Unlock()

	glog.V(4).Infof("deleting snapshot %s", snapshotId)
	path := getSnapshotPath(snapshotId, hps.DataRoot)
	os.RemoveAll(path)
	delete(hps.HostPathSnapshots, snapshotId)
	return nil
}

func (hps *HostPathDriverState) GetSnapshotLocker() *sync.RWMutex {
	return hps.SnapshotsFileRWLock
}

func getSnapshotID(file string) (bool, string) {
	glog.V(4).Infof("file: %s", file)
	// Files with .snap extension are volumesnapshot files.
	// e.g. foo.snap, foo.bar.snap
	if filepath.Ext(file) == SnapshotExt {
		return true, strings.TrimSuffix(file, SnapshotExt)
	}
	return false, ""
}

// getSnapshotPath returns the full path to where the snapshot is stored
func getSnapshotPath(snapshotID, dataRoot string) string {
	return filepath.Join(dataRoot, fmt.Sprintf("%s%s", snapshotID, SnapshotExt))
}

func LoadStateFromFile(dataRoot string) (*HostPathDriverState, error) {
	driverState := NewHostPathDriverState(dataRoot)

	if err := driverState.loadVolumesFromFile(); err != nil {
		return nil, err
	}

	if err := driverState.loadSnapshotsFromFile(); err != nil {
		return nil, err
	}

	return driverState, nil
}

func (hps *HostPathDriverState) loadVolumesFromFile() error {
	hps.VolumesFileRWLock.RLock()
	defer hps.VolumesFileRWLock.RUnlock()
	glog.V(4).Infof("discovering existing volume data in %s", hps.VolumeDataFilePath)
	data, err := ioutil.ReadFile(hps.VolumeDataFilePath)
	if err != nil {
		glog.Errorf("failed to discover existing volumes under %s: %v", hps.VolumeDataFilePath, err)
		return err
	}

	volumes := make(HostPathVolumes)
	err = json.Unmarshal(data, &volumes)
	if err != nil {
		glog.Errorf("failed to unmarshal existing volumes: %s", err)
		return err
	}

	hps.HostPathVolumes = volumes
	glog.V(4).Info("discover existing volumes successfully")
	return nil
}

func (hps *HostPathDriverState) loadSnapshotsFromFile() error {
	glog.V(4).Infof("discovering existing snapshots in %s", hps.DataRoot)
	files, err := ioutil.ReadDir(hps.DataRoot)
	if err != nil {
		glog.Errorf("failed to discover snapshots under %s: %v", hps.DataRoot, err)
		return err
	}

	for _, file := range files {
		isSnapshot, snapshotID := getSnapshotID(file.Name())
		if isSnapshot {
			glog.V(4).Infof("adding snapshot %s from file %s", snapshotID, getSnapshotPath(hps.DataRoot, snapshotID))
			hps.HostPathSnapshots[snapshotID] = HostPathSnapshot{
				Id:         snapshotID,
				Path:       getSnapshotPath(hps.DataRoot, snapshotID),
				ReadyToUse: true,
			}
		}
	}

	return nil
}

// Capacity simulates linear storage of certain types ("fast",
// "slow"). To calculate the amount of allocated space, the size of
// all currently existing volumes of the same kind is summed up.
//
// Available capacity is configurable with a command line flag
// -capacity <type>=<size> where <type> is a string and <size>
// is a quantity (1T, 1Gi). More than one of those
// flags can be used.
//
// The underlying map will be initialized if needed by Set,
// which makes it possible to define and use a Capacity instance
// without explicit initialization (`var capacity Capacity` or as
// member in a struct).
type Capacity map[string]resource.Quantity

// Set is an implementation of flag.Value.Set.
func (c *Capacity) Set(arg string) error {
	parts := strings.SplitN(arg, "=", 2)
	if len(parts) != 2 {
		return errors.New("must be of format <type>=<size>")
	}
	quantity, err := resource.ParseQuantity(parts[1])
	if err != nil {
		return err
	}

	// We overwrite any previous value.
	if *c == nil {
		*c = Capacity{}
	}
	(*c)[parts[0]] = quantity
	return nil
}

func (c *Capacity) String() string {
	return fmt.Sprintf("%v", map[string]resource.Quantity(*c))
}

var _ flag.Value = &Capacity{}

// Enabled returns true if capacities are configured.
func (c *Capacity) Enabled() bool {
	return len(*c) > 0
}
