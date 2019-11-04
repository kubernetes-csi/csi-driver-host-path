The deployment for Kubernetes 1.16 enables ephemeral inline volumes via
its CSIDriverInfo and thus is incompatible with Kubernetes < 1.16
because the `VolumeLifecycleModes` field is rejected by those release.

The following canary images are known to be incompatible with this
deployment:
- csi-snapshotter (canary uses VolumeSnapshot v1beta)
