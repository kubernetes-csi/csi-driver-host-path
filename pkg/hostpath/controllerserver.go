/*
Copyright 2017 The Kubernetes Authors.

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
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/golang/protobuf/ptypes"

	"github.com/golang/glog"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/container-storage-interface/spec/lib/go/csi"
	utilexec "k8s.io/utils/exec"
)

const (
	deviceID           = "deviceID"
	maxStorageCapacity = tib
)

type accessType int

const (
	mountAccess accessType = iota
	blockAccess
)

func (hp *hostPath) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (resp *csi.CreateVolumeResponse, finalErr error) {
	if err := hp.validateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); err != nil {
		glog.V(3).Infof("invalid create volume req: %v", req)
		return nil, err
	}

	// Check arguments
	if len(req.GetName()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Name missing in request")
	}
	caps := req.GetVolumeCapabilities()
	if caps == nil {
		return nil, status.Error(codes.InvalidArgument, "Volume Capabilities missing in request")
	}

	// Keep a record of the requested access types.
	var accessTypeMount, accessTypeBlock bool

	for _, cap := range caps {
		if cap.GetBlock() != nil {
			accessTypeBlock = true
		}
		if cap.GetMount() != nil {
			accessTypeMount = true
		}
	}
	// A real driver would also need to check that the other
	// fields in VolumeCapabilities are sane. The check above is
	// just enough to pass the "[Testpattern: Dynamic PV (block
	// volmode)] volumeMode should fail in binding dynamic
	// provisioned PV to PVC" storage E2E test.

	if accessTypeBlock && accessTypeMount {
		return nil, status.Error(codes.InvalidArgument, "cannot have both block and mount access type")
	}

	var requestedAccessType accessType

	if accessTypeBlock {
		requestedAccessType = blockAccess
	} else {
		// Default to mount.
		requestedAccessType = mountAccess
	}

	// Lock before acting on global state. A production-quality
	// driver might use more fine-grained locking.
	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	// Check for maximum available capacity
	capacity := int64(req.GetCapacityRange().GetRequiredBytes())
	if capacity >= maxStorageCapacity {
		return nil, status.Errorf(codes.OutOfRange, "Requested capacity %d exceeds maximum allowed %d", capacity, maxStorageCapacity)
	}

	topologies := []*csi.Topology{&csi.Topology{
		Segments: map[string]string{TopologyKeyNode: hp.nodeID},
	}}

	// Need to check for already existing volume name, and if found
	// check for the requested capacity and already allocated capacity
	if exVol, err := hp.getVolumeByName(req.GetName()); err == nil {
		// Since err is nil, it means the volume with the same name already exists
		// need to check if the size of existing volume is the same as in new
		// request
		if exVol.VolSize < capacity {
			return nil, status.Errorf(codes.AlreadyExists, "Volume with the same name: %s but with different size already exist", req.GetName())
		}
		if req.GetVolumeContentSource() != nil {
			volumeSource := req.VolumeContentSource
			switch volumeSource.Type.(type) {
			case *csi.VolumeContentSource_Snapshot:
				if volumeSource.GetSnapshot() != nil && exVol.ParentSnapID != "" && exVol.ParentSnapID != volumeSource.GetSnapshot().GetSnapshotId() {
					return nil, status.Error(codes.AlreadyExists, "existing volume source snapshot id not matching")
				}
			case *csi.VolumeContentSource_Volume:
				if volumeSource.GetVolume() != nil && exVol.ParentVolID != volumeSource.GetVolume().GetVolumeId() {
					return nil, status.Error(codes.AlreadyExists, "existing volume source volume id not matching")
				}
			default:
				return nil, status.Errorf(codes.InvalidArgument, "%v not a proper volume source", volumeSource)
			}
		}
		// TODO (sbezverk) Do I need to make sure that volume still exists?
		return &csi.CreateVolumeResponse{
			Volume: &csi.Volume{
				VolumeId:           exVol.VolID,
				CapacityBytes:      int64(exVol.VolSize),
				VolumeContext:      req.GetParameters(),
				ContentSource:      req.GetVolumeContentSource(),
				AccessibleTopology: topologies,
			},
		}, nil
	}

	volumeID := uuid.NewUUID().String()

	vol, err := createHostpathVolume(volumeID, req.GetName(), capacity, requestedAccessType, false /* ephemeral */)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create volume %v: %v", volumeID, err)
	}
	glog.V(4).Infof("created volume %s at path %s", vol.VolID, vol.VolPath)

	if req.GetVolumeContentSource() != nil {
		path := getVolumePath(volumeID)
		volumeSource := req.VolumeContentSource
		switch volumeSource.Type.(type) {
		case *csi.VolumeContentSource_Snapshot:
			if snapshot := volumeSource.GetSnapshot(); snapshot != nil {
				err = loadFromSnapshot(capacity, snapshot.GetSnapshotId(), path, requestedAccessType)
				vol.ParentSnapID = snapshot.GetSnapshotId()
			}
		case *csi.VolumeContentSource_Volume:
			if srcVolume := volumeSource.GetVolume(); srcVolume != nil {
				err = loadFromVolume(capacity, srcVolume.GetVolumeId(), path, requestedAccessType)
				vol.ParentVolID = srcVolume.GetVolumeId()
			}
		default:
			err = status.Errorf(codes.InvalidArgument, "%v not a proper volume source", volumeSource)
		}
		if err != nil {
			glog.V(4).Infof("VolumeSource error: %v", err)
			if delErr := hp.deleteVolume(volumeID); delErr != nil {
				glog.V(2).Infof("deleting hostpath volume %v failed: %v", volumeID, delErr)
			}
			return nil, err
		}
		glog.V(4).Infof("successfully populated volume %s", vol.VolID)
	}

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:           volumeID,
			CapacityBytes:      req.GetCapacityRange().GetRequiredBytes(),
			VolumeContext:      req.GetParameters(),
			ContentSource:      req.GetVolumeContentSource(),
			AccessibleTopology: topologies,
		},
	}, nil
}

