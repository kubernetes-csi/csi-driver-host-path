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

// Package state manages the internal state of the driver which needs to be maintained
// across driver restarts.
package state

import (
	"encoding/json"
	"errors"
	"os"
	"sort"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccessType int

const (
	MountAccess AccessType = iota
	BlockAccess
)

type Volume struct {
	VolName        string
	VolID          string
	VolSize        int64
	VolPath        string
	VolAccessType  AccessType
	ParentVolID    string
	ParentSnapID   string
	Ephemeral      bool
	NodeID         string
	Kind           string
	ReadOnlyAttach bool
	Attached       bool
	// Staged contains the staging target path at which the volume
	// was staged. A set of paths is used for consistency
	// with Published.
	Staged Strings
	// Published contains the target paths where the volume
	// was published.
	Published Strings
}

type Snapshot struct {
	Name            string
	Id              string
	VolID           string
	Path            string
	CreationTime    *timestamp.Timestamp
	SizeBytes       int64
	ReadyToUse      bool
	GroupSnapshotID string
}

type GroupSnapshot struct {
	Name            string
	Id              string
	SnapshotIDs     []string
	SourceVolumeIDs []string
	CreationTime    *timestamp.Timestamp
	ReadyToUse      bool
}

// State is the interface that the rest of the code has to use to
// access and change state. All error messages contain gRPC
// status codes and can be returned without wrapping.
type State interface {
	// GetVolumeByID retrieves a volume by its unique ID or returns
	// an error including that ID when not found.
	GetVolumeByID(volID string) (Volume, error)

	// GetVolumeByName retrieves a volume by its name or returns
	// an error including that name when not found.
	GetVolumeByName(volName string) (Volume, error)

	// GetVolumes returns all currently existing volumes.
	GetVolumes() []Volume

	// UpdateVolume updates the existing hostpath volume,
	// identified by its volume ID, or adds it if it does
	// not exist yet.
	UpdateVolume(volume Volume) error

	// DeleteVolume deletes the volume with the given
	// volume ID. It is not an error when such a volume
	// does not exist.
	DeleteVolume(volID string) error

	// GetSnapshotByID retrieves a snapshot by its unique ID or returns
	// an error including that ID when not found.
	GetSnapshotByID(snapshotID string) (Snapshot, error)

	// GetSnapshotByName retrieves a snapshot by its name or returns
	// an error including that name when not found.
	GetSnapshotByName(volName string) (Snapshot, error)

	// GetSnapshots returns all currently existing snapshots.
	GetSnapshots() []Snapshot

	// UpdateSnapshot updates the existing hostpath snapshot,
	// identified by its snapshot ID, or adds it if it does
	// not exist yet.
	UpdateSnapshot(snapshot Snapshot) error

	// DeleteSnapshot deletes the snapshot with the given
	// snapshot ID. It is not an error when such a snapshot
	// does not exist.
	DeleteSnapshot(snapshotID string) error

	// GetGroupSnapshotByID retrieves a groupsnapshot by its unique ID or
	// returns an error including that ID when not found.
	GetGroupSnapshotByID(vgsID string) (GroupSnapshot, error)

	// GetGroupSnapshotByName retrieves a groupsnapshot by its name or
	// returns an error including that name when not found.
	GetGroupSnapshotByName(volName string) (GroupSnapshot, error)

	// GetGroupSnapshots returns all currently existing groupsnapshots.
	GetGroupSnapshots() []GroupSnapshot

	// UpdateGroupSnapshot updates the existing hostpath groupsnapshot,
	// identified by its snapshot ID, or adds it if it does not exist yet.
	UpdateGroupSnapshot(snapshot GroupSnapshot) error

	// DeleteGroupSnapshot deletes the groupsnapshot with the given
	// groupsnapshot ID. It is not an error when such a groupsnapshot does
	// not exist.
	DeleteGroupSnapshot(groupSnapshotID string) error
}

type resources struct {
	Volumes        []Volume
	Snapshots      []Snapshot
	GroupSnapshots []GroupSnapshot
}

type state struct {
	resources

	statefilePath string
}

var _ State = &state{}

// New retrieves the complete state of the driver from the file if given
// and then ensures that all changes are mirrored immediately in the
// given file. If not given, the initial state is empty and changes
// are not saved.
func New(statefilePath string) (State, error) {
	s := &state{
		statefilePath: statefilePath,
	}

	return s, s.restore()
}

func (s *state) dump() error {
	data, err := json.Marshal(&s.resources)
	if err != nil {
		return status.Errorf(codes.Internal, "error encoding volumes and snapshots: %v", err)
	}
	if err := os.WriteFile(s.statefilePath, data, 0600); err != nil {
		return status.Errorf(codes.Internal, "error writing state file: %v", err)
	}
	return nil
}

func (s *state) restore() error {
	s.Volumes = nil
	s.Snapshots = nil

	data, err := os.ReadFile(s.statefilePath)
	switch {
	case errors.Is(err, os.ErrNotExist):
		// Nothing to do.
		return nil
	case err != nil:
		return status.Errorf(codes.Internal, "error reading state file: %v", err)
	}
	if err := json.Unmarshal(data, &s.resources); err != nil {
		return status.Errorf(codes.Internal, "error encoding volumes and snapshots from state file %q: %v", s.statefilePath, err)
	}
	return nil
}

func (s *state) GetVolumeByID(volID string) (Volume, error) {
	for _, volume := range s.Volumes {
		if volume.VolID == volID {
			return volume, nil
		}
	}
	return Volume{}, status.Errorf(codes.NotFound, "volume id %s does not exist in the volumes list", volID)
}

func (s *state) GetVolumeByName(volName string) (Volume, error) {
	for _, volume := range s.Volumes {
		if volume.VolName == volName {
			return volume, nil
		}
	}
	return Volume{}, status.Errorf(codes.NotFound, "volume name %s does not exist in the volumes list", volName)
}

func (s *state) GetVolumes() []Volume {
	volumes := make([]Volume, len(s.Volumes))
	copy(volumes, s.Volumes)
	return volumes
}

func (s *state) UpdateVolume(update Volume) error {
	for i, volume := range s.Volumes {
		if volume.VolID == update.VolID {
			s.Volumes[i] = update
			return s.dump()
		}
	}
	s.Volumes = append(s.Volumes, update)
	return s.dump()
}

func (s *state) DeleteVolume(volID string) error {
	for i, volume := range s.Volumes {
		if volume.VolID == volID {
			s.Volumes = append(s.Volumes[:i], s.Volumes[i+1:]...)
			return s.dump()
		}
	}
	return nil
}

func (s *state) GetSnapshotByID(snapshotID string) (Snapshot, error) {
	for _, snapshot := range s.Snapshots {
		if snapshot.Id == snapshotID {
			return snapshot, nil
		}
	}
	return Snapshot{}, status.Errorf(codes.NotFound, "snapshot id %s does not exist in the snapshots list", snapshotID)
}

func (s *state) GetSnapshotByName(name string) (Snapshot, error) {
	for _, snapshot := range s.Snapshots {
		if snapshot.Name == name {
			return snapshot, nil
		}
	}
	return Snapshot{}, status.Errorf(codes.NotFound, "snapshot name %s does not exist in the snapshots list", name)
}

func (s *state) GetSnapshots() []Snapshot {
	snapshots := make([]Snapshot, len(s.Snapshots))
	copy(snapshots, s.Snapshots)
	return snapshots
}

func (s *state) UpdateSnapshot(update Snapshot) error {
	for i, snapshot := range s.Snapshots {
		if snapshot.Id == update.Id {
			s.Snapshots[i] = update
			return s.dump()
		}
	}
	s.Snapshots = append(s.Snapshots, update)
	return s.dump()
}

func (s *state) DeleteSnapshot(snapshotID string) error {
	for i, snapshot := range s.Snapshots {
		if snapshot.Id == snapshotID {
			s.Snapshots = append(s.Snapshots[:i], s.Snapshots[i+1:]...)
			return s.dump()
		}
	}
	return nil
}

func (s *state) GetGroupSnapshotByID(groupSnapshotID string) (GroupSnapshot, error) {
	for _, groupSnapshot := range s.GroupSnapshots {
		if groupSnapshot.Id == groupSnapshotID {
			return groupSnapshot, nil
		}
	}
	return GroupSnapshot{}, status.Errorf(codes.NotFound, "groupsnapshot id %s does not exist in the groupsnapshots list", groupSnapshotID)
}

func (s *state) GetGroupSnapshotByName(name string) (GroupSnapshot, error) {
	for _, groupSnapshot := range s.GroupSnapshots {
		if groupSnapshot.Name == name {
			return groupSnapshot, nil
		}
	}
	return GroupSnapshot{}, status.Errorf(codes.NotFound, "groupsnapshot name %s does not exist in the groupsnapshots list", name)
}

func (s *state) GetGroupSnapshots() []GroupSnapshot {
	groupSnapshots := make([]GroupSnapshot, len(s.GroupSnapshots))
	for i, groupSnapshot := range s.GroupSnapshots {
		groupSnapshots[i] = groupSnapshot
	}
	return groupSnapshots
}

func (s *state) UpdateGroupSnapshot(update GroupSnapshot) error {
	for i, groupSnapshot := range s.GroupSnapshots {
		if groupSnapshot.Id == update.Id {
			s.GroupSnapshots[i] = update
			return s.dump()
		}
	}
	s.GroupSnapshots = append(s.GroupSnapshots, update)
	return s.dump()
}

func (s *state) DeleteGroupSnapshot(groupSnapshotID string) error {
	for i, groupSnapshot := range s.GroupSnapshots {
		if groupSnapshot.Id == groupSnapshotID {
			s.GroupSnapshots = append(s.GroupSnapshots[:i], s.GroupSnapshots[i+1:]...)
			return s.dump()
		}
	}
	return nil
}

func (gs *GroupSnapshot) MatchesSourceVolumeIDs(sourceVolumeIDs []string) bool {
	return equalIDs(gs.SourceVolumeIDs, sourceVolumeIDs)
}

func (gs *GroupSnapshot) MatchesSnapshotIDs(snapshotIDs []string) bool {
	return equalIDs(gs.SnapshotIDs, snapshotIDs)
}

func equalIDs(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// sort slices so that values are at the same location
	sort.Strings(a)
	sort.Strings(b)

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
