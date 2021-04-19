/*
Copyright 2021 The Kubernetes Authors.

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
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/golang/glog"
	fs "k8s.io/kubernetes/pkg/volume/util/fs"
)

const (
	podVolumeTargetPath       = "/var/lib/kubelet/pods"
	csiSignOfVolumeTargetPath = "kubernetes.io~csi/pvc"
)

type MountPointInfo struct {
	Target              string           `json:"target"`
	Source              string           `json:"source"`
	FsType              string           `json:"fstype"`
	Options             string           `json:"options"`
	ContainerFileSystem []MountPointInfo `json:"children,omitempty"`
}

type ContainerFileSystem struct {
	Children []MountPointInfo `json:"children"`
}

type FileSystems struct {
	Filsystem []ContainerFileSystem `json:"filesystems"`
}

func getSourcePath(volumeHandle string) string {
	return fmt.Sprintf("%s/%s", dataRoot, volumeHandle)
}

func checkSourcePathExist(volumeHandle string) (bool, error) {
	sourcePath := getSourcePath(volumeHandle)
	glog.V(3).Infof("Volume: %s Source path is: %s", volumeHandle, sourcePath)
	_, err := os.Stat(sourcePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func parseMountInfo(originalMountInfo []byte) ([]MountPointInfo, error) {
	fs := FileSystems{
		Filsystem: make([]ContainerFileSystem, 0),
	}

	if err := json.Unmarshal(originalMountInfo, &fs); err != nil {
		return nil, err
	}

	if len(fs.Filsystem) <= 0 {
		return nil, fmt.Errorf("failed to get mount info")
	}

	return fs.Filsystem[0].Children, nil
}

func checkMountPointExist(sourcePath string) (bool, error) {
	cmdPath, err := exec.LookPath("findmnt")
	if err != nil {
		return false, fmt.Errorf("findmnt not found: %w", err)
	}

	out, err := exec.Command(cmdPath, "--json").CombinedOutput()
	if err != nil {
		glog.V(3).Infof("failed to execute command: %+v", cmdPath)
		return false, err
	}

	if len(out) < 1 {
		return false, fmt.Errorf("mount point info is nil")
	}

	mountInfos, err := parseMountInfo([]byte(out))
	if err != nil {
		return false, fmt.Errorf("failed to parse the mount infos: %+v", err)
	}

	mountInfosOfPod := MountPointInfo{}
	for _, mountInfo := range mountInfos {
		if mountInfo.Target == podVolumeTargetPath {
			mountInfosOfPod = mountInfo
			break
		}
	}

	for _, mountInfo := range mountInfosOfPod.ContainerFileSystem {
		if !strings.Contains(mountInfo.Source, sourcePath) {
			continue
		}

		_, err = os.Stat(mountInfo.Target)
		if err != nil {
			if os.IsNotExist(err) {
				return false, nil
			}

			return false, err
		}

		return true, nil
	}

	return false, nil
}

func (hp *hostPath) checkPVCapacityValid(volumeHandle string) (bool, error) {
	sourcePath := getSourcePath(volumeHandle)
	_, fscapacity, _, _, _, _, err := fs.FsInfo(sourcePath)
	if err != nil {
		return false, fmt.Errorf("failed to get capacity info: %+v", err)
	}

	volume, err := hp.hostPathDriverState.GetVolumeByID(volumeHandle)
	if err != nil {
		return false, err
	}
	volumeCapacity := volume.VolSize
	glog.V(3).Infof("volume capacity: %+v fs capacity:%+v", volumeCapacity, fscapacity)
	return fscapacity >= volumeCapacity, nil
}

func getPVStats(volumePath string) (available int64, capacity int64, used int64, inodes int64, inodesFree int64, inodesUsed int64, err error) {
	return fs.FsInfo(volumePath)
}

func checkPVUsage(volumeHandle string) (bool, error) {
	sourcePath := getSourcePath(volumeHandle)
	fsavailable, _, _, _, _, _, err := fs.FsInfo(sourcePath)
	if err != nil {
		return false, err
	}

	glog.V(3).Infof("fs available: %+v", fsavailable)
	return fsavailable > 0, nil
}

func (hp *hostPath) doHealthCheckInControllerSide(volumeHandle string) (bool, string) {
	spExist, err := checkSourcePathExist(volumeHandle)
	if err != nil {
		return false, err.Error()
	}

	if !spExist {
		return false, "The source path of the volume doesn't exist"
	}

	capValid, err := hp.checkPVCapacityValid(volumeHandle)
	if err != nil {
		return false, err.Error()
	}

	if !capValid {
		return false, "The capacity of volume is greater than actual storage"
	}

	available, err := checkPVUsage(volumeHandle)
	if err != nil {
		return false, err.Error()
	}

	if !available {
		return false, "The free space of the volume is insufficient"
	}

	return true, ""
}

func doHealthCheckInNodeSide(volumeHandle string) (bool, string) {
	mpExist, err := checkMountPointExist(volumeHandle)
	if err != nil {
		return false, err.Error()
	}

	if !mpExist {
		return false, "The volume isn't mounted"
	}

	return true, ""
}
