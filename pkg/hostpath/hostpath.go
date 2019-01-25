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
	"fmt"

	"github.com/golang/glog"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

const (
	kib    int64 = 1024
	mib    int64 = kib * 1024
	gib    int64 = mib * 1024
	gib100 int64 = gib * 100
	tib    int64 = gib * 1024
	tib100 int64 = tib * 100
)

type hostPath struct {
	name     string
	nodeID   string
	version  string
	endpoint string

	ids *identityServer
	ns  *nodeServer
	cs  *controllerServer
}

type hostPathVolume struct {
	VolName       string     `json:"volName"`
	VolID         string     `json:"volID"`
	VolSize       int64      `json:"volSize"`
	VolPath       string     `json:"volPath"`
	VolAccessType accessType `json:"volAccessType"`
}

type hostPathSnapshot struct {
	Name         string              `json:"name"`
	Id           string              `json:"id"`
	VolID        string              `json:"volID"`
	Path         string              `json:"path"`
	CreationTime timestamp.Timestamp `json:"creationTime"`
	SizeBytes    int64               `json:"sizeBytes"`
	ReadyToUse   bool                `json:"readyToUse"`
}

var hostPathVolumes map[string]hostPathVolume
var hostPathVolumeSnapshots map[string]hostPathSnapshot

var (
	vendorVersion = "dev"
)

func init() {
	hostPathVolumes = map[string]hostPathVolume{}
	hostPathVolumeSnapshots = map[string]hostPathSnapshot{}
}

func NewHostPathDriver(driverName, nodeID, endpoint string) (*hostPath, error) {
	if driverName == "" {
		return nil, fmt.Errorf("No driver name provided")
	}

	if nodeID == "" {
		return nil, fmt.Errorf("No node id provided")
	}

	if endpoint == "" {
		return nil, fmt.Errorf("No driver endpoint provided")
	}

	glog.Infof("Driver: %v ", driverName)
	glog.Infof("Version: %s", vendorVersion)

	return &hostPath{
		name:     driverName,
		version:  vendorVersion,
		nodeID:   nodeID,
		endpoint: endpoint,
	}, nil
}

func (hp *hostPath) Run() {

	// Create GRPC servers
	hp.ids = NewIdentityServer(hp.name, hp.version)
	hp.ns = NewNodeServer(hp.nodeID)
	hp.cs = NewControllerServer()

	s := NewNonBlockingGRPCServer()
	s.Start(hp.endpoint, hp.ids, hp.cs, hp.ns)
	s.Wait()
}

func getVolumeByID(volumeID string) (hostPathVolume, error) {
	if hostPathVol, ok := hostPathVolumes[volumeID]; ok {
		return hostPathVol, nil
	}
	return hostPathVolume{}, fmt.Errorf("volume id %s does not exit in the volumes list", volumeID)
}

func getVolumeByName(volName string) (hostPathVolume, error) {
	for _, hostPathVol := range hostPathVolumes {
		if hostPathVol.VolName == volName {
			return hostPathVol, nil
		}
	}
	return hostPathVolume{}, fmt.Errorf("volume name %s does not exit in the volumes list", volName)
}

func getSnapshotByName(name string) (hostPathSnapshot, error) {
	for _, snapshot := range hostPathVolumeSnapshots {
		if snapshot.Name == name {
			return snapshot, nil
		}
	}
	return hostPathSnapshot{}, fmt.Errorf("snapshot name %s does not exit in the snapshots list", name)
}
