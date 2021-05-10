The deployment for Kubernetes 1.20 uses the CSI snapshotter sidecar
4.x and thus is incompatible with Kubernetes clusters where older
snapshotter CRDs are installed.

It uses separate pods and service accounts for each sidecar. This is
not how they would normally be deployed. It gets done this way to test
that the individual RBAC rules are correct.