func (hp *hostPath) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	// Check arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}

	if err := hp.validateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); err != nil {
		glog.V(3).Infof("invalid delete volume req: %v", req)
		return nil, err
	}

	// Lock before acting on global state. A production-quality
	// driver might use more fine-grained locking.
	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	volId := req.GetVolumeId()
	if err := hp.deleteVolume(volId); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete volume %v: %v", volId, err)
	}

	glog.V(4).Infof("volume %v successfully deleted", volId)

	return &csi.DeleteVolumeResponse{}, nil
}

func (hp *hostPath) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: hp.getControllerServiceCapabilities(),
	}, nil
}

func (hp *hostPath) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {

	// Check arguments
	if len(req.GetVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID cannot be empty")
	}
	if len(req.VolumeCapabilities) == 0 {
		return nil, status.Error(codes.InvalidArgument, req.VolumeId)
	}

	// Lock before acting on global state. A production-quality
	// driver might use more fine-grained locking.
	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	if _, err := hp.getVolumeByID(req.GetVolumeId()); err != nil {
		return nil, status.Error(codes.NotFound, req.GetVolumeId())
	}

	for _, cap := range req.GetVolumeCapabilities() {
		if cap.GetMount() == nil && cap.GetBlock() == nil {
			return nil, status.Error(codes.InvalidArgument, "cannot have both mount and block access type be undefined")
		}

		// A real driver would check the capabilities of the given volume with
		// the set of requested capabilities.
	}

	return &csi.ValidateVolumeCapabilitiesResponse{
		Confirmed: &csi.ValidateVolumeCapabilitiesResponse_Confirmed{
			VolumeContext:      req.GetVolumeContext(),
			VolumeCapabilities: req.GetVolumeCapabilities(),
			Parameters:         req.GetParameters(),
		},
	}, nil
}

func (hp *hostPath) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (hp *hostPath) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (hp *hostPath) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

