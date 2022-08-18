The deployment for Kubernetes 1.21 uses the CSI snapshotter sidecar
4.x and thus is incompatible with Kubernetes clusters where older
snapshotter CRDs are installed.

The health-monitor-agent is no longer getting deployed because its
functionality was moved into kubelet in Kubernetes 1.21.
