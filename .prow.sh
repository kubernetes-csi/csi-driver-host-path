#! /bin/bash

# Copyright 2021 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# TODO: move into job
CSI_PROW_GINKO_PARALLEL="-p -nodes 40" # default was 7

# Simulates canary test job.
# TODO: canary periodic job
#CSI_PROW_BUILD_JOB=false
#CSI_PROW_KUBERNETES_VERSION=latest
#CSI_PROW_HOSTPATH_CANARY=canary
CSI_PROW_HOSTPATH_DRIVER_NAME="hostpath.csi.k8s.io"

CSI_PROW_TESTS_SANITY="sanity"

# We need to hardcode e2e version for resizer for now, because
# we need fixes from latest release-1.31 branch for all e2es to pass.
#
# See: https://github.com/kubernetes-csi/csi-driver-host-path/pull/581#issuecomment-2634529098
# See: https://github.com/kubernetes-csi/external-resizer/blob/20072c0fdf8baaf919ef95d6e918538ba9d84eaf/.prow.sh
export CSI_PROW_E2E_VERSION="release-1.31"

. release-tools/prow.sh

main