func (hp *hostPath) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	volumeRes := &csi.ListVolumesResponse{
		Entries: []*csi.ListVolumesResponse_Entry{},
	}

	var (
		startIdx, volumesLength, maxLength int64
		hpVolume                           hostPathVolume
	)

	// Lock before acting on global state. A production-quality
	// driver might use more fine-grained locking.
	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	volumeIds := getSortedVolumeIDs()
	if req.StartingToken == "" {
		req.StartingToken = "1"
	}

	startIdx, err := strconv.ParseInt(req.StartingToken, 10, 32)
	if err != nil {
		return nil, status.Error(codes.Aborted, "The type of startingToken should be integer")
	}

	volumesLength = int64(len(volumeIds))
	maxLength = int64(req.MaxEntries)

	if maxLength > volumesLength || maxLength <= 0 {
		maxLength = volumesLength
	}

	for index := startIdx - 1; index < volumesLength && index < maxLength; index++ {
		hpVolume = hp.volumes[volumeIds[index]]
		healthy, msg := doHealthCheckInControllerSide(volumeIds[index])
		glog.V(3).Infof("Healthy state: %s Volume: %t", hpVolume.VolName, healthy)
		volumeRes.Entries = append(volumeRes.Entries, &csi.ListVolumesResponse_Entry{
			Volume: &csi.Volume{
				VolumeId:      hpVolume.VolID,
				CapacityBytes: hpVolume.VolSize,
			},
			Status: &csi.ListVolumesResponse_VolumeStatus{
				PublishedNodeIds: []string{hpVolume.NodeID},
				VolumeCondition: &csi.VolumeCondition{
					Abnormal: !healthy,
					Message:  msg,
				},
			},
		})
	}

	glog.V(5).Infof("Volumes are: %+v", *volumeRes)
	return volumeRes, nil
}

func (hp *hostPath) ControllerGetVolume(ctx context.Context, req *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	// Lock before acting on global state. A production-quality
	// driver might use more fine-grained locking.
	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	volume, ok := hp.volumes[req.GetVolumeId()]
	if !ok {
		return nil, status.Error(codes.NotFound, "The volume not found")
	}

	healthy, msg := doHealthCheckInControllerSide(req.GetVolumeId())
	glog.V(3).Infof("Healthy state: %s Volume: %t", volume.VolName, healthy)
	return &csi.ControllerGetVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      volume.VolID,
			CapacityBytes: volume.VolSize,
		},
		Status: &csi.ControllerGetVolumeResponse_VolumeStatus{
			PublishedNodeIds: []string{volume.NodeID},
			VolumeCondition: &csi.VolumeCondition{
				Abnormal: !healthy,
				Message:  msg,
			},
		},
	}, nil
}

// getSnapshotPath returns the full path to where the snapshot is stored
func getSnapshotPath(snapshotID string) string {
	return filepath.Join(dataRoot, fmt.Sprintf("%s%s", snapshotID, snapshotExt))
}

// CreateSnapshot uses tar command to create snapshot for hostpath volume. The tar command can quickly create
// archives of entire directories. The host image must have "tar" binaries in /bin, /usr/sbin, or /usr/bin.
func (hp *hostPath) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	if err := hp.validateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT); err != nil {
		glog.V(3).Infof("invalid create snapshot req: %v", req)
		return nil, err
	}

	if len(req.GetName()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Name missing in request")
	}
	// Check arguments
	if len(req.GetSourceVolumeId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "SourceVolumeId missing in request")
	}

	// Lock before acting on global state. A production-quality
	// driver might use more fine-grained locking.
	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	// Need to check for already existing snapshot name, and if found check for the
	// requested sourceVolumeId and sourceVolumeId of snapshot that has been created.
	if exSnap, err := hp.getSnapshotByName(req.GetName()); err == nil {
		// Since err is nil, it means the snapshot with the same name already exists need
		// to check if the sourceVolumeId of existing snapshot is the same as in new request.
		if exSnap.VolID == req.GetSourceVolumeId() {
			// same snapshot has been created.
			return &csi.CreateSnapshotResponse{
				Snapshot: &csi.Snapshot{
					SnapshotId:     exSnap.Id,
					SourceVolumeId: exSnap.VolID,
					CreationTime:   exSnap.CreationTime,
					SizeBytes:      exSnap.SizeBytes,
					ReadyToUse:     exSnap.ReadyToUse,
				},
			}, nil
		}
		return nil, status.Errorf(codes.AlreadyExists, "snapshot with the same name: %s but with different SourceVolumeId already exist", req.GetName())
	}

	volumeID := req.GetSourceVolumeId()
	hostPathVolume, ok := hp.volumes[volumeID]
	if !ok {
		return nil, status.Error(codes.Internal, "volumeID is not exist")
	}

	snapshotID := uuid.NewUUID().String()
	creationTime := ptypes.TimestampNow()
	volPath := hostPathVolume.VolPath
	file := getSnapshotPath(snapshotID)

	var cmd []string
	if hostPathVolume.VolAccessType == blockAccess {
		glog.V(4).Infof("Creating snapshot of Raw Block Mode Volume")
		cmd = []string{"cp", volPath, file}
	} else {
		glog.V(4).Infof("Creating snapshot of Filsystem Mode Volume")
		cmd = []string{"tar", "czf", file, "-C", volPath, "."}
	}
	executor := utilexec.New()
	out, err := executor.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed create snapshot: %v: %s", err, out)
	}

	glog.V(4).Infof("create volume snapshot %s", file)
	snapshot := hostPathSnapshot{}
	snapshot.Name = req.GetName()
	snapshot.Id = snapshotID
	snapshot.VolID = volumeID
	snapshot.Path = file
	snapshot.CreationTime = creationTime
	snapshot.SizeBytes = hostPathVolume.VolSize
	snapshot.ReadyToUse = true

	hostPathVolumeSnapshots[snapshotID] = snapshot

	return &csi.CreateSnapshotResponse{
		Snapshot: &csi.Snapshot{
			SnapshotId:     snapshot.Id,
			SourceVolumeId: snapshot.VolID,
			CreationTime:   snapshot.CreationTime,
			SizeBytes:      snapshot.SizeBytes,
			ReadyToUse:     snapshot.ReadyToUse,
		},
	}, nil
}

