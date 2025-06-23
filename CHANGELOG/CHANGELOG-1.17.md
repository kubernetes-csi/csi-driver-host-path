# Release notes for v1.17.0

# Changelog since v1.17.0

## Changes by Kind

### Feature

- Add a new `--automaxprocs` flag to set the `GOMAXPROCS` environment variable to match the configured Linux container CPU quota. ([#605](https://github.com/kubernetes-csi/csi-driver-host-path/pull/605), [@nixpanic](https://github.com/nixpanic))
- Introduce support for the `--enable-list-snapshots=<bool>` command line option, allowing users to explicitly enable or disable the ControllerServiceCapability_RPC_LIST_SNAPSHOTS capability. By default, this option is set to true to preserve existing behavior. ([#597](https://github.com/kubernetes-csi/csi-driver-host-path/pull/597), [@leonardoce](https://github.com/leonardoce))
- Set a random MaxVolumesPerNode if config.AttachLimit == -1. ([#587](https://github.com/kubernetes-csi/csi-driver-host-path/pull/587), [@torredil](https://github.com/torredil))
- Update kubernetes dependencies to v1.33.0 ([#604](https://github.com/kubernetes-csi/csi-driver-host-path/pull/604), [@Aishwarya-Hebbar](https://github.com/Aishwarya-Hebbar))
- Update to csi-snapshot-metadata alpha image. ([#595](https://github.com/kubernetes-csi/csi-driver-host-path/pull/595), [@carlbraganza](https://github.com/carlbraganza))

### Bug or Regression

- Fix the ControllerExpandVolume response returns the wrong NodeExpansionRequired when the `--disable-node-expansion` is true. ([#599](https://github.com/kubernetes-csi/csi-driver-host-path/pull/599), [@carlory](https://github.com/carlory))
- Fixed: Considered only non-zero blocks for GetMetadataAllocated. ([#579](https://github.com/kubernetes-csi/csi-driver-host-path/pull/579), [@iPraveenParihar](https://github.com/iPraveenParihar))
- The `ignoreFailedRead` VolumeSnapshotClass parameter is now honored by using the GNU implementation of `tar`. ([#586](https://github.com/kubernetes-csi/csi-driver-host-path/pull/586), [@leonardoce](https://github.com/leonardoce))

### Uncategorized

- Add better validation of the CapacityRange parameter to CreateVolume. ([#577](https://github.com/kubernetes-csi/csi-driver-host-path/pull/577), [@ebblake](https://github.com/ebblake))

## Dependencies

### Added
- github.com/containerd/errdefs/pkg: [v0.3.0](https://github.com/containerd/errdefs/pkg/tree/v0.3.0)
- github.com/containerd/typeurl/v2: [v2.2.2](https://github.com/containerd/typeurl/v2/tree/v2.2.2)
- github.com/opencontainers/cgroups: [v0.0.1](https://github.com/opencontainers/cgroups/tree/v0.0.1)
- github.com/opencontainers/image-spec: [v1.1.1](https://github.com/opencontainers/image-spec/tree/v1.1.1)
- github.com/prashantv/gostub: [v1.1.0](https://github.com/prashantv/gostub/tree/v1.1.0)
- go.uber.org/automaxprocs: v1.6.0
- gopkg.in/go-jose/go-jose.v2: v2.6.3
- sigs.k8s.io/randfill: v1.0.0

### Changed
- cel.dev/expr: v0.18.0 → v0.19.1
- github.com/containerd/containerd/api: [v1.7.19 → v1.8.0](https://github.com/containerd/containerd/api/compare/v1.7.19...v1.8.0)
- github.com/containerd/errdefs: [v0.1.0 → v1.0.0](https://github.com/containerd/errdefs/compare/v0.1.0...v1.0.0)
- github.com/containerd/ttrpc: [v1.2.5 → v1.2.6](https://github.com/containerd/ttrpc/compare/v1.2.5...v1.2.6)
- github.com/coredns/corefile-migration: [v1.0.24 → v1.0.25](https://github.com/coredns/corefile-migration/compare/v1.0.24...v1.0.25)
- github.com/coreos/go-oidc: [v2.2.1+incompatible → v2.3.0+incompatible](https://github.com/coreos/go-oidc/compare/v2.2.1...v2.3.0)
- github.com/cyphar/filepath-securejoin: [v0.3.4 → v0.4.1](https://github.com/cyphar/filepath-securejoin/compare/v0.3.4...v0.4.1)
- github.com/golang-jwt/jwt/v4: [v4.5.0 → v4.5.2](https://github.com/golang-jwt/jwt/v4/compare/v4.5.0...v4.5.2)
- github.com/google/btree: [v1.0.1 → v1.1.3](https://github.com/google/btree/compare/v1.0.1...v1.1.3)
- github.com/google/cadvisor: [v0.51.0 → v0.52.1](https://github.com/google/cadvisor/compare/v0.51.0...v0.52.1)
- github.com/google/cel-go: [v0.22.0 → v0.23.2](https://github.com/google/cel-go/compare/v0.22.0...v0.23.2)
- github.com/google/go-cmp: [v0.6.0 → v0.7.0](https://github.com/google/go-cmp/compare/v0.6.0...v0.7.0)
- github.com/gorilla/websocket: [v1.5.0 → e064f32](https://github.com/gorilla/websocket/compare/v1.5.0...e064f32)
- github.com/grpc-ecosystem/grpc-gateway/v2: [v2.20.0 → v2.24.0](https://github.com/grpc-ecosystem/grpc-gateway/v2/compare/v2.20.0...v2.24.0)
- github.com/klauspost/compress: [v1.17.11 → v1.18.0](https://github.com/klauspost/compress/compare/v1.17.11...v1.18.0)
- github.com/kubernetes-csi/csi-lib-utils: [v0.20.0 → v0.22.0](https://github.com/kubernetes-csi/csi-lib-utils/compare/v0.20.0...v0.22.0)
- github.com/prometheus/client_golang: [v1.20.5 → v1.22.0](https://github.com/prometheus/client_golang/compare/v1.20.5...v1.22.0)
- github.com/prometheus/common: [v0.61.0 → v0.62.0](https://github.com/prometheus/common/compare/v0.61.0...v0.62.0)
- github.com/rogpeppe/go-internal: [v1.12.0 → v1.13.1](https://github.com/rogpeppe/go-internal/compare/v1.12.0...v1.13.1)
- github.com/vishvananda/netlink: [b1ce50c → 62fb240](https://github.com/vishvananda/netlink/compare/b1ce50c...62fb240)
- go.etcd.io/etcd/api/v3: v3.5.16 → v3.5.21
- go.etcd.io/etcd/client/pkg/v3: v3.5.16 → v3.5.21
- go.etcd.io/etcd/client/v2: v2.305.16 → v2.305.21
- go.etcd.io/etcd/client/v3: v3.5.16 → v3.5.21
- go.etcd.io/etcd/pkg/v3: v3.5.16 → v3.5.21
- go.etcd.io/etcd/raft/v3: v3.5.16 → v3.5.21
- go.etcd.io/etcd/server/v3: v3.5.16 → v3.5.21
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp: v0.53.0 → v0.58.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc: v1.27.0 → v1.33.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace: v1.28.0 → v1.33.0
- go.opentelemetry.io/otel/sdk: v1.31.0 → v1.33.0
- go.opentelemetry.io/proto/otlp: v1.3.1 → v1.4.0
- golang.org/x/crypto: v0.30.0 → v0.36.0
- golang.org/x/net: v0.32.0 → v0.38.0
- golang.org/x/oauth2: v0.24.0 → v0.27.0
- golang.org/x/sync: v0.10.0 → v0.12.0
- golang.org/x/sys: v0.28.0 → v0.31.0
- golang.org/x/term: v0.27.0 → v0.30.0
- golang.org/x/text: v0.21.0 → v0.23.0
- golang.org/x/time: v0.8.0 → v0.9.0
- google.golang.org/genproto/googleapis/api: 796eee8 → e6fa225
- google.golang.org/protobuf: v1.36.0 → v1.36.5
- k8s.io/api: v0.32.0 → v0.33.0
- k8s.io/apiextensions-apiserver: v0.32.0 → v0.33.0
- k8s.io/apimachinery: v0.32.0 → v0.33.0
- k8s.io/apiserver: v0.32.0 → v0.33.0
- k8s.io/cli-runtime: v0.32.0 → v0.33.0
- k8s.io/client-go: v0.32.0 → v0.33.0
- k8s.io/cloud-provider: v0.32.0 → v0.33.0
- k8s.io/cluster-bootstrap: v0.32.0 → v0.33.0
- k8s.io/code-generator: v0.32.0 → v0.33.0
- k8s.io/component-base: v0.32.0 → v0.33.0
- k8s.io/component-helpers: v0.32.0 → v0.33.0
- k8s.io/controller-manager: v0.32.0 → v0.33.0
- k8s.io/cri-api: v0.32.0 → v0.33.0
- k8s.io/cri-client: v0.32.0 → v0.33.0
- k8s.io/csi-translation-lib: v0.32.0 → v0.33.0
- k8s.io/dynamic-resource-allocation: v0.32.0 → v0.33.0
- k8s.io/endpointslice: v0.32.0 → v0.33.0
- k8s.io/externaljwt: v0.32.0 → v0.33.0
- k8s.io/gengo/v2: 2b36238 → 1244d31
- k8s.io/kms: v0.32.0 → v0.33.0
- k8s.io/kube-aggregator: v0.32.0 → v0.33.0
- k8s.io/kube-controller-manager: v0.32.0 → v0.33.0
- k8s.io/kube-openapi: 2c72e55 → c8a335a
- k8s.io/kube-proxy: v0.32.0 → v0.33.0
- k8s.io/kube-scheduler: v0.32.0 → v0.33.0
- k8s.io/kubectl: v0.32.0 → v0.33.0
- k8s.io/kubelet: v0.32.0 → v0.33.0
- k8s.io/kubernetes: v1.32.0 → v1.33.0
- k8s.io/metrics: v0.32.0 → v0.33.0
- k8s.io/mount-utils: v0.32.0 → v0.33.0
- k8s.io/pod-security-admission: v0.32.0 → v0.33.0
- k8s.io/sample-apiserver: v0.32.0 → v0.33.0
- sigs.k8s.io/apiserver-network-proxy/konnectivity-client: v0.31.0 → v0.31.2
- sigs.k8s.io/kustomize/api: v0.18.0 → v0.19.0
- sigs.k8s.io/kustomize/kustomize/v5: v5.5.0 → v5.6.0
- sigs.k8s.io/kustomize/kyaml: v0.18.1 → v0.19.0
- sigs.k8s.io/structured-merge-diff/v4: v4.5.0 → v4.6.0

### Removed
- github.com/asaskevich/govalidator: [f61b66f](https://github.com/asaskevich/govalidator/tree/f61b66f)
- github.com/go-kit/log: [v0.2.1](https://github.com/go-kit/log/tree/v0.2.1)
- github.com/go-logfmt/logfmt: [v0.5.1](https://github.com/go-logfmt/logfmt/tree/v0.5.1)
- github.com/google/gofuzz: [v1.2.0](https://github.com/google/gofuzz/tree/v1.2.0)
- github.com/opencontainers/runc: [v1.2.1](https://github.com/opencontainers/runc/tree/v1.2.1)
- gopkg.in/square/go-jose.v2: v2.6.0
