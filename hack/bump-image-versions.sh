#!/bin/sh

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

# This script will determine the latest release of all sidecars and
# update to that version in *all* deployment files. Do not commit
# everything! Sometimes sidecars must be locked on an older release
# for older Kubernetes releases.

set -e
set -x

# List of images under registry.k8s.io/sig-storage
images="
csi-attacher
csi-node-driver-registrar
csi-provisioner
csi-resizer
csi-snapshotter
csi-external-health-monitor-agent
csi-external-health-monitor-controller
livenessprobe
"

if ! command -v skopeo > /dev/null 2>&1; then
    echo "skopeo not found. For installation instructions, see https://github.com/containers/skopeo/blob/main/install.md"
    exit 1
fi

if ! command -v jq > /dev/null 2>&1; then
    echo "jq not found."
    exit 1
fi

for image in $images; do
    latest=$(skopeo list-tags --retry-times 3 docker://registry.k8s.io/sig-storage/"$image" | jq -r '.Tags[]' | grep -e '^v.*$' | sort -V | tail -n 1)
    find deploy -type f -exec sed -i '' "s;\(image: registry.k8s.io/sig-storage/$image:\).*;\1$latest;" {} \;
done
