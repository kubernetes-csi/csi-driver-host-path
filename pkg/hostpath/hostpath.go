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
	"os"

	"github.com/golang/glog"
	"github.com/kubernetes-csi/csi-driver-host-path/pkg/state"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	utilexec "k8s.io/utils/exec"
)

type hostPath struct {
	config              Config
	hostPathDriverState state.DriverState
}

type Config struct {
	DriverName        string
	Endpoint          string
	ProxyEndpoint     string
	NodeID            string
	VendorVersion     string
	MaxVolumesPerNode int64
	MaxVolumeSize     int64
	Capacity          state.Capacity
	Ephemeral         bool
	ShowVersion       bool
	EnableAttach      bool
	DataRoot          string
}

var (
	vendorVersion = "dev"
)

const (
	// Directory where data for volumes and snapshots are persisted.
	// This can be ephemeral within the container or persisted if
	// backed by a Pod volume.
	dataRoot = "/csi-data-dir"
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

	// build up volumes and snapshot data in memory
	hostPathDriverState, err := state.LoadStateFromFile(cfg.DataRoot)
	if err != nil {
		return nil, err
	}

	hp := &hostPath{
		config:              cfg,
		hostPathDriverState: hostPathDriverState,
	}

	return hp, nil
}

func (hp *hostPath) Run() error {
	s := NewNonBlockingGRPCServer()
	// hp itself implements ControllerServer, NodeServer, and IdentityServer.
	s.Start(hp.config.Endpoint, hp, hp, hp)
	s.Wait()

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
func (hp *hostPath) loadFromSnapshot(size int64, snapshotId, destPath string, mode state.AccessType) error {
	snapshot, err := hp.hostPathDriverState.GetSnapshotByID(snapshotId)
	if err != nil {
		return err
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
	case state.MountAccess:
		cmd = []string{"tar", "zxvf", snapshotPath, "-C", destPath}
	case state.BlockAccess:
		cmd = []string{"dd", "if=" + snapshotPath, "of=" + destPath}
	default:
		return status.Errorf(codes.InvalidArgument, "unknown state.AccessType: %d", mode)
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
func (hp *hostPath) loadFromVolume(size int64, srcVolumeId, destPath string, mode state.AccessType) error {
	hostPathVolume, err := hp.hostPathDriverState.GetVolumeByID(srcVolumeId)
	if err != nil {
		return err
	}
	if hostPathVolume.VolSize > size {
		return status.Errorf(codes.InvalidArgument, "volume %v size %v is greater than requested volume size %v", srcVolumeId, hostPathVolume.VolSize, size)
	}

	if mode != hostPathVolume.VolAccessType {
		return status.Errorf(codes.InvalidArgument, "volume %v mode is not compatible with requested mode", srcVolumeId)
	}

	switch mode {
	case state.MountAccess:
		return loadFromFilesystemVolume(hostPathVolume, destPath)
	case state.BlockAccess:
		return loadFromBlockVolume(hostPathVolume, destPath)
	default:
		return status.Errorf(codes.InvalidArgument, "unknown state.AccessType: %d", mode)
	}
}

func loadFromFilesystemVolume(hostPathVolume state.HostPathVolume, destPath string) error {
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

func loadFromBlockVolume(hostPathVolume state.HostPathVolume, destPath string) error {
	srcPath := hostPathVolume.VolPath
	args := []string{"if=" + srcPath, "of=" + destPath}
	executor := utilexec.New()
	out, err := executor.Command("dd", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed pre-populate data from volume %v: %w: %s", hostPathVolume.VolID, err, out)
	}
	return nil
}
