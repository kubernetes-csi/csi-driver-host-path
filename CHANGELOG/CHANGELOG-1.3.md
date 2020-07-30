# Changelog since v1.2.0

## New Features

- Adds deployment for Kubernetes 1.17. Starting with this deployment, the snapshot CRD is no longer part of the driver installation and must be installed as part of the cluster. ([#98](https://github.com/kubernetes-csi/csi-driver-host-path/pull/98), [@xing-yang](https://github.com/xing-yang))
- Add option to limit number of volumes per node. ([#110](https://github.com/kubernetes-csi/csi-driver-host-path/pull/110), [@bertinatto](https://github.com/bertinatto))
- updated sidecars to latest stable releases ([#102](https://github.com/kubernetes-csi/csi-driver-host-path/pull/102), [@pohly](https://github.com/pohly))
- The -ephemeral parameter (currently alpha) is still supported, but only needed for Kubernetes 1.15 and will be removed once Kubernetes 1.15 stops being supported. On Kubernetes 1.16, the same deployment supports normal persistent volumes and inline ephemeral volumes. ([#67](https://github.com/kubernetes-csi/csi-driver-host-path/pull/67), [@pohly](https://github.com/pohly))


## Bug Fixes

- fixed raw block volumes on hosts that don't have /dev/loop devices pre-defined ([#109](https://github.com/kubernetes-csi/csi-driver-host-path/pull/109), [@pohly](https://github.com/pohly))
- NodeVolumeUnpublish: tolerate repeated requests ([#139](https://github.com/kubernetes-csi/csi-driver-host-path/pull/139), [@okartau](https://github.com/okartau))
- NodePublishVolume: return error when mount fails ([#146](https://github.com/kubernetes-csi/csi-driver-host-path/pull/146), [@c3y1huang](https://github.com/c3y1huang))
- CreateVolume: validate requested size ([#151](https://github.com/kubernetes-csi/csi-driver-host-path/pull/151), [@Madhu-1](https://github.com/Madhu-1))
- CreateVolume: check for volume content source ([#148](https://github.com/kubernetes-csi/csi-driver-host-path/pull/148), [@Madhu-1](https://github.com/Madhu-1))

## Other Notable Changes

- Remove deployment for the unsupported Kubernetes 1.14 ([#138](https://github.com/kubernetes-csi/csi-driver-host-path/pull/138), [@msau42](https://github.com/msau42))
