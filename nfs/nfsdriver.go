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

package nfs

import (
	"sync"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"

	"github.com/kubernetes-csi/drivers/csi-common"
)

type nfsDriver struct {
	csiDriver *csicommon.CSIDriver

	ids *identityServer
	ns  *nodeServer
	cs  *controllerServer

	cap   []*csi.VolumeCapability_AccessMode
	cscap []*csi.ControllerServiceCapability
}

const (
	driverName = "NFS"
)

var (
	driver  *nfsDriver
	runOnce sync.Once
	version = csi.Version{
		Minor: 1,
	}
)

func GetSupportedVersions() []*csi.Version {
	return []*csi.Version{&version}
}

func GetNFSDriver() *nfsDriver {
	runOnce.Do(func() {
		driver = &nfsDriver{}
	})
	return driver
}

func NewIdentityServer(d *csicommon.CSIDriver) *identityServer {
	return &identityServer{
		DefaultIdentityServer: csicommon.NewDefaultIdentityServer(d),
	}
}

func NewControllerServer(d *csicommon.CSIDriver) *controllerServer {
	return &controllerServer{
		DefaultControllerServer: csicommon.NewDefaultControllerServer(d),
	}
}

func NewNodeServer(d *csicommon.CSIDriver) *nodeServer {
	return &nodeServer{
		DefaultNodeServer: csicommon.NewDefaultNodeServer(d),
	}
}

func (f *nfsDriver) Run(driverPath, nodeID, endpoint string) {

	glog.Infof("Driver: %v version: %v", driverName, csicommon.GetVersionString(&version))

	// Initialize default library driver
	driver.csiDriver = csicommon.NewCSIDriver(driverName, &version, GetSupportedVersions(), nodeID)
	driver.csiDriver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER})

	// Create GRPC servers
	f.ids = NewIdentityServer(driver.csiDriver)
	f.ns = NewNodeServer(driver.csiDriver)
	f.cs = NewControllerServer(driver.csiDriver)

	csicommon.Serve(endpoint, f.ids, f.cs, f.ns)
}
