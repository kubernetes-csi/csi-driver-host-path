This deployment is meant for Kubernetes clusters with
CSIStorageCapacity enabled. It deploys the hostpath driver on each
node, using distributed provisioning, and configures it so that it has
10Gi of "fast" storage and 100Gi of "slow" storage.

The "kind" storage class parameter can selected between the two. If
not set, an arbitrary kind with enough capacity is picked.

