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
	"reflect"
	"testing"

	"github.com/kubernetes-csi/csi-driver-host-path/pkg/state"
)

func TestOptionsFromParameters(t *testing.T) {
	mountedVolume := state.Volume{
		VolAccessType: state.MountAccess,
	}
	blockVolume := state.Volume{
		VolAccessType: state.BlockAccess,
	}

	cases := []struct {
		testName string
		volume   state.Volume
		params   map[string]string
		result   []string
		success  bool
	}{
		{
			testName: "no options passed (mounted)",
			params:   nil,
			result:   nil,
			success:  true,
			volume:   mountedVolume,
		},
		{
			testName: "no options passed (block)",
			params:   nil,
			result:   nil,
			success:  true,
			volume:   blockVolume,
		},
		{
			testName: "ignoreFailedRead=false (mounted)",
			params: map[string]string{
				ignoreFailedReadParameterName: "false",
			},
			result:  nil,
			success: true,
			volume:  mountedVolume,
		},
		{
			testName: "ignoreFailedRead=true (mounted)",
			params: map[string]string{
				ignoreFailedReadParameterName: "true",
			},
			result: []string{
				"--ignore-failed-read",
			},
			success: true,
			volume:  mountedVolume,
		},
		{
			testName: "invalid ignoreFailedRead (mounted)",
			params: map[string]string{
				ignoreFailedReadParameterName: "ABC",
			},
			result:  nil,
			success: false,
			volume:  mountedVolume,
		},
		{
			testName: "ignoreFailedRead=true (block)",
			params: map[string]string{
				ignoreFailedReadParameterName: "true",
			},
			result:  nil,
			success: true,
			volume:  blockVolume,
		},
		{
			testName: "invalid ignoreFailedRead (block)",
			params: map[string]string{
				ignoreFailedReadParameterName: "ABC",
			},
			result:  nil,
			success: true,
			volume:  blockVolume,
		},
	}

	for _, tt := range cases {
		t.Run(tt.testName, func(t *testing.T) {
			result, err := optionsFromParameters(tt.volume, tt.params)
			if tt.success && err != nil {
				t.Fatalf("expected success, found error: %s", err.Error())
			}
			if !tt.success && err == nil {
				t.Fatalf("expected failure, but succeeded with %q", result)
			}

			if tt.success && !reflect.DeepEqual(tt.result, result) {
				t.Fatalf("expected %q but received %q", tt.result, result)
			}
		})
	}
}
