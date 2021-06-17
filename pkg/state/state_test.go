/*
Copyright 2021 The Kubernetes Authors.

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

package state

import (
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var grpcError = status.Errorf(codes.NotFound, "not found")

func TestVolumes(t *testing.T) {
	tmp := t.TempDir()
	statefileName := path.Join(tmp, "state.json")

	s, err := New(statefileName)
	require.NoError(t, err, "construct state")
	require.Empty(t, s.GetVolumes(), "initial volumes")

	_, err = s.GetVolumeByID("foo")
	require.Equal(t, codes.NotFound, status.Convert(err).Code(), "GetVolumeByID of non-existent volume")
	require.Contains(t, status.Convert(err).Message(), "foo")

	err = s.UpdateVolume(Volume{VolID: "foo", VolName: "bar"})
	require.NoError(t, err, "add volume")

	s, err = New(statefileName)
	require.NoError(t, err, "reconstruct state")
	_, err = s.GetVolumeByID("foo")
	require.NoError(t, err, "get existing volume by ID")
	_, err = s.GetVolumeByName("bar")
	require.NoError(t, err, "get existing volume by name")

	err = s.DeleteVolume("foo")
	require.NoError(t, err, "delete existing volume")

	err = s.DeleteVolume("foo")
	require.NoError(t, err, "delete non-existent volume")

	require.Empty(t, s.GetVolumes(), "final volumes")
}

func TestSnapshots(t *testing.T) {
	tmp := t.TempDir()
	statefileName := path.Join(tmp, "state.json")

	s, err := New(statefileName)
	require.NoError(t, err, "construct state")
	require.Empty(t, s.GetSnapshots(), "initial snapshots")

	_, err = s.GetSnapshotByID("foo")
	require.Equal(t, codes.NotFound, status.Convert(err).Code(), "GetSnapshotByID of non-existent snapshot")
	require.Contains(t, status.Convert(err).Message(), "foo")

	err = s.UpdateSnapshot(Snapshot{Id: "foo", Name: "bar"})
	require.NoError(t, err, "add snapshot")

	s, err = New(statefileName)
	require.NoError(t, err, "reconstruct state")
	_, err = s.GetSnapshotByID("foo")
	require.NoError(t, err, "get existing snapshot by ID")
	_, err = s.GetSnapshotByName("bar")
	require.NoError(t, err, "get existing snapshot by name")

	err = s.DeleteSnapshot("foo")
	require.NoError(t, err, "delete existing snapshot")

	err = s.DeleteSnapshot("foo")
	require.NoError(t, err, "delete non-existent snapshot")

	require.Empty(t, s.GetSnapshots(), "final snapshots")
}

// TestSnapshotsFromSameSource tests that multiple snapshots from the same
// source can exist at the same time.
func TestSnapshotsFromSameSource(t *testing.T) {
	tmp := t.TempDir()
	statefileName := path.Join(tmp, "state.json")

	s, err := New(statefileName)
	require.NoError(t, err, "construct state")

	err = s.UpdateSnapshot(Snapshot{Id: "foo", Name: "foo-name", VolID: "source"})
	require.NoError(t, err, "add snapshot")
	err = s.UpdateSnapshot(Snapshot{Id: "bar", Name: "bar-name", VolID: "source"})
	require.NoError(t, err, "add snapshot")

	_, err = s.GetSnapshotByID("foo")
	require.NoError(t, err, "get existing snapshot by ID 'foo'")
	_, err = s.GetSnapshotByName("foo-name")
	require.NoError(t, err, "get existing snapshot by name 'foo-name'")
	_, err = s.GetSnapshotByID("bar")
	require.NoError(t, err, "get existing snapshot by ID 'bar'")
	_, err = s.GetSnapshotByName("bar-name")
	require.NoError(t, err, "get existing snapshot by name 'bar-name'")

	// Make sure it still works after reconstruction
	s, err = New(statefileName)
	require.NoError(t, err, "reconstruct state")
	_, err = s.GetSnapshotByID("foo")
	require.NoError(t, err, "get existing snapshot by ID 'foo'")
	_, err = s.GetSnapshotByName("foo-name")
	require.NoError(t, err, "get existing snapshot by name 'foo-name'")
	_, err = s.GetSnapshotByID("bar")
	require.NoError(t, err, "get existing snapshot by ID 'bar'")
	_, err = s.GetSnapshotByName("bar-name")
	require.NoError(t, err, "get existing snapshot by name 'bar-name'")
}
