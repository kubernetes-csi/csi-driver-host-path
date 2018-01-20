#!/bin/bash

## This file is for app/hostpathplugin
## It could be used for other apps in this repo, but
## those applications may or may not take the same
## arguments

## Must be run from the root of the repo

UDS="/tmp/e2e-csi-sanity.sock"
CSI_ENDPOINT="unix://${UDS}"
APP=hostpathplugin

# Get csi-sanity
./hack/get-sanity.sh

# Build
cd app/hostpathplugin
  go install || exit 1
cd ../..

# Cleanup
rm -f $UDS

# Start the application in the background
sudo $GOPATH/bin/$APP --endpoint=$CSI_ENDPOINT --nodeid=1 &
pid=$!

# Need to skip Capacity testing since hostpath does not support it
sudo $GOPATH/bin/csi-sanity $@ --ginkgo.skip="WithCapacity" --csi.endpoint=$CSI_ENDPOINT ; ret=$?
sudo kill -9 $pid
sudo rm -f $UDS

if [ $ret -ne 0 ] ; then
	exit $ret
fi

exit 0
