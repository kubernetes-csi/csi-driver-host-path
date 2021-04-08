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
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/golang/glog"
	"github.com/kubernetes-csi/csi-driver-host-path/internal/proxy"
	"github.com/kubernetes-csi/csi-driver-host-path/pkg/hostpath"
)

func init() {
	flag.Set("logtostderr", "true")
}

var (
	csiEndpoint       = flag.String("endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	driverName        = flag.String("drivername", "hostpath.csi.k8s.io", "name of the driver")
	nodeID            = flag.String("nodeid", "", "node id")
	ephemeral         = flag.Bool("ephemeral", false, "publish volumes in ephemeral mode even if kubelet did not ask for it (only needed for Kubernetes 1.15)")
	maxVolumesPerNode = flag.Int64("maxvolumespernode", 0, "limit of volumes per node")
	showVersion       = flag.Bool("version", false, "Show version.")
	capacity          = func() hostpath.Capacity {
		c := hostpath.Capacity{}
		flag.Var(&c, "capacity", "Simulate storage capacity. The parameter is <kind>=<quantity> where <kind> is the value of a 'kind' storage class parameter and <quantity> is the total amount of bytes for that kind. The flag may be used multiple times to configure different kinds.")
		return c
	}()
	enableAttach = flag.Bool("enable-attach", false, "Enables RPC_PUBLISH_UNPUBLISH_VOLUME capability.")
	// The proxy-endpoint option is intended to used by the Kubernetes E2E test suite
	// for proxying incoming calls to the embedded mock CSI driver.
	proxyEndpoint = flag.String("proxy-endpoint", "", "Instead of running the CSI driver code, just proxy connections from csiEndpoint to the given listening socket.")
	// Set by the build process
	version = ""
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

	if *proxyEndpoint != "" {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		closer, err := proxy.Run(ctx, *csiEndpoint, *proxyEndpoint)
		if err != nil {
			glog.Fatalf("failed to run proxy: %v", err)
		}
		defer closer.Close()

		// Wait for signal
		sigc := make(chan os.Signal, 1)
		sigs := []os.Signal{
			syscall.SIGTERM,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGQUIT,
		}
		signal.Notify(sigc, sigs...)

		<-sigc
		return
	}

	driver, err := hostpath.NewHostPathDriver(*driverName, *nodeID, *csiEndpoint, *ephemeral, *maxVolumesPerNode, version, capacity, *enableAttach)
	if err != nil {
		fmt.Printf("Failed to initialize driver: %s", err.Error())
		os.Exit(1)
	}

	if err := driver.Run(); err != nil {
		fmt.Printf("Failed to run driver: %s", err.Error())
		os.Exit(1)

	}
}
