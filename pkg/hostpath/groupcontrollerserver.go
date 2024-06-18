/*
Copyright 2023 The Kubernetes Authors.

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

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/protobuf/ptypes"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"

	"github.com/kubernetes-csi/csi-driver-host-path/pkg/state"
)

func (hp *hostPath) GroupControllerGetCapabilities(context.Context, *csi.GroupControllerGetCapabilitiesRequest) (*csi.GroupControllerGetCapabilitiesResponse, error) {
	return &csi.GroupControllerGetCapabilitiesResponse{
		Capabilities: []*csi.GroupControllerServiceCapability{{
			Type: &csi.GroupControllerServiceCapability_Rpc{
				Rpc: &csi.GroupControllerServiceCapability_RPC{
					Type: csi.GroupControllerServiceCapability_RPC_CREATE_DELETE_GET_VOLUME_GROUP_SNAPSHOT,
				},
			},
		}},
	}, nil
}

func (hp *hostPath) CreateVolumeGroupSnapshot(ctx context.Context, req *csi.CreateVolumeGroupSnapshotRequest) (*csi.CreateVolumeGroupSnapshotResponse, error) {
	if err := hp.validateGroupControllerServiceRequest(csi.GroupControllerServiceCapability_RPC_CREATE_DELETE_GET_VOLUME_GROUP_SNAPSHOT); err != nil {
		klog.V(3).Infof("invalid create volume group snapshot req: %v", req)
		return nil, err
	}

	if len(req.GetName()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Name missing in request")
	}
	// Check arguments
	if len(req.GetSourceVolumeIds()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "SourceVolumeIds missing in request")
	}

	// Lock before acting on global state. A production-quality
	// driver might use more fine-grained locking.
	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	// Need to check for already existing groupsnapshot name, and if found check for the
	// requested sourceVolumeIds and sourceVolumeIds of groupsnapshot that has been created.
	if exGS, err := hp.state.GetGroupSnapshotByName(req.GetName()); err == nil {
		// Since err is nil, it means the groupsnapshot with the same name already exists. Need
		// to check if the sourceVolumeIds of existing groupsnapshot is the same as in new request.

		if !exGS.MatchesSourceVolumeIDs(req.GetSourceVolumeIds()) {
			return nil, status.Errorf(codes.AlreadyExists, "group snapshot with the same name: %s but with different SourceVolumeIds already exist", req.GetName())
		}

		// same groupsnapshot has been created.
		snapshots := make([]*csi.Snapshot, len(exGS.SnapshotIDs))
		readyToUse := true

		for i, snapshotID := range exGS.SnapshotIDs {
			snapshot, err := hp.state.GetSnapshotByID(snapshotID)
			if err != nil {
				return nil, err
			}

			snapshots[i] = &csi.Snapshot{
				SizeBytes:       snapshot.SizeBytes,
				CreationTime:    snapshot.CreationTime,
				ReadyToUse:      snapshot.ReadyToUse,
				GroupSnapshotId: snapshot.GroupSnapshotID,
			}

			readyToUse = readyToUse && snapshot.ReadyToUse
		}

		return &csi.CreateVolumeGroupSnapshotResponse{
			GroupSnapshot: &csi.VolumeGroupSnapshot{
				GroupSnapshotId: exGS.Id,
				Snapshots:       snapshots,
				CreationTime:    exGS.CreationTime,
				ReadyToUse:      readyToUse,
			},
		}, nil
	}

	groupSnapshot := state.GroupSnapshot{
		Name:            req.GetName(),
		Id:              uuid.NewUUID().String(),
		CreationTime:    ptypes.TimestampNow(),
		SnapshotIDs:     make([]string, len(req.GetSourceVolumeIds())),
		SourceVolumeIDs: make([]string, len(req.GetSourceVolumeIds())),
		ReadyToUse:      true,
	}

	copy(groupSnapshot.SourceVolumeIDs, req.GetSourceVolumeIds())

	snapshots := make([]*csi.Snapshot, len(req.GetSourceVolumeIds()))

	// TODO: defer a cleanup function to remove snapshots in case of a failure

	for i, volumeID := range req.GetSourceVolumeIds() {
		hostPathVolume, err := hp.state.GetVolumeByID(volumeID)
		if err != nil {
			return nil, err
		}

		snapshotID := uuid.NewUUID().String()
		file := hp.getSnapshotPath(snapshotID)

		if err := hp.createSnapshotFromVolume(hostPathVolume, file); err != nil {
			return nil, err
		}

		klog.V(4).Infof("create volume snapshot %s", file)
		snapshot := state.Snapshot{}
		snapshot.Name = req.GetName() + "-" + volumeID
		snapshot.Id = snapshotID
		snapshot.VolID = volumeID
		snapshot.Path = file
		snapshot.CreationTime = groupSnapshot.CreationTime
		snapshot.SizeBytes = hostPathVolume.VolSize
		snapshot.ReadyToUse = true
		snapshot.GroupSnapshotID = groupSnapshot.Id

		hp.state.UpdateSnapshot(snapshot)

		groupSnapshot.SnapshotIDs[i] = snapshotID

		snapshots[i] = &csi.Snapshot{
			SizeBytes:       hostPathVolume.VolSize,
			SnapshotId:      snapshotID,
			SourceVolumeId:  volumeID,
			CreationTime:    groupSnapshot.CreationTime,
			ReadyToUse:      true,
			GroupSnapshotId: groupSnapshot.Id,
		}
	}

	if err := hp.state.UpdateGroupSnapshot(groupSnapshot); err != nil {
		return nil, err
	}

	return &csi.CreateVolumeGroupSnapshotResponse{
		GroupSnapshot: &csi.VolumeGroupSnapshot{
			GroupSnapshotId: groupSnapshot.Id,
			Snapshots:       snapshots,
			CreationTime:    groupSnapshot.CreationTime,
			ReadyToUse:      groupSnapshot.ReadyToUse,
		},
	}, nil
}

func (hp *hostPath) DeleteVolumeGroupSnapshot(ctx context.Context, req *csi.DeleteVolumeGroupSnapshotRequest) (*csi.DeleteVolumeGroupSnapshotResponse, error) {
	if err := hp.validateGroupControllerServiceRequest(csi.GroupControllerServiceCapability_RPC_CREATE_DELETE_GET_VOLUME_GROUP_SNAPSHOT); err != nil {
		klog.V(3).Infof("invalid delete volume group snapshot req: %v", req)
		return nil, err
	}

	// Check arguments
	if len(req.GetGroupSnapshotId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "GroupSnapshot ID missing in request")
	}

	groupSnapshotID := req.GetGroupSnapshotId()

	// Lock before acting on global state. A production-quality
	// driver might use more fine-grained locking.
	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	groupSnapshot, err := hp.state.GetGroupSnapshotByID(groupSnapshotID)
	if err != nil {
		// ok if NotFound, the VolumeGroupSnapshot was deleted already
		if status.Code(err) == codes.NotFound {
			return &csi.DeleteVolumeGroupSnapshotResponse{}, nil
		}

		return nil, err
	}

	for _, snapshotID := range groupSnapshot.SnapshotIDs {
		klog.V(4).Infof("deleting snapshot %s", snapshotID)
		path := hp.getSnapshotPath(snapshotID)
		os.RemoveAll(path)

		if err := hp.state.DeleteSnapshot(snapshotID); err != nil {
			return nil, err
		}
	}

	klog.V(4).Infof("deleting groupsnapshot %s", groupSnapshotID)
	if err := hp.state.DeleteGroupSnapshot(groupSnapshotID); err != nil {
		return nil, err
	}

	return &csi.DeleteVolumeGroupSnapshotResponse{}, nil
}

func (hp *hostPath) GetVolumeGroupSnapshot(ctx context.Context, req *csi.GetVolumeGroupSnapshotRequest) (*csi.GetVolumeGroupSnapshotResponse, error) {
	if err := hp.validateGroupControllerServiceRequest(csi.GroupControllerServiceCapability_RPC_CREATE_DELETE_GET_VOLUME_GROUP_SNAPSHOT); err != nil {
		klog.V(3).Infof("invalid get volume group snapshot req: %v", req)
		return nil, err
	}

	groupSnapshotID := req.GetGroupSnapshotId()

	// Check arguments
	if len(groupSnapshotID) == 0 {
		return nil, status.Error(codes.InvalidArgument, "GroupSnapshot ID missing in request")
	}

	// Lock before acting on global state. A production-quality
	// driver might use more fine-grained locking.
	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	groupSnapshot, err := hp.state.GetGroupSnapshotByID(groupSnapshotID)
	if err != nil {
		return nil, err
	}

	if !groupSnapshot.MatchesSnapshotIDs(req.GetSnapshotIds()) {
		return nil, status.Error(codes.InvalidArgument, "Snapshot IDs do not match the GroupSnapshot IDs")
	}

	snapshots := make([]*csi.Snapshot, len(groupSnapshot.SnapshotIDs))
	for i, snapshotID := range groupSnapshot.SnapshotIDs {
		snapshot, err := hp.state.GetSnapshotByID(snapshotID)
		if err != nil {
			return nil, err
		}

		snapshots[i] = &csi.Snapshot{
			SizeBytes:       snapshot.SizeBytes,
			SnapshotId:      snapshotID,
			SourceVolumeId:  snapshot.VolID,
			CreationTime:    snapshot.CreationTime,
			ReadyToUse:      snapshot.ReadyToUse,
			GroupSnapshotId: snapshot.GroupSnapshotID,
		}
	}

	return &csi.GetVolumeGroupSnapshotResponse{
		GroupSnapshot: &csi.VolumeGroupSnapshot{
			GroupSnapshotId: groupSnapshotID,
			Snapshots:       snapshots,
			CreationTime:    groupSnapshot.CreationTime,
			ReadyToUse:      groupSnapshot.ReadyToUse,
		},
	}, nil
}

func (hp *hostPath) validateGroupControllerServiceRequest(c csi.GroupControllerServiceCapability_RPC_Type) error {
	if c == csi.GroupControllerServiceCapability_RPC_UNKNOWN {
		return nil
	}

	if c == csi.GroupControllerServiceCapability_RPC_CREATE_DELETE_GET_VOLUME_GROUP_SNAPSHOT {
		return nil
	}

	return status.Errorf(codes.InvalidArgument, "unsupported capability %s", c)
}
