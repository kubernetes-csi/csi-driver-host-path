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
	"strconv"

	"github.com/kubernetes-csi/csi-driver-host-path/pkg/state"
)

// ignoreFailedReadParameterName is a parameter that, when set to true,
// causes the `--ignore-failed-read` option to be passed to `tar`.
const ignoreFailedReadParameterName = "ignoreFailedRead"

func optionsFromParameters(vol state.Volume, parameters map[string]string) ([]string, error) {
	// We do not support options for snapshots of block volumes
	if vol.VolAccessType == state.BlockAccess {
		return nil, nil
	}

	ignoreFailedReadString := parameters[ignoreFailedReadParameterName]
	if len(ignoreFailedReadString) == 0 {
		return nil, nil
	}

	if ok, err := strconv.ParseBool(ignoreFailedReadString); err != nil {
		return nil, fmt.Errorf(
			"invalid value for %q, expected boolean but was %q",
			ignoreFailedReadParameterName,
			ignoreFailedReadString,
		)
	} else if ok {
		return []string{"--ignore-failed-read"}, nil
	}

	return nil, nil
}
