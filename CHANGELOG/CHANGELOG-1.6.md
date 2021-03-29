# Release notes for v1.6.2

[Documentation](https://kubernetes-csi.github.io)

# Changelog since v1.6.1

## Changes by Kind

### Bug or Regression
 - The YAML files in v1.6.1 should have used the v1.6.1 image, but were still using v1.6.0. ([#266](https://github.com/kubernetes-csi/csi-driver-host-path/pull/266), [@pohly](https://github.com/pohly))

# Release notes for v1.6.1

# Changelog since v1.6.0

## Changes by Kind

### Bug or Regression
 - CSI ephemeral inline volumes failed to get created with an error saying `MountVolume.SetUp failed for volume "ephemeral-volume" : rpc error: code = OutOfRange desc = Requested capacity 1099511627776 exceeds maximum allowed 1099511627776` ([#254](https://github.com/kubernetes-csi/csi-driver-host-path/pull/254), [@pohly](https://github.com/pohly))
 - Deploying `kubernetes-distributed` did not always correctly detect whether the right CSIStorageCapacity API was supported. Incorrectly enabling the feature then prevented scheduling of pods with late binding volumes. ([#262](https://github.com/kubernetes-csi/csi-driver-host-path/pull/262) and [#266](https://github.com/kubernetes-csi/csi-driver-host-path/pull/266), [@pohly](https://github.com/pohly))

### Other (Cleanup or Flake)
 - Updated runtime (Go 1.16) and dependencies ([#259](https://github.com/kubernetes-csi/csi-driver-host-path/pull/259), [@pohly](https://github.com/pohly))
 - Updated external-provisioner from v2.1.0 to v2.1.1 to include some bug fixes ([#268](https://github.com/kubernetes-csi/csi-driver-host-path/pull/268), [@pohly](https://github.com/pohly))

## Dependencies

### Added
_Nothing has changed._

### Changed
_Nothing has changed._

### Removed
_Nothing has changed._

# Release notes for v1.6.0

# Changelog since v1.5.0

## Changes by Kind

### Feature
 - Simulate storage capacity constraints and publish available capacity information ([#248](https://github.com/kubernetes-csi/csi-driver-host-path/pull/248), [@pohly](https://github.com/pohly))

### Other (Cleanup or Flake)
 - Image updates and deployment for Kubernetes 1.20 with external-snapshotter v4.0.0. ([#238](https://github.com/kubernetes-csi/csi-driver-host-path/pull/238), [@pohly](https://github.com/pohly))

## Dependencies

### Added
_Nothing has changed._

### Changed
_Nothing has changed._

### Removed
_Nothing has changed._
