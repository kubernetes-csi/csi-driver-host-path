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

package hostpath

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

const (
	// volumeHealthSuffix is the suffix for the "both-sides" marker file
	// (e.g. <stateDir>/<volID>.health), read by both the controller and
	// the node RPCs as a fallback.
	volumeHealthSuffix = ".health"
	// volumeHealthControllerSuffix is the suffix for the controller-only
	// marker file (e.g. <stateDir>/<volID>.controller.health).
	volumeHealthControllerSuffix = ".controller.health"
	// volumeHealthNodeSuffix is the suffix for the node-only marker file
	// (e.g. <stateDir>/<volID>.node.health).
	volumeHealthNodeSuffix = ".node.health"
	// storageHealthFile is the name of the node-level storage health
	// marker file inside the state directory.
	storageHealthFile = "storage.health"
)

// VolumeHealthScope selects which side(s) a per-volume health marker
// applies to. ScopeBoth writes the generic <volID>.health marker that
// both controller and node RPCs read as a fallback; ScopeController and
// ScopeNode write side-specific markers that take precedence over the
// generic one.
type VolumeHealthScope string

const (
	ScopeBoth       VolumeHealthScope = "both"
	ScopeController VolumeHealthScope = "controller"
	ScopeNode       VolumeHealthScope = "node"
)

// Validate returns an error if the scope is not a recognized value.
func (s VolumeHealthScope) Validate() error {
	switch s {
	case ScopeBoth, ScopeController, ScopeNode:
		return nil
	default:
		return fmt.Errorf("invalid scope %q: expected one of both, controller, node", s)
	}
}

// suffix returns the on-disk file suffix for the scope.
func (s VolumeHealthScope) suffix() string {
	switch s {
	case ScopeController:
		return volumeHealthControllerSuffix
	case ScopeNode:
		return volumeHealthNodeSuffix
	default:
		return volumeHealthSuffix
	}
}

