# Release notes for v1.11.0

# Changelog since v1.10.0

## Changes by Kind

### Feature

- Enforces volumes the with SINGLE_NODE_SINGLE_WRITER access mode can only be mounted at one target path at a time ([#381](https://github.com/kubernetes-csi/csi-driver-host-path/pull/381), [@chrishenzie](https://github.com/chrishenzie)) [SIG Storage]
- Provide option to enable volume mode conversion flag to HostPath driver deployment script ([#379](https://github.com/kubernetes-csi/csi-driver-host-path/pull/379), [@RaunakShah](https://github.com/RaunakShah))

## Dependencies

### Added
- github.com/cenkalti/backoff/v4: [v4.1.3](https://github.com/cenkalti/backoff/v4/tree/v4.1.3)
- github.com/go-logr/stdr: [v1.2.2](https://github.com/go-logr/stdr/tree/v1.2.2)
- github.com/grpc-ecosystem/grpc-gateway/v2: [v2.7.0](https://github.com/grpc-ecosystem/grpc-gateway/v2/tree/v2.7.0)
- go.opentelemetry.io/otel/exporters/otlp/internal/retry: v1.10.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc: v1.10.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace: v1.10.0
- k8s.io/dynamic-resource-allocation: v0.26.0
- k8s.io/kms: v0.26.0

### Changed
- github.com/antlr/antlr4/runtime/Go/antlr: [f25a4f6 → v1.4.10](https://github.com/antlr/antlr4/runtime/Go/antlr/compare/f25a4f6...v1.4.10)
- github.com/aws/aws-sdk-go: [v1.38.49 → v1.44.116](https://github.com/aws/aws-sdk-go/compare/v1.38.49...v1.44.116)
- github.com/container-storage-interface/spec: [v1.6.0 → v1.7.0](https://github.com/container-storage-interface/spec/compare/v1.6.0...v1.7.0)
- github.com/containerd/ttrpc: [v1.0.2 → v1.1.0](https://github.com/containerd/ttrpc/compare/v1.0.2...v1.1.0)
- github.com/cpuguy83/go-md2man/v2: [v2.0.1 → v2.0.2](https://github.com/cpuguy83/go-md2man/v2/compare/v2.0.1...v2.0.2)
- github.com/docker/go-units: [v0.4.0 → v0.5.0](https://github.com/docker/go-units/compare/v0.4.0...v0.5.0)
- github.com/emicklei/go-restful/v3: [v3.8.0 → v3.9.0](https://github.com/emicklei/go-restful/v3/compare/v3.8.0...v3.9.0)
- github.com/felixge/httpsnoop: [v1.0.1 → v1.0.3](https://github.com/felixge/httpsnoop/compare/v1.0.1...v1.0.3)
- github.com/fsnotify/fsnotify: [v1.4.9 → v1.6.0](https://github.com/fsnotify/fsnotify/compare/v1.4.9...v1.6.0)
- github.com/go-openapi/jsonreference: [v0.19.5 → v0.20.0](https://github.com/go-openapi/jsonreference/compare/v0.19.5...v0.20.0)
- github.com/google/cadvisor: [v0.45.0 → v0.46.0](https://github.com/google/cadvisor/compare/v0.45.0...v0.46.0)
- github.com/google/go-cmp: [v0.5.8 → v0.5.9](https://github.com/google/go-cmp/compare/v0.5.8...v0.5.9)
- github.com/google/martian/v3: [v3.2.1 → v3.0.0](https://github.com/google/martian/v3/compare/v3.2.1...v3.0.0)
- github.com/ianlancetaylor/demangle: [28f6c0f → 5e5cf60](https://github.com/ianlancetaylor/demangle/compare/28f6c0f...5e5cf60)
- github.com/inconshreveable/mousetrap: [v1.0.0 → v1.0.1](https://github.com/inconshreveable/mousetrap/compare/v1.0.0...v1.0.1)
- github.com/karrick/godirwalk: [v1.16.1 → v1.17.0](https://github.com/karrick/godirwalk/compare/v1.16.1...v1.17.0)
- github.com/kr/pretty: [v0.2.0 → v0.1.0](https://github.com/kr/pretty/compare/v0.2.0...v0.1.0)
- github.com/kubernetes-csi/csi-lib-utils: [v0.11.0 → v0.12.0](https://github.com/kubernetes-csi/csi-lib-utils/compare/v0.11.0...v0.12.0)
- github.com/matttproud/golang_protobuf_extensions: [c182aff → v1.0.2](https://github.com/matttproud/golang_protobuf_extensions/compare/c182aff...v1.0.2)
- github.com/moby/term: [3f7ff69 → 39b0c02](https://github.com/moby/term/compare/3f7ff69...39b0c02)
- github.com/onsi/ginkgo/v2: [v2.1.6 → v2.4.0](https://github.com/onsi/ginkgo/v2/compare/v2.1.6...v2.4.0)
- github.com/onsi/gomega: [v1.20.1 → v1.23.0](https://github.com/onsi/gomega/compare/v1.20.1...v1.23.0)
- github.com/opencontainers/runc: [v1.1.3 → v1.1.4](https://github.com/opencontainers/runc/compare/v1.1.3...v1.1.4)
- github.com/prometheus/client_golang: [v1.13.0 → v1.14.0](https://github.com/prometheus/client_golang/compare/v1.13.0...v1.14.0)
- github.com/prometheus/client_model: [v0.2.0 → v0.3.0](https://github.com/prometheus/client_model/compare/v0.2.0...v0.3.0)
- github.com/spf13/cobra: [v1.4.0 → v1.6.0](https://github.com/spf13/cobra/compare/v1.4.0...v1.6.0)
- github.com/stretchr/objx: [v0.2.0 → v0.4.0](https://github.com/stretchr/objx/compare/v0.2.0...v0.4.0)
- github.com/stretchr/testify: [v1.7.0 → v1.8.0](https://github.com/stretchr/testify/compare/v1.7.0...v1.8.0)
- github.com/yuin/goldmark: [v1.4.13 → v1.2.1](https://github.com/yuin/goldmark/compare/v1.4.13...v1.2.1)
- go.etcd.io/etcd/api/v3: v3.5.4 → v3.5.5
- go.etcd.io/etcd/client/pkg/v3: v3.5.4 → v3.5.5
- go.etcd.io/etcd/client/v2: v2.305.4 → v2.305.5
- go.etcd.io/etcd/client/v3: v3.5.4 → v3.5.5
- go.etcd.io/etcd/pkg/v3: v3.5.4 → v3.5.5
- go.etcd.io/etcd/raft/v3: v3.5.4 → v3.5.5
- go.etcd.io/etcd/server/v3: v3.5.4 → v3.5.5
- go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful: v0.20.0 → v0.35.0
- go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc: v0.20.0 → v0.35.0
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp: v0.20.0 → v0.35.0
- go.opentelemetry.io/otel/metric: v0.20.0 → v0.31.0
- go.opentelemetry.io/otel/sdk: v0.20.0 → v1.10.0
- go.opentelemetry.io/otel/trace: v0.20.0 → v1.10.0
- go.opentelemetry.io/otel: v0.20.0 → v1.10.0
- go.opentelemetry.io/proto/otlp: v0.7.0 → v0.19.0
- go.uber.org/goleak: v1.1.10 → v1.2.0
- golang.org/x/crypto: 3147a52 → v0.1.0
- golang.org/x/exp: 85be41e → 6cc2880
- golang.org/x/lint: 6edffad → 738671d
- golang.org/x/mod: 86c51ed → v0.6.0
- golang.org/x/net: d300de1 → v0.4.0
- golang.org/x/sys: fb04ddd → v0.3.0
- golang.org/x/term: 03fcf44 → v0.3.0
- golang.org/x/text: v0.3.7 → v0.5.0
- golang.org/x/tools: v0.1.12 → v0.2.0
- k8s.io/api: v0.25.2 → v0.26.0
- k8s.io/apiextensions-apiserver: v0.25.2 → v0.26.0
- k8s.io/apimachinery: v0.25.2 → v0.26.0
- k8s.io/apiserver: v0.25.2 → v0.26.0
- k8s.io/cli-runtime: v0.25.2 → v0.26.0
- k8s.io/client-go: v0.25.2 → v0.26.0
- k8s.io/cloud-provider: v0.25.2 → v0.26.0
- k8s.io/cluster-bootstrap: v0.25.2 → v0.26.0
- k8s.io/code-generator: v0.25.2 → v0.26.0
- k8s.io/component-base: v0.25.2 → v0.26.0
- k8s.io/component-helpers: v0.25.2 → v0.26.0
- k8s.io/controller-manager: v0.25.2 → v0.26.0
- k8s.io/cri-api: v0.25.2 → v0.26.0
- k8s.io/csi-translation-lib: v0.25.2 → v0.26.0
- k8s.io/gengo: c02415c → c0856e2
- k8s.io/kube-aggregator: v0.25.2 → v0.26.0
- k8s.io/kube-controller-manager: v0.25.2 → v0.26.0
- k8s.io/kube-openapi: 67bda5d → 172d655
- k8s.io/kube-proxy: v0.25.2 → v0.26.0
- k8s.io/kube-scheduler: v0.25.2 → v0.26.0
- k8s.io/kubectl: v0.25.2 → v0.26.0
- k8s.io/kubelet: v0.25.2 → v0.26.0
- k8s.io/kubernetes: v1.25.2 → v1.26.0
- k8s.io/legacy-cloud-providers: v0.25.2 → v0.26.0
- k8s.io/metrics: v0.25.2 → v0.26.0
- k8s.io/mount-utils: v0.25.2-rc.0 → v0.26.0
- k8s.io/pod-security-admission: v0.25.2 → v0.26.0
- k8s.io/sample-apiserver: v0.25.2 → v0.26.0
- k8s.io/system-validators: v1.7.0 → v1.8.0
- k8s.io/utils: 7796b5f → 1a15be2
- sigs.k8s.io/apiserver-network-proxy/konnectivity-client: v0.0.32 → v0.0.33
- sigs.k8s.io/yaml: v1.2.0 → v1.3.0

### Removed
- github.com/OneOfOne/xxhash: [v1.2.2](https://github.com/OneOfOne/xxhash/tree/v1.2.2)
- github.com/PuerkitoBio/purell: [v1.1.1](https://github.com/PuerkitoBio/purell/tree/v1.1.1)
- github.com/PuerkitoBio/urlesc: [de5bf2a](https://github.com/PuerkitoBio/urlesc/tree/de5bf2a)
- github.com/antihax/optional: [v1.0.0](https://github.com/antihax/optional/tree/v1.0.0)
- github.com/auth0/go-jwt-middleware: [v1.0.1](https://github.com/auth0/go-jwt-middleware/tree/v1.0.1)
- github.com/benbjohnson/clock: [v1.1.0](https://github.com/benbjohnson/clock/tree/v1.1.0)
- github.com/boltdb/bolt: [v1.3.1](https://github.com/boltdb/bolt/tree/v1.3.1)
- github.com/cespare/xxhash: [v1.1.0](https://github.com/cespare/xxhash/tree/v1.1.0)
- github.com/creack/pty: [v1.1.11](https://github.com/creack/pty/tree/v1.1.11)
- github.com/docopt/docopt-go: [ee0de3b](https://github.com/docopt/docopt-go/tree/ee0de3b)
- github.com/getkin/kin-openapi: [v0.76.0](https://github.com/getkin/kin-openapi/tree/v0.76.0)
- github.com/ghodss/yaml: [v1.0.0](https://github.com/ghodss/yaml/tree/v1.0.0)
- github.com/go-ozzo/ozzo-validation: [v3.5.0+incompatible](https://github.com/go-ozzo/ozzo-validation/tree/v3.5.0)
- github.com/golang/snappy: [v0.0.3](https://github.com/golang/snappy/tree/v0.0.3)
- github.com/gophercloud/gophercloud: [v0.1.0](https://github.com/gophercloud/gophercloud/tree/v0.1.0)
- github.com/gorilla/mux: [v1.8.0](https://github.com/gorilla/mux/tree/v1.8.0)
- github.com/heketi/heketi: [v10.3.0+incompatible](https://github.com/heketi/heketi/tree/v10.3.0)
- github.com/heketi/tests: [f3775cb](https://github.com/heketi/tests/tree/f3775cb)
- github.com/hpcloud/tail: [v1.0.0](https://github.com/hpcloud/tail/tree/v1.0.0)
- github.com/lpabon/godbc: [v0.1.1](https://github.com/lpabon/godbc/tree/v0.1.1)
- github.com/mvdan/xurls: [v1.1.0](https://github.com/mvdan/xurls/tree/v1.1.0)
- github.com/nxadm/tail: [v1.4.8](https://github.com/nxadm/tail/tree/v1.4.8)
- github.com/onsi/ginkgo: [v1.16.4](https://github.com/onsi/ginkgo/tree/v1.16.4)
- github.com/rogpeppe/fastuuid: [v1.2.0](https://github.com/rogpeppe/fastuuid/tree/v1.2.0)
- github.com/russross/blackfriday: [v1.5.2](https://github.com/russross/blackfriday/tree/v1.5.2)
- github.com/spaolacci/murmur3: [f09979e](https://github.com/spaolacci/murmur3/tree/f09979e)
- github.com/spf13/afero: [v1.9.2](https://github.com/spf13/afero/tree/v1.9.2)
- github.com/urfave/negroni: [v1.0.0](https://github.com/urfave/negroni/tree/v1.0.0)
- go.opentelemetry.io/contrib: v0.20.0
- go.opentelemetry.io/otel/exporters/otlp: v0.20.0
- go.opentelemetry.io/otel/oteltest: v0.20.0
- go.opentelemetry.io/otel/sdk/export/metric: v0.20.0
- go.opentelemetry.io/otel/sdk/metric: v0.20.0
- gonum.org/v1/gonum: v0.6.2
- gonum.org/v1/netlib: 7672324
- google.golang.org/grpc/cmd/protoc-gen-go-grpc: v1.1.0
- gopkg.in/fsnotify.v1: v1.4.7
- gopkg.in/tomb.v1: dd63297
