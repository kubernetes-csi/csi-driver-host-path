#!/usr/bin/env bash

# This script captures the steps required to successfully
# deploy the hostpath plugin driver.  This should be considered
# authoritative and all updates for this process should be
# done here and referenced elsewhere.

# The script assumes that kubectl is available on the OS path 
# where it is executed.

set -e
set -o pipefail

function image_version () {
    yaml="$1"
    image="$2"

    # get version from `image: quay.io/k8scsi/csi-attacher:v1.0.1`
    grep "image:.*$image" "$yaml" | sed -e 's/.*:v/v/'
}

BASE_DIR=$(dirname "$0")
K8S_RELEASE=${K8S_RELEASE:-"release-1.13"}
PROVISIONER_RELEASE=${PROVISIONER_RELEASE:-$(image_version "${BASE_DIR}/hostpath/csi-hostpath-provisioner.yaml" csi-provisioner)}
ATTACHER_RELEASE=${ATTACHER_RELEASE:-$(image_version "${BASE_DIR}/hostpath/csi-hostpath-attacher.yaml" csi-attacher)}
SNAPSHOTTER_RELEASE=${SNAPSHOTTER_RELEASE:-$(image_version "${BASE_DIR}/snapshotter/csi-hostpath-snpshotter.yaml" csi-snapshotter)}
INSTALL_CRD=${INSTALL_CRD:-"false"}

# apply CSIDriver and CSINodeInfo API objects
if [[ "${INSTALL_CRD}" =~ ^(y|Y|yes|true)$ ]] ; then
    echo "installing CRDs"
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/csi-api/${K8S_RELEASE}/pkg/crd/manifests/csidriver.yaml --validate=false
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/csi-api/${K8S_RELEASE}/pkg/crd/manifests/csinodeinfo.yaml --validate=false
fi

# rbac rules
echo "applying RBAC rules"
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-provisioner/${PROVISIONER_RELEASE}/deploy/kubernetes/rbac.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-attacher/${ATTACHER_RELEASE}/deploy/kubernetes/rbac.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${SNAPSHOTTER_RELEASE}/deploy/kubernetes/rbac.yaml

# deploy hostpath plugin and registrar sidecar
echo "deploying hostpath components"
kubectl apply -f ${BASE_DIR}/hostpath

# deploy snapshotter
echo "deploying snapshotter"
kubectl apply -f ${BASE_DIR}/snapshotter
