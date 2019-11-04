The deployment for Kubernetes 1.14 uses CSI 1.0 and thus is
incompatible with Kubernetes < 1.13.

It uses the APIs for CSIDriverInfo and CSINodeInfo that were
introduced in Kubernetes 1.14, so features depending on those (like
topology) will not work on Kubernetes 1.13. But because this example
deployment does not enable those features, it can run on Kubernetes
1.13.

The following canary images are known to be incompatible with this
deployment:
- csi-snapshotter (canary uses VolumeSnapshot v1beta)
