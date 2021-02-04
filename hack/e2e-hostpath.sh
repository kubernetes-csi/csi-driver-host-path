#!/bin/bash

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

## This file is for app/hostpathplugin
## It could be used for other apps in this repo, but
## those applications may or may not take the same
## arguments

## Must be run from the root of the repo

UDS="/tmp/e2e-csi-sanity.sock"
CSI_ENDPOINT="unix://${UDS}"
CSI_MOUNTPOINT="/mnt"
APP=hostpathplugin

SKIP="WithCapacity"

# Get csi-sanity
./hack/get-sanity.sh

# Build
make hostpath

# Cleanup
rm -f $UDS

# Start the application in the background
sudo _output/$APP --endpoint=$CSI_ENDPOINT --nodeid=1 &
pid=$!

# Need to skip Capacity testing since hostpath does not support it
sudo $GOPATH/bin/csi-sanity $@ \
    --ginkgo.skip=${SKIP} \
    --csi.mountdir=$CSI_MOUNTPOINT \
    --csi.endpoint=$CSI_ENDPOINT ; ret=$?
sudo kill -9 $pid
sudo rm -f $UDS

if [ $ret -ne 0 ] ; then
	exit $ret
fi

exit 0
