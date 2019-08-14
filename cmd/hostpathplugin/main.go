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
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/kubernetes-csi/csi-driver-host-path/pkg/hostpath"
)

func init() {
	flag.Set("logtostderr", "true")
}

var (
	endpoint    = flag.String("endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	driverName  = flag.String("drivername", "hostpath.csi.k8s.io", "name of the driver")
	nodeID      = flag.String("nodeid", "", "node id")
	ephemeral   = flag.Bool("ephemeral", false, "publish volumes in ephemeral mode even if kubelet did not ask for it (only needed for Kubernetes 1.15)")
	showVersion = flag.Bool("version", false, "Show version.")
	// Set by the build process
	version    = ""
	controller = flag.Bool("controller", false, "Open ControllerPublishVolume and ControllerUnpublishVolume")
)

func main() {
	flag.Parse()

	if *showVersion {
		baseName := path.Base(os.Args[0])
		fmt.Println(baseName, version)
		return
	}

	if *ephemeral {
		fmt.Fprintln(os.Stderr, "Deprecation warning: The ephemeral flag is deprecated and should only be used when deploying on Kubernetes 1.15. It will be removed in the future.")
	}

	handle()
	os.Exit(0)
}

func handle() {
	driver, err := hostpath.NewHostPathDriver(*driverName, *nodeID, *endpoint, version, *ephemeral, *controller)
	if err != nil {
		fmt.Printf("Failed to initialize driver: %s", err.Error())
		os.Exit(1)
	}
	driver.Run()
}
