#! /bin/bash

# shellcheck disable=SC1091
. release-tools/prow.sh

if find . -name Dockerfile | grep -v ^./vendor | xargs --no-run-if-empty cat | grep -q ^RUN; then
    # Needed for "RUN apk" on non linux/amd64 platforms.
    (set -x; docker run --rm --privileged multiarch/qemu-user-static --reset -p yes)
fi

gcr_cloud_build
