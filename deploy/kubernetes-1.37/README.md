The deployment for Kubernetes 1.37 uses the CSI snapshotter sidecar
4.x and thus is incompatible with Kubernetes clusters where older
snapshotter CRDs are installed.

This deployment includes the external-health-monitor-controller
sidecar. It requires a hostpath driver that implements the
ControllerGetVolumeHealth and ControllerListVolumeHealth RPCs (CSI spec
1.13) and exits otherwise.

TODO: update the hostpathplugin image to a release that implements these
RPCs. v1.17.1 does not, so this deployment only works when the driver
image is overridden (as the CI jobs do).

The health-monitor-agent is no longer getting deployed because its
functionality was moved into kubelet in Kubernetes 1.21.
