# Changelog since v1.1.0

## Breaking Changes

- The deployment uses `hostpath.csi.k8s.io` as driver name ([#64](https://github.com/kubernetes-csi/csi-driver-host-path/pull/64), [@pohly](https://github.com/pohly)).
  Make sure that there are no persistent or ephemeral volumes using the old `csi-hostpath` name before updating because otherwise
  those volumes cannot be removed. Pods with such ephemeral volumes will be stuck in "terminating" state. New pods
  will not be able to start if they reference a volume that uses the old name.
  Any storage class that references the driver must be updated together with the driver.

## New Features

- normal deployment supports ephemeral inline volumes on Kubernetes 1.15 and 1.16 ([#67](https://github.com/kubernetes-csi/csi-driver-host-path/pull/67), [@pohly](https://github.com/pohly), [#97](https://github.com/kubernetes-csi/csi-driver-host-path/pull/97), [@pohly](https://github.com/pohly))
- volume expansion support ([#87](https://github.com/kubernetes-csi/csi-driver-host-path/pull/87), [@bertinatto](https://github.com/bertinatto), [#90](https://github.com/kubernetes-csi/csi-driver-host-path/pull/90), [@bertinatto](https://github.com/bertinatto))
- topology support ([#88](https://github.com/kubernetes-csi/csi-driver-host-path/pull/88), [@mucahitkurt](https://github.com/mucahitkurt))
- cloning support ([#58](https://github.com/kubernetes-csi/csi-driver-host-path/pull/58), [@j-griffith](https://github.com/j-griffith))

## Bug Fixes

- /csi-data-dir optional ([#73](https://github.com/kubernetes-csi/csi-driver-host-path/pull/73), [@msau42](https://github.com/msau42))
- Set volume content source if creating volume from snapshot ([#51](https://github.com/kubernetes-csi/csi-driver-host-path/pull/51), [@zhucan](https://github.com/zhucan)).
- Fixes cp expansion issue for volume cloning ([#82](https://github.com/kubernetes-csi/csi-driver-host-path/pull/82), [@j-griffith](https://github.com/j-griffith)).

## Other Notable Changes

- Added deployment specs for K8s 1.15 ([#63](https://github.com/kubernetes-csi/csi-driver-host-path/pull/63), [@msau42](https://github.com/msau42)) and 1.16 ([#97](https://github.com/kubernetes-csi/csi-driver-host-path/pull/97), [@pohly](https://github.com/pohly)).
- Removed deployment specs for K8s 1.13 because it is no longer supported ([#102](https://github.com/kubernetes-csi/csi-driver-host-path/pull/102), [@pohly](https://github.com/pohly)).
- Updated sidecars to latest stable releases ([#102](https://github.com/kubernetes-csi/csi-driver-host-path/pull/102), [@pohly](https://github.com/pohly)).
