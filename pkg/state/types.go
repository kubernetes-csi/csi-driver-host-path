package state

import (
	"path"
	"sync"

	"github.com/golang/protobuf/ptypes/timestamp"
)

type AccessType int

const (
	MountAccess AccessType = iota
	BlockAccess
)

const (
	kib    int64 = 1024
	mib    int64 = kib * 1024
	gib    int64 = mib * 1024
	gib100 int64 = gib * 100
	tib    int64 = gib * 1024
	tib100 int64 = tib * 100

	// StorageKind is the special parameter which requests
	// storage of a certain kind (only affects capacity checks).
	StorageKind = "kind"
)

var (
	// Extension with which snapshot files will be saved.
	SnapshotExt = ".snap"
)

const (
	deviceID           = "deviceID"
	maxStorageCapacity = tib
)

type HostPathVolume struct {
	VolName        string     `json:"volName"`
	VolID          string     `json:"volID"`
	VolSize        int64      `json:"volSize"`
	VolPath        string     `json:"volPath"`
	VolAccessType  AccessType `json:"volAccessType"`
	ParentVolID    string     `json:"parentVolID,omitempty"`
	ParentSnapID   string     `json:"parentSnapID,omitempty"`
	Ephemeral      bool       `json:"ephemeral"`
	NodeID         string     `json:"nodeID"`
	Kind           string     `json:"kind"`
	ReadOnlyAttach bool       `json:"readOnlyAttach"`
	IsAttached     bool       `json:"isAttached"`
	IsStaged       bool       `json:"isStaged"`
	IsPublished    bool       `json:"isPublished"`
}

type HostPathVolumes map[string]HostPathVolume

type HostPathSnapshot struct {
	Name         string
	Id           string
	VolID        string
	Path         string
	CreationTime *timestamp.Timestamp
	SizeBytes    int64
	ReadyToUse   bool
}

type HostPathSnapshots map[string]HostPathSnapshot

type HostPathDriverState struct {
	HostPathVolumes   HostPathVolumes
	HostPathSnapshots HostPathSnapshots

	VolumesFileRWLock   *sync.RWMutex
	SnapshotsFileRWLock *sync.RWMutex

	DataRoot           string
	VolumeDataFilePath string
}

func NewHostPathDriverState(dataRoot string) *HostPathDriverState {
	return &HostPathDriverState{
		HostPathVolumes:   make(HostPathVolumes),
		HostPathSnapshots: make(HostPathSnapshots),

		VolumesFileRWLock:   &sync.RWMutex{},
		SnapshotsFileRWLock: &sync.RWMutex{},
		DataRoot:            dataRoot,
		VolumeDataFilePath:  path.Join(dataRoot, "volumes.json"),
	}
}

type DriverState interface {
	VolumeState
	SnapshotState
}

type VolumeState interface {
	ListVolumes() (HostPathVolumes, error)
	GetVolumeByID(volID string) (HostPathVolume, error)
	UpdateVolume(volID string, volume HostPathVolume) error
	GetVolumeByName(volName string) (HostPathVolume, error)
	DeleteVolume(volID string, capacity Capacity) error
	CreateVolume(volID, name string, cap int64, volAccessType AccessType, ephemeral bool, kind string, maxVolumeSize int64, capacity Capacity) (hpv *HostPathVolume, finalErr error)
	GetVolumeLocker() *sync.RWMutex
	GetSortedVolumeIDs() []string
	SumVolumeSizes(kind string) (sum int64)
}

type SnapshotState interface {
	ListSnapshots() (HostPathSnapshots, error)
	GetSnapshotByName(name string) (HostPathSnapshot, error)
	GetSnapshotByID(id string) (HostPathSnapshot, error)
	CreateSnapshot(name, snapshotId, volumeId, snapshotFilePath string, creationTime *timestamp.Timestamp, size int64, readyToUse bool) (HostPathSnapshot, error)
	DeleteSnapshot(snapshotId string) error
	GetSnapshotLocker() *sync.RWMutex
}
