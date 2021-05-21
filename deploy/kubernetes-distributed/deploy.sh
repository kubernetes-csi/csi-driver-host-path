#!/usr/bin/env bash

# This script captures the steps required to successfully
# deploy the hostpath plugin driver.  This should be considered
# authoritative and all updates for this process should be
# done here and referenced elsewhere.

# The script assumes that kubectl is available on the OS path
# where it is executed.

set -e
set -o pipefail

BASE_DIR=$(dirname "$0")

TEMP_DIR="$( mktemp -d )"
trap 'rm -rf ${TEMP_DIR}' EXIT

# KUBELET_DATA_DIR can be set to replace the default /var/lib/kubelet.
# All nodes must use the same directory.
default_kubelet_data_dir=/var/lib/kubelet
: ${KUBELET_DATA_DIR:=${default_kubelet_data_dir}}

# If set, the following env variables override image registry and/or tag for each of the images.
# They are named after the image name, with hyphen replaced by underscore and in upper case.
#
# - CSI_ATTACHER_REGISTRY
# - CSI_ATTACHER_TAG
# - CSI_NODE_DRIVER_REGISTRAR_REGISTRY
# - CSI_NODE_DRIVER_REGISTRAR_TAG
# - CSI_PROVISIONER_REGISTRY
# - CSI_PROVISIONER_TAG
# - CSI_SNAPSHOTTER_REGISTRY
# - CSI_SNAPSHOTTER_TAG
# - HOSTPATHPLUGIN_REGISTRY
# - HOSTPATHPLUGIN_TAG
#
# Alternatively, it is possible to override all registries or tags with:
# - IMAGE_REGISTRY
# - IMAGE_TAG
# These are used as fallback when the more specific variables are unset or empty.
#
# IMAGE_TAG=canary is ignored for images that are blacklisted in the
# deployment's optional canary-blacklist.txt file. This is meant for
# images which have known API breakages and thus cannot work in those
# deployments anymore. That text file must have the name of the blacklisted
# image on a line by itself, other lines are ignored. Example:
#
#     # The following canary images are known to be incompatible with this
#     # deployment:
#     csi-snapshotter
#
# Beware that the .yaml files do not have "imagePullPolicy: Always". That means that
# also the "canary" images will only be pulled once. This is good for testing
# (starting a pod multiple times will always run with the same canary image), but
# implies that refreshing that image has to be done manually.
#
# As a special case, 'none' as registry removes the registry name.


# The default is to use the RBAC rules that match the image that is
# being used, also in the case that the image gets overridden. This
# way if there are breaking changes in the RBAC rules, the deployment
# will continue to work.
#
# However, such breaking changes should be rare and only occur when updating
# to a new major version of a sidecar. Nonetheless, to allow testing the scenario
# where the image gets overridden but not the RBAC rules, updating the RBAC
# rules can be disabled.
: ${UPDATE_RBAC_RULES:=true}
function rbac_version () {
    yaml="$1"
    image="$2"
    update_rbac="$3"

    # get version from `image: quay.io/k8scsi/csi-attacher:v1.0.1`, ignoring comments
    version="$(sed -e 's/ *#.*$//' "$yaml" | grep "image:.*$image" | sed -e 's/ *#.*//' -e 's/.*://')"

    if $update_rbac; then
        # apply overrides
        varname=$(echo $image | tr - _ | tr a-z A-Z)
        eval version=\${${varname}_TAG:-\${IMAGE_TAG:-\$version}}
    fi

    # When using canary images, we have to assume that the
    # canary images were built from the corresponding branch.
    case "$version" in canary) version=master;;
                        *-canary) version="$(echo "$version" | sed -e 's/\(.*\)-canary/release-\1/')";;
    esac

    echo "$version"
}

