# Release notes for v1.19.0

# Changelog since v1.18.0

## Changes by Kind

### Feature

- Implement volume health support (`ControllerGetVolume` and volume condition reporting) in the CSI driver, and add a `healthctl` tool to simulate volume health states. Bump the CSI spec to v1.13.0. ([#678](https://github.com/kubernetes-csi/csi-driver-host-path/pull/678), [@gnufied](https://github.com/gnufied))

### Uncategorized

- Update deployment structure for Kubernetes 1.37 and update sidecar images. ([#700](https://github.com/kubernetes-csi/csi-driver-host-path/pull/700), [@xing-yang](https://github.com/xing-yang))


## Dependencies

### Added
- go.yaml.in/yaml/v3: v3.0.5

### Changed
- github.com/container-storage-interface/spec: [v1.12.0 → v1.13.0](https://github.com/container-storage-interface/spec/compare/v1.12.0...v1.13.0)
- github.com/fxamacker/cbor/v2: [v2.9.0 → v2.9.3](https://github.com/fxamacker/cbor/v2/compare/v2.9.0...v2.9.3)
- github.com/go-logr/logr: [v1.4.3 → v1.4.4](https://github.com/go-logr/logr/compare/v1.4.3...v1.4.4)
- github.com/kubernetes-csi/csi-lib-utils: [v0.24.0 → v0.25.0](https://github.com/kubernetes-csi/csi-lib-utils/compare/v0.24.0...v0.25.0)
- github.com/prometheus/client_golang: [v1.23.2 → v1.24.1](https://github.com/prometheus/client_golang/compare/v1.23.2...v1.24.1)
- github.com/prometheus/common: [v0.67.5 → v0.70.1](https://github.com/prometheus/common/compare/v0.67.5...v0.70.1)
- github.com/prometheus/procfs: [v0.19.2 → v0.21.1](https://github.com/prometheus/procfs/compare/v0.19.2...v0.21.1)
- github.com/stretchr/testify: [v1.11.1 → v1.12.1](https://github.com/stretchr/testify/compare/v1.11.1...v1.12.1)
- go.opentelemetry.io/otel/trace: v1.43.0 → v1.46.0
- go.opentelemetry.io/otel: v1.43.0 → v1.46.0
- go.yaml.in/yaml/v2: v2.4.3 → v2.4.4
- golang.org/x/net: v0.49.0 → v0.61.0
- golang.org/x/sys: v0.40.0 → v0.49.0
- golang.org/x/text: v0.33.0 → v0.43.0
- google.golang.org/genproto/googleapis/rpc: 8636f87 → da73d73
- google.golang.org/grpc: v1.79.3 → v1.84.0
- google.golang.org/protobuf: v1.36.12-0.20260120151049-f2248ac996af → v1.36.12
- k8s.io/apimachinery: v0.36.1 → v0.37.0
- k8s.io/component-base: v0.36.1 → v0.37.0
- k8s.io/kubernetes: v1.36.1 → v1.36.3
- k8s.io/utils: b8788ab → cf1189d

### Removed
- github.com/davecgh/go-spew: [d8f796a](https://github.com/davecgh/go-spew/tree/d8f796a)
- github.com/pmezard/go-difflib: [5d4384e](https://github.com/pmezard/go-difflib/tree/5d4384e)
- gopkg.in/yaml.v3: v3.0.1
