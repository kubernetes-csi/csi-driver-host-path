The deployment for Kubernetes 1.13 uses CSI 1.0 and thus is
incompatible with older Kubernetes releases.

The sidecar images rely on the CRDs for CSIDriverInfo and CSINodeInfo,
which were replaced with builtin APIs in Kubernetes 1.14. They can be
deployed on Kubernetes 1.14 if the CRDs are installed, but features
relying on these CRDs (like topology) are unlikely to work.