# version_gt returns true if arg1 is greater than arg2.
# 
# This function expects versions to be one of the following formats:
#   X.Y.Z, release-X.Y.Z, vX.Y.Z
#
#   where X,Y, and Z are any number.
#
# Partial versions (1.2, release-1.2) work as well.
# The follow substrings are stripped before version comparison:
#   - "v"
#   - "release-"
#
# Usage:
# version_gt release-1.3 v1.2.0  (returns true)
# version_gt v1.1.1 v1.2.0  (returns false)
# version_gt 1.1.1 v1.2.0  (returns false)
# version_gt 1.3.1 v1.2.0  (returns true)
# version_gt 1.1.1 release-1.2.0  (returns false)
# version_gt 1.2.0 1.2.2  (returns false)
function version_gt() { 
    versions=$(for ver in "$@"; do ver=${ver#release-}; ver=${ver#kubernetes-}; echo ${ver#v}; done)
    greaterVersion=${1#"release-"};  
    greaterVersion=${greaterVersion#"kubernetes-"};
    greaterVersion=${greaterVersion#"v"}; 
    test "$(printf '%s' "$versions" | sort -V | head -n 1)" != "$greaterVersion"
}


CSI_PROVISIONER_RBAC_YAML="https://raw.githubusercontent.com/kubernetes-csi/external-provisioner/$(rbac_version "${BASE_DIR}/hostpath/csi-hostpath-plugin.yaml" csi-provisioner false)/deploy/kubernetes/rbac.yaml"
: ${CSI_PROVISIONER_RBAC:=https://raw.githubusercontent.com/kubernetes-csi/external-provisioner/$(rbac_version "${BASE_DIR}/hostpath/csi-hostpath-plugin.yaml" csi-provisioner "${UPDATE_RBAC_RULES}")/deploy/kubernetes/rbac.yaml}

# Some images are not affected by *_REGISTRY/*_TAG and IMAGE_* variables.
# The default is to update unless explicitly excluded.
update_image () {
    case "$1" in socat) return 1;; esac
}

run () {
    echo "$@" >&2
    "$@"
}

# rbac rules
echo "applying RBAC rules"
for component in CSI_PROVISIONER; do
    eval current="\${${component}_RBAC}"
    eval original="\${${component}_RBAC_YAML}"
    if [ "$current" != "$original" ]; then
        echo "Using non-default RBAC rules for $component. Changes from $original to $current are:"
        diff -c <(wget --quiet -O - "$original") <(if [[ "$current" =~ ^http ]]; then wget --quiet -O - "$current"; else cat "$current"; fi) || true
    fi

    # using kustomize kubectl plugin to add labels to he rbac files.
    # since we are deploying rbas directly with the url, the kustomize plugin only works with the local files
    # we need to add the files locally in temp folder and using kustomize adding labels it will be applied
    if [[ "${current}" =~ ^http:// ]] || [[ "${current}" =~ ^https:// ]]; then
      run curl "${current}" --output "${TEMP_DIR}"/rbac.yaml --silent --location
    else
        # Even for local files we need to copy because kustomize only supports files inside
        # the root of a kustomization.
        cp "${current}" "${TEMP_DIR}"/rbac.yaml
    fi

    cat <<- EOF > "${TEMP_DIR}"/kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

commonLabels:
  app.kubernetes.io/instance: hostpath.csi.k8s.io
  app.kubernetes.io/part-of: csi-driver-host-path

resources:
- ./rbac.yaml
EOF

    run kubectl apply --kustomize "${TEMP_DIR}"
done

# The cluster must support exactly the version that the external-provisioner supports.
# The problem then becomes that the version of the external-provisioner might get
# changed via CSI_PROVISIONER_TAG, so we cannot just check for the version currently
# listed in the YAML file.
case "$CSI_PROVISIONER_TAG" in
    "") csistoragecapacities_api=v1alpha1;; # unchanged, assume version from YAML
    *) csistoragecapacities_api=v1beta1;; # set, assume that it is more recent *and* a version that uses v1beta1 (https://github.com/kubernetes-csi/external-provisioner/pull/584)
esac
get_csistoragecapacities=$(kubectl get csistoragecapacities.${csistoragecapacities_api}.storage.k8s.io 2>&1 || true)
if  echo "$get_csistoragecapacities" | grep -q "the server doesn't have a resource type"; then
    have_csistoragecapacity=false
else
    have_csistoragecapacity=true
    echo "csistoragecapacities.${csistoragecapacities_api}.storage.k8s.io:"
    show_lines=4
    echo "$get_csistoragecapacities" | head -n $show_lines | sed -e 's/^/   /'
    if [ $(echo "$get_csistoragecapacities" | wc -l) -gt $show_lines ]; then
        echo "   ..."
    fi
fi
echo "deploying with CSIStorageCapacity $csistoragecapacities_api: $have_csistoragecapacity"

# deploy hostpath plugin and registrar sidecar
echo "deploying hostpath components"
for i in $(ls ${BASE_DIR}/hostpath/*.yaml | sort); do
    echo "   $i"
    modified="$(cat "$i" | sed -e "s;${default_kubelet_data_dir}/;${KUBELET_DATA_DIR}/;" | while IFS= read -r line; do
        nocomments="$(echo "$line" | sed -e 's/ *#.*$//')"
        if echo "$nocomments" | grep -q '^[[:space:]]*image:[[:space:]]*'; then
            # Split 'image: quay.io/k8scsi/csi-attacher:v1.0.1'
            # into image (quay.io/k8scsi/csi-attacher:v1.0.1),
            # registry (quay.io/k8scsi),
            # name (csi-attacher),
            # tag (v1.0.1).
            image=$(echo "$nocomments" | sed -e 's;.*image:[[:space:]]*;;')
            registry=$(echo "$image" | sed -e 's;\(.*\)/.*;\1;')
            name=$(echo "$image" | sed -e 's;.*/\([^:]*\).*;\1;')
            tag=$(echo "$image" | sed -e 's;.*:;;')

            # Variables are with underscores and upper case.
            varname=$(echo $name | tr - _ | tr a-z A-Z)

            # Now replace registry and/or tag, if set as env variables.
            # If not set, the replacement is the same as the original value.
            # Only do this for the images which are meant to be configurable.
            if update_image "$name"; then
                prefix=$(eval echo \${${varname}_REGISTRY:-${IMAGE_REGISTRY:-${registry}}}/ | sed -e 's;none/;;')
                if [ "$IMAGE_TAG" = "canary" ] &&
                   [ -f ${BASE_DIR}/canary-blacklist.txt ] &&
                   grep -q "^$name\$" ${BASE_DIR}/canary-blacklist.txt; then
                    # Ignore IMAGE_TAG=canary for this particular image because its
                    # canary image is blacklisted in the deployment blacklist.
                    suffix=$(eval echo :\${${varname}_TAG:-${tag}})
                else
                    suffix=$(eval echo :\${${varname}_TAG:-${IMAGE_TAG:-${tag}}})
                fi
                line="$(echo "$nocomments" | sed -e "s;$image;${prefix}${name}${suffix};")"
            fi
            echo "        using $line" >&2
        fi
        if ! $have_csistoragecapacity; then
            line="$(echo "$line" | grep -v -e 'storageCapacity: true' -e '--enable-capacity')"
        fi
        echo "$line"
    done)"
    if ! echo "$modified" | kubectl apply -f -; then
        echo "modified version of $i:"
        echo "$modified"
        exit 1
    fi
done

wait_for_daemonset () {
    retries=10
    while [ $retries -ge 0 ]; do
        ready=$(kubectl get -n $1 daemonset $2 -o jsonpath="{.status.numberReady}")
        required=$(kubectl get -n $1 daemonset $2 -o jsonpath="{.status.desiredNumberScheduled}")
        if [ $ready -gt 0 ] && [ $ready -eq $required ]; then
            return 0
        fi
        retries=$((retries - 1))
        sleep 3
    done
    return 1
}


# Wait until the DaemonSet is running on all nodes.
if ! wait_for_daemonset default csi-hostpathplugin; then
    echo
    echo "driver not ready"
    echo "Deployment:"
    (set +e; set -x; kubectl describe all,role,clusterrole,rolebinding,clusterrolebinding,serviceaccount,storageclass,csidriver --all-namespaces -l app.kubernetes.io/instance=hostpath.csi.k8s.io)
    echo
    echo "Pod logs:"
    kubectl get pods -l app.kubernetes.io/instance=hostpath.csi.k8s.io --all-namespaces -o=jsonpath='{range .items[*]}{.metadata.name}{" "}{range .spec.containers[*]}{.name}{" "}{end}{"\n"}{end}' | while read -r pod containers; do
        for c in $containers; do
            echo
            (set +e; set -x; kubectl logs $pod $c)
        done
    done
    exit 1
fi

# Create a test driver configuration in the place where the prow job
# expects it?
if [ "${CSI_PROW_TEST_DRIVER}" ]; then
    sed -e "s/capacity: true/capacity: ${have_csistoragecapacity}/" "${BASE_DIR}/test-driver.yaml" >"${CSI_PROW_TEST_DRIVER}"
fi
