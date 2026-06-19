/*
Copyright 2026 The Kubernetes Authors.

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

// Command healthctl is a development/testing utility for the hostpath CSI
// driver. It flips per-volume and node-level storage health by writing and
// removing marker files in the driver's state directory, so that the
// ControllerGetVolumeHealth, ControllerListVolumeHealth, NodeGetVolumeHealth
// and NodeGetStorageHealth RPCs report unhealthy conditions.
//
// This binary is intended for manual use against a running hostpath driver
// and is not part of the normal CSI driver image. Build it with:
//
//	make
//
// and run it on the host that holds the driver's state directory, e.g.:
//
//	./bin/healthctl -statedir /csi-data-dir mark-volume-unhealthy <volID> --scope controller --status INACCESSIBLE
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/kubernetes-csi/csi-driver-host-path/pkg/hostpath"
)

// globalStateDir is populated from the -statedir global flag parsed before
// the subcommand. All subcommands operate on this directory.
var globalStateDir string

func main() {
	// Global flags (-statedir) come before the subcommand. The global flag
	// set parses them and stops at the first non-flag argument (the
	// subcommand name). The remaining args (subcommand + its flags) are then
	// dispatched.
	globalFS := flag.NewFlagSet("healthctl", flag.ExitOnError)
	globalFS.StringVar(&globalStateDir, "statedir", "/csi-data-dir", "directory the hostpath driver uses for its state")
	globalFS.Usage = usage
	globalFS.Parse(os.Args[1:])

	if globalFS.NArg() < 1 {
		usage()
		os.Exit(2)
	}

	subcommand := globalFS.Arg(0)
	args := globalFS.Args()[1:]

	switch subcommand {
	case "mark-volume-unhealthy":
		os.Exit(runMarkVolumeUnhealthy(args))
	case "mark-volume-healthy":
		os.Exit(runMarkVolumeHealthy(args))
	case "mark-storage-unhealthy":
		os.Exit(runMarkStorageUnhealthy(args))
	case "mark-storage-healthy":
		os.Exit(runMarkStorageHealthy(args))
	case "list":
		os.Exit(runList(args))
	case "-h", "--help", "help":
		usage()
		return
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n\n", subcommand)
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `healthctl - manipulate hostpath CSI driver health markers (development only)

Usage:
  healthctl [global flags] <subcommand> [flags] [args]

Global flags:
  -statedir string   directory the hostpath driver uses for its state (default "/csi-data-dir")

Subcommands:
  mark-volume-unhealthy <volID>   Mark a volume as unhealthy
  mark-volume-healthy <volID>     Clear an unhealthy marker from a volume
  mark-storage-unhealthy          Mark node storage as unhealthy
  mark-storage-healthy            Clear the node storage unhealthy marker
  list                            List all current health markers

Volume subcommands accept --scope (both|controller|node, default both) to
control which side(s) report the unhealthy condition.

Run "healthctl <subcommand> -h" for subcommand-specific flags.`)
}

// reorderArgs moves all flag-style arguments (and their values) to the front
// and positional arguments to the back, so that the standard library flag
// package can parse flags that appear after a positional argument like a
// volume ID. This supports both "--flag value" and "--flag=value" forms and
// assumes non-boolean flags (all subcommand flags here take a value).
func reorderArgs(args []string) []string {
	var flags, positionals []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if strings.HasPrefix(a, "-") {
			flags = append(flags, a)
			// If the flag is not in "--flag=value" form, the next arg is its value.
			if !strings.Contains(a, "=") && i+1 < len(args) {
				flags = append(flags, args[i+1])
				i++
			}
			continue
		}
		positionals = append(positionals, a)
	}
	return append(flags, positionals...)
}

// parseScope validates and returns the scope string as a typed
// VolumeHealthScope, exiting the process on invalid input.
func parseScope(s string) hostpath.VolumeHealthScope {
	scope := hostpath.VolumeHealthScope(s)
	if err := scope.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
	return scope
}

func runMarkVolumeUnhealthy(args []string) int {
	fs := flag.NewFlagSet("mark-volume-unhealthy", flag.ExitOnError)
	scope := fs.String("scope", "both", "which side reports unhealthy: both, controller, or node")
	status := fs.String("status", "DEGRADED", "volume health status (one of: DEGRADED, INACCESSIBLE, DATA_LOSS)")
	reason := fs.String("reason", "manually-marked", "short CamelCase reason for the condition")
	message := fs.String("message", "", "human-readable description (defaults to a templated string)")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: healthctl -statedir <dir> mark-volume-unhealthy <volID> [flags]")
		fs.PrintDefaults()
	}
	fs.Parse(reorderArgs(args))

	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "error: exactly one volume ID argument is required")
		fs.Usage()
		return 2
	}
	volID := fs.Arg(0)

	sc := parseScope(*scope)

	if _, err := hostpath.ValidateVolumeHealthStatus(*status); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}

	msg := *message
	if msg == "" {
		msg = fmt.Sprintf("volume %s marked %s by healthctl", volID, *status)
	}

	m := hostpath.VolumeHealthMarker{
		Status:  *status,
		Reason:  *reason,
		Message: msg,
	}
	if err := hostpath.WriteVolumeHealthMarker(globalStateDir, volID, sc, m); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	fmt.Printf("volume %s marked unhealthy (scope=%s status=%s reason=%q)\n", volID, sc, *status, *reason)
	return 0
}

func runMarkVolumeHealthy(args []string) int {
	fs := flag.NewFlagSet("mark-volume-healthy", flag.ExitOnError)
	scope := fs.String("scope", "both", "which side's marker to clear: both, controller, or node (both clears all scopes)")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: healthctl -statedir <dir> mark-volume-healthy <volID> [flags]")
		fs.PrintDefaults()
	}
	fs.Parse(reorderArgs(args))

	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "error: exactly one volume ID argument is required")
		fs.Usage()
		return 2
	}
	volID := fs.Arg(0)

	sc := parseScope(*scope)

	// For "both", clear all scope markers to be safe; for a specific scope,
	// clear only that one.
	if sc == hostpath.ScopeBoth {
		if err := hostpath.ClearAllVolumeHealthMarkers(globalStateDir, volID); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		fmt.Printf("volume %s marked healthy (all scope markers cleared)\n", volID)
		return 0
	}
	if err := hostpath.ClearVolumeHealthMarker(globalStateDir, volID, sc); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	fmt.Printf("volume %s %s-scope marker cleared\n", volID, sc)
	return 0
}

func runMarkStorageUnhealthy(args []string) int {
	fs := flag.NewFlagSet("mark-storage-unhealthy", flag.ExitOnError)
	status := fs.String("status", "STORAGE_DEGRADED", "storage health status (one of: STORAGE_UNREACHABLE, STORAGE_DEGRADED)")
	reason := fs.String("reason", "manually-marked", "short CamelCase reason for the condition")
	message := fs.String("message", "", "human-readable description (defaults to a templated string)")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: healthctl -statedir <dir> mark-storage-unhealthy [flags]")
		fs.PrintDefaults()
	}
	fs.Parse(reorderArgs(args))

	if _, err := hostpath.ValidateStorageHealthStatus(*status); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}

	msg := *message
	if msg == "" {
		msg = fmt.Sprintf("node storage marked %s by healthctl", *status)
	}

	m := hostpath.StorageHealthMarker{
		Status:  *status,
		Reason:  *reason,
		Message: msg,
	}
	if err := hostpath.WriteStorageHealthMarker(globalStateDir, m); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	fmt.Printf("node storage marked unhealthy (status=%s reason=%q)\n", *status, *reason)
	return 0
}

func runMarkStorageHealthy(args []string) int {
	fs := flag.NewFlagSet("mark-storage-healthy", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: healthctl -statedir <dir> mark-storage-healthy")
		fs.PrintDefaults()
	}
	fs.Parse(reorderArgs(args))

	if err := hostpath.ClearStorageHealthMarker(globalStateDir); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	fmt.Println("node storage marked healthy (marker cleared)")
	return 0
}

func runList(args []string) int {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: healthctl -statedir <dir> list")
		fs.PrintDefaults()
	}
	fs.Parse(reorderArgs(args))

	volMarkers, err := hostpath.ListVolumeHealthMarkers(globalStateDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing volume markers: %v\n", err)
		return 1
	}
	storageMarker, err := hostpath.ReadStorageHealthMarker(globalStateDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading storage marker: %v\n", err)
		return 1
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "VOLUME HEALTH MARKERS")
	fmt.Fprintln(w, "VOLUME ID\tSCOPE\tSTATUS\tREASON\tMESSAGE")
	volIDs := make([]string, 0, len(volMarkers))
	for volID := range volMarkers {
		volIDs = append(volIDs, volID)
	}
	sort.Strings(volIDs)
	if len(volIDs) == 0 {
		fmt.Fprintln(w, "(none - all volumes healthy)\t\t\t\t")
	}
	for _, volID := range volIDs {
		m := volMarkers[volID]
		for _, e := range []struct {
			scope  string
			marker *hostpath.VolumeHealthMarker
		}{
			{"both", m.Both},
			{"controller", m.Controller},
			{"node", m.Node},
		} {
			if e.marker == nil {
				continue
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", volID, e.scope, e.marker.Status, e.marker.Reason, e.marker.Message)
		}
	}
	w.Flush()

	fmt.Fprintln(w)
	fmt.Fprintln(w, "NODE STORAGE HEALTH")
	fmt.Fprintln(w, "STATUS\tREASON\tMESSAGE")
	if storageMarker == nil {
		fmt.Fprintln(w, "(healthy - no marker)\t\t")
	} else {
		fmt.Fprintf(w, "%s\t%s\t%s\n", storageMarker.Status, storageMarker.Reason, storageMarker.Message)
	}
	w.Flush()

	return 0
}
