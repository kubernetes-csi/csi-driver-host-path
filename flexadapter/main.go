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

package main

import (
	"fmt"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/kubernetes-csi/drivers/lib"
)

type flexAdapter struct {
	driver *lib.CSIDriver

	flexDriver *flexVolumeDriver

	cap   []*csi.VolumeCapability_AccessMode
	cscap []*csi.ControllerServiceCapability
}

var (
	adapter    flexAdapter
	endpoint   string
	driverName string
	driverPath string
	nodeID     string
	version    = csi.Version{
		Minor: 1,
	}
)

func GetSupportedVersions() []*csi.Version {
	return []*csi.Version{&version}
}

func main() {
	cmd := &cobra.Command{
		Use:   "flexadapter",
		Short: "Flex volume adapter for CSI",
		Run: func(cmd *cobra.Command, args []string) {
			handle()
		},
	}

	cmd.PersistentFlags().StringVar(&nodeID, "nodeid", "", "node id")
	cmd.MarkPersistentFlagRequired("nodeid")

	cmd.PersistentFlags().StringVar(&endpoint, "endpoint", "", "CSI endpoint")
	cmd.MarkPersistentFlagRequired("endpoint")

	cmd.PersistentFlags().StringVar(&driverPath, "driverpath", "", "path to flexvolume driver path")
	cmd.MarkPersistentFlagRequired("driverpath")

	cmd.PersistentFlags().StringVar(&driverName, "drivername", "", "name of the driver")
	cmd.MarkPersistentFlagRequired("drivername")

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func handle() {

	var err error

	adapter.flexDriver, err = NewFlexVolumeDriver(driverPath)
	if err != nil {
		glog.Fatalf("Failed to initialize flex volume driver, error: %v", err.Error())
	}

	glog.Infof("Driver: %v version: %v", driverName, GetVersionString(&version))

	adapter.driver = lib.NewCSIDriver(driverName, &version, GetSupportedVersions(), nodeID)
	ids := &identityServer{
		IdentityServerDefaults: lib.IdentityServerDefaults{
			Driver: adapter.driver,
		}}
	ns := &nodeServer{
		NodeServerDefaults: lib.NodeServerDefaults{
			Driver: adapter.driver,
		},
	}
	cs := &controllerServer{
		ControllerServerDefaults: lib.ControllerServerDefaults{
			Driver: adapter.driver,
		}}

	if adapter.flexDriver.capabilities.Attach {
		adapter.driver.AddControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME})
	}
	adapter.driver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER})

	lib.Serve(endpoint, ids, cs, ns)
}
