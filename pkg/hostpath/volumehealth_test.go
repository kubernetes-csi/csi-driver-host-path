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
	"context"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-driver-host-path/pkg/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// newTestHostPath returns a hostPath driver backed by a temp state dir.
func newTestHostPath(t *testing.T) *hostPath {
	t.Helper()
	stateDir := t.TempDir()
	cfg := Config{
		StateDir:       stateDir,
		Endpoint:       "unix://tmp/csi.sock",
		DriverName:     "hostpath.csi.k8s.io",
		NodeID:         "fakeNodeID",
		MaxVolumeSize:  1024 * 1024 * 1024 * 1024,
		EnableTopology: true,
	}
	hp, err := NewHostPathDriver(cfg)
	require.NoError(t, err)
	return hp
}

func mustCreateVolume(t *testing.T, hp *hostPath, volID string) {
	t.Helper()
	_, err := hp.createVolume(volID, volID, 1024*1024, state.MountAccess, false, "")
	require.NoError(t, err)
}

func TestControllerGetVolumeHealth_Healthy(t *testing.T) {
	hp := newTestHostPath(t)
	mustCreateVolume(t, hp, "vol-healthy")

	resp, err := hp.ControllerGetVolumeHealth(context.TODO(), &csi.ControllerGetVolumeHealthRequest{VolumeId: "vol-healthy"})
	require.NoError(t, err)
	require.NotNil(t, resp.VolumeHealth)
	assert.Empty(t, resp.VolumeHealth.GetHealthStatuses(), "expected no adverse health entries")
}

func TestControllerGetVolumeHealth_Unhealthy(t *testing.T) {
	hp := newTestHostPath(t)
	mustCreateVolume(t, hp, "vol-unhealthy")

	require.NoError(t, WriteVolumeHealthMarker(hp.config.StateDir, "vol-unhealthy", ScopeBoth, VolumeHealthMarker{
		Status:  "INACCESSIBLE",
		Reason:  "simulated",
		Message: "marked by test",
	}))

	resp, err := hp.ControllerGetVolumeHealth(context.TODO(), &csi.ControllerGetVolumeHealthRequest{VolumeId: "vol-unhealthy"})
	require.NoError(t, err)
	require.NotNil(t, resp.VolumeHealth)
	assert.Equal(t, "vol-unhealthy", resp.VolumeHealth.GetVolumeId())

	entries := resp.VolumeHealth.GetHealthStatuses()
	require.Len(t, entries, 1)
	assert.Equal(t, csi.VolumeHealthErrorType_INACCESSIBLE, entries[0].GetStatus())
	assert.Equal(t, "simulated", entries[0].GetReason())
	assert.Equal(t, "marked by test", entries[0].GetMessage())

	// Clearing the marker restores healthy state.
	require.NoError(t, ClearVolumeHealthMarker(hp.config.StateDir, "vol-unhealthy", ScopeBoth))
	resp, err = hp.ControllerGetVolumeHealth(context.TODO(), &csi.ControllerGetVolumeHealthRequest{VolumeId: "vol-unhealthy"})
	require.NoError(t, err)
	assert.Empty(t, resp.VolumeHealth.GetHealthStatuses())
}

func TestControllerGetVolumeHealth_ControllerScopedOnly(t *testing.T) {
	hp := newTestHostPath(t)
	mustCreateVolume(t, hp, "vol-ctrl-only")

	// A controller-scoped marker is visible to the controller RPC...
	require.NoError(t, WriteVolumeHealthMarker(hp.config.StateDir, "vol-ctrl-only", ScopeController, VolumeHealthMarker{
		Status: "DEGRADED", Reason: "ctrl-only",
	}))
	resp, err := hp.ControllerGetVolumeHealth(context.TODO(), &csi.ControllerGetVolumeHealthRequest{VolumeId: "vol-ctrl-only"})
	require.NoError(t, err)
	entries := resp.VolumeHealth.GetHealthStatuses()
	require.Len(t, entries, 1)
	assert.Equal(t, csi.VolumeHealthErrorType_DEGRADED, entries[0].GetStatus())

	// ...but NOT to the node RPC.
	respNode, err := hp.NodeGetVolumeHealth(context.TODO(), &csi.NodeGetVolumeHealthRequest{VolumeId: "vol-ctrl-only"})
	require.NoError(t, err)
	assert.Empty(t, respNode.VolumeHealth.GetHealthStatuses(), "controller-scoped marker should not show on node side")
}

