# Release notes for v1.7.3

[Documentation](https://kubernetes-csi.github.io)

# Changelog since v1.7.2

## Changes by Kind

### Bug or Regression

- Fix build on 32bit platforms ([#320](https://github.com/kubernetes-csi/csi-driver-host-path/pull/320), [@c0va23](https://github.com/c0va23))
- Fix rescheduling with distributed provisioning when simulated storage capacity is exhausted ([#331](https://github.com/kubernetes-csi/csi-driver-host-path/pull/331), [@pohly](https://github.com/pohly))
- Fixed bug where `UpdateSnapshot` checks the volume ID when it should check the snapshot ID before updating. ([#315](https://github.com/kubernetes-csi/csi-driver-host-path/pull/315), [@verult](https://github.com/verult))

### Other (Cleanup or Flake)

- Added Kubernetes 1.21 deployments without the obsolete health-monitor-agent ([#313](https://github.com/kubernetes-csi/csi-driver-host-path/pull/313), [@pohly](https://github.com/pohly))
- Changed deployments to use the latest CSI sidecar versions, released during the k8s 1.22 cycle. ([#324](https://github.com/kubernetes-csi/csi-driver-host-path/pull/324), [@verult](https://github.com/verult))
- Livenessprobe version is now 2.4.0 in all deployments. ([#323](https://github.com/kubernetes-csi/csi-driver-host-path/pull/323), [@verult](https://github.com/verult))
- NodePublishVolume now requires that the CO creates the parents of the target directory, as implied by the CSI spec. ([#307](https://github.com/kubernetes-csi/csi-driver-host-path/pull/307), [@pohly](https://github.com/pohly))

## Dependencies

### Added
_Nothing has changed._

### Changed
_Nothing has changed._

### Removed
_Nothing has changed._

# Release notes for v1.7.2

# Changelog since v1.7.1

## Changes by Kind

### Bug or Regression
 - The deploy scripts did not work when used to test sidecars because kustomize didn't allow rbac.yaml files with absolute paths. ([#302](https://github.com/kubernetes-csi/csi-driver-host-path/pull/302), [@pohly](https://github.com/pohly))

## Dependencies

### Added
_Nothing has changed._

### Changed
_Nothing has changed._

### Removed
_Nothing has changed._

# Changelog since v1.7.0

## Changes by Kind

### Uncategorized
 - Fixed image building. ([#298](https://github.com/kubernetes-csi/csi-driver-host-path/pull/298), [@pohly](https://github.com/pohly))

## Dependencies

### Added
_Nothing has changed._

### Changed
_Nothing has changed._

### Removed
_Nothing has changed._

# Release notes for v1.7.0

# Changelog since v1.6.2

## Changes by Kind

### Feature
 - Added 'proxy' support where the host-path driver offloads the CSI requests transparently to the external CSI driver running at `--proxy-endpoint`. ([#260](https://github.com/kubernetes-csi/csi-driver-host-path/pull/260), [@avalluri](https://github.com/avalluri))
 - Implemented Controller{Publish,Unpublish}Volume calls. Introduced a new command-line argument `--enable-attach` (defaults to `false`) which controls if the driver should add `RPC_PUBLISH_UNPUBLISH_VOLUME` to its controller capablilities. ([#260](https://github.com/kubernetes-csi/csi-driver-host-path/pull/260), [@avalluri](https://github.com/avalluri))
 - GRPC calls are logged in a unified(request, reply, and error) JSON format. ([#260](https://github.com/kubernetes-csi/csi-driver-host-path/pull/260), [@avalluri](https://github.com/avalluri))
 - /var/lib/kubelet will be replaced on-the-fly by the deploy.sh scripts with the content of the KUBELET_DATA_DIR env variable. ([#286](https://github.com/kubernetes-csi/csi-driver-host-path/pull/286), [@pohly](https://github.com/pohly))
 - Added support for configuring the maximum number of volumes that could be attached on a node using `--attach-limit`. ([#269](https://github.com/kubernetes-csi/csi-driver-host-path/pull/269), [@avalluri](https://github.com/avalluri))
 - New command-line option `--enable-topology` for enabling/disabling driver topology. ([#269](https://github.com/kubernetes-csi/csi-driver-host-path/pull/269), [@avalluri](https://github.com/avalluri))
 - New command-line option `--node-expand-required` for enabling/disabling volume expansion feature. ([#269](https://github.com/kubernetes-csi/csi-driver-host-path/pull/269), [@avalluri](https://github.com/avalluri))
 - The hostpath driver now has a configurable fixed maximum volume size. It reports the minimum of that and the remaining capacity as `GetCapacityResponse.MaximumVolumeSize`. `GetCapacityResponse.MinimumVolumeSize` is always zero. ([#253](https://github.com/kubernetes-csi/csi-driver-host-path/pull/253), [@pohly](https://github.com/pohly))
 - User now can get volume stats data with Prometheus ([#275](https://github.com/kubernetes-csi/csi-driver-host-path/pull/275), [@stoneshi-yunify](https://github.com/stoneshi-yunify))

### Failing Test
 - Some violations of the volume lifecycle (specifically, VolumeDeleted without NodeUnpublishVolume+NodeUnstageVolume) are not fatal (the behavior in csi-driver-host-path <1.7) and merely cause a warning (new). `--check-volume-lifecycle` can be used to turn such violations into errors. ([#293](https://github.com/kubernetes-csi/csi-driver-host-path/pull/293), [@pohly](https://github.com/pohly))

### Bug or Regression
 - During startup, the driver may have restored internal state incorrectly (volumes added to internal list that belong to some other driver) or failed to start completely (`failed to get capacity info: no such file or directory`). ([#277](https://github.com/kubernetes-csi/csi-driver-host-path/pull/277), [@pohly](https://github.com/pohly))
 - Added the Fix for path resolution in deploy script. Now the user can run script from anywhere, not necessarily from the base path as `./deploy/kubernetes-1.1x/depploy-script.sh` ([#218](https://github.com/kubernetes-csi/csi-driver-host-path/pull/218), [@aayushrangwala](https://github.com/aayushrangwala))

### Other (Cleanup or Flake)
 - All the csi resources from k8s 1.17 will be having a common label `partof: csi-hostpath-driver` which will be used to collectively identify all the resources under csi installation and also can be used to delete them in bulk ([#216](https://github.com/kubernetes-csi/csi-driver-host-path/pull/216), [@aayushrangwala](https://github.com/aayushrangwala))
 - Deployments use "serviceAccountName", the official name for the field, instead of the deprecated "serviceAccount" alias. ([#271](https://github.com/kubernetes-csi/csi-driver-host-path/pull/271), [@pohly](https://github.com/pohly))
 - Testing of CSIStorageCapacity publishing with https://github.com/kubernetes/kubernetes/pull/100537 is enabled. ([#265](https://github.com/kubernetes-csi/csi-driver-host-path/pull/265), [@pohly](https://github.com/pohly))
 - The deploy/kubernetes-x.yy deployments use one pod for all sidecars and the driver. Deploying with separate pods is still supported though the deploy/kubernetes-x.yy-prow deployments, for testing of RBAC rules with Prow. ([#282](https://github.com/kubernetes-csi/csi-driver-host-path/pull/282), [@pohly](https://github.com/pohly))
 - Updated sidecar versions. ([#280](https://github.com/kubernetes-csi/csi-driver-host-path/pull/280), [@pohly](https://github.com/pohly)), ([#294](https://github.com/kubernetes-csi/csi-driver-host-path/pull/294), [@pohly](https://github.com/pohly))
 - Use exec.LookPath function ([#240](https://github.com/kubernetes-csi/csi-driver-host-path/pull/240), [@guilhem](https://github.com/guilhem))
 - Optimize sparse file cloning ([#291](https://github.com/kubernetes-csi/csi-driver-host-path/pull/291), [@stoneshi-yunify](https://github.com/stoneshi-yunify))

## Dependencies

### Added
_Nothing has changed._

### Changed
- github.com/container-storage-interface/spec: [v1.3.0 → v1.4.0](https://github.com/container-storage-interface/spec/compare/v1.3.0...v1.4.0)
- github.com/stretchr/testify: [v1.6.1 → v1.7.0](https://github.com/stretchr/testify/compare/v1.6.1...v1.7.0)

### Removed
_Nothing has changed._
