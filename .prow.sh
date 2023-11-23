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

# FIXME: do not merge this, only for testing PR #399
CSI_PROW_SANITY_VERSION="master"

. release-tools/prow.sh

main
