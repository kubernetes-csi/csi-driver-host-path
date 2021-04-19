package state

import (
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/assert"
)

const (
	tmpDataRoot = "/tmp/csi-data-dir"
)

var (
	testHostPathDriverState = NewHostPathDriverState(tmpDataRoot)
	testVolumes             = make(HostPathVolumes)
	testSnapshots           = make(HostPathSnapshots)
)

func TestCreateVolumeByID(t *testing.T) {
	maxVolumeSize := 100 * 1024 * 1024 // 10Mi
	testVolume := HostPathVolume{
		VolID:         "volume-1",
		VolName:       "volume-1",
		VolSize:       int64(maxVolumeSize),
		VolAccessType: MountAccess,
		Ephemeral:     false,
		Kind:          "",
		VolPath:       "/tmp/csi-data-dir/volume-1",
	}

	cap := make(Capacity)
	_, err := testHostPathDriverState.CreateVolume(
		testVolume.VolID, testVolume.VolName, int64(maxVolumeSize),
		testVolume.VolAccessType, false, "",
		int64(maxVolumeSize), cap)

	assert.Nil(t, err)
	testVolumes[testVolume.VolID] = testVolume

	// Get volume from memory
	volumeresult, err := testHostPathDriverState.GetVolumeByID(testVolume.VolID)
	assert.Nil(t, err)
	assert.EqualValues(t, testVolume, volumeresult)
}

func TestListVolumes(t *testing.T) {
	maxVolumeSize := 100 * 1024 * 1024 * 1024 // 1Gi
	testVolume := HostPathVolume{
		VolID:         "volume-2",
		VolName:       "volume-2",
		VolSize:       int64(maxVolumeSize),
		VolAccessType: MountAccess,
		Ephemeral:     false,
		Kind:          "",
		VolPath:       "/tmp/csi-data-dir/volume-2",
	}

	cap := make(Capacity)
	_, err := testHostPathDriverState.CreateVolume(
		testVolume.VolID, testVolume.VolName, int64(maxVolumeSize),
		testVolume.VolAccessType, false, "",
		int64(maxVolumeSize), cap)

	assert.Nil(t, err)
	testVolumes[testVolume.VolID] = testVolume

	volumesResult, err := testHostPathDriverState.ListVolumes()
	assert.Nil(t, err)
	assert.EqualValues(t, testVolumes, volumesResult)
}

func TestUpdateVolumes(t *testing.T) {
	maxVolumeSize := 100 * 1024 // 1Ki
	testVolume := HostPathVolume{
		VolID:         "volume-3",
		VolName:       "volume-3",
		VolSize:       int64(maxVolumeSize),
		VolAccessType: MountAccess,
		Ephemeral:     false,
		Kind:          "",
		VolPath:       "/tmp/csi-data-dir/volume-3",
	}
	cap := make(Capacity)
	_, err := testHostPathDriverState.CreateVolume(
		testVolume.VolID, testVolume.VolName, int64(maxVolumeSize),
		testVolume.VolAccessType, false, "",
		int64(maxVolumeSize), cap)
	assert.Nil(t, err)
	testVolumes[testVolume.VolID] = testVolume

	testVolume.Ephemeral = true
	err = testHostPathDriverState.UpdateVolume(testVolume.VolID, testVolume)
	assert.Nil(t, err)
	// Get volume from memory
	volumeresult, err := testHostPathDriverState.GetVolumeByID(testVolume.VolID)
	assert.Nil(t, err)
	assert.EqualValues(t, testVolume, volumeresult)
}

func TestDeleteVolumes(t *testing.T) {
	err := testHostPathDriverState.DeleteVolume("volume-3", make(Capacity))
	assert.Nil(t, err)

	delete(testVolumes, "volume-3")

	volumesResult, err := testHostPathDriverState.ListVolumes()
	assert.Nil(t, err)
	assert.EqualValues(t, testVolumes, volumesResult)
}

func TestLoadVolumesFromFile(t *testing.T) {
	err := testHostPathDriverState.loadVolumesFromFile()
	assert.Nil(t, err)
	assert.EqualValues(t, testVolumes, testHostPathDriverState.HostPathVolumes)
}

