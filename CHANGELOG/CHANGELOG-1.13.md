# Release notes for v1.13.0

# Changelog since v1.12.0

## Changes by Kind

### Feature

- Add socat binary into registry.k8s.io/sig-storage/hostpathplugin image ([#498](https://github.com/kubernetes-csi/csi-driver-host-path/pull/498), [@carlory](https://github.com/carlory))
- Hostpath driver implements the ControllerModityVolume feature ([#481](https://github.com/kubernetes-csi/csi-driver-host-path/pull/481), [@carlory](https://github.com/carlory))
- Together with the Kubernetes CSI external-snapshotter, the csi-driver-host-path provides ALPHA support for the VolumeGroupSnapshot functionality. ([#399](https://github.com/kubernetes-csi/csi-driver-host-path/pull/399), [@nixpanic](https://github.com/nixpanic))

### Bug or Regression

- Fix missing published target paths when republish the ephemeral volume ([#480](https://github.com/kubernetes-csi/csi-driver-host-path/pull/480), [@carlory](https://github.com/carlory))

### Other (Cleanup or Flake)

- Update deploy directory structure and update sidecar releases, corresponding to Kubernetes 1.29. ([#501](https://github.com/kubernetes-csi/csi-driver-host-path/pull/501), [@RaunakShah](https://github.com/RaunakShah))

### Uncategorized

- Fix an issue with GetVolumeGroupSnapshot where SourceVolumeIDs were compared with SnapshotIDs. ([#494](https://github.com/kubernetes-csi/csi-driver-host-path/pull/494), [@nixpanic](https://github.com/nixpanic))
- Update kubernetes dependencies to v1.29.0 ([#487](https://github.com/kubernetes-csi/csi-driver-host-path/pull/487), [@sunnylovestiramisu](https://github.com/sunnylovestiramisu))

## Dependencies

### Added
- github.com/danwinship/knftables: [v0.0.13](https://github.com/danwinship/knftables/tree/v0.0.13)
- github.com/distribution/reference: [v0.5.0](https://github.com/distribution/reference/tree/v0.5.0)
- github.com/google/s2a-go: [v0.1.7](https://github.com/google/s2a-go/tree/v0.1.7)

### Changed
- cloud.google.com/go/compute: v1.21.0 → v1.23.3
- github.com/cncf/xds/go: [e9ce688 → 0fa0005](https://github.com/cncf/xds/go/compare/e9ce688...0fa0005)
- github.com/container-storage-interface/spec: [v1.8.0 → v1.9.0](https://github.com/container-storage-interface/spec/compare/v1.8.0...v1.9.0)
- github.com/coredns/corefile-migration: [v1.0.20 → v1.0.21](https://github.com/coredns/corefile-migration/compare/v1.0.20...v1.0.21)
- github.com/cyphar/filepath-securejoin: [v0.2.3 → v0.2.4](https://github.com/cyphar/filepath-securejoin/compare/v0.2.3...v0.2.4)
- github.com/emicklei/go-restful/v3: [v3.9.0 → v3.11.0](https://github.com/emicklei/go-restful/v3/compare/v3.9.0...v3.11.0)
- github.com/envoyproxy/go-control-plane: [v0.11.1 → v0.12.0](https://github.com/envoyproxy/go-control-plane/compare/v0.11.1...v0.12.0)
- github.com/envoyproxy/protoc-gen-validate: [v1.0.2 → v1.0.4](https://github.com/envoyproxy/protoc-gen-validate/compare/v1.0.2...v1.0.4)
- github.com/fsnotify/fsnotify: [v1.6.0 → v1.7.0](https://github.com/fsnotify/fsnotify/compare/v1.6.0...v1.7.0)
- github.com/go-logr/logr: [v1.2.4 → v1.4.1](https://github.com/go-logr/logr/compare/v1.2.4...v1.4.1)
- github.com/godbus/dbus/v5: [v5.0.6 → v5.1.0](https://github.com/godbus/dbus/v5/compare/v5.0.6...v5.1.0)
- github.com/golang/glog: [v1.1.2 → v1.2.0](https://github.com/golang/glog/compare/v1.1.2...v1.2.0)
- github.com/google/cadvisor: [v0.47.3 → v0.48.1](https://github.com/google/cadvisor/compare/v0.47.3...v0.48.1)
- github.com/google/cel-go: [v0.16.1 → v0.17.7](https://github.com/google/cel-go/compare/v0.16.1...v0.17.7)
- github.com/google/go-cmp: [v0.5.9 → v0.6.0](https://github.com/google/go-cmp/compare/v0.5.9...v0.6.0)
- github.com/google/uuid: [v1.3.0 → v1.6.0](https://github.com/google/uuid/compare/v1.3.0...v1.6.0)
- github.com/googleapis/gax-go/v2: [v2.7.1 → v2.11.0](https://github.com/googleapis/gax-go/v2/compare/v2.7.1...v2.11.0)
- github.com/gorilla/websocket: [v1.4.2 → v1.5.0](https://github.com/gorilla/websocket/compare/v1.4.2...v1.5.0)
- github.com/grpc-ecosystem/grpc-gateway/v2: [v2.7.0 → v2.16.0](https://github.com/grpc-ecosystem/grpc-gateway/v2/compare/v2.7.0...v2.16.0)
- github.com/ishidawataru/sctp: [7c296d4 → 7ff4192](https://github.com/ishidawataru/sctp/compare/7c296d4...7ff4192)
- github.com/kubernetes-csi/csi-lib-utils: [v0.15.0 → v0.17.0](https://github.com/kubernetes-csi/csi-lib-utils/compare/v0.15.0...v0.17.0)
- github.com/mrunalp/fileutils: [v0.5.0 → v0.5.1](https://github.com/mrunalp/fileutils/compare/v0.5.0...v0.5.1)
- github.com/onsi/ginkgo/v2: [v2.9.4 → v2.13.0](https://github.com/onsi/ginkgo/v2/compare/v2.9.4...v2.13.0)
- github.com/onsi/gomega: [v1.27.6 → v1.29.0](https://github.com/onsi/gomega/compare/v1.27.6...v1.29.0)
- github.com/opencontainers/runc: [v1.1.7 → v1.1.10](https://github.com/opencontainers/runc/compare/v1.1.7...v1.1.10)
- github.com/opencontainers/selinux: [v1.10.0 → v1.11.0](https://github.com/opencontainers/selinux/compare/v1.10.0...v1.11.0)
- github.com/vmware/govmomi: [v0.30.0 → v0.30.6](https://github.com/vmware/govmomi/compare/v0.30.0...v0.30.6)
- go.etcd.io/bbolt: v1.3.7 → v1.3.8
- go.etcd.io/etcd/api/v3: v3.5.9 → v3.5.10
- go.etcd.io/etcd/client/pkg/v3: v3.5.9 → v3.5.10
- go.etcd.io/etcd/client/v2: v2.305.9 → v2.305.10
- go.etcd.io/etcd/client/v3: v3.5.9 → v3.5.10
- go.etcd.io/etcd/pkg/v3: v3.5.9 → v3.5.10
- go.etcd.io/etcd/raft/v3: v3.5.9 → v3.5.10
- go.etcd.io/etcd/server/v3: v3.5.9 → v3.5.10
- go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful: v0.35.0 → v0.42.0
- go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc: v0.41.0 → v0.44.0
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp: v0.35.1 → v0.44.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc: v1.10.0 → v1.19.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace: v1.10.0 → v1.19.0
- go.opentelemetry.io/otel/metric: v0.38.0 → v1.19.0
- go.opentelemetry.io/otel/sdk: v1.10.0 → v1.19.0
- go.opentelemetry.io/otel/trace: v1.15.0 → v1.19.0
- go.opentelemetry.io/otel: v1.15.0 → v1.19.0
- go.opentelemetry.io/proto/otlp: v0.19.0 → v1.0.0
- golang.org/x/crypto: v0.14.0 → v0.19.0
- golang.org/x/mod: v0.10.0 → v0.12.0
- golang.org/x/net: v0.17.0 → v0.21.0
- golang.org/x/oauth2: v0.10.0 → v0.16.0
- golang.org/x/sync: v0.3.0 → v0.6.0
- golang.org/x/sys: v0.13.0 → v0.17.0
- golang.org/x/term: v0.13.0 → v0.17.0
- golang.org/x/text: v0.13.0 → v0.14.0
- golang.org/x/tools: v0.8.0 → v0.12.0
- google.golang.org/api: v0.114.0 → v0.126.0
- google.golang.org/appengine: v1.6.7 → v1.6.8
- google.golang.org/genproto/googleapis/api: 782d3b1 → ef43131
- google.golang.org/genproto/googleapis/rpc: 782d3b1 → ef43131
- google.golang.org/genproto: 782d3b1 → ef43131
- google.golang.org/grpc: v1.58.2 → v1.62.0
- google.golang.org/protobuf: v1.31.0 → v1.32.0
- k8s.io/api: v0.28.2 → v0.29.0
- k8s.io/apiextensions-apiserver: v0.28.2 → v0.29.0
- k8s.io/apimachinery: v0.28.2 → v0.29.0
- k8s.io/apiserver: v0.28.2 → v0.29.0
- k8s.io/cli-runtime: v0.28.2 → v0.29.0
- k8s.io/client-go: v0.28.2 → v0.29.0
- k8s.io/cloud-provider: v0.28.2 → v0.29.0
- k8s.io/cluster-bootstrap: v0.28.2 → v0.29.0
- k8s.io/code-generator: v0.28.2 → v0.29.0
- k8s.io/component-base: v0.28.2 → v0.29.0
- k8s.io/component-helpers: v0.28.2 → v0.29.0
- k8s.io/controller-manager: v0.28.2 → v0.29.0
- k8s.io/cri-api: v0.28.2 → v0.29.0
- k8s.io/csi-translation-lib: v0.28.2 → v0.29.0
- k8s.io/dynamic-resource-allocation: v0.28.2 → v0.29.0
- k8s.io/endpointslice: v0.28.2 → v0.29.0
- k8s.io/gengo: c0856e2 → 9cce18d
- k8s.io/klog/v2: v2.100.1 → v2.120.1
- k8s.io/kms: v0.28.2 → v0.29.0
- k8s.io/kube-aggregator: v0.28.2 → v0.29.0
- k8s.io/kube-controller-manager: v0.28.2 → v0.29.0
- k8s.io/kube-openapi: 2695361 → 2dd684a
- k8s.io/kube-proxy: v0.28.2 → v0.29.0
- k8s.io/kube-scheduler: v0.28.2 → v0.29.0
- k8s.io/kubectl: v0.28.2 → v0.29.0
- k8s.io/kubelet: v0.28.2 → v0.29.0
- k8s.io/kubernetes: v1.28.2 → v1.29.0
- k8s.io/legacy-cloud-providers: v0.28.2 → v0.29.0
- k8s.io/metrics: v0.28.2 → v0.29.0
- k8s.io/mount-utils: v0.28.2 → v0.29.0
- k8s.io/pod-security-admission: v0.28.2 → v0.29.0
- k8s.io/sample-apiserver: v0.28.2 → v0.29.0
- k8s.io/utils: d93618c → 3b25d92
- sigs.k8s.io/apiserver-network-proxy/konnectivity-client: v0.1.2 → v0.28.0
- sigs.k8s.io/structured-merge-diff/v4: v4.2.3 → v4.4.1

### Removed
- github.com/docker/distribution: [v2.8.2+incompatible](https://github.com/docker/distribution/tree/v2.8.2)
- go.opentelemetry.io/otel/exporters/otlp/internal/retry: v1.10.0
