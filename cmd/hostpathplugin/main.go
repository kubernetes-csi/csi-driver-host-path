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
	// Set by the build process
	version = ""
)

func main() {
	var cfg hostpath.Config

	flag.StringVar(&cfg.Endpoint, "endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	flag.StringVar(&cfg.DriverName, "drivername", "hostpath.csi.k8s.io", "name of the driver")
	flag.StringVar(&cfg.NodeID, "nodeid", "", "node id")
	flag.BoolVar(&cfg.Ephemeral, "ephemeral", false, "publish volumes in ephemeral mode even if kubelet did not ask for it (only needed for Kubernetes 1.15)")
	flag.Int64Var(&cfg.MaxVolumesPerNode, "maxvolumespernode", 0, "limit of volumes per node")
	flag.Var(&cfg.Capacity, "capacity", "Simulate storage capacity. The parameter is <kind>=<quantity> where <kind> is the value of a 'kind' storage class parameter and <quantity> is the total amount of bytes for that kind. The flag may be used multiple times to configure different kinds.")
	flag.BoolVar(&cfg.EnableAttach, "enable-attach", false, "Enables RPC_PUBLISH_UNPUBLISH_VOLUME capability.")
	showVersion := flag.Bool("version", false, "Show version.")
	// The proxy-endpoint option is intended to used by the Kubernetes E2E test suite
	// for proxying incoming calls to the embedded mock CSI driver.
	proxyEndpoint := flag.String("proxy-endpoint", "", "Instead of running the CSI driver code, just proxy connections from csiEndpoint to the given listening socket.")

	flag.Parse()

	if *showVersion {
		baseName := path.Base(os.Args[0])
		fmt.Println(baseName, version)
		return
	}

	if cfg.Ephemeral {
		fmt.Fprintln(os.Stderr, "Deprecation warning: The ephemeral flag is deprecated and should only be used when deploying on Kubernetes 1.15. It will be removed in the future.")
	}

	if *proxyEndpoint != "" {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		closer, err := proxy.Run(ctx, cfg.Endpoint, *proxyEndpoint)
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

	if version != "" {
		cfg.VendorVersion = version
	}

	driver, err := hostpath.NewHostPathDriver(cfg)
	if err != nil {
		fmt.Printf("Failed to initialize driver: %s", err.Error())
		os.Exit(1)
	}

	if err := driver.Run(); err != nil {
		fmt.Printf("Failed to run driver: %s", err.Error())
		os.Exit(1)

	}
}
