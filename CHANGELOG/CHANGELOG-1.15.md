## Changes by Kind

### Bug or Regression

- The `ignoreFailedRead` parameter in the VolumeSnapshotClass, when set to `true`, results in the `--ignore-failed-read` option to be passed to `tar`. ([#543](https://github.com/kubernetes-csi/csi-driver-host-path/pull/543), [@leonardoce](https://github.com/leonardoce))

### Other (Cleanup or Flake)

- Bump external-resizer to v1.11.2 ([#545](https://github.com/kubernetes-csi/csi-driver-host-path/pull/545), [@AndrewSirenko](https://github.com/AndrewSirenko))

## Dependencies

### Added
- cel.dev/expr: v0.15.0

### Changed
- cloud.google.com/go/compute/metadata: v0.2.3 → v0.3.0
- cloud.google.com/go/compute: v1.25.1 → v1.23.0
- github.com/cespare/xxhash/v2: [v2.2.0 → v2.3.0](https://github.com/cespare/xxhash/compare/v2.2.0...v2.3.0)
- github.com/cncf/xds/go: [8a4994d → 555b57e](https://github.com/cncf/xds/compare/8a4994d...555b57e)
- github.com/container-storage-interface/spec: [v1.9.0 → v1.10.0](https://github.com/container-storage-interface/spec/compare/v1.9.0...v1.10.0)
- github.com/golang/glog: [v1.2.0 → v1.2.1](https://github.com/golang/glog/compare/v1.2.0...v1.2.1)
- golang.org/x/crypto: v0.24.0 → v0.25.0
- golang.org/x/net: v0.26.0 → v0.27.0
- golang.org/x/oauth2: v0.18.0 → v0.20.0
- golang.org/x/sys: v0.21.0 → v0.22.0
- golang.org/x/term: v0.21.0 → v0.22.0
- google.golang.org/appengine: v1.6.8 → v1.6.7
- google.golang.org/genproto/googleapis/api: 94a12d6 → 5315273
- google.golang.org/genproto/googleapis/rpc: 94a12d6 → 5315273
- google.golang.org/grpc: v1.64.0 → v1.65.0
- google.golang.org/protobuf: v1.33.0 → v1.34.2
- k8s.io/klog/v2: v2.130.0 → v2.130.1

### Removed
_Nothing has changed._
