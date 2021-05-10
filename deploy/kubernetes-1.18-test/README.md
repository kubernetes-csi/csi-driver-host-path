The deployment for Kubernetes 1.18 uses CSIDriver v1 and
thus is incompatible with Kubernetes < 1.18.

It uses separate pods and service accounts for each sidecar. This is
not how they would normally be deployed. It gets done this way to test
that the individual RBAC rules are correct.
