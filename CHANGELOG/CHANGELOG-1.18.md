# Release notes for v1.18.0

# Changelog since v1.17.0

## Changes by Kind

### Feature

- Add `ControllerServer.GetSnapshot` CSI procedure. ([#606](https://github.com/kubernetes-csi/csi-driver-host-path/pull/606), [@nixpanic](https://github.com/nixpanic))
- Add support for VolumeAttributesClasses (VACs). ([#618](https://github.com/kubernetes-csi/csi-driver-host-path/pull/618), [@gnufied](https://github.com/gnufied))

### Bug or Regression

- Fix inverted error check in `DeleteSnapshot` so errors from `GetSnapshotByID` are handled correctly for group snapshots. ([#658](https://github.com/kubernetes-csi/csi-driver-host-path/pull/658), [@carterpewpew](https://github.com/carterpewpew))
- Fix the node-driver-registrar liveness probe to use `--http-endpoint`. ([#650](https://github.com/kubernetes-csi/csi-driver-host-path/pull/650), [@humblec](https://github.com/humblec))
- Fix idempotency of `CreateVolumeGroupSnapshot`. ([#627](https://github.com/kubernetes-csi/csi-driver-host-path/pull/627), [@jsafrane](https://github.com/jsafrane))
- Prevent duplicate sidecars by using temporary copies instead of `sed -i` when updating deployment manifests. ([#621](https://github.com/kubernetes-csi/csi-driver-host-path/pull/621), [@kaovilai](https://github.com/kaovilai))
- Add signal catch to stop the server gracefully. ([#455](https://github.com/kubernetes-csi/csi-driver-host-path/pull/455), [@astraw99](https://github.com/astraw99))

### Uncategorized

- Update csi-snapshot-metadata image version to v1.0.0. ([#645](https://github.com/kubernetes-csi/csi-driver-host-path/pull/645), [@PrasadG193](https://github.com/PrasadG193))
- Bump Kubernetes dependencies to v1.36.1. ([#663](https://github.com/kubernetes-csi/csi-driver-host-path/pull/663), [@dfajmon](https://github.com/dfajmon))
- Update deployment structure for Kubernetes 1.36. ([#668](https://github.com/kubernetes-csi/csi-driver-host-path/pull/668), [@xing-yang](https://github.com/xing-yang))


## Dependencies

### Added
- go.yaml.in/yaml/v2: v2.4.3
- sigs.k8s.io/json: v0.0.0-20250730193827-2d320260d730

### Changed
- github.com/container-storage-interface/spec: [v1.11.0 → v1.12.0](https://github.com/container-storage-interface/spec/compare/v1.11.0...v1.12.0)
- github.com/fxamacker/cbor/v2: [v2.7.0 → v2.9.0](https://github.com/fxamacker/cbor/v2/compare/v2.7.0...v2.9.0)
- github.com/go-logr/logr: [v1.4.2 → v1.4.3](https://github.com/go-logr/logr/compare/v1.4.2...v1.4.3)
- github.com/kubernetes-csi/csi-lib-utils: [v0.22.0 → v0.24.0](https://github.com/kubernetes-csi/csi-lib-utils/compare/v0.22.0...v0.24.0)
- github.com/prometheus/client_golang: [v1.22.0 → v1.23.2](https://github.com/prometheus/client_golang/compare/v1.22.0...v1.23.2)
- github.com/prometheus/client_model: [v0.6.1 → v0.6.2](https://github.com/prometheus/client_model/compare/v0.6.1...v0.6.2)
- github.com/prometheus/common: [v0.62.0 → v0.67.5](https://github.com/prometheus/common/compare/v0.62.0...v0.67.5)
- github.com/prometheus/procfs: [v0.15.1 → v0.19.2](https://github.com/prometheus/procfs/compare/v0.15.1...v0.19.2)
- github.com/spf13/pflag: [v1.0.5 → v1.0.10](https://github.com/spf13/pflag/compare/v1.0.5...v1.0.10)
- github.com/stretchr/testify: [v1.10.0 → v1.11.1](https://github.com/stretchr/testify/compare/v1.10.0...v1.11.1)
- go.opentelemetry.io/otel/trace: v1.33.0 → v1.43.0
- go.opentelemetry.io/otel: v1.33.0 → v1.43.0
- golang.org/x/net: v0.38.0 → v0.49.0
- golang.org/x/sys: v0.31.0 → v0.40.0
- golang.org/x/text: v0.23.0 → v0.33.0
- google.golang.org/genproto/googleapis/rpc: 9240e9c → 8636f87
- google.golang.org/grpc: v1.69.0 → v1.79.3
- google.golang.org/protobuf: v1.36.5 → v1.36.12-0.20260120151049-f2248ac996af
- k8s.io/api: v0.33.0 → v0.36.1
- k8s.io/apiextensions-apiserver: v0.33.0 → v0.36.1
- k8s.io/apimachinery: v0.33.0 → v0.36.1
- k8s.io/apiserver: v0.33.0 → v0.36.1
- k8s.io/cli-runtime: v0.33.0 → v0.36.1
- k8s.io/client-go: v0.33.0 → v0.36.1
- k8s.io/cloud-provider: v0.33.0 → v0.36.1
- k8s.io/cluster-bootstrap: v0.33.0 → v0.36.1
- k8s.io/code-generator: v0.33.0 → v0.36.1
- k8s.io/component-base: v0.33.0 → v0.36.1
- k8s.io/component-helpers: v0.33.0 → v0.36.1
- k8s.io/controller-manager: v0.33.0 → v0.36.1
- k8s.io/cri-api: v0.33.0 → v0.36.1
- k8s.io/cri-client: v0.33.0 → v0.36.1
- k8s.io/csi-translation-lib: v0.33.0 → v0.36.1
- k8s.io/dynamic-resource-allocation: v0.33.0 → v0.36.1
- k8s.io/endpointslice: v0.33.0 → v0.36.1
- k8s.io/externaljwt: v0.33.0 → v0.36.1
- k8s.io/klog/v2: v2.130.1 → v2.140.0
- k8s.io/kms: v0.33.0 → v0.36.1
- k8s.io/kube-aggregator: v0.33.0 → v0.36.1
- k8s.io/kube-controller-manager: v0.33.0 → v0.36.1
- k8s.io/kube-proxy: v0.33.0 → v0.36.1
- k8s.io/kube-scheduler: v0.33.0 → v0.36.1
- k8s.io/kubectl: v0.33.0 → v0.36.1
- k8s.io/kubelet: v0.33.0 → v0.36.1
- k8s.io/kubernetes: v1.33.0 → v1.36.1
- k8s.io/metrics: v0.33.0 → v0.36.1
- k8s.io/mount-utils: v0.33.0 → v0.36.1
- k8s.io/pod-security-admission: v0.33.0 → v0.36.1
- k8s.io/sample-apiserver: v0.33.0 → v0.36.1
- k8s.io/utils: 24370be → b8788ab

### Removed
- github.com/gogo/protobuf: v1.3.2
- sigs.k8s.io/yaml: v1.4.0
