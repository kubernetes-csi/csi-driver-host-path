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

VERSION="v4.0.2"
SANITYTGZ="csi-sanity-${VERSION}.linux.amd64.tar.gz"

echo "Downloading csi-test from https://github.com/kubernetes-csi/csi-test/releases/download/${VERSION}/${SANITYTGZ}"
curl -s -L "https://github.com/kubernetes-csi/csi-test/releases/download/${VERSION}/${SANITYTGZ}" -o ${SANITYTGZ}
tar xzvf ${SANITYTGZ} -C /tmp && \
rm -f ${SANITYTGZ} && \
rm -f $GOPATH/bin/csi-sanity
cp /tmp/csi-sanity/csi-sanity $GOPATH/bin/csi-sanity && \
rm -rf /tmp/csi-sanity
