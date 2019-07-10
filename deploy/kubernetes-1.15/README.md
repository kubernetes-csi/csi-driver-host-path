The deployment for Kubernetes 1.15 uses CSI 1.0 and thus is
incompatible with Kubernetes < 1.13.

The sidecars depend on 1.15 API changes for migration and resizing,
and 1.14 API changes for CSIDriver and CSINode.
However the hostpath driver doesn't use those features, so this
deployment can work on older Kubernetes versions.
