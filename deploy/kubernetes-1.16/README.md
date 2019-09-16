The deployment for Kubernetes 1.16 enables ephemeral inline volumes via
its CSIDriverInfo and thus is incompatible with Kubernetes < 1.16
because the `VolumeLifecycleModes` field is rejected by those release.
