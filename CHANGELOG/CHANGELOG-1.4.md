# Release notes for v1.4.0

[Documentation](https://kubernetes-csi.github.io)

# Changelog since v1.3.0

## Changes by Kind

### Feature

- Added cloning support to raw block volumes. ([#157](https://github.com/kubernetes-csi/csi-driver-host-path/pull/157), [@jsafrane](https://github.com/jsafrane))
- Added support for volume expansion. ([#87](https://github.com/kubernetes-csi/csi-driver-host-path/pull/87), [@bertinatto](https://github.com/bertinatto))
- Added volume expansion support. ([#90](https://github.com/kubernetes-csi/csi-driver-host-path/pull/90), [@bertinatto](https://github.com/bertinatto))
- Adds ability to discover and import existing on-disk tarballs as snapshots. ([#161](https://github.com/kubernetes-csi/csi-driver-host-path/pull/161), [@ashish-amarnath](https://github.com/ashish-amarnath))
- Enable topology support of CSI hostpath driver. ([#88](https://github.com/kubernetes-csi/csi-driver-host-path/pull/88), [@mucahitkurt](https://github.com/mucahitkurt))

### Bug or Regression

- Fixed snapshots of block volumes. ([#162](https://github.com/kubernetes-csi/csi-driver-host-path/pull/162), [@jsafrane](https://github.com/jsafrane))
- Removed useless preStop hooks from the driver pod. ([#176](https://github.com/kubernetes-csi/csi-driver-host-path/pull/176), [@jsafrane](https://github.com/jsafrane))
- Use absolute path name of snapshot file to import ([#166](https://github.com/kubernetes-csi/csi-driver-host-path/pull/166), [@ashish-amarnath](https://github.com/ashish-amarnath))

### Other (Cleanup or Flake)

- Add 1.18 deployment specs and update sidecar and driver versions. ([#169](https://github.com/kubernetes-csi/csi-driver-host-path/pull/169), [@msau42](https://github.com/msau42))

### Uncategorized

- Publishing of images on k8s.gcr.io ([#180](https://github.com/kubernetes-csi/csi-driver-host-path/pull/180), [@pohly](https://github.com/pohly))

## Dependencies

### Added
- github.com/BurntSushi/toml: [v0.3.1](https://github.com/BurntSushi/toml/tree/v0.3.1)
- github.com/census-instrumentation/opencensus-proto: [v0.2.1](https://github.com/census-instrumentation/opencensus-proto/tree/v0.2.1)
- github.com/envoyproxy/go-control-plane: [5f8ba28](https://github.com/envoyproxy/go-control-plane/tree/5f8ba28)
- github.com/envoyproxy/protoc-gen-validate: [v0.1.0](https://github.com/envoyproxy/protoc-gen-validate/tree/v0.1.0)
- github.com/google/go-cmp: [v0.2.0](https://github.com/google/go-cmp/tree/v0.2.0)
- github.com/prometheus/client_model: [14fe0d1](https://github.com/prometheus/client_model/tree/14fe0d1)
- golang.org/x/crypto: c2843e0
- golang.org/x/exp: 509febe

### Changed
- github.com/container-storage-interface/spec: [v1.1.0 → v1.2.0](https://github.com/container-storage-interface/spec/compare/v1.1.0...v1.2.0)
- github.com/golang/protobuf: [v1.2.0 → v1.3.2](https://github.com/golang/protobuf/compare/v1.2.0...v1.3.2)
- golang.org/x/lint: 06c8688 → d0100b6
- golang.org/x/net: 88d92db → d888771
- golang.org/x/sync: 1d60e46 → 1122301
- golang.org/x/sys: 66b7b13 → d0b11bd
- golang.org/x/tools: 6cd1fce → 2c0ae70
- google.golang.org/appengine: v1.1.0 → v1.4.0
- google.golang.org/genproto: b5d4398 → 24fa4b2
- google.golang.org/grpc: v1.16.0 → v1.26.0
- honnef.co/go/tools: 8849700 → ea95bdf

### Removed
- github.com/golang/lint: [06c8688](https://github.com/golang/lint/tree/06c8688)
- github.com/kisielk/gotool: [v1.0.0](https://github.com/kisielk/gotool/tree/v1.0.0)
