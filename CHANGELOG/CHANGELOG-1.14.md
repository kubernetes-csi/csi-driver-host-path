# Release notes for v1.14.0

# Changelog since v1.13.0

## Changes by Kind

### Bug or Regression

- CSI ephemeral inline volumes failed to get created with an error saying `MountVolume.SetUp failed for volume "ephemeral-volume" : rpc error: code = OutOfRange desc = Requested capacity 1099511627776 exceeds maximum allowed 1099511627776` ([#254](https://github.com/kubernetes-csi/csi-driver-host-path/pull/254), [@pohly](https://github.com/pohly))
- Fix broken symbolic links in the deploy dir ([#510](https://github.com/kubernetes-csi/csi-driver-host-path/pull/510), [@carlory](https://github.com/carlory))

### Other (Cleanup or Flake)

- Bump image versions ([#528](https://github.com/kubernetes-csi/csi-driver-host-path/pull/528), [@carlory](https://github.com/carlory))
- Replace socat image with hostpathplugin image ([#499](https://github.com/kubernetes-csi/csi-driver-host-path/pull/499), [@carlory](https://github.com/carlory))
- Switched logging library from glog to klog ([#534](https://github.com/kubernetes-csi/csi-driver-host-path/pull/534), [@huww98](https://github.com/huww98))

## Dependencies

### Added
_Nothing has changed._

### Changed
- cloud.google.com/go/compute: v1.23.3 → v1.25.1
- github.com/cncf/xds/go: [0fa0005 → 8a4994d](https://github.com/cncf/xds/go/compare/0fa0005...8a4994d)
- github.com/golang/protobuf: [v1.5.3 → v1.5.4](https://github.com/golang/protobuf/compare/v1.5.3...v1.5.4)
- github.com/stretchr/objx: [v0.5.0 → v0.5.2](https://github.com/stretchr/objx/compare/v0.5.0...v0.5.2)
- github.com/stretchr/testify: [v1.8.4 → v1.9.0](https://github.com/stretchr/testify/compare/v1.8.4...v1.9.0)
- golang.org/x/crypto: v0.19.0 → v0.24.0
- golang.org/x/mod: v0.12.0 → v0.17.0
- golang.org/x/net: v0.21.0 → v0.26.0
- golang.org/x/oauth2: v0.16.0 → v0.18.0
- golang.org/x/sync: v0.6.0 → v0.7.0
- golang.org/x/sys: v0.17.0 → v0.21.0
- golang.org/x/term: v0.17.0 → v0.21.0
- golang.org/x/text: v0.14.0 → v0.16.0
- golang.org/x/tools: v0.12.0 → e35e4cc
- google.golang.org/genproto/googleapis/api: ef43131 → 94a12d6
- google.golang.org/genproto/googleapis/rpc: ef43131 → 94a12d6
- google.golang.org/genproto: ef43131 → f966b18
- google.golang.org/grpc: v1.62.0 → v1.64.0
- google.golang.org/protobuf: v1.32.0 → v1.33.0
- k8s.io/klog/v2: v2.120.1 → v2.130.0
- k8s.io/kubernetes: v1.29.0 → v1.29.2

### Removed
- github.com/cncf/udpa/go: [c52dc94](https://github.com/cncf/udpa/go/tree/c52dc94)
