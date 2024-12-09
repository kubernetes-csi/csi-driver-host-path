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
	"context"
	"errors"
	"fmt"
	"io"

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
		return status.Error(codes.NotFound, "cannot find the snapshot")
	}
	if !source.ReadyToUse {
		return status.Error(codes.Unavailable, fmt.Sprintf("snapshot %v is not yet ready to use", snapID))
	}

	vol, err := hp.state.GetVolumeByID(source.VolID)
	if err != nil {
		return err
	}
	if vol.VolAccessType != state.BlockAccess {
		return status.Error(codes.InvalidArgument, "source volume does not have block mode access type")
	}

	br, err := newFileBlockReader(
		"",
		hp.getSnapshotPath(snapID),
		req.StartingOffset,
		state.BlockSizeBytes,
		hp.config.SnapshotMetadataBlockType,
		req.MaxResults,
	)
	if err != nil {
		klog.Errorf("failed initialize file block reader: %v", err)
		return status.Error(codes.Internal, "failed initialize file block reader")
	}
	defer br.Close()
	if err := br.seekToStartingOffset(); err != nil {
		return status.Error(codes.OutOfRange, fmt.Sprintf("failed to seek to starting offset: %v", err.Error()))
	}

	// Read all blocks can find allocated blocks till EOF in chunks of size == maxSize
	for {
		cb, cbErr := br.getChangedBlockMetadata(ctx)
		if cbErr != nil {
			if errors.Is(cbErr, context.Canceled) {
				klog.V(4).Info("context canceled while getting allocated block metadata, returning")
				return nil
			}
			if errors.Is(cbErr, context.DeadlineExceeded) {
				klog.V(4).Info("context deadline exceeded while getting allocated block metadata, returning")
				return nil
			}
			if errors.Is(cbErr, io.EOF) {
				klog.V(4).Info("reached EOF while getting allocated block metadata, returning")
				// send allocated blocks found till EOF
				if err := sendGetMetadataAllocatedResponse(stream, vol.VolSize, hp.config.SnapshotMetadataBlockType, cb); err != nil {
					return err
				}
				return nil
			}
			klog.Errorf("Failed to get allocated block metadata: %v", cbErr)
			return status.Error(codes.Internal, "failed to get allocated block metadata")
		}
		// stream response to client
		if err := sendGetMetadataAllocatedResponse(stream, vol.VolSize, hp.config.SnapshotMetadataBlockType, cb); err != nil {
			return err
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
		return status.Error(codes.NotFound, "cannot find the source snapshot")
	}
	target, err := hp.state.GetSnapshotByID(targetSnapID)
	if err != nil {
		return status.Error(codes.NotFound, "cannot find the target snapshot")
	}

	if !source.ReadyToUse {
		return status.Error(codes.Unavailable, fmt.Sprintf("snapshot %v is not yet ready to use", baseSnapID))
	}
	if !target.ReadyToUse {
		return status.Error(codes.Unavailable, fmt.Sprintf("snapshot %v is not yet ready to use", targetSnapID))
	}

	if source.VolID != target.VolID {
		return status.Error(codes.InvalidArgument, "snapshots don't belong to the same Volume")
	}
	vol, err := hp.state.GetVolumeByID(source.VolID)
	if err != nil {
		return err
	}
	if vol.VolAccessType != state.BlockAccess {
		return status.Error(codes.InvalidArgument, "source volume does not have block mode access type")
	}

	br, err := newFileBlockReader(
		hp.getSnapshotPath(baseSnapID),
		hp.getSnapshotPath(targetSnapID),
		req.StartingOffset,
		state.BlockSizeBytes,
		hp.config.SnapshotMetadataBlockType,
		req.MaxResults,
	)
	if err != nil {
		klog.Errorf("failed initialize file block reader: %v", err)
		return status.Error(codes.Internal, "failed initialize file block reader")
	}
	defer br.Close()
	if err := br.seekToStartingOffset(); err != nil {
		return status.Error(codes.OutOfRange, fmt.Sprintf("failed to seek to starting offset: %v", err.Error()))
	}

	// Read all blocks can find changed blocks till EOF in chunks of size == maxSize
	for {
		cb, cbErr := br.getChangedBlockMetadata(ctx)
		if cbErr != nil {
			if errors.Is(cbErr, context.Canceled) {
				klog.V(4).Info("context canceled while getting changed block metadata, returning")
				return nil
			}
			if errors.Is(cbErr, context.DeadlineExceeded) {
				klog.V(4).Info("context deadline exceeded while getting changed block metadata, returning")
				return nil
			}
			if errors.Is(cbErr, io.EOF) {
				klog.V(4).Info("reached EOF while getting changed block metadata, returning")
				// send changed blocks found till EOF
				if err := sendGetMetadataDeltaResponse(stream, vol.VolSize, hp.config.SnapshotMetadataBlockType, cb); err != nil {
					return err
				}
				return nil
			}
			klog.Errorf("failed to get changed block metadata: %v", cbErr)
			return status.Error(codes.Internal, "failed to get changed block metadata")
		}
		// stream response to client
		if err := sendGetMetadataDeltaResponse(stream, vol.VolSize, hp.config.SnapshotMetadataBlockType, cb); err != nil {
			return err
		}
	}
}

func sendGetMetadataDeltaResponse(
	stream csi.SnapshotMetadata_GetMetadataDeltaServer,
	volSize int64,
	blockMetadataType csi.BlockMetadataType,
	cb []*csi.BlockMetadata,
) error {
	if len(cb) == 0 {
		return nil
	}
	resp := csi.GetMetadataDeltaResponse{
		BlockMetadataType:   blockMetadataType,
		VolumeCapacityBytes: volSize,
		BlockMetadata:       cb,
	}
	return stream.Send(&resp)
}

func sendGetMetadataAllocatedResponse(
	stream csi.SnapshotMetadata_GetMetadataAllocatedServer,
	volSize int64,
	blockMetadataType csi.BlockMetadataType,
	cb []*csi.BlockMetadata,
) error {
	if len(cb) == 0 {
		return nil
	}
	resp := csi.GetMetadataAllocatedResponse{
		BlockMetadataType:   blockMetadataType,
		VolumeCapacityBytes: volSize,
		BlockMetadata:       cb,
	}
	return stream.Send(&resp)
}