func (hp *hostPath) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	// Check arguments
	if len(req.GetSnapshotId()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Snapshot ID missing in request")
	}

	if err := hp.validateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT); err != nil {
		glog.V(3).Infof("invalid delete snapshot req: %v", req)
		return nil, err
	}
	snapshotID := req.GetSnapshotId()

	// Lock before acting on global state. A production-quality
	// driver might use more fine-grained locking.
	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	glog.V(4).Infof("deleting snapshot %s", snapshotID)
	path := getSnapshotPath(snapshotID)
	os.RemoveAll(path)
	delete(hostPathVolumeSnapshots, snapshotID)
	return &csi.DeleteSnapshotResponse{}, nil
}

func (hp *hostPath) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	if err := hp.validateControllerServiceRequest(csi.ControllerServiceCapability_RPC_LIST_SNAPSHOTS); err != nil {
		glog.V(3).Infof("invalid list snapshot req: %v", req)
		return nil, err
	}

	// Lock before acting on global state. A production-quality
	// driver might use more fine-grained locking.
	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	// case 1: SnapshotId is not empty, return snapshots that match the snapshot id.
	if len(req.GetSnapshotId()) != 0 {
		snapshotID := req.SnapshotId
		if snapshot, ok := hostPathVolumeSnapshots[snapshotID]; ok {
			return convertSnapshot(snapshot), nil
		}
	}

	// case 2: SourceVolumeId is not empty, return snapshots that match the source volume id.
	if len(req.GetSourceVolumeId()) != 0 {
		for _, snapshot := range hostPathVolumeSnapshots {
			if snapshot.VolID == req.SourceVolumeId {
				return convertSnapshot(snapshot), nil
			}
		}
	}

	var snapshots []csi.Snapshot
	// case 3: no parameter is set, so we return all the snapshots.
	sortedKeys := make([]string, 0)
	for k := range hostPathVolumeSnapshots {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		snap := hostPathVolumeSnapshots[key]
		snapshot := csi.Snapshot{
			SnapshotId:     snap.Id,
			SourceVolumeId: snap.VolID,
			CreationTime:   snap.CreationTime,
			SizeBytes:      snap.SizeBytes,
			ReadyToUse:     snap.ReadyToUse,
		}
		snapshots = append(snapshots, snapshot)
	}

	var (
		ulenSnapshots = int32(len(snapshots))
		maxEntries    = req.MaxEntries
		startingToken int32
	)

	if v := req.StartingToken; v != "" {
		i, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return nil, status.Errorf(
				codes.Aborted,
				"startingToken=%d !< int32=%d",
				startingToken, math.MaxUint32)
		}
		startingToken = int32(i)
	}

	if startingToken > ulenSnapshots {
		return nil, status.Errorf(
			codes.Aborted,
			"startingToken=%d > len(snapshots)=%d",
			startingToken, ulenSnapshots)
	}

	// Discern the number of remaining entries.
	rem := ulenSnapshots - startingToken

	// If maxEntries is 0 or greater than the number of remaining entries then
	// set maxEntries to the number of remaining entries.
	if maxEntries == 0 || maxEntries > rem {
		maxEntries = rem
	}

	var (
		i       int
		j       = startingToken
		entries = make(
			[]*csi.ListSnapshotsResponse_Entry,
			maxEntries)
	)

	for i = 0; i < len(entries); i++ {
		entries[i] = &csi.ListSnapshotsResponse_Entry{
			Snapshot: &snapshots[j],
		}
		j++
	}

	var nextToken string
	if j < ulenSnapshots {
		nextToken = fmt.Sprintf("%d", j)
	}

	return &csi.ListSnapshotsResponse{
		Entries:   entries,
		NextToken: nextToken,
	}, nil
}