func TestNodeGetVolumeHealth_NodeScopedOnly(t *testing.T) {
	hp := newTestHostPath(t)
	mustCreateVolume(t, hp, "vol-node-only")

	// A node-scoped marker is visible to the node RPC...
	require.NoError(t, WriteVolumeHealthMarker(hp.config.StateDir, "vol-node-only", ScopeNode, VolumeHealthMarker{
		Status: "DATA_LOSS", Reason: "node-only",
	}))
	respNode, err := hp.NodeGetVolumeHealth(context.TODO(), &csi.NodeGetVolumeHealthRequest{VolumeId: "vol-node-only"})
	require.NoError(t, err)
	entries := respNode.VolumeHealth.GetHealthStatuses()
	require.Len(t, entries, 1)
	assert.Equal(t, csi.VolumeHealthErrorType_DATA_LOSS, entries[0].GetStatus())

	// ...but NOT to the controller RPC.
	respCtrl, err := hp.ControllerGetVolumeHealth(context.TODO(), &csi.ControllerGetVolumeHealthRequest{VolumeId: "vol-node-only"})
	require.NoError(t, err)
	assert.Empty(t, respCtrl.VolumeHealth.GetHealthStatuses(), "node-scoped marker should not show on controller side")
}

func TestControllerGetVolumeHealth_NotFound(t *testing.T) {
	hp := newTestHostPath(t)
	_, err := hp.ControllerGetVolumeHealth(context.TODO(), &csi.ControllerGetVolumeHealthRequest{VolumeId: "nope"})
	assert.Error(t, err)
}

func TestControllerGetVolume_NotFound(t *testing.T) {
	hp := newTestHostPath(t)
	_, err := hp.ControllerGetVolume(context.TODO(), &csi.ControllerGetVolumeRequest{VolumeId: "nope"})
	require.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
}

func TestNodeGetVolumeHealth_Unhealthy(t *testing.T) {
	hp := newTestHostPath(t)
	mustCreateVolume(t, hp, "vol-node")

	require.NoError(t, WriteVolumeHealthMarker(hp.config.StateDir, "vol-node", ScopeBoth, VolumeHealthMarker{
		Status: "DATA_LOSS", Reason: "r", Message: "m",
	}))

	resp, err := hp.NodeGetVolumeHealth(context.TODO(), &csi.NodeGetVolumeHealthRequest{VolumeId: "vol-node"})
	require.NoError(t, err)
	require.NotNil(t, resp.VolumeHealth)
	entries := resp.VolumeHealth.GetHealthStatuses()
	require.Len(t, entries, 1)
	assert.Equal(t, csi.VolumeHealthErrorType_DATA_LOSS, entries[0].GetStatus())
}

func TestNodeGetStorageHealth_Healthy(t *testing.T) {
	hp := newTestHostPath(t)
	resp, err := hp.NodeGetStorageHealth(context.TODO(), &csi.NodeGetStorageHealthRequest{})
	require.NoError(t, err)
	assert.Empty(t, resp.GetBackendHealth(), "expected healthy (empty) backend health")
}

func TestNodeGetStorageHealth_Unhealthy(t *testing.T) {
	hp := newTestHostPath(t)
	require.NoError(t, WriteStorageHealthMarker(hp.config.StateDir, StorageHealthMarker{
		Status: "STORAGE_UNREACHABLE", Reason: "disk gone", Message: "backend down",
	}))

	resp, err := hp.NodeGetStorageHealth(context.TODO(), &csi.NodeGetStorageHealthRequest{})
	require.NoError(t, err)
	bh := resp.GetBackendHealth()
	require.Len(t, bh, 1)
	assert.Equal(t, csi.StorageHealthErrorType_STORAGE_UNREACHABLE, bh[0].GetStatus())
	assert.Equal(t, "disk gone", bh[0].GetReason())
	assert.Equal(t, "backend down", bh[0].GetMessage())

	// Clearing restores healthy.
	require.NoError(t, ClearStorageHealthMarker(hp.config.StateDir))
	resp, err = hp.NodeGetStorageHealth(context.TODO(), &csi.NodeGetStorageHealthRequest{})
	require.NoError(t, err)
	assert.Empty(t, resp.GetBackendHealth())
}

