package hostpath

import (
	"context"
	"os"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-driver-host-path/pkg/state"
	"github.com/stretchr/testify/assert"
)

func TestCreateVolume(t *testing.T) {
	// TODO: add more test cases
	testCases := []struct {
		name                          string
		enableControllerModifyVolume  bool
		acceptedMutableParameterNames StringArray

		req       *csi.CreateVolumeRequest
		resp      *csi.CreateVolumeResponse
		expectErr bool
	}{
		{
			name: "failure case - no volume name",
			req: &csi.CreateVolumeRequest{
				Name:               "",
				VolumeCapabilities: []*csi.VolumeCapability{},
			},
			expectErr: true,
		},
		{
			name: "failure case - no volume capabilities",
			req: &csi.CreateVolumeRequest{
				Name:               "fakeVolume",
				VolumeCapabilities: nil,
			},
			expectErr: true,
		},
		{
			name: "failure case - have both block and mount access type",
			req: &csi.CreateVolumeRequest{
				Name: "fakeVolume",
				VolumeCapabilities: []*csi.VolumeCapability{
					{
						AccessType: &csi.VolumeCapability_Mount{
							Mount: &csi.VolumeCapability_MountVolume{},
						},
						AccessMode: &csi.VolumeCapability_AccessMode{
							Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_SINGLE_WRITER,
						},
					},
					{
						AccessType: &csi.VolumeCapability_Block{
							Block: &csi.VolumeCapability_BlockVolume{},
						},
						AccessMode: &csi.VolumeCapability_AccessMode{
							Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_SINGLE_WRITER,
						},
					},
				},
			},
			expectErr: true,
		},
		{
			name:                         "failure case - no cap",
			enableControllerModifyVolume: false,
			req: &csi.CreateVolumeRequest{
				Name: "fakeVolume",
				MutableParameters: map[string]string{
					"fakeKey": "fakeValue",
				},
			},
			expectErr: true,
		},
		{
			name:                          "failure case - contains unaccepted parameters",
			enableControllerModifyVolume:  true,
			acceptedMutableParameterNames: StringArray{"barKey"},
			req: &csi.CreateVolumeRequest{
				Name: "fakeVolume",
				MutableParameters: map[string]string{
					"fakeKey": "fakeValue",
				},
				VolumeCapabilities: []*csi.VolumeCapability{},
			},
			expectErr: true,
		},
		{
			name:                         "success case - valid request 1",
			enableControllerModifyVolume: true,
			req: &csi.CreateVolumeRequest{
				Name: "fakeVolume",
				MutableParameters: map[string]string{
					"fakeKey": "fakeValue",
				},
				VolumeCapabilities: []*csi.VolumeCapability{},
			},
			resp: &csi.CreateVolumeResponse{
				Volume: &csi.Volume{
					AccessibleTopology: []*csi.Topology{
						{
							Segments: map[string]string{
								TopologyKeyNode: "fakeNodeID",
							},
						},
					},
				},
			},
			expectErr: false,
		},
		{
			name:                          "success case - valid request 2",
			enableControllerModifyVolume:  true,
			acceptedMutableParameterNames: StringArray{"fakeKey"},
			req: &csi.CreateVolumeRequest{
				Name: "fakeVolume",
				MutableParameters: map[string]string{
					"fakeKey": "fakeValue",
				},
				VolumeCapabilities: []*csi.VolumeCapability{},
			},
			resp: &csi.CreateVolumeResponse{
				Volume: &csi.Volume{
					AccessibleTopology: []*csi.Topology{
						{
							Segments: map[string]string{
								TopologyKeyNode: "fakeNodeID",
							},
						},
					},
				},
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stateDir, err := os.MkdirTemp(os.TempDir(), "csi-data-dir")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(stateDir)

			cfg := Config{
				StateDir:                      stateDir,
				Endpoint:                      "unix://tmp/csi.sock",
				DriverName:                    "hostpath.csi.k8s.io",
				NodeID:                        "fakeNodeID",
				MaxVolumeSize:                 1024 * 1024 * 1024 * 1024,
				EnableTopology:                true,
				EnableControllerModifyVolume:  tc.enableControllerModifyVolume,
				AcceptedMutableParameterNames: tc.acceptedMutableParameterNames,
			}
			hp, err := NewHostPathDriver(cfg)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := hp.CreateVolume(context.TODO(), tc.req)
			if tc.expectErr && err == nil {
				t.Fatalf("expected error, got none")
			}
			if !tc.expectErr && err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			// Clear out the volume ID so we can compare the rest of the response
			if resp != nil {
				resp.Volume.VolumeId = ""
			}
			assert.Equal(t, tc.resp, resp)
		})
	}
}

