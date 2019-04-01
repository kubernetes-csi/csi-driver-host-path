The deployment for master uses CSI 1.0 and thus is incompatible with
Kubernetes < 1.13.

It uses the APIs for CSIDriverInfo and CSINodeInfo that were
introduced in Kubernetes 1.14, so features depending on those (like
topology) will not work on Kubernetes 1.13. But because this example
deployment does not enable those features, it can run on Kubernetes
1.13.

WARNING: this example uses the "canary" images. It can break at any
time.