func TestControllerListVolumeHealth(t *testing.T) {
	hp := newTestHostPath(t)
	mustCreateVolume(t, hp, "vol-a")
	mustCreateVolume(t, hp, "vol-b")
	mustCreateVolume(t, hp, "vol-c")

	// Mark vol-a (both) and vol-c (controller-only) unhealthy; vol-b stays healthy.
	require.NoError(t, WriteVolumeHealthMarker(hp.config.StateDir, "vol-a", ScopeBoth, VolumeHealthMarker{Status: "DEGRADED", Reason: "a"}))
	require.NoError(t, WriteVolumeHealthMarker(hp.config.StateDir, "vol-c", ScopeController, VolumeHealthMarker{Status: "INACCESSIBLE", Reason: "c"}))
	// A node-only marker for vol-b must NOT appear in the controller list.
	require.NoError(t, WriteVolumeHealthMarker(hp.config.StateDir, "vol-b", ScopeNode, VolumeHealthMarker{Status: "DEGRADED"}))
	// A marker for a non-existent volume must be skipped.
	require.NoError(t, WriteVolumeHealthMarker(hp.config.StateDir, "vol-stale", ScopeBoth, VolumeHealthMarker{Status: "DEGRADED"}))

	resp, err := hp.ControllerListVolumeHealth(context.TODO(), &csi.ControllerListVolumeHealthRequest{})
	require.NoError(t, err)
	// Sorted by volID: vol-a, vol-c.
	require.Len(t, resp.GetEntries(), 2)
	assert.Equal(t, "vol-a", resp.GetEntries()[0].GetVolumeId())
	assert.Len(t, resp.GetEntries()[0].GetHealthStatuses(), 1)
	assert.Equal(t, "vol-c", resp.GetEntries()[1].GetVolumeId())
	assert.Empty(t, resp.GetNextToken(), "all entries fit in one page")
}

func TestControllerListVolumeHealth_Pagination(t *testing.T) {
	hp := newTestHostPath(t)
	for _, id := range []string{"v1", "v2", "v3", "v4"} {
		mustCreateVolume(t, hp, id)
		require.NoError(t, WriteVolumeHealthMarker(hp.config.StateDir, id, ScopeBoth, VolumeHealthMarker{Status: "DEGRADED"}))
	}

	// Page size 2.
	resp, err := hp.ControllerListVolumeHealth(context.TODO(), &csi.ControllerListVolumeHealthRequest{MaxEntries: 2})
	require.NoError(t, err)
	require.Len(t, resp.GetEntries(), 2)
	assert.Equal(t, "v1", resp.GetEntries()[0].GetVolumeId())
	assert.Equal(t, "v2", resp.GetEntries()[1].GetVolumeId())
	assert.Equal(t, "3", resp.GetNextToken())

	// Second page.
	resp, err = hp.ControllerListVolumeHealth(context.TODO(), &csi.ControllerListVolumeHealthRequest{MaxEntries: 2, StartingToken: resp.GetNextToken()})
	require.NoError(t, err)
	require.Len(t, resp.GetEntries(), 2)
	assert.Equal(t, "v3", resp.GetEntries()[0].GetVolumeId())
	assert.Equal(t, "v4", resp.GetEntries()[1].GetVolumeId())
	assert.Empty(t, resp.GetNextToken(), "no more pages")
}

func TestDeleteVolumeClearsHealthMarker(t *testing.T) {
	hp := newTestHostPath(t)
	mustCreateVolume(t, hp, "vol-del")

	// Write markers for all three scopes.
	for _, scope := range []VolumeHealthScope{ScopeBoth, ScopeController, ScopeNode} {
		require.NoError(t, WriteVolumeHealthMarker(hp.config.StateDir, "vol-del", scope, VolumeHealthMarker{Status: "DEGRADED"}))
	}

	require.NoError(t, hp.deleteVolume("vol-del"))

	for _, scope := range []VolumeHealthScope{ScopeBoth, ScopeController, ScopeNode} {
		marker, err := ReadVolumeHealthMarker(hp.config.StateDir, "vol-del", scope)
		require.NoError(t, err)
		assert.Nil(t, marker, "scope %s marker should be removed after volume deletion", scope)
	}
}
