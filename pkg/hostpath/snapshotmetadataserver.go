/*
Copyright 2024 The Kubernetes Authors.

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
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-driver-host-path/pkg/state"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

func (hp *hostPath) GetMetadataAllocated(req *csi.GetMetadataAllocatedRequest, stream csi.SnapshotMetadata_GetMetadataAllocatedServer) error {
	ctx := stream.Context()
	// Check arguments
	snapID := req.GetSnapshotId()
	if len(snapID) == 0 {
		return status.Error(codes.InvalidArgument, "SnapshotID missing in request")
	}

	// Load snapshots
	source, err := hp.state.GetSnapshotByID(snapID)
	if err != nil {
		return status.Error(codes.Internal, "Cannot find the snapshot")
	}
	if !source.ReadyToUse {
		return status.Error(codes.Unavailable, fmt.Sprintf("snapshot %v is not yet ready to use", snapID))
	}

	vol, err := hp.state.GetVolumeByID(source.VolID)
	if err != nil {
		return err
	}
	if vol.VolAccessType != state.BlockAccess {
		return status.Error(codes.InvalidArgument, "Source volume does not have block mode access type")
	}

	allocatedBlocks := make(chan []*csi.BlockMetadata, 100)
	go func() {
		err := hp.getAllocatedBlockMetadata(ctx, hp.getSnapshotPath(snapID), req.StartingOffset, state.BlockSizeBytes, req.MaxResults, allocatedBlocks)
		if err != nil {
			klog.Errorf("failed to get allocated block metadata: %v", err)
		}
	}()

	for {
		select {
		case cb, ok := <-allocatedBlocks:
			if !ok {
				klog.V(4).Info("channel closed, returning")
				return nil
			}
			resp := csi.GetMetadataAllocatedResponse{
				BlockMetadataType:   csi.BlockMetadataType_FIXED_LENGTH,
				VolumeCapacityBytes: vol.VolSize,
				BlockMetadata:       cb,
			}
			if err := stream.Send(&resp); err != nil {
				return err
			}
		case <-ctx.Done():
			klog.V(4).Info("received cancellation signal, returning")
			return nil
		}
	}
}

func (hp *hostPath) GetMetadataDelta(req *csi.GetMetadataDeltaRequest, stream csi.SnapshotMetadata_GetMetadataDeltaServer) error {
	ctx := stream.Context()
	// Check arguments
	baseSnapID := req.GetBaseSnapshotId()
	targetSnapID := req.GetTargetSnapshotId()
	if len(baseSnapID) == 0 {
		return status.Error(codes.InvalidArgument, "BaseSnapshotID missing in request")
	}
	if len(targetSnapID) == 0 {
		return status.Error(codes.InvalidArgument, "TargetSnapshotID missing in request")
	}

	// Load snapshots
	source, err := hp.state.GetSnapshotByID(baseSnapID)
	if err != nil {
		return status.Error(codes.Internal, "Cannot find the source snapshot")
	}
	target, err := hp.state.GetSnapshotByID(targetSnapID)
	if err != nil {
		return status.Error(codes.Internal, "Cannot find the target snapshot")
	}

	if !source.ReadyToUse {
		return status.Error(codes.Unavailable, fmt.Sprintf("snapshot %v is not yet ready to use", baseSnapID))
	}
	if !target.ReadyToUse {
		return status.Error(codes.Unavailable, fmt.Sprintf("snapshot %v is not yet ready to use", targetSnapID))
	}

	if source.VolID != target.VolID {
		return status.Error(codes.InvalidArgument, "Snapshots don't belong to the same Volume")
	}
	vol, err := hp.state.GetVolumeByID(source.VolID)
	if err != nil {
		return err
	}
	if vol.VolAccessType != state.BlockAccess {
		return status.Error(codes.InvalidArgument, "Source volume does not have block mode access type")
	}

	changedBlocks := make(chan []*csi.BlockMetadata, 100)
	go func() {
		err := hp.getChangedBlockMetadata(ctx, hp.getSnapshotPath(baseSnapID), hp.getSnapshotPath(targetSnapID), req.StartingOffset, state.BlockSizeBytes, req.MaxResults, changedBlocks)
		if err != nil {
			klog.Errorf("failed to get changed block metadata: %v", err)
		}
	}()

	for {
		select {
		case cb, ok := <-changedBlocks:
			if !ok {
				klog.V(4).Info("channel closed, returning")
				return nil
			}
			resp := csi.GetMetadataDeltaResponse{
				BlockMetadataType:   csi.BlockMetadataType_FIXED_LENGTH,
				VolumeCapacityBytes: vol.VolSize,
				BlockMetadata:       cb,
			}
			if err := stream.Send(&resp); err != nil {
				return err
			}
		case <-ctx.Done():
			klog.V(4).Info("received cancellation signal, returning")
			return nil
		}
	}
}
