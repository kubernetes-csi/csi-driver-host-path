#!/usr/bin/env bash

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

# This script captures the steps required to successfully
# deploy the hostpath plugin driver.  This should be considered
# authoritative and all updates for this process should be
# done here and referenced elsewhere.

# The script assumes that kubectl is available on the OS path
# where it is executed.

set -euxo pipefail

BASE_DIR="$( cd "$( dirname "$0" )" && pwd )"

TEMP_DIR="$( mktemp -d )"

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
# - CSI_SNAPSHOT_METADATA_REGISTRY
# - CSI_SNAPSHOT_METADATA_TAG
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
# Set VOLUME_MODE_CONVERSION_TESTS to "true" to enable the feature in external-provisioner.

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

    if ! [ -f "$yaml" ]; then
        # Fall back to csi-hostpath-plugin.yaml for those deployments which do not
        # have individual pods for the sidecars.
        yaml="$(dirname "$yaml")/csi-hostpath-plugin.yaml"
    fi

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

function volume_mode_conversion () {
    [ "${VOLUME_MODE_CONVERSION_TESTS}" == "true" ]
}

function snapshot_metadata () {
    [ "${SNAPSHOT_METADATA_TESTS}" == "true" ]
}

# In addition, the RBAC rules can be overridden separately.
# For snapshotter 2.0+, the directory has changed.
SNAPSHOTTER_RBAC_RELATIVE_PATH="rbac.yaml"
if version_gt $(rbac_version "${BASE_DIR}/hostpath/csi-hostpath-snapshotter.yaml" csi-snapshotter "${UPDATE_RBAC_RULES}") "v1.255.255"; then
	SNAPSHOTTER_RBAC_RELATIVE_PATH="csi-snapshotter/rbac-csi-snapshotter.yaml"
fi
SNAPSHOT_METADATA_RBAC_RELATIVE_PATH="snapshot-metadata-cluster-role.yaml"
SNAPSHOT_METADATA_SIDECAR_PATCH_RELATIVE_PATH="${BASE_DIR}/hostpath/csi-snapshot-metadata-sidecar.patch"

