The deployment for Kubernetes 1.15 uses CSI 1.0 and thus is
incompatible with Kubernetes < 1.13.

The sidecars depend on 1.15 API changes for migration and resizing,
and 1.14 API changes for CSIDriver and CSINode.

In contrast to the normal deployment, this deployment does not use the
external-attacher sidecar. Instead, it creates a CSIDriver instance
that tells Kubernetes not to use attach (ControllerPublishVolume) at all.