func TestGetSnapshotID(t *testing.T) {
	testCases := []struct {
		name               string
		inputPath          string
		expectedIsSnapshot bool
		expectedSnapshotID string
	}{
		{
			name:               "should recognize foo.snap as a valid snapshot with ID foo",
			inputPath:          "foo.snap",
			expectedIsSnapshot: true,
			expectedSnapshotID: "foo",
		},
		{
			name:               "should recognize baz.tar.gz as an invalid snapshot",
			inputPath:          "baz.tar.gz",
			expectedIsSnapshot: false,
			expectedSnapshotID: "",
		},
		{
			name:               "should recognize baz.tar.snap as a valid snapshot with ID baz.tar",
			inputPath:          "baz.tar.snap",
			expectedIsSnapshot: true,
			expectedSnapshotID: "baz.tar",
		},
		{
			name:               "should recognize baz.tar.snap.snap as a valid snapshot with ID baz.tar.snap",
			inputPath:          "baz.tar.snap.snap",
			expectedIsSnapshot: true,
			expectedSnapshotID: "baz.tar.snap",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualIsSnapshot, actualSnapshotID := getSnapshotID(tc.inputPath)
			if actualIsSnapshot != tc.expectedIsSnapshot {
				t.Errorf("unexpected result for path %s, Want: %t, Got: %t", tc.inputPath, tc.expectedIsSnapshot, actualIsSnapshot)
			}
			if actualSnapshotID != tc.expectedSnapshotID {
				t.Errorf("unexpected snapshotID for path %s, Want: %s; Got :%s", tc.inputPath, tc.expectedSnapshotID, actualSnapshotID)
			}
		})
	}
}

func TestCreateSnapshot(t *testing.T) {
	testSnapshot := HostPathSnapshot{
		Name:         "snapshot-1",
		Id:           "snapshot-1",
		VolID:        "volume-1",
		Path:         getSnapshotPath("snapshot-1", testHostPathDriverState.DataRoot),
		CreationTime: ptypes.TimestampNow(),
		SizeBytes:    10 * 1024 * 1024,
		ReadyToUse:   true,
	}

	snapshot, err := testHostPathDriverState.CreateSnapshot(
		testSnapshot.Name, testSnapshot.Id,
		testSnapshot.VolID, testSnapshot.Path,
		testSnapshot.CreationTime,
		testSnapshot.SizeBytes, testSnapshot.ReadyToUse)

	assert.Nil(t, err)
	assert.EqualValues(t, testSnapshot, snapshot)
	testSnapshots[testSnapshot.Id] = snapshot
}

func TestListSnapshots(t *testing.T) {
	testSnapshot := HostPathSnapshot{
		Name:         "snapshot-2",
		Id:           "snapshot-2",
		VolID:        "volume-2",
		Path:         getSnapshotPath("snapshot-2", testHostPathDriverState.DataRoot),
		CreationTime: ptypes.TimestampNow(),
		SizeBytes:    10 * 1024 * 1024,
		ReadyToUse:   true,
	}

	snapshot, err := testHostPathDriverState.CreateSnapshot(
		testSnapshot.Name, testSnapshot.Id,
		testSnapshot.VolID, testSnapshot.Path,
		testSnapshot.CreationTime,
		testSnapshot.SizeBytes, testSnapshot.ReadyToUse)

	assert.Nil(t, err)
	assert.EqualValues(t, testSnapshot, snapshot)
	testSnapshots[testSnapshot.Id] = snapshot

	snapshots, err := testHostPathDriverState.ListSnapshots()
	assert.Nil(t, err)
	assert.EqualValues(t, testSnapshots, snapshots)
}

func TestGetSnapshotByID(t *testing.T) {
	expectedSnapshot, ok := testSnapshots["snapshot-1"]
	assert.True(t, ok)

	actualSnapshot, err := testHostPathDriverState.GetSnapshotByID("snapshot-1")
	assert.Nil(t, err)

	assert.EqualValues(t, expectedSnapshot, actualSnapshot)
}

func TestGetSnapshotByName(t *testing.T) {
	expectedSnapshot, ok := testSnapshots["snapshot-2"]
	assert.True(t, ok)

	actualSnapshot, err := testHostPathDriverState.GetSnapshotByName("snapshot-2")
	assert.Nil(t, err)

	assert.EqualValues(t, expectedSnapshot, actualSnapshot)
}

func TestDeleteSnapshot(t *testing.T) {
	delete(testSnapshots, "snapshot-2")

	err := testHostPathDriverState.DeleteSnapshot("snapshot-2")
	assert.Nil(t, err)

	snapshots, err := testHostPathDriverState.ListSnapshots()
	assert.Nil(t, err)
	assert.EqualValues(t, testSnapshots, snapshots)
}

func TestLoadSnapshotsFromFile(t *testing.T) {
	err := testHostPathDriverState.loadSnapshotsFromFile()
	assert.Nil(t, err)

	snapshots, err := testHostPathDriverState.ListSnapshots()
	assert.Nil(t, err)
	assert.EqualValues(t, testSnapshots, snapshots)
}

// func TestLoadStateFromFile(t *testing.T) {
// 	newdriverState, err := LoadStateFromFile(tmpDataRoot)
// 	assert.Nil(t, err)
// 	assert.EqualValues(t, *testHostPathDriverState, *newdriverState)
// }
