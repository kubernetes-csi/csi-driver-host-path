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
	cfg := hostpath.Config{
		VendorVersion: version,
	}

	flag.StringVar(&cfg.Endpoint, "endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	flag.StringVar(&cfg.DriverName, "drivername", "hostpath.csi.k8s.io", "name of the driver")
	flag.StringVar(&cfg.StateDir, "statedir", "/csi-data-dir", "directory for storing state information across driver restarts, volumes and snapshots")
	flag.StringVar(&cfg.NodeID, "nodeid", "", "node id")
	flag.BoolVar(&cfg.Ephemeral, "ephemeral", false, "publish volumes in ephemeral mode even if kubelet did not ask for it (only needed for Kubernetes 1.15)")
	flag.Int64Var(&cfg.MaxVolumesPerNode, "maxvolumespernode", 0, "limit of volumes per node")
	flag.Var(&cfg.Capacity, "capacity", "Simulate storage capacity. The parameter is <kind>=<quantity> where <kind> is the value of a 'kind' storage class parameter and <quantity> is the total amount of bytes for that kind. The flag may be used multiple times to configure different kinds.")
	flag.BoolVar(&cfg.EnableAttach, "enable-attach", false, "Enables RPC_PUBLISH_UNPUBLISH_VOLUME capability.")
	flag.BoolVar(&cfg.CheckVolumeLifecycle, "check-volume-lifecycle", false, "Can be used to turn some violations of the volume lifecycle into warnings instead of failing the incorrect gRPC call. Disabled by default because of https://github.com/kubernetes/kubernetes/issues/101911.")
	flag.Int64Var(&cfg.MaxVolumeSize, "max-volume-size", 1024*1024*1024*1024, "maximum size of volumes in bytes (inclusive)")
	flag.BoolVar(&cfg.EnableTopology, "enable-topology", true, "Enables PluginCapability_Service_VOLUME_ACCESSIBILITY_CONSTRAINTS capability.")
	flag.BoolVar(&cfg.EnableVolumeExpansion, "node-expand-required", true, "Enables volume expansion capability of the plugin(Deprecated). Please use enable-volume-expansion flag.")

	flag.BoolVar(&cfg.EnableVolumeExpansion, "enable-volume-expansion", true, "Enables volume expansion feature.")
	flag.BoolVar(&cfg.EnableControllerModifyVolume, "enable-controller-modify-volume", false, "Enables Controller modify volume feature.")
	flag.Var(&cfg.AcceptedMutableParameterNames, "accepted-mutable-parameter-names", "Comma separated list of parameter names that can be modified on a persistent volume. This is only used when enable-controller-modify-volume is true. If unset, all parameters are mutable.")
	flag.BoolVar(&cfg.DisableControllerExpansion, "disable-controller-expansion", false, "Disables Controller volume expansion capability.")
	flag.BoolVar(&cfg.DisableNodeExpansion, "disable-node-expansion", false, "Disables Node volume expansion capability.")
	flag.Int64Var(&cfg.MaxVolumeExpansionSizeNode, "max-volume-size-node", 0, "Maximum allowed size of volume when expanded on the node. Defaults to same size as max-volume-size.")

	flag.Int64Var(&cfg.AttachLimit, "attach-limit", 0, "Maximum number of attachable volumes on a node. Zero refers to no limit.")
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

	if cfg.MaxVolumeExpansionSizeNode == 0 {
		cfg.MaxVolumeExpansionSizeNode = cfg.MaxVolumeSize
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
