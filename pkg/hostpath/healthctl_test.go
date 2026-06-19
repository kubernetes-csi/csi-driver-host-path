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
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVolumeHealthPath(t *testing.T) {
	t.Run("valid both", func(t *testing.T) {
		p, err := VolumeHealthPath("/state", "vol-1", ScopeBoth)
		require.NoError(t, err)
		assert.Equal(t, filepath.Join("/state", "vol-1.health"), p)
	})
	t.Run("valid controller", func(t *testing.T) {
		p, err := VolumeHealthPath("/state", "vol-1", ScopeController)
		require.NoError(t, err)
		assert.Equal(t, filepath.Join("/state", "vol-1.controller.health"), p)
	})
	t.Run("valid node", func(t *testing.T) {
		p, err := VolumeHealthPath("/state", "vol-1", ScopeNode)
		require.NoError(t, err)
		assert.Equal(t, filepath.Join("/state", "vol-1.node.health"), p)
	})
	t.Run("empty id", func(t *testing.T) {
		_, err := VolumeHealthPath("/state", "", ScopeBoth)
		assert.Error(t, err)
	})
	t.Run("slash in id", func(t *testing.T) {
		_, err := VolumeHealthPath("/state", "a/b", ScopeBoth)
		assert.Error(t, err)
	})
	t.Run("parent reference", func(t *testing.T) {
		_, err := VolumeHealthPath("/state", "..", ScopeBoth)
		assert.Error(t, err)
	})
	t.Run("invalid scope", func(t *testing.T) {
		_, err := VolumeHealthPath("/state", "vol-1", VolumeHealthScope("bogus"))
		assert.Error(t, err)
	})
}

func TestVolumeHealthScopeValidate(t *testing.T) {
	for _, s := range []VolumeHealthScope{ScopeBoth, ScopeController, ScopeNode} {
		assert.NoError(t, s.Validate(), s)
	}
	assert.Error(t, VolumeHealthScope("bogus").Validate())
}

func TestVolumeHealthMarkerRoundTrip(t *testing.T) {
	for _, scope := range []VolumeHealthScope{ScopeBoth, ScopeController, ScopeNode} {
		t.Run(string(scope), func(t *testing.T) {
			stateDir := t.TempDir()
			volID := "vol-roundtrip"

			got, err := ReadVolumeHealthMarker(stateDir, volID, scope)
			require.NoError(t, err)
			assert.Nil(t, got, "expected no marker initially")

			m := VolumeHealthMarker{Status: "DEGRADED", Reason: "simulated", Message: "manually marked unhealthy"}
			require.NoError(t, WriteVolumeHealthMarker(stateDir, volID, scope, m))

			got, err = ReadVolumeHealthMarker(stateDir, volID, scope)
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, m, *got)

			require.NoError(t, ClearVolumeHealthMarker(stateDir, volID, scope))

			got, err = ReadVolumeHealthMarker(stateDir, volID, scope)
			require.NoError(t, err)
			assert.Nil(t, got)

			// Clearing again is idempotent.
			require.NoError(t, ClearVolumeHealthMarker(stateDir, volID, scope))
		})
	}
}

func TestReadEffectiveVolumeHealthMarker(t *testing.T) {
	stateDir := t.TempDir()
	volID := "vol-eff"

	// No markers at all.
	got, err := ReadEffectiveVolumeHealthMarker(stateDir, volID, ScopeController)
	require.NoError(t, err)
	assert.Nil(t, got)

	// Generic "both" marker is used as fallback by controller and node.
	both := VolumeHealthMarker{Status: "DEGRADED", Reason: "both"}
	require.NoError(t, WriteVolumeHealthMarker(stateDir, volID, ScopeBoth, both))

	got, err = ReadEffectiveVolumeHealthMarker(stateDir, volID, ScopeController)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "both", got.Reason)

	got, err = ReadEffectiveVolumeHealthMarker(stateDir, volID, ScopeNode)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "both", got.Reason)

	// A scope-specific marker takes precedence over the generic one.
	ctrl := VolumeHealthMarker{Status: "INACCESSIBLE", Reason: "controller-specific"}
	require.NoError(t, WriteVolumeHealthMarker(stateDir, volID, ScopeController, ctrl))

	got, err = ReadEffectiveVolumeHealthMarker(stateDir, volID, ScopeController)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "controller-specific", got.Reason, "controller marker should take precedence")

	// Node still falls back to the generic marker.
	got, err = ReadEffectiveVolumeHealthMarker(stateDir, volID, ScopeNode)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "both", got.Reason)
}

func TestClearAllVolumeHealthMarkers(t *testing.T) {
	stateDir := t.TempDir()
	volID := "vol-clearall"

	for _, scope := range []VolumeHealthScope{ScopeBoth, ScopeController, ScopeNode} {
		require.NoError(t, WriteVolumeHealthMarker(stateDir, volID, scope, VolumeHealthMarker{Status: "DEGRADED"}))
	}

	require.NoError(t, ClearAllVolumeHealthMarkers(stateDir, volID))

	for _, scope := range []VolumeHealthScope{ScopeBoth, ScopeController, ScopeNode} {
		got, err := ReadVolumeHealthMarker(stateDir, volID, scope)
		require.NoError(t, err)
		assert.Nil(t, got, "scope %s marker should be cleared", scope)
	}

	// Idempotent.
	require.NoError(t, ClearAllVolumeHealthMarkers(stateDir, volID))
}