CSI_PROVISIONER_RBAC_YAML="https://raw.githubusercontent.com/kubernetes-csi/external-provisioner/$(rbac_version "${BASE_DIR}/hostpath/csi-hostpath-provisioner.yaml" csi-provisioner false)/deploy/kubernetes/rbac.yaml"
: ${CSI_PROVISIONER_RBAC:=https://raw.githubusercontent.com/kubernetes-csi/external-provisioner/$(rbac_version "${BASE_DIR}/hostpath/csi-hostpath-provisioner.yaml" csi-provisioner "${UPDATE_RBAC_RULES}")/deploy/kubernetes/rbac.yaml}
CSI_ATTACHER_RBAC_YAML="https://raw.githubusercontent.com/kubernetes-csi/external-attacher/$(rbac_version "${BASE_DIR}/hostpath/csi-hostpath-attacher.yaml" csi-attacher false)/deploy/kubernetes/rbac.yaml"
: ${CSI_ATTACHER_RBAC:=https://raw.githubusercontent.com/kubernetes-csi/external-attacher/$(rbac_version "${BASE_DIR}/hostpath/csi-hostpath-attacher.yaml" csi-attacher "${UPDATE_RBAC_RULES}")/deploy/kubernetes/rbac.yaml}
CSI_SNAPSHOTTER_RBAC_YAML="https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/$(rbac_version "${BASE_DIR}/hostpath/csi-hostpath-snapshotter.yaml" csi-snapshotter false)/deploy/kubernetes/${SNAPSHOTTER_RBAC_RELATIVE_PATH}"
: ${CSI_SNAPSHOTTER_RBAC:=https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/$(rbac_version "${BASE_DIR}/hostpath/csi-hostpath-snapshotter.yaml" csi-snapshotter "${UPDATE_RBAC_RULES}")/deploy/kubernetes/${SNAPSHOTTER_RBAC_RELATIVE_PATH}}
CSI_SNAPSHOT_METADATA_RBAC_YAML="https://raw.githubusercontent.com/kubernetes-csi/external-snapshot-metadata/$(rbac_version "${BASE_DIR}/hostpath/csi-snapshot-metadata-sidecar.patch" csi-snapshot-metadata false)/deploy/${SNAPSHOT_METADATA_RBAC_RELATIVE_PATH}"
: ${CSI_SNAPSHOT_METADATA_RBAC:=https://raw.githubusercontent.com/kubernetes-csi/external-snapshot-metadata/$(rbac_version "${BASE_DIR}/hostpath/csi-snapshot-metadata-sidecar.patch" csi-snapshot-metadata "${UPDATE_RBAC_RULES}")/deploy/${SNAPSHOT_METADATA_RBAC_RELATIVE_PATH}}
CSI_RESIZER_RBAC_YAML="https://raw.githubusercontent.com/kubernetes-csi/external-resizer/$(rbac_version "${BASE_DIR}/hostpath/csi-hostpath-resizer.yaml" csi-resizer false)/deploy/kubernetes/rbac.yaml"
: ${CSI_RESIZER_RBAC:=https://raw.githubusercontent.com/kubernetes-csi/external-resizer/$(rbac_version "${BASE_DIR}/hostpath/csi-hostpath-resizer.yaml" csi-resizer "${UPDATE_RBAC_RULES}")/deploy/kubernetes/rbac.yaml}

CSI_EXTERNALHEALTH_MONITOR_RBAC_YAML="https://raw.githubusercontent.com/kubernetes-csi/external-health-monitor/$(rbac_version "${BASE_DIR}/hostpath/csi-hostpath-plugin.yaml" csi-external-health-monitor-controller false)/deploy/kubernetes/external-health-monitor-controller/rbac.yaml"
: ${CSI_EXTERNALHEALTH_MONITOR_RBAC:=https://raw.githubusercontent.com/kubernetes-csi/external-health-monitor/$(rbac_version "${BASE_DIR}/hostpath/csi-hostpath-plugin.yaml" csi-external-health-monitor-controller "${UPDATE_RBAC_RULES}")/deploy/kubernetes/external-health-monitor-controller/rbac.yaml}

CSI_SNAPSHOT_METADATA_TLS_CERT_YAML="https://raw.githubusercontent.com/kubernetes-csi/external-snapshot-metadata/$(rbac_version "${BASE_DIR}/hostpath/csi-snapshot-metadata-sidecar.patch" csi-snapshot-metadata false)/deploy/example/csi-driver/testdata/csi-snapshot-metadata-tls-secret.yaml"
SNAPSHOT_METADATA_SERVICE_CR_YAML="https://raw.githubusercontent.com/kubernetes-csi/external-snapshot-metadata/$(rbac_version "${BASE_DIR}/hostpath/csi-snapshot-metadata-sidecar.patch" csi-snapshot-metadata false)/deploy/example/csi-driver/testdata/snapshotmetadataservice.yaml"
CSI_SNAPSHOT_METADATA_SERVICE_YAML="https://raw.githubusercontent.com/kubernetes-csi/external-snapshot-metadata/$(rbac_version "${BASE_DIR}/hostpath/csi-snapshot-metadata-sidecar.patch" csi-snapshot-metadata false)/deploy/example/csi-driver/testdata/csi-snapshot-metadata-service.yaml"

INSTALL_CRD=${INSTALL_CRD:-"false"}
VOLUME_MODE_CONVERSION_TESTS=${VOLUME_MODE_CONVERSION_TESTS:-"false"}
SNAPSHOT_METADATA_TESTS=${SNAPSHOT_METADATA_TESTS:-"false"}

# Some images are not affected by *_REGISTRY/*_TAG and IMAGE_* variables.
# The default is to update unless explicitly excluded.
update_image () {
    case "$1" in socat) return 1;; esac
}

run () {
    echo "$@" >&2
    "$@"
}

CSI_HOSTPATHPLUGIN_TEMPLATES=${BASE_DIR}/../../charts/csi-hostpathplugin/templates
CSI_HOSTPATHPLUGIN_YAML=${CSI_HOSTPATHPLUGIN_TEMPLATES}/csi-hostpath-plugin.yaml
CSI_STORAGECLASS_YAML=${CSI_HOSTPATHPLUGIN_TEMPLATES}/csi-storageclass.yaml
CSI_VOLUMESNAPSHOTCLASS_YAML=${CSI_HOSTPATHPLUGIN_TEMPLATES}/csi-volumesnapshotclass.yaml
CSI_GROUPSNAPSHOTCLASS_YAML=${CSI_HOSTPATHPLUGIN_TEMPLATES}/csi-groupsnapshotclass-v1beta1.yaml

# rbac rules
echo "applying RBAC rules"
components=(CSI_PROVISIONER CSI_ATTACHER CSI_SNAPSHOTTER CSI_RESIZER CSI_EXTERNALHEALTH_MONITOR)
if snapshot_metadata; then
    components+=(CSI_SNAPSHOT_METADATA)
fi
for component in "${components[@]}"; do
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

    cp -rfp ${TEMP_DIR}/rbac.yaml ${CSI_HOSTPATHPLUGIN_TEMPLATES}/$component-rbac.yaml
done

# deploy hostpath plugin and registrar sidecar
echo "deploying hostpath components"
for i in $(ls ${BASE_DIR}/hostpath/csi-hostpath-driverinfo.yaml ${BASE_DIR}/hostpath/csi-hostpath-plugin.yaml | sort); do
    echo "   $i"
    if volume_mode_conversion; then
      sed -i -e 's/# end csi-provisioner args/- \"--prevent-volume-mode-conversion=true\"\n            # end csi-provisioner args/' $i
    fi

    # Add external-snapshot-metadata sidecar to the driver, mount TLS certs,
    # and enable snapshot-metadata service
    if snapshot_metadata; then
      sed -i -e "/# end csi containers/r ${SNAPSHOT_METADATA_SIDECAR_PATCH_RELATIVE_PATH}" $i
      sed -i -e 's/# end csi volumes/- name: csi-snapshot-metadata-server-certs\n          secret:\n            secretName: csi-snapshot-metadata-certs\n        # end csi volumes/' $i
      sed -i -e 's/# end hostpath args/- \"--enable-snapshot-metadata\"\n            # end hostpath args/' $i
    fi
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
                suffix=$(eval echo :\${${varname}_TAG:-${IMAGE_TAG:-${tag}}})
                line="$(echo "$nocomments" | sed -e "s;$image;${prefix}${name}${suffix};")"
            fi
            echo "        using $line" >&2
        fi
        echo "$line"
    done)"
    echo "$modified" > ${CSI_HOSTPATHPLUGIN_TEMPLATES}/$(basename $i)
done

cp -rfp examples/*class*.yaml ${CSI_HOSTPATHPLUGIN_TEMPLATES}/

# Convert csi-hostpathplugin as Heal Charts.
find ${CSI_HOSTPATHPLUGIN_TEMPLATES} -type f -name '*.yaml' | while read -r _YAML
do
    yq -i -P -I 2 '... comments=""' $_YAML
    yq -i e '(select(.metadata | has("namespace")) | .metadata.namespace) = "{{ .Release.Namespace }}"' $_YAML
    yq -i e '(select(.subjects.[] | has("namespace")) | .subjects.[].namespace) = "{{ .Release.Namespace }}"' $_YAML
done

yq -i e '(select(.spec | has("replicas")) | .spec.replicas) = "{{ .Values.csiHostpathplugin.replicaCount }}"' $CSI_HOSTPATHPLUGIN_YAML

yq -i e '(select(.spec | has("template")) | .spec.template.spec.containers.[] | select(.name == "hostpath") | .image) = "{{ .Values.csiHostpathplugin.hostpath.image.repository }}:{{ .Values.csiHostpathplugin.hostpath.image.tag }}"' $CSI_HOSTPATHPLUGIN_YAML
yq -i e '(select(.spec | has("template")) | .spec.template.spec.containers.[] | select(.name == "hostpath") | .imagePullPolicy) = "{{ .Values.csiHostpathplugin.hostpath.image.pullPolicy }}"' $CSI_HOSTPATHPLUGIN_YAML

yq -i e '(select(.spec | has("template")) | .spec.template.spec.containers.[] | select(.name == "csi-external-health-monitor-controller") | .image) = "{{ .Values.csiHostpathplugin.csiExternalHealthMonitorController.image.repository }}:{{ .Values.csiHostpathplugin.csiExternalHealthMonitorController.image.tag }}"' $CSI_HOSTPATHPLUGIN_YAML
yq -i e '(select(.spec | has("template")) | .spec.template.spec.containers.[] | select(.name == "csi-external-health-monitor-controller") | .imagePullPolicy) = "{{ .Values.csiHostpathplugin.csiExternalHealthMonitorController.image.pullPolicy }}"' $CSI_HOSTPATHPLUGIN_YAML

yq -i e '(select(.spec | has("template")) | .spec.template.spec.containers.[] | select(.name == "node-driver-registrar") | .image) = "{{ .Values.csiHostpathplugin.nodeDriverRegistrar.image.repository }}:{{ .Values.csiHostpathplugin.nodeDriverRegistrar.image.tag }}"' $CSI_HOSTPATHPLUGIN_YAML
yq -i e '(select(.spec | has("template")) | .spec.template.spec.containers.[] | select(.name == "node-driver-registrar") | .imagePullPolicy) = "{{ .Values.csiHostpathplugin.nodeDriverRegistrar.image.pullPolicy }}"' $CSI_HOSTPATHPLUGIN_YAML

yq -i e '(select(.spec | has("template")) | .spec.template.spec.containers.[] | select(.name == "liveness-probe") | .image) = "{{ .Values.csiHostpathplugin.livenessProbe.image.repository }}:{{ .Values.csiHostpathplugin.livenessProbe.image.tag }}"' $CSI_HOSTPATHPLUGIN_YAML
yq -i e '(select(.spec | has("template")) | .spec.template.spec.containers.[] | select(.name == "liveness-probe") | .imagePullPolicy) = "{{ .Values.csiHostpathplugin.livenessProbe.image.pullPolicy }}"' $CSI_HOSTPATHPLUGIN_YAML

yq -i e '(select(.spec | has("template")) | .spec.template.spec.containers.[] | select(.name == "csi-attacher") | .image) = "{{ .Values.csiHostpathplugin.csiAttacher.image.repository }}:{{ .Values.csiHostpathplugin.csiAttacher.image.tag }}"' $CSI_HOSTPATHPLUGIN_YAML
yq -i e '(select(.spec | has("template")) | .spec.template.spec.containers.[] | select(.name == "csi-attacher") | .imagePullPolicy) = "{{ .Values.csiHostpathplugin.csiAttacher.image.pullPolicy }}"' $CSI_HOSTPATHPLUGIN_YAML

yq -i e '(select(.spec | has("template")) | .spec.template.spec.containers.[] | select(.name == "csi-provisioner") | .image) = "{{ .Values.csiHostpathplugin.csiProvisioner.image.repository }}:{{ .Values.csiHostpathplugin.csiProvisioner.image.tag }}"' $CSI_HOSTPATHPLUGIN_YAML
yq -i e '(select(.spec | has("template")) | .spec.template.spec.containers.[] | select(.name == "csi-provisioner") | .imagePullPolicy) = "{{ .Values.csiHostpathplugin.csiProvisioner.image.pullPolicy }}"' $CSI_HOSTPATHPLUGIN_YAML

yq -i e '(select(.spec | has("template")) | .spec.template.spec.containers.[] | select(.name == "csi-resizer") | .image) = "{{ .Values.csiHostpathplugin.csiResizer.image.repository }}:{{ .Values.csiHostpathplugin.csiResizer.image.tag }}"' $CSI_HOSTPATHPLUGIN_YAML
yq -i e '(select(.spec | has("template")) | .spec.template.spec.containers.[] | select(.name == "csi-resizer") | .imagePullPolicy) = "{{ .Values.csiHostpathplugin.csiResizer.image.pullPolicy }}"' $CSI_HOSTPATHPLUGIN_YAML

yq -i e '(select(.spec | has("template")) | .spec.template.spec.containers.[] | select(.name == "csi-snapshotter") | .image) = "{{ .Values.csiHostpathplugin.csiSnapshotter.image.repository }}:{{ .Values.csiHostpathplugin.csiSnapshotter.image.tag }}"' $CSI_HOSTPATHPLUGIN_YAML
yq -i e '(select(.spec | has("template")) | .spec.template.spec.containers.[] | select(.name == "csi-snapshotter") | .imagePullPolicy) = "{{ .Values.csiHostpathplugin.csiSnapshotter.image.pullPolicy }}"' $CSI_HOSTPATHPLUGIN_YAML

yq -i e '(select(.spec | has("template")) | .spec.template.spec.volumes.[] | select(.name == "csi-data-dir") | .hostPath.path) = "{{ .Values.csiHostpathplugin.volumes.csiDataDir.path }}"' $CSI_HOSTPATHPLUGIN_YAML

yq -i e '(select(.metadata | has("name")) | .metadata.name) = "{{ .Values.storageClass.name }}"' $CSI_STORAGECLASS_YAML
sed -i '1s/^/{{- if .Values.storageClass.create -}}\n/g' $CSI_STORAGECLASS_YAML
sed -i '$s/$/\n{{- end -}}/g' $CSI_STORAGECLASS_YAML

yq -i e '(select(.metadata | has("name")) | .metadata.name) = "{{ .Values.volumeSnapshotClass.name }}"' $CSI_VOLUMESNAPSHOTCLASS_YAML
sed -i '1s/^/{{- if .Values.volumeSnapshotClass.create -}}\n/g' $CSI_VOLUMESNAPSHOTCLASS_YAML
sed -i '$s/$/\n{{- end -}}/g' $CSI_VOLUMESNAPSHOTCLASS_YAML

yq -i e '(select(.metadata | has("name")) | .metadata.name) = "{{ .Values.volumeGroupSnapshotClass.name }}"' $CSI_GROUPSNAPSHOTCLASS_YAML
sed -i '1s/^/{{- if .Values.volumeGroupSnapshotClass.create -}}\n/g' $CSI_GROUPSNAPSHOTCLASS_YAML
sed -i '$s/$/\n{{- end -}}/g' $CSI_GROUPSNAPSHOTCLASS_YAML

find ${CSI_HOSTPATHPLUGIN_TEMPLATES} -type f -name '*.yaml' | xargs sed -i "s/'{{/{{/g;s/}}'/}}/g"