func (hp *hostPath) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {

	volID := req.GetVolumeId()
	if len(volID) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Volume ID missing in request")
	}

	capRange := req.GetCapacityRange()
	if capRange == nil {
		return nil, status.Error(codes.InvalidArgument, "Capacity range not provided")
	}

	capacity := int64(capRange.GetRequiredBytes())
	if capacity >= maxStorageCapacity {
		return nil, status.Errorf(codes.OutOfRange, "Requested capacity %d exceeds maximum allowed %d", capacity, maxStorageCapacity)
	}

	// Lock before acting on global state. A production-quality
	// driver might use more fine-grained locking.
	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	exVol, err := hp.getVolumeByID(volID)
	if err != nil {
		// Assume not found error
		return nil, status.Errorf(codes.NotFound, "Could not get volume %s: %v", volID, err)
	}

	if exVol.VolSize < capacity {
		exVol.VolSize = capacity
		if err := hp.updateVolume(volID, exVol); err != nil {
			return nil, status.Errorf(codes.Internal, "Could not update volume %s: %v", volID, err)
		}
	}

	return &csi.ControllerExpandVolumeResponse{
		CapacityBytes:         exVol.VolSize,
		NodeExpansionRequired: true,
	}, nil
}

func convertSnapshot(snap hostPathSnapshot) *csi.ListSnapshotsResponse {
	entries := []*csi.ListSnapshotsResponse_Entry{
		{
			Snapshot: &csi.Snapshot{
				SnapshotId:     snap.Id,
				SourceVolumeId: snap.VolID,
				CreationTime:   snap.CreationTime,
				SizeBytes:      snap.SizeBytes,
				ReadyToUse:     snap.ReadyToUse,
			},
		},
	}

	rsp := &csi.ListSnapshotsResponse{
		Entries: entries,
	}

	return rsp
}

func (hp *hostPath) validateControllerServiceRequest(c csi.ControllerServiceCapability_RPC_Type) error {
	if c == csi.ControllerServiceCapability_RPC_UNKNOWN {
		return nil
	}

	for _, cap := range hp.getControllerServiceCapabilities() {
		if c == cap.GetRpc().GetType() {
			return nil
		}
	}
	return status.Errorf(codes.InvalidArgument, "unsupported capability %s", c)
}

func (hp *hostPath) getControllerServiceCapabilities() []*csi.ControllerServiceCapability {
	var cl []csi.ControllerServiceCapability_RPC_Type
	if !hp.ephemeral {
		cl = []csi.ControllerServiceCapability_RPC_Type{
			csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
			csi.ControllerServiceCapability_RPC_GET_VOLUME,
			csi.ControllerServiceCapability_RPC_CREATE_DELETE_SNAPSHOT,
			csi.ControllerServiceCapability_RPC_LIST_SNAPSHOTS,
			csi.ControllerServiceCapability_RPC_LIST_VOLUMES,
			csi.ControllerServiceCapability_RPC_CLONE_VOLUME,
			csi.ControllerServiceCapability_RPC_EXPAND_VOLUME,
			csi.ControllerServiceCapability_RPC_VOLUME_CONDITION,
		}
	}

	var csc []*csi.ControllerServiceCapability

	for _, cap := range cl {
		csc = append(csc, &csi.ControllerServiceCapability{
			Type: &csi.ControllerServiceCapability_Rpc{
				Rpc: &csi.ControllerServiceCapability_RPC{
					Type: cap,
				},
			},
		})
	}

	return csc
}