// VolumeHealthMarker is the JSON shape persisted on disk for a per-volume
// health marker. The Status field holds the string name of a
// csi.VolumeHealthErrorType value (e.g. "DEGRADED").
type VolumeHealthMarker struct {
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

// StorageHealthMarker is the JSON shape persisted on disk for the
// node-level storage health marker. The Status field holds the string
// name of a csi.StorageHealthErrorType value (e.g. "STORAGE_DEGRADED").
type StorageHealthMarker struct {
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

// validateVolumeID rejects volume IDs that could escape the state
// directory (containing path separators or "..").
func validateVolumeID(volID string) error {
	if volID == "" {
		return fmt.Errorf("volume ID is empty")
	}
	if strings.ContainsAny(volID, `/\`) || volID == ".." || strings.Contains(volID, "..") {
		return fmt.Errorf("invalid volume ID %q: must not contain path separators or parent references", volID)
	}
	return nil
}

// VolumeHealthPath returns the path to the per-volume health marker file
// for the given volume ID and scope.
func VolumeHealthPath(stateDir, volID string, scope VolumeHealthScope) (string, error) {
	if err := validateVolumeID(volID); err != nil {
		return "", err
	}
	if err := scope.Validate(); err != nil {
		return "", err
	}
	return filepath.Join(stateDir, volID+scope.suffix()), nil
}

// StorageHealthPath returns the path to the node-level storage health
// marker file.
func StorageHealthPath(stateDir string) string {
	return filepath.Join(stateDir, storageHealthFile)
}

// WriteVolumeHealthMarker atomically writes a per-volume health marker for
// the given scope. An existing marker for the same volume and scope is
// overwritten.
func WriteVolumeHealthMarker(stateDir, volID string, scope VolumeHealthScope, m VolumeHealthMarker) error {
	path, err := VolumeHealthPath(stateDir, volID, scope)
	if err != nil {
		return err
	}
	return writeJSONAtomic(path, m)
}

// ReadVolumeHealthMarker reads a per-volume health marker for exactly the
// given scope (no fallback). It returns (nil, nil) when no marker exists
// for that scope.
func ReadVolumeHealthMarker(stateDir, volID string, scope VolumeHealthScope) (*VolumeHealthMarker, error) {
	path, err := VolumeHealthPath(stateDir, volID, scope)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var m VolumeHealthMarker
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse volume health marker %q: %w", path, err)
	}
	return &m, nil
}

// ReadEffectiveVolumeHealthMarker reads the marker that should be reported
// for the given scope, falling back to the generic "both" marker when no
// scope-specific marker exists. For ScopeController it reads
// <volID>.controller.health then falls back to <volID>.health; for
// ScopeNode it reads <volID>.node.health then falls back to
// <volID>.health; for ScopeBoth it reads <volID>.health directly.
func ReadEffectiveVolumeHealthMarker(stateDir, volID string, scope VolumeHealthScope) (*VolumeHealthMarker, error) {
	if scope != ScopeBoth {
		m, err := ReadVolumeHealthMarker(stateDir, volID, scope)
		if err != nil {
			return nil, err
		}
		if m != nil {
			return m, nil
		}
	}
	return ReadVolumeHealthMarker(stateDir, volID, ScopeBoth)
}

// ClearVolumeHealthMarker removes the per-volume health marker for the
// given scope. It is not an error when the marker does not exist.
func ClearVolumeHealthMarker(stateDir, volID string, scope VolumeHealthScope) error {
	path, err := VolumeHealthPath(stateDir, volID, scope)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// ClearAllVolumeHealthMarkers removes every per-volume health marker
// (both, controller and node scopes) for the given volume ID. It is not an
// error when none of the markers exist.
func ClearAllVolumeHealthMarkers(stateDir, volID string) error {
	for _, scope := range []VolumeHealthScope{ScopeBoth, ScopeController, ScopeNode} {
		if err := ClearVolumeHealthMarker(stateDir, volID, scope); err != nil {
			return err
		}
	}
	return nil
}

// VolumeHealthMarkers holds all per-scope markers that exist for a single
// volume. Any field may be nil when the corresponding marker file is
// absent.
type VolumeHealthMarkers struct {
	Both       *VolumeHealthMarker
	Controller *VolumeHealthMarker
	Node       *VolumeHealthMarker
}

// EffectiveForController returns the marker a controller-side RPC should
// report: the controller-specific marker if present, otherwise the generic
// "both" marker.
func (v VolumeHealthMarkers) EffectiveForController() *VolumeHealthMarker {
	if v.Controller != nil {
		return v.Controller
	}
	return v.Both
}

// EffectiveForNode returns the marker a node-side RPC should report: the
// node-specific marker if present, otherwise the generic "both" marker.
func (v VolumeHealthMarkers) EffectiveForNode() *VolumeHealthMarker {
	if v.Node != nil {
		return v.Node
	}
	return v.Both
}

// ListVolumeHealthMarkers returns all per-volume health markers keyed by
// volume ID, with each present scope populated. The node-level storage
// marker is excluded. Malformed marker files are skipped.
func ListVolumeHealthMarkers(stateDir string) (map[string]VolumeHealthMarkers, error) {
	entries, err := os.ReadDir(stateDir)
	if err != nil {
		return nil, err
	}
	result := make(map[string]VolumeHealthMarkers)
	for _, entry := range entries {
		name := entry.Name()
		if name == storageHealthFile {
			continue
		}
		var scope VolumeHealthScope
		var volID string
		switch {
		case strings.HasSuffix(name, volumeHealthControllerSuffix):
			scope = ScopeController
			volID = strings.TrimSuffix(name, volumeHealthControllerSuffix)
		case strings.HasSuffix(name, volumeHealthNodeSuffix):
			scope = ScopeNode
			volID = strings.TrimSuffix(name, volumeHealthNodeSuffix)
		case strings.HasSuffix(name, volumeHealthSuffix):
			scope = ScopeBoth
			volID = strings.TrimSuffix(name, volumeHealthSuffix)
		default:
			continue
		}
		data, err := os.ReadFile(filepath.Join(stateDir, name))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		var m VolumeHealthMarker
		if err := json.Unmarshal(data, &m); err != nil {
			// Skip malformed markers rather than failing the whole list.
			continue
		}
		entry := result[volID]
		switch scope {
		case ScopeBoth:
			entry.Both = &m
		case ScopeController:
			entry.Controller = &m
		case ScopeNode:
			entry.Node = &m
		}
		result[volID] = entry
	}
	return result, nil
}

// WriteStorageHealthMarker atomically writes the node-level storage
// health marker, overwriting any existing marker.
func WriteStorageHealthMarker(stateDir string, m StorageHealthMarker) error {
	return writeJSONAtomic(StorageHealthPath(stateDir), m)
}

// ReadStorageHealthMarker reads the node-level storage health marker. It
// returns (nil, nil) when no marker exists.
func ReadStorageHealthMarker(stateDir string) (*StorageHealthMarker, error) {
	data, err := os.ReadFile(StorageHealthPath(stateDir))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var m StorageHealthMarker
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse storage health marker: %w", err)
	}
	return &m, nil
}

// ClearStorageHealthMarker removes the node-level storage health marker.
// It is not an error when the marker does not exist.
func ClearStorageHealthMarker(stateDir string) error {
	if err := os.Remove(StorageHealthPath(stateDir)); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// ValidateVolumeHealthStatus maps a status string (e.g. "DEGRADED") to
// the corresponding csi.VolumeHealthErrorType. It rejects the
// UNKNOWN_VOLUME_HEALTH_TYPE placeholder and any unrecognized value.
func ValidateVolumeHealthStatus(s string) (csi.VolumeHealthErrorType, error) {
	v, ok := csi.VolumeHealthErrorType_value[s]
	if !ok {
		return 0, fmt.Errorf("invalid volume health status %q: expected one of %v", s, validVolumeHealthStatusNames())
	}
	t := csi.VolumeHealthErrorType(v)
	if t == csi.VolumeHealthErrorType_UNKNOWN_VOLUME_HEALTH_TYPE {
		return 0, fmt.Errorf("invalid volume health status %q: must not be the UNKNOWN placeholder", s)
	}
	return t, nil
}

// ValidateStorageHealthStatus maps a status string (e.g.
// "STORAGE_DEGRADED") to the corresponding csi.StorageHealthErrorType.
// It rejects the UNKNOWN_STORAGE_HEALTH_ERROR_TYPE placeholder and any
// unrecognized value.
func ValidateStorageHealthStatus(s string) (csi.StorageHealthErrorType, error) {
	v, ok := csi.StorageHealthErrorType_value[s]
	if !ok {
		return 0, fmt.Errorf("invalid storage health status %q: expected one of %v", s, validStorageHealthStatusNames())
	}
	t := csi.StorageHealthErrorType(v)
	if t == csi.StorageHealthErrorType_UNKNOWN_STORAGE_HEALTH_ERROR_TYPE {
		return 0, fmt.Errorf("invalid storage health status %q: must not be the UNKNOWN placeholder", s)
	}
	return t, nil
}

func validVolumeHealthStatusNames() []string {
	names := make([]string, 0, len(csi.VolumeHealthErrorType_value))
	for name, val := range csi.VolumeHealthErrorType_value {
		if csi.VolumeHealthErrorType(val) != csi.VolumeHealthErrorType_UNKNOWN_VOLUME_HEALTH_TYPE {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

func validStorageHealthStatusNames() []string {
	names := make([]string, 0, len(csi.StorageHealthErrorType_value))
	for name, val := range csi.StorageHealthErrorType_value {
		if csi.StorageHealthErrorType(val) != csi.StorageHealthErrorType_UNKNOWN_STORAGE_HEALTH_ERROR_TYPE {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

// writeJSONAtomic marshals the value to JSON and writes it to path via a
// temp file + rename so that readers never observe a partially-written
// marker.
// getVolumeHealthEntries returns the CSI health entries for a volume based on
// its on-disk marker for the given scope. An empty slice (nil) means no
// adverse condition is known (healthy). Errors reading the marker are treated
// as healthy so a corrupt marker cannot break the RPC.
func (hp *hostPath) getVolumeHealthEntries(volID string, scope VolumeHealthScope) []*csi.VolumeHealth_VolumeHealthEntry {
	m, err := ReadEffectiveVolumeHealthMarker(hp.config.StateDir, volID, scope)
	if err != nil || m == nil {
		return nil
	}
	return []*csi.VolumeHealth_VolumeHealthEntry{toVolumeHealthEntry(m)}
}

// toVolumeHealthEntry converts a stored marker into a CSI VolumeHealth entry.
// An invalid status string is treated as DEGRADED so a corrupt marker still
// surfaces an adverse condition rather than being silently healthy.
func toVolumeHealthEntry(m *VolumeHealthMarker) *csi.VolumeHealth_VolumeHealthEntry {
	status, err := ValidateVolumeHealthStatus(m.Status)
	if err != nil {
		status = csi.VolumeHealthErrorType_DEGRADED
	}
	return &csi.VolumeHealth_VolumeHealthEntry{
		Status:  status,
		Reason:  m.Reason,
		Message: m.Message,
	}
}

// getStorageBackendHealth returns the CSI storage backend health entries based
// on the node-level on-disk marker. An empty slice (nil) means healthy.
func (hp *hostPath) getStorageBackendHealth() []*csi.NodeGetStorageHealthResponse_StorageBackendHealth {
	m, err := ReadStorageHealthMarker(hp.config.StateDir)
	if err != nil || m == nil {
		return nil
	}
	return []*csi.NodeGetStorageHealthResponse_StorageBackendHealth{toStorageBackendHealth(m)}
}

func toStorageBackendHealth(m *StorageHealthMarker) *csi.NodeGetStorageHealthResponse_StorageBackendHealth {
	status, err := ValidateStorageHealthStatus(m.Status)
	if err != nil {
		status = csi.StorageHealthErrorType_STORAGE_DEGRADED
	}
	return &csi.NodeGetStorageHealthResponse_StorageBackendHealth{
		Status:  status,
		Reason:  m.Reason,
		Message: m.Message,
	}
}

func writeJSONAtomic(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode marker: %w", err)
	}
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".health-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpName := tmp.Name()
	cleanup := func() {
		tmp.Close()
		os.Remove(tmpName)
	}
	if _, err := tmp.Write(data); err != nil {
		cleanup()
		return fmt.Errorf("failed to write marker: %w", err)
	}
	if err := tmp.Chmod(0600); err != nil {
		cleanup()
		return fmt.Errorf("failed to set marker permissions: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("failed to close marker: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return fmt.Errorf("failed to commit marker: %w", err)
	}
	return nil
}
