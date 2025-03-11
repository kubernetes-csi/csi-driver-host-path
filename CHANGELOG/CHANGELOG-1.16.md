# Release notes for v1.16.1

# Changelog since v1.16.0

## Changes by Kind

### Feature

- Set a random MaxVolumesPerNode if config.AttachLimit == -1. ([#587](https://github.com/kubernetes-csi/csi-driver-host-path/pull/587), [@torredil](https://github.com/torredil))

### Bug or Regression

- Fixed: Considered only non-zero blocks for GetMetadataAllocated. ([#579](https://github.com/kubernetes-csi/csi-driver-host-path/pull/579), [@iPraveenParihar](https://github.com/iPraveenParihar))

### Uncategorized

- Add better validation of the CapacityRange parameter to CreateVolume. ([#577](https://github.com/kubernetes-csi/csi-driver-host-path/pull/577), [@ebblake](https://github.com/ebblake))

## Dependencies

### Added
_Nothing has changed._

### Changed
_Nothing has changed._

### Removed
_Nothing has changed._

# Release notes for v1.16.0

# Changelog since v1.15.0

## Changes by Kind

### Bug or Regression

- Updated external-provisioner from v2.1.0 to v2.1.1 to include some bug fixes ([#268](https://github.com/kubernetes-csi/csi-driver-host-path/pull/268), [@pohly](https://github.com/pohly))

### Other (Cleanup or Flake)

- Bump hostpath image to v1.15.0 in /deploy ([#559](https://github.com/kubernetes-csi/csi-driver-host-path/pull/559), [@carlory](https://github.com/carlory))
- Update Kubernetes dependencies to 1.32.0 ([#576](https://github.com/kubernetes-csi/csi-driver-host-path/pull/576), [@dfajmon](https://github.com/dfajmon))
- Update deploy directory structure and update sidecar releases, corresponding to Kubernetes 1.32. ([#581](https://github.com/kubernetes-csi/csi-driver-host-path/pull/581), [@leonardoce](https://github.com/leonardoce))

### Uncategorized

- Add sample implementation of CSI SnapshotMetadata service ([#569](https://github.com/kubernetes-csi/csi-driver-host-path/pull/569), [@PrasadG193](https://github.com/PrasadG193))

## Dependencies

### Added
- github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp: [v1.24.2](https://github.com/GoogleCloudPlatform/opentelemetry-operations-go/tree/detectors/gcp/v1.24.2)
- github.com/Microsoft/hnslib: [v0.0.8](https://github.com/Microsoft/hnslib/tree/v0.0.8)
- github.com/antlr4-go/antlr/v4: [v4.13.0](https://github.com/antlr4-go/antlr/tree/v4.13.0)
- github.com/containerd/containerd/api: [v1.7.19](https://github.com/containerd/containerd/tree/api/v1.7.19)
- github.com/containerd/errdefs: [v0.1.0](https://github.com/containerd/errdefs/tree/v0.1.0)
- github.com/containerd/log: [v0.1.0](https://github.com/containerd/log/tree/v0.1.0)
- github.com/fxamacker/cbor/v2: [v2.7.0](https://github.com/fxamacker/cbor/tree/v2.7.0)
- github.com/go-task/slim-sprig/v3: [v3.0.0](https://github.com/go-task/slim-sprig/tree/v3.0.0)
- github.com/klauspost/compress: [v1.17.11](https://github.com/klauspost/compress/tree/v1.17.11)
- github.com/kylelemons/godebug: [v1.1.0](https://github.com/kylelemons/godebug/tree/v1.1.0)
- github.com/moby/sys/userns: [v0.1.0](https://github.com/moby/sys/tree/userns/v0.1.0)
- github.com/planetscale/vtprotobuf: [0393e58](https://github.com/planetscale/vtprotobuf/tree/0393e58)
- github.com/x448/float16: [v0.8.4](https://github.com/x448/float16/tree/v0.8.4)
- go.opentelemetry.io/auto/sdk: v1.1.0
- go.opentelemetry.io/contrib/detectors/gcp: v1.31.0
- go.opentelemetry.io/otel/sdk/metric: v1.31.0
- gopkg.in/evanphx/json-patch.v4: v4.12.0
- k8s.io/cri-client: v0.32.0
- k8s.io/externaljwt: v0.32.0
- k8s.io/gengo/v2: 2b36238
- sigs.k8s.io/knftables: v0.0.17

### Changed
- cel.dev/expr: v0.15.0 → v0.18.0
- cloud.google.com/go/compute/metadata: v0.3.0 → v0.5.2
- github.com/Azure/go-ansiterm: [d185dfc → 306776e](https://github.com/Azure/go-ansiterm/compare/d185dfc...306776e)
- github.com/Microsoft/go-winio: [v0.6.0 → v0.6.2](https://github.com/Microsoft/go-winio/compare/v0.6.0...v0.6.2)
- github.com/alecthomas/kingpin/v2: [v2.3.2 → v2.4.0](https://github.com/alecthomas/kingpin/compare/v2.3.2...v2.4.0)
- github.com/armon/circbuf: [bbbad09 → 5111143](https://github.com/armon/circbuf/compare/bbbad09...5111143)
- github.com/cenkalti/backoff/v4: [v4.2.1 → v4.3.0](https://github.com/cenkalti/backoff/compare/v4.2.1...v4.3.0)
- github.com/cncf/xds/go: [555b57e → b4127c9](https://github.com/cncf/xds/compare/555b57e...b4127c9)
- github.com/container-storage-interface/spec: [v1.10.0 → v1.11.0](https://github.com/container-storage-interface/spec/compare/v1.10.0...v1.11.0)
- github.com/containerd/ttrpc: [v1.2.2 → v1.2.5](https://github.com/containerd/ttrpc/compare/v1.2.2...v1.2.5)
- github.com/coredns/corefile-migration: [v1.0.21 → v1.0.24](https://github.com/coredns/corefile-migration/compare/v1.0.21...v1.0.24)
- github.com/cpuguy83/go-md2man/v2: [v2.0.2 → v2.0.4](https://github.com/cpuguy83/go-md2man/compare/v2.0.2...v2.0.4)
- github.com/cyphar/filepath-securejoin: [v0.2.4 → v0.3.4](https://github.com/cyphar/filepath-securejoin/compare/v0.2.4...v0.3.4)
- github.com/davecgh/go-spew: [v1.1.1 → d8f796a](https://github.com/davecgh/go-spew/compare/v1.1.1...d8f796a)
- github.com/distribution/reference: [v0.5.0 → v0.6.0](https://github.com/distribution/reference/compare/v0.5.0...v0.6.0)
- github.com/emicklei/go-restful/v3: [v3.11.0 → v3.12.1](https://github.com/emicklei/go-restful/compare/v3.11.0...v3.12.1)
- github.com/envoyproxy/go-control-plane: [v0.12.0 → v0.13.1](https://github.com/envoyproxy/go-control-plane/compare/v0.12.0...v0.13.1)
- github.com/envoyproxy/protoc-gen-validate: [v1.0.4 → v1.1.0](https://github.com/envoyproxy/protoc-gen-validate/compare/v1.0.4...v1.1.0)
- github.com/exponent-io/jsonpath: [d6023ce → 1de76d7](https://github.com/exponent-io/jsonpath/compare/d6023ce...1de76d7)
- github.com/felixge/httpsnoop: [v1.0.3 → v1.0.4](https://github.com/felixge/httpsnoop/compare/v1.0.3...v1.0.4)
- github.com/go-logr/logr: [v1.4.1 → v1.4.2](https://github.com/go-logr/logr/compare/v1.4.1...v1.4.2)
- github.com/go-logr/zapr: [v1.2.3 → v1.3.0](https://github.com/go-logr/zapr/compare/v1.2.3...v1.3.0)
- github.com/go-openapi/jsonpointer: [v0.19.6 → v0.21.0](https://github.com/go-openapi/jsonpointer/compare/v0.19.6...v0.21.0)
- github.com/go-openapi/jsonreference: [v0.20.2 → v0.21.0](https://github.com/go-openapi/jsonreference/compare/v0.20.2...v0.21.0)
- github.com/go-openapi/swag: [v0.22.3 → v0.23.0](https://github.com/go-openapi/swag/compare/v0.22.3...v0.23.0)
- github.com/golang/glog: [v1.2.1 → v1.2.2](https://github.com/golang/glog/compare/v1.2.1...v1.2.2)
- github.com/google/cadvisor: [v0.48.1 → v0.51.0](https://github.com/google/cadvisor/compare/v0.48.1...v0.51.0)
- github.com/google/cel-go: [v0.17.7 → v0.22.0](https://github.com/google/cel-go/compare/v0.17.7...v0.22.0)
- github.com/google/gnostic-models: [v0.6.8 → v0.6.9](https://github.com/google/gnostic-models/compare/v0.6.8...v0.6.9)
- github.com/google/pprof: [4bb14d4 → d1b30fe](https://github.com/google/pprof/compare/4bb14d4...d1b30fe)
- github.com/gregjones/httpcache: [9cad4c3 → 901d907](https://github.com/gregjones/httpcache/compare/9cad4c3...901d907)
- github.com/grpc-ecosystem/grpc-gateway/v2: [v2.16.0 → v2.20.0](https://github.com/grpc-ecosystem/grpc-gateway/compare/v2.16.0...v2.20.0)
- github.com/jonboulle/clockwork: [v0.2.2 → v0.4.0](https://github.com/jonboulle/clockwork/compare/v0.2.2...v0.4.0)
- github.com/kubernetes-csi/csi-lib-utils: [v0.17.0 → v0.20.0](https://github.com/kubernetes-csi/csi-lib-utils/compare/v0.17.0...v0.20.0)
- github.com/mailru/easyjson: [v0.7.7 → v0.9.0](https://github.com/mailru/easyjson/compare/v0.7.7...v0.9.0)
- github.com/moby/spdystream: [v0.2.0 → v0.5.0](https://github.com/moby/spdystream/compare/v0.2.0...v0.5.0)
- github.com/moby/sys/mountinfo: [v0.6.2 → v0.7.2](https://github.com/moby/sys/compare/mountinfo/v0.6.2...mountinfo/v0.7.2)
- github.com/moby/term: [1aeaba8 → v0.5.0](https://github.com/moby/term/compare/1aeaba8...v0.5.0)
- github.com/mohae/deepcopy: [491d360 → c48cc78](https://github.com/mohae/deepcopy/compare/491d360...c48cc78)
- github.com/onsi/ginkgo/v2: [v2.13.0 → v2.21.0](https://github.com/onsi/ginkgo/compare/v2.13.0...v2.21.0)
- github.com/onsi/gomega: [v1.29.0 → v1.35.1](https://github.com/onsi/gomega/compare/v1.29.0...v1.35.1)
- github.com/opencontainers/runc: [v1.1.10 → v1.2.1](https://github.com/opencontainers/runc/compare/v1.1.10...v1.2.1)
- github.com/opencontainers/runtime-spec: [494a5a6 → v1.2.0](https://github.com/opencontainers/runtime-spec/compare/494a5a6...v1.2.0)
- github.com/opencontainers/selinux: [v1.11.0 → v1.11.1](https://github.com/opencontainers/selinux/compare/v1.11.0...v1.11.1)
- github.com/pmezard/go-difflib: [v1.0.0 → 5d4384e](https://github.com/pmezard/go-difflib/compare/v1.0.0...5d4384e)
- github.com/prometheus/client_golang: [v1.16.0 → v1.20.5](https://github.com/prometheus/client_golang/compare/v1.16.0...v1.20.5)
- github.com/prometheus/client_model: [v0.4.0 → v0.6.1](https://github.com/prometheus/client_model/compare/v0.4.0...v0.6.1)
- github.com/prometheus/common: [v0.44.0 → v0.61.0](https://github.com/prometheus/common/compare/v0.44.0...v0.61.0)
- github.com/prometheus/procfs: [v0.10.1 → v0.15.1](https://github.com/prometheus/procfs/compare/v0.10.1...v0.15.1)
- github.com/rogpeppe/go-internal: [v1.10.0 → v1.12.0](https://github.com/rogpeppe/go-internal/compare/v1.10.0...v1.12.0)
- github.com/sirupsen/logrus: [v1.9.0 → v1.9.3](https://github.com/sirupsen/logrus/compare/v1.9.0...v1.9.3)
- github.com/spf13/cobra: [v1.7.0 → v1.8.1](https://github.com/spf13/cobra/compare/v1.7.0...v1.8.1)
- github.com/stoewer/go-strcase: [v1.2.0 → v1.3.0](https://github.com/stoewer/go-strcase/compare/v1.2.0...v1.3.0)
- github.com/stretchr/testify: [v1.9.0 → v1.10.0](https://github.com/stretchr/testify/compare/v1.9.0...v1.10.0)
- github.com/vishvananda/netlink: [v1.1.0 → b1ce50c](https://github.com/vishvananda/netlink/compare/v1.1.0...b1ce50c)
- github.com/xiang90/probing: [43a291a → a49e3df](https://github.com/xiang90/probing/compare/43a291a...a49e3df)
- go.etcd.io/bbolt: v1.3.8 → v1.3.11
- go.etcd.io/etcd/api/v3: v3.5.10 → v3.5.16
- go.etcd.io/etcd/client/pkg/v3: v3.5.10 → v3.5.16
- go.etcd.io/etcd/client/v2: v2.305.10 → v2.305.16
- go.etcd.io/etcd/client/v3: v3.5.10 → v3.5.16
- go.etcd.io/etcd/pkg/v3: v3.5.10 → v3.5.16
- go.etcd.io/etcd/raft/v3: v3.5.10 → v3.5.16
- go.etcd.io/etcd/server/v3: v3.5.10 → v3.5.16
- go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc: v0.44.0 → v0.58.0
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp: v0.44.0 → v0.53.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc: v1.19.0 → v1.27.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace: v1.19.0 → v1.28.0
- go.opentelemetry.io/otel/metric: v1.19.0 → v1.33.0
- go.opentelemetry.io/otel/sdk: v1.19.0 → v1.31.0
- go.opentelemetry.io/otel/trace: v1.19.0 → v1.33.0
- go.opentelemetry.io/otel: v1.19.0 → v1.33.0
- go.opentelemetry.io/proto/otlp: v1.0.0 → v1.3.1
- go.uber.org/goleak: v1.2.1 → v1.3.0
- go.uber.org/zap: v1.19.0 → v1.27.0
- golang.org/x/crypto: v0.25.0 → v0.30.0
- golang.org/x/exp: a9213ee → 8a7402a
- golang.org/x/mod: v0.17.0 → v0.21.0
- golang.org/x/net: v0.27.0 → v0.32.0
- golang.org/x/oauth2: v0.20.0 → v0.24.0
- golang.org/x/sync: v0.7.0 → v0.10.0
- golang.org/x/sys: v0.22.0 → v0.28.0
- golang.org/x/term: v0.22.0 → v0.27.0
- golang.org/x/text: v0.16.0 → v0.21.0
- golang.org/x/time: v0.3.0 → v0.8.0
- golang.org/x/tools: e35e4cc → v0.26.0
- google.golang.org/genproto/googleapis/api: 5315273 → 796eee8
- google.golang.org/genproto/googleapis/rpc: 5315273 → 9240e9c
- google.golang.org/genproto: f966b18 → ef43131
- google.golang.org/grpc: v1.65.0 → v1.69.0
- google.golang.org/protobuf: v1.34.2 → v1.36.0
- k8s.io/api: v0.29.0 → v0.32.0
- k8s.io/apiextensions-apiserver: v0.29.0 → v0.32.0
- k8s.io/apimachinery: v0.29.0 → v0.32.0
- k8s.io/apiserver: v0.29.0 → v0.32.0
- k8s.io/cli-runtime: v0.29.0 → v0.32.0
- k8s.io/client-go: v0.29.0 → v0.32.0
- k8s.io/cloud-provider: v0.29.0 → v0.32.0
- k8s.io/cluster-bootstrap: v0.29.0 → v0.32.0
- k8s.io/code-generator: v0.29.0 → v0.32.0
- k8s.io/component-base: v0.29.0 → v0.32.0
- k8s.io/component-helpers: v0.29.0 → v0.32.0
- k8s.io/controller-manager: v0.29.0 → v0.32.0
- k8s.io/cri-api: v0.29.0 → v0.32.0
- k8s.io/csi-translation-lib: v0.29.0 → v0.32.0
- k8s.io/dynamic-resource-allocation: v0.29.0 → v0.32.0
- k8s.io/endpointslice: v0.29.0 → v0.32.0
- k8s.io/kms: v0.29.0 → v0.32.0
- k8s.io/kube-aggregator: v0.29.0 → v0.32.0
- k8s.io/kube-controller-manager: v0.29.0 → v0.32.0
- k8s.io/kube-openapi: 2dd684a → 2c72e55
- k8s.io/kube-proxy: v0.29.0 → v0.32.0
- k8s.io/kube-scheduler: v0.29.0 → v0.32.0
- k8s.io/kubectl: v0.29.0 → v0.32.0
- k8s.io/kubelet: v0.29.0 → v0.32.0
- k8s.io/kubernetes: v1.29.2 → v1.32.0
- k8s.io/metrics: v0.29.0 → v0.32.0
- k8s.io/mount-utils: v0.29.0 → v0.32.0
- k8s.io/pod-security-admission: v0.29.0 → v0.32.0
- k8s.io/sample-apiserver: v0.29.0 → v0.32.0
- k8s.io/system-validators: v1.8.0 → v1.9.1
- k8s.io/utils: 3b25d92 → 24370be
- sigs.k8s.io/apiserver-network-proxy/konnectivity-client: v0.28.0 → v0.31.0
- sigs.k8s.io/json: bc3834c → cfa47c3
- sigs.k8s.io/kustomize/api: 6ce0bf3 → v0.18.0
- sigs.k8s.io/kustomize/kustomize/v5: 6ce0bf3 → v5.5.0
- sigs.k8s.io/kustomize/kyaml: 6ce0bf3 → v0.18.1
- sigs.k8s.io/structured-merge-diff/v4: v4.4.1 → v4.5.0
- sigs.k8s.io/yaml: v1.3.0 → v1.4.0

### Removed
- cloud.google.com/go/compute: v1.23.0
- github.com/Azure/azure-sdk-for-go: [v68.0.0+incompatible](https://github.com/Azure/azure-sdk-for-go/tree/v68.0.0)
- github.com/Azure/go-autorest/autorest/adal: [v0.9.23](https://github.com/Azure/go-autorest/tree/autorest/adal/v0.9.23)
- github.com/Azure/go-autorest/autorest/date: [v0.3.0](https://github.com/Azure/go-autorest/tree/autorest/date/v0.3.0)
- github.com/Azure/go-autorest/autorest/mocks: [v0.4.2](https://github.com/Azure/go-autorest/tree/autorest/mocks/v0.4.2)
- github.com/Azure/go-autorest/autorest/to: [v0.4.0](https://github.com/Azure/go-autorest/tree/autorest/to/v0.4.0)
- github.com/Azure/go-autorest/autorest/validation: [v0.3.1](https://github.com/Azure/go-autorest/tree/autorest/validation/v0.3.1)
- github.com/Azure/go-autorest/autorest: [v0.11.29](https://github.com/Azure/go-autorest/tree/autorest/v0.11.29)
- github.com/Azure/go-autorest/logger: [v0.2.1](https://github.com/Azure/go-autorest/tree/logger/v0.2.1)
- github.com/Azure/go-autorest/tracing: [v0.6.0](https://github.com/Azure/go-autorest/tree/tracing/v0.6.0)
- github.com/Azure/go-autorest: [v14.2.0+incompatible](https://github.com/Azure/go-autorest/tree/v14.2.0)
- github.com/GoogleCloudPlatform/k8s-cloud-provider: [f118173](https://github.com/GoogleCloudPlatform/k8s-cloud-provider/tree/f118173)
- github.com/Microsoft/hcsshim: [v0.8.25](https://github.com/Microsoft/hcsshim/tree/v0.8.25)
- github.com/antlr/antlr4/runtime/Go/antlr/v4: [8188dc5](https://github.com/antlr/antlr4/tree/runtime/Go/antlr/v4/8188dc5)
- github.com/checkpoint-restore/go-criu/v5: [v5.3.0](https://github.com/checkpoint-restore/go-criu/tree/v5.3.0)
- github.com/cilium/ebpf: [v0.9.1](https://github.com/cilium/ebpf/tree/v0.9.1)
- github.com/containerd/cgroups: [v1.1.0](https://github.com/containerd/cgroups/tree/v1.1.0)
- github.com/containerd/console: [v1.0.3](https://github.com/containerd/console/tree/v1.0.3)
- github.com/creack/pty: [v1.1.9](https://github.com/creack/pty/tree/v1.1.9)
- github.com/danwinship/knftables: [v0.0.13](https://github.com/danwinship/knftables/tree/v0.0.13)
- github.com/daviddengcn/go-colortext: [v1.0.0](https://github.com/daviddengcn/go-colortext/tree/v1.0.0)
- github.com/evanphx/json-patch: [v4.12.0+incompatible](https://github.com/evanphx/json-patch/tree/v4.12.0)
- github.com/fvbommel/sortorder: [v1.1.0](https://github.com/fvbommel/sortorder/tree/v1.1.0)
- github.com/go-task/slim-sprig: [52ccab3](https://github.com/go-task/slim-sprig/tree/52ccab3)
- github.com/gofrs/uuid: [v4.4.0+incompatible](https://github.com/gofrs/uuid/tree/v4.4.0)
- github.com/golang/groupcache: [41bb18b](https://github.com/golang/groupcache/tree/41bb18b)
- github.com/golang/mock: [v1.6.0](https://github.com/golang/mock/tree/v1.6.0)
- github.com/google/s2a-go: [v0.1.7](https://github.com/google/s2a-go/tree/v0.1.7)
- github.com/googleapis/enterprise-certificate-proxy: [v0.2.3](https://github.com/googleapis/enterprise-certificate-proxy/tree/v0.2.3)
- github.com/googleapis/gax-go/v2: [v2.11.0](https://github.com/googleapis/gax-go/tree/v2.11.0)
- github.com/imdario/mergo: [v0.3.6](https://github.com/imdario/mergo/tree/v0.3.6)
- github.com/matttproud/golang_protobuf_extensions: [v1.0.4](https://github.com/matttproud/golang_protobuf_extensions/tree/v1.0.4)
- github.com/rubiojr/go-vhd: [02e2102](https://github.com/rubiojr/go-vhd/tree/02e2102)
- github.com/seccomp/libseccomp-golang: [v0.10.0](https://github.com/seccomp/libseccomp-golang/tree/v0.10.0)
- github.com/syndtr/gocapability: [42c35b4](https://github.com/syndtr/gocapability/tree/42c35b4)
- github.com/vmware/govmomi: [v0.30.6](https://github.com/vmware/govmomi/tree/v0.30.6)
- go.opencensus.io: v0.24.0
- go.starlark.net: a134d8f
- go.uber.org/atomic: v1.10.0
- google.golang.org/api: v0.126.0
- google.golang.org/appengine: v1.6.7
- gopkg.in/gcfg.v1: v1.2.3
- gopkg.in/warnings.v0: v0.1.2
- k8s.io/gengo: 9cce18d
- k8s.io/legacy-cloud-providers: v0.29.0
