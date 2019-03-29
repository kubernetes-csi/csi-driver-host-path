The deployment for Kubernetes 1.13 uses CSI 1.0 and this is
incompatible with older Kubernetes releases.

It relies on the CRDs for CSIDriverInfo and CSINodeInfo, which are
about to be replaced with builtin APIs in Kubernetes 1.14. It can be
deployed on Kubernetes 1.14 if the CRDs are installed, but features
relying on these CRDs (like topology) are unlikely to work.

Kubernetes 1.14 will need a different deployment.
