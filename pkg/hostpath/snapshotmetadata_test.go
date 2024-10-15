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
	"math"
	"os"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-driver-host-path/pkg/state"
)

func TestGetChangedBlockMetadata(t *testing.T) {
	testCases := []struct {
		name             string
		sourceFileBlocks int
		targetFileBlocks int
		changedBlocks    []int
		startingOffset   int64
		maxResult        int32
		expectedResponse []*csi.BlockMetadata
		expectErr        bool
	}{
		{
			name:             "success case",
			sourceFileBlocks: 100,
			targetFileBlocks: 100,
			changedBlocks:    []int{2, 4, 7, 30, 70},
			maxResult:        100,
			expectedResponse: []*csi.BlockMetadata{
				{
					ByteOffset: 2 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 4 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 7 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 30 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 70 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
			},
			expectErr: false,
		},
		{
			name:             "success case with max result",
			sourceFileBlocks: 100,
			targetFileBlocks: 100,
			changedBlocks:    []int{2, 4, 7, 10, 30, 65, 70},
			maxResult:        3,
			expectedResponse: []*csi.BlockMetadata{
				{
					ByteOffset: 2 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 4 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 7 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 10 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 30 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 65 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 70 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
			},
			expectErr: false,
		},
		{
			name:             "success case with starting offset",
			sourceFileBlocks: 100,
			targetFileBlocks: 100,
			changedBlocks:    []int{2, 4, 7, 10, 30, 70, 65},
			startingOffset:   9 * state.BlockSizeBytes,
			maxResult:        3,
			expectedResponse: []*csi.BlockMetadata{
				{
					ByteOffset: 10 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 30 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 65 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 70 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
			},
			expectErr: false,
		},
		{
			name:             "sucess case empty response",
			sourceFileBlocks: 100,
			targetFileBlocks: 100,
			startingOffset:   9 * state.BlockSizeBytes,
			maxResult:        3,
			expectedResponse: []*csi.BlockMetadata{},
			expectErr:        false,
		},
		{
			name:             "sucess case different sizes",
			sourceFileBlocks: 95,
			targetFileBlocks: 100,
			changedBlocks:    []int{70, 97},
			startingOffset:   9 * state.BlockSizeBytes,
			maxResult:        3,
			expectedResponse: []*csi.BlockMetadata{
				{
					ByteOffset: 70 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 95 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 96 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 97 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 98 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 99 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
			},
			expectErr: false,
		},
		{
			name:             "sucess case different sizes",
			sourceFileBlocks: 100,
			targetFileBlocks: 95,
			changedBlocks:    []int{70, 97},
			startingOffset:   9 * state.BlockSizeBytes,
			maxResult:        3,
			expectedResponse: []*csi.BlockMetadata{
				{
					ByteOffset: 70 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 95 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 96 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 97 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 98 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
				},
				{
					ByteOffset: 99 * state.BlockSizeBytes,
					SizeBytes:  state.BlockSizeBytes,
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

			// Create test files with data
			sourceFile := createTempFile(t, tc.sourceFileBlocks)
			defer sourceFile.Close()
			targetFile := createTempFile(t, tc.targetFileBlocks)
			defer targetFile.Close()
			for _, i := range tc.changedBlocks {
				modifyBlock(t, targetFile, i, []byte("changed block"))
			}

			cfg := Config{
				StateDir:                     stateDir,
				Endpoint:                     "unix://tmp/csi.sock",
				DriverName:                   "hostpath.csi.k8s.io",
				NodeID:                       "fakeNodeID",
				MaxVolumeSize:                1024 * 1024 * 1024 * 1024,
				EnableTopology:               true,
				EnableControllerModifyVolume: true,
			}

			hp, err := NewHostPathDriver(cfg)
			if err != nil {
				t.Fatal(err)
			}
			cb := make(chan []*csi.BlockMetadata, 100)
			err1 := hp.getChangedBlockMetadata(context.Background(), sourceFile.Name(), targetFile.Name(), tc.startingOffset, state.BlockSizeBytes, tc.maxResult, cb)
			if tc.expectErr {
				if err1 == nil {
					t.Fatalf("expected error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			response := []*csi.BlockMetadata{}
			responsePages := 0
			for c := range cb {
				responsePages++
				response = append(response, c...)
			}
			// Validate max result limit
			expPages := int(math.Ceil(float64(len(tc.expectedResponse)) / float64(tc.maxResult)))
			if responsePages != expPages {
				t.Fatalf("expected %d pages of response, got: %d", expPages, responsePages)
			}
			// Validate response content
			if len(tc.expectedResponse) != len(response) {
				t.Fatalf("expected %d changed blocks metadata, got: %d", len(tc.changedBlocks), len(response))
			}
			for i := 0; i < len(response); i++ {
				if response[i].String() != tc.expectedResponse[i].String() {
					t.Fatalf("received unexpected block metadata, expected: %s\n, got %s", tc.expectedResponse[i].String(), response[i].String())
				}
			}

		})
	}
}

// createTempFile creates a file with given number of blocks
func createTempFile(t *testing.T, blocks int) *os.File {
	f, err := os.CreateTemp("", "test-*.img")
	if err != nil {
		t.Fatal(err)
	}
	// Create n blocks of default block size; declared by BlockSizeBytes
	for i := 0; i < blocks; i++ {
		data := make([]byte, state.BlockSizeBytes)
		// Set different content in each block
		for j := 0; j < state.BlockSizeBytes; j++ {
			data[j] = byte(i + 1)
		}
		_, err = f.Write(data)
		if err != nil {
			t.Fatal(err)
		}
	}
	return f
}

// modifyBlock modifies the content of a specific block in the file
func modifyBlock(t *testing.T, file *os.File, blockNumber int, newContent []byte) {
	offset := int64(blockNumber) * state.BlockSizeBytes
	// Seek to the start of the block
	_, err := file.Seek(offset, 0)
	if err != nil {
		t.Fatal(err)
	}

	// Create a buffer with the same size and copy new content into it
	data := make([]byte, state.BlockSizeBytes)
	copy(data, newContent)

	_, err = file.Write(data)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAllocatedBlockMetadata(t *testing.T) {
	testCases := []struct {
		name           string
		fileBlocks     int
		startingOffset int64
		maxResult      int32
		expectedBlocks []int
		expectErr      bool
	}{
		{
			name:           "success case",
			fileBlocks:     5,
			maxResult:      100,
			expectedBlocks: []int{0, 1, 2, 3, 4},
			expectErr:      false,
		},
		{
			name:           "success case with max result",
			fileBlocks:     10,
			expectedBlocks: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			maxResult:      3,
			expectErr:      false,
		},
		{
			name:           "success case with starting offset",
			fileBlocks:     10,
			startingOffset: 4 * state.BlockSizeBytes,
			expectedBlocks: []int{4, 5, 6, 7, 8, 9},
			maxResult:      3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stateDir, err := os.MkdirTemp(os.TempDir(), "csi-data-dir")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(stateDir)

			file := createTempFile(t, tc.fileBlocks)
			defer file.Close()

			cfg := Config{
				StateDir:                     stateDir,
				Endpoint:                     "unix://tmp/csi.sock",
				DriverName:                   "hostpath.csi.k8s.io",
				NodeID:                       "fakeNodeID",
				MaxVolumeSize:                1024 * 1024 * 1024 * 1024,
				EnableTopology:               true,
				EnableControllerModifyVolume: true,
			}

			hp, err := NewHostPathDriver(cfg)
			if err != nil {
				t.Fatal(err)
			}
			cb := make(chan []*csi.BlockMetadata, 100)
			err1 := hp.getAllocatedBlockMetadata(context.Background(), file.Name(), tc.startingOffset, state.BlockSizeBytes, tc.maxResult, cb)
			if tc.expectErr {
				if err1 == nil {
					t.Fatalf("expected error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got: %v", err)
			}
			response := []*csi.BlockMetadata{}
			responsePages := 0
			for c := range cb {
				responsePages++
				response = append(response, c...)
			}
			// Validate max result limit
			expPages := int(math.Ceil(float64(len(tc.expectedBlocks)) / float64(tc.maxResult)))
			if responsePages != expPages {
				t.Fatalf("expected %d pages of response, got: %d", expPages, responsePages)
			}
			// Validate response content
			if len(tc.expectedBlocks) != len(response) {
				t.Fatalf("expected %d changed blocks metadata, got: %d", tc.fileBlocks, len(response))
			}
			for i := 0; i < len(tc.expectedBlocks); i++ {
				expCB := createBlockMetadata(int64(tc.expectedBlocks[i]), state.BlockSizeBytes)
				if response[i].String() != expCB.String() {
					t.Fatalf("received unexpected block metadata, expected: %s\n, got %s", expCB.String(), response[i].String())
				}
			}
		})
	}
}