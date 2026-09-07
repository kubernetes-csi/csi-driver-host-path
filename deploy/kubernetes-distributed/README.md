This deployment is meant for Kubernetes clusters with
CSIStorageCapacity enabled. It deploys the hostpath driver on each
node, using distributed provisioning, and configures it so that it has
10Gi of "fast" storage and 100Gi of "slow" storage.

The "kind" storage class parameter can selected between the two. If
not set, an arbitrary kind with enough capacity is picked.

## Prerequisites

Snapshot support in this deployment uses per-node `csi-snapshotter`
sidecars (`--node-deployment=true`). For that to work, the cluster
must run an external `snapshot-controller` started with
`--enable-distributed-snapshotting=true`, plus the matching
`VolumeSnapshot*` CRDs. The `deploy.sh` script in this directory does
not install snapshot-controller; install it separately before
deploying the driver. See the
[external-snapshotter documentation](https://github.com/kubernetes-csi/external-snapshotter#distributed-snapshotting)
for the install steps.

Cross-node snapshot restore (provisioning a PVC from a VolumeSnapshot
whose source data lives on a different node) is not yet covered by
this deployment alone. The companion `csi-topology-coordinator`
component handles that case.