func TestWriteVolumeHealthMarkerIsAtomic(t *testing.T) {
	stateDir := t.TempDir()
	volID := "vol-atomic"

	require.NoError(t, WriteVolumeHealthMarker(stateDir, volID, ScopeBoth, VolumeHealthMarker{Status: "DEGRADED"}))

	entries, err := os.ReadDir(stateDir)
	require.NoError(t, err)
	for _, e := range entries {
		assert.False(t, filepath.Ext(e.Name()) == ".tmp", "leftover temp file %q", e.Name())
	}
}

func TestListVolumeHealthMarkers(t *testing.T) {
	stateDir := t.TempDir()

	require.NoError(t, WriteVolumeHealthMarker(stateDir, "vol-a", ScopeBoth, VolumeHealthMarker{Status: "DEGRADED", Reason: "a-both"}))
	require.NoError(t, WriteVolumeHealthMarker(stateDir, "vol-b", ScopeController, VolumeHealthMarker{Status: "INACCESSIBLE", Reason: "b-ctrl"}))
	require.NoError(t, WriteVolumeHealthMarker(stateDir, "vol-b", ScopeNode, VolumeHealthMarker{Status: "DATA_LOSS", Reason: "b-node"}))
	// Node-level storage marker must be excluded.
	require.NoError(t, WriteStorageHealthMarker(stateDir, StorageHealthMarker{Status: "STORAGE_DEGRADED"}))
	// Unrelated file must be ignored.
	require.NoError(t, os.WriteFile(filepath.Join(stateDir, "state.json"), []byte("{}"), 0600))

	got, err := ListVolumeHealthMarkers(stateDir)
	require.NoError(t, err)
	require.Contains(t, got, "vol-a")
	require.Contains(t, got, "vol-b")

	assert.NotNil(t, got["vol-a"].Both)
	assert.Equal(t, "a-both", got["vol-a"].Both.Reason)
	assert.Nil(t, got["vol-a"].Controller)
	assert.Nil(t, got["vol-a"].Node)

	assert.Nil(t, got["vol-b"].Both)
	assert.NotNil(t, got["vol-b"].Controller)
	assert.Equal(t, "b-ctrl", got["vol-b"].Controller.Reason)
	assert.NotNil(t, got["vol-b"].Node)
	assert.Equal(t, "b-node", got["vol-b"].Node.Reason)

	// EffectiveForController/Node helpers.
	assert.Equal(t, "a-both", got["vol-a"].EffectiveForController().Reason)
	assert.Equal(t, "b-ctrl", got["vol-b"].EffectiveForController().Reason)
	assert.Equal(t, "b-node", got["vol-b"].EffectiveForNode().Reason)

	// A malformed marker file is skipped, not fatal.
	require.NoError(t, os.WriteFile(filepath.Join(stateDir, "vol-x.health"), []byte("not json"), 0600))
	got, err = ListVolumeHealthMarkers(stateDir)
	require.NoError(t, err)
	assert.Len(t, got, 2, "malformed marker should be skipped")
}

func TestStorageHealthMarkerRoundTrip(t *testing.T) {
	stateDir := t.TempDir()

	got, err := ReadStorageHealthMarker(stateDir)
	require.NoError(t, err)
	assert.Nil(t, got)

	m := StorageHealthMarker{Status: "STORAGE_UNREACHABLE", Reason: "disk gone", Message: "backend down"}
	require.NoError(t, WriteStorageHealthMarker(stateDir, m))

	got, err = ReadStorageHealthMarker(stateDir)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, m, *got)

	require.NoError(t, ClearStorageHealthMarker(stateDir))

	got, err = ReadStorageHealthMarker(stateDir)
	require.NoError(t, err)
	assert.Nil(t, got)

	require.NoError(t, ClearStorageHealthMarker(stateDir), "clearing again should be idempotent")
}

func TestValidateVolumeHealthStatus(t *testing.T) {
	for _, s := range []string{"INACCESSIBLE", "DATA_LOSS", "DEGRADED"} {
		v, err := ValidateVolumeHealthStatus(s)
		require.NoError(t, err, s)
		assert.NotEqual(t, csi.VolumeHealthErrorType_UNKNOWN_VOLUME_HEALTH_TYPE, v, s)
	}
	_, err := ValidateVolumeHealthStatus("UNKNOWN_VOLUME_HEALTH_TYPE")
	assert.Error(t, err)
	_, err = ValidateVolumeHealthStatus("bogus")
	assert.Error(t, err)
}

func TestValidateStorageHealthStatus(t *testing.T) {
	for _, s := range []string{"STORAGE_UNREACHABLE", "STORAGE_DEGRADED"} {
		v, err := ValidateStorageHealthStatus(s)
		require.NoError(t, err, s)
		assert.NotEqual(t, csi.StorageHealthErrorType_UNKNOWN_STORAGE_HEALTH_ERROR_TYPE, v, s)
	}
	_, err := ValidateStorageHealthStatus("UNKNOWN_STORAGE_HEALTH_ERROR_TYPE")
	assert.Error(t, err)
	_, err = ValidateStorageHealthStatus("bogus")
	assert.Error(t, err)
}

func TestValidStatusNamesAreSortedAndExcludeUnknown(t *testing.T) {
	vol := validVolumeHealthStatusNames()
	assert.Equal(t, []string{"DATA_LOSS", "DEGRADED", "INACCESSIBLE"}, vol)
	stor := validStorageHealthStatusNames()
	assert.Equal(t, []string{"STORAGE_DEGRADED", "STORAGE_UNREACHABLE"}, stor)
	// Sanity check sort (already implicitly checked above).
	assert.True(t, sort.StringsAreSorted(vol))
	assert.True(t, sort.StringsAreSorted(stor))
}