func TestControllerModifyVolume(t *testing.T) {
	volume1 := state.Volume{VolID: "fakeVolumeID", VolName: "fakeVolume", VolSize: int64(1), VolAccessType: state.MountAccess}

	testCases := []struct {
		name                          string
		volumes                       []state.Volume
		enableControllerModifyVolume  bool
		acceptedMutableParameterNames StringArray

		req       *csi.ControllerModifyVolumeRequest
		resp      *csi.ControllerModifyVolumeResponse
		expectErr bool
	}{
		{
			name: "failure case - no volume id",
			req: &csi.ControllerModifyVolumeRequest{
				VolumeId:          "",
				MutableParameters: map[string]string{},
			},
			expectErr: true,
		},
		{
			name: "failure case - no mutable parameters",
			req: &csi.ControllerModifyVolumeRequest{
				VolumeId:          "fakeVolumeID",
				MutableParameters: map[string]string{},
			},
			expectErr: true,
		},
		{
			name: "failure case - no volume found",
			req: &csi.ControllerModifyVolumeRequest{
				VolumeId: "fakeVolumeID",
				MutableParameters: map[string]string{
					"fakeKey": "fakeValue",
				},
			},
			expectErr: true,
		},
		{
			name:                         "failure case - no cap",
			volumes:                      []state.Volume{volume1},
			enableControllerModifyVolume: false,
			req: &csi.ControllerModifyVolumeRequest{
				VolumeId: "fakeVolumeID",
				MutableParameters: map[string]string{
					"fakeKey": "fakeValue",
				},
			},
			expectErr: true,
		},
		{
			name:                          "failure case - contains unaccepted parameters",
			volumes:                       []state.Volume{volume1},
			enableControllerModifyVolume:  true,
			acceptedMutableParameterNames: StringArray{"barKey"},
			req: &csi.ControllerModifyVolumeRequest{
				VolumeId: "fakeVolumeID",
				MutableParameters: map[string]string{
					"fakeKey": "fakeValue",
				},
			},
			expectErr: true,
		},
		{
			name:                         "success case - valid request 1",
			volumes:                      []state.Volume{volume1},
			enableControllerModifyVolume: true,
			req: &csi.ControllerModifyVolumeRequest{
				VolumeId: "fakeVolumeID",
				MutableParameters: map[string]string{
					"fakeKey": "fakeValue",
				},
			},
			resp:      &csi.ControllerModifyVolumeResponse{},
			expectErr: false,
		},
		{
			name:                          "success case - valid request 2",
			volumes:                       []state.Volume{volume1},
			enableControllerModifyVolume:  true,
			acceptedMutableParameterNames: StringArray{"fakeKey"},
			req: &csi.ControllerModifyVolumeRequest{
				VolumeId: "fakeVolumeID",
				MutableParameters: map[string]string{
					"fakeKey": "fakeValue",
				},
			},
			resp:      &csi.ControllerModifyVolumeResponse{},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stateDir, err := os.MkdirTemp(os.TempDir(), "csi-data-dir")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(stateDir)

			cfg := Config{
				StateDir:                      stateDir,
				Endpoint:                      "unix://tmp/csi.sock",
				DriverName:                    "hostpath.csi.k8s.io",
				NodeID:                        "fakeNodeID",
				MaxVolumeSize:                 1024 * 1024 * 1024 * 1024,
				EnableTopology:                true,
				EnableControllerModifyVolume:  tc.enableControllerModifyVolume,
				AcceptedMutableParameterNames: tc.acceptedMutableParameterNames,
			}
			hp, err := NewHostPathDriver(cfg)
			if err != nil {
				t.Fatal(err)
			}
			for _, volume := range tc.volumes {
				_, err := hp.createVolume(volume.VolID, volume.VolName, volume.VolSize, volume.VolAccessType, volume.Ephemeral, volume.Kind)
				if err != nil {
					t.Fatal(err)
				}
			}

			resp, err := hp.ControllerModifyVolume(context.TODO(), tc.req)
			if tc.expectErr && err == nil {
				t.Fatalf("expected error, got none")
			}
			if !tc.expectErr && err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			assert.Equal(t, tc.resp, resp)
		})
	}
}
