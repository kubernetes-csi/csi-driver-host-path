#! /bin/bash

# TODO: move into job
CSI_PROW_GINKO_PARALLEL="-p -nodes 40" # default was 7

# Simulates canary test job.
# TODO: canary periodic job
#CSI_PROW_BUILD_JOB=false
#CSI_PROW_KUBERNETES_VERSION=latest
#CSI_PROW_HOSTPATH_CANARY=canary
CSI_PROW_HOSTPATH_DRIVER_NAME="hostpath.csi.k8s.io"

CSI_PROW_TESTS_SANITY="sanity"
CSI_PROW_KIND_VERSION="86bc23d84ac12dcb56a0528890736e2c347c2dc3"
CSI_PROW_KUBERNETES_VERSION="release-1.15"
CSI_PROW_DEPLOYMENT="kubernetes-1.15"
CSI_PROW_TESTS="parallel"

. release-tools/prow.sh

main
