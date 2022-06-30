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

# List of images under https://console.cloud.google.com/gcr/images/k8s-staging-sig-storage/GLOBAL
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

for image in $images; do
    latest=$(gcloud container images list-tags k8s.gcr.io/sig-storage/$image --format='get(tags)' --filter='tags~^v AND NOT tags~v2020 AND NOT tags~-rc' --sort-by=tags | tail -n 1)

    sed -i -e "s;\(image: registry.k8s.io/sig-storage/$image:\).*;\1$latest;" $(find deploy -type f)
done
