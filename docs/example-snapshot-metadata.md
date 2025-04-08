## Snapshot Changed Block Metadata Support

The CSI hostpath driver now includes support for the [CSI SnapshotMetadata](https://github.com/container-storage-interface/spec/blob/master/csi.proto#L130) service.
This service provides APIs to retrieve metadata about the allocated blocks of a CSI VolumeSnapshot or the changed blocks between any two CSI VolumeSnapshot objects of the same PersistentVolume.

A Kubernetes application cannot directly access a CSI SnapshotMetadata service but instead
communicates with a
[Kubernetes SnapshotMetadataService](https://github.com/kubernetes-csi/external-snapshot-metadata/blob/main/proto/schema.proto#L11)
provided by an
[external-snapshot-metadata sidecar](https://github.com/kubernetes/enhancements/tree/master/keps/sig-storage/3314-csi-changed-block-tracking#the-external-snapshot-metadata-sidecar)
that fronts the CSI service.
Access to the Kubernetes SnapshotMetadata service is advertised by a
[SnapshotMetadata CR](https://github.com/kubernetes/enhancements/tree/master/keps/sig-storage/3314-csi-changed-block-tracking#snapshot-metadata-service-custom-resource).

This document outlines the steps to test this feature on a Kubernetes cluster using the CSI hostpath driver and contains the following sections:
- [Deploy the CSI hostpath driver with a Kubernetes SnapshotMetadata service](#deploy-the-csi-hostpath-driver-with-a-kubernetes-snapshotmetadata-service)
- [Setup a Kubernetes SnapshotMetadata client](#setup-a-kubernetes-snapshotmetadata-client)
- [Create a stateful application](#create-a-stateful-application)
- [Test the GetMetadataAllocated RPC](#test-the-getmetadataallocated-rpc)
- [Test the GetMetadataDelta RPC](#test-the-getmetadatadelta-rpc)

### Deploy the CSI hostpath driver with a Kubernetes SnapshotMetadata service

Setting up the CSI hostpath driver with a Kubernetes SnapshotMetadata service requires provisioning TLS certificates, creating TLS secrets, a SnapshotMetadata custom resource, patching the csi-hostpathplugin deployments, etc.
These steps are automated in the `deploy.sh` script used to deploy CSI Hostpath driver into the current namespace when invoked with the
appropriate environment variable.

Follow the steps below to deploy the CSI hostpath driver with a Kubernetes SnapshotMetadata service:

1. Create the `SnapshotMetadataService` CRD using the definition in the 
   [external-snapshot-metadata](https://github.com/kubernetes-csi/external-snapshot-metadata/tree/main/examples/snapshot-metadata-lister) repository.
   ```
    kubectl create -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshot-metadata/main/client/config/crd/cbt.storage.k8s.io_snapshotmetadataservices.yaml
   ```

2. Execute the deploy script to setup the hostpath plugin driver with the Kubernetes SnapshotMetadata service:
   ```
   SNAPSHOT_METADATA_TESTS=true HOSTPATHPLUGIN_REGISTRY=registry.k8s.io/sig-storage HOSTPATHPLUGIN_TAG=v1.16.1 ./deploy/kubernetes-latest/deploy.sh
   ```
   Specifying the `SNAPSHOT_METADATA_TESTS=true` environment variable causes it to configure the `external-snapshot-metadata` sidecar
   in the CSI hostpath driver Pod.
   The `HOSTPATHPLUGIN_REGISTRY` and `HOSTPATHPLUGIN_TAG` environment variables are used to override defaults for the CSI hostpath driver image.

### Setup a Kubernetes SnapshotMetadata client

The `SnapshotMetadata` service implements gRPC APIs. A gRPC client can query these APIs to retrieve metadata about the allocated blocks of a CSI VolumeSnapshot or the changed blocks between any two CSI VolumeSnapshot objects.

For our testing, we will be using a sample client implementation in Go provided as a example in the
[external-snapshot-metadata](https://github.com/kubernetes-csi/external-snapshot-metadata/tree/main/examples/snapshot-metadata-lister) repository.
This client performs the actions that would be taken by a real backup application to fetch snapshot metadata
on allocated or changed blocks.

The following steps sets up the client with all the required permissions:

1. Setup RBAC for the client

   a. Create a ClusterRole containing all the required permissions for the client

   ```bash
   kubectl create -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshot-metadata/main/deploy/snapshot-metadata-client-cluster-role.yaml
   ```

   b. Create a Namespace to deploy the client

   ```
   kubectl create namespace csi-client
   ```
   Note that the `csi-client` name is hardcoded in the Pod YAML used below.

   c. Create a ServiceAccount in the namespace

   ```
   kubectl create serviceaccount csi-client-sa -n csi-client
   ```

   d. Bind the ClusterRole to this ServiceAccount

   ```
   kubectl create clusterrolebinding csi-client-cluster-role-binding --clusterrole=external-snapshot-metadata-client-runner --serviceaccount=csi-client:csi-client-sa
   ```

2. Deploy a Pod with the
   [snapshot-metadata-lister](https://github.com/kubernetes-csi/external-snapshot-metadata) tool,
   using the sample YAML from the
   [external-snapshot-metadata](https://github.com/kubernetes-csi/external-snapshot-metadata/tree/main/examples/snapshot-metadata-lister) repository:
    ```
    kubectl create -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshot-metadata/main/examples/snapshot-metadata-lister/deploy/snapshot-metadata-lister-pod.yaml -n csi-client
    ```
    This takes some time to complete because an init container is used to compile the tool.
    Wait for the Pod to become ready before continuing:
    ```
    kubectl wait -n csi-client --for=jsonpath='{.status.phase}'=Running --timeout=5m pod/csi-client
    ```

### Create a stateful application

We now create a stateful application that writes data to a PersistentVolume provided by the CSI hostpath driver.

The application is installed into the current Namespace (assumed to be `default`) for convenience, but could as well be installed in a separate namespace.
If this is not the case then substitute the explicit use of `-n default` with the appropriate value in the `snapshot-metadata-lister` invocations that follow.

1. Create the CSI hostpath driver's StorageClass:
    ```
    kubectl create -f examples/csi-storageclass.yaml
    ```

2. Create a dynamically provisioned PVC with Block mode access for this StorageClass in the current Namespace:
    ```
    kubectl apply -f - <<EOF
    kind: PersistentVolumeClaim
    apiVersion: v1
    metadata:
      name: pvc-raw
    spec:
      accessModes:
        - ReadWriteOnce
      storageClassName: csi-hostpath-sc
      volumeMode: Block
      resources:
        requests:
          storage: 10Mi
    EOF
    ```
    Note: the CSI hostpath driver only supports the CSI SnapshotMetadata service with volumes in Block mode.
    Real CSI drivers that support this service are expected do so with any volume mode.

3. Mount the PVC in a Pod:
    ```
    kubectl apply -f - <<EOF
    apiVersion: v1
    kind: Pod
    metadata:
      name: pod-raw
      labels:
        name: busybox-test
    spec:
      restartPolicy: Always
      containers:
        - image: gcr.io/google_containers/busybox
          command: ["/bin/sh", "-c"]
          args: [ "tail -f /dev/null" ]
          name: busybox
          volumeDevices:
            - name: vol
              devicePath: /dev/loop3
      volumes:
        - name: vol
          persistentVolumeClaim:
            claimName: pvc-raw
    EOF
    ```
    Here the volume is mounted at the `/dev/loop3` device path.
    Wait for the Pod to become ready before continuing:
    ```
    kubectl wait --for=jsonpath='{.status.phase}'=Running pod/pod-raw
    ```

### Test the GetMetadataAllocated RPC
In this test we add data to the volume, take a snapshot and then fetch
metadata on the allocated block of the snapshot.

1. Add some initial data to the mounted volume:
    ```
    kubectl exec -it pod-raw -- dd if=/dev/urandom of=/dev/loop3 bs=4K count=1 seek=1 conv=notrunc
    kubectl exec -it pod-raw -- dd if=/dev/urandom of=/dev/loop3 bs=4K count=1 seek=3 conv=notrunc
    kubectl exec -it pod-raw -- dd if=/dev/urandom of=/dev/loop3 bs=4K count=1 seek=5 conv=notrunc
    kubectl exec -it pod-raw -- dd if=/dev/urandom of=/dev/loop3 bs=4K count=1 seek=7 conv=notrunc
    kubectl exec -it pod-raw -- dd if=/dev/urandom of=/dev/loop3 bs=4K count=1 seek=9 conv=notrunc
    ```
    The calls above add 4Ki sized blocks of random data at block addresses 1, 3, 5, 7 and 9.

2. Snapshot the `pvc-raw` PVC creating VolumeSnapshot `raw-pvc-snap-1`
    ```
    kubectl apply -f - <<EOF
    apiVersion: snapshot.storage.k8s.io/v1
    kind: VolumeSnapshot
    metadata:
      name: raw-pvc-snap-1
    spec:
      volumeSnapshotClassName: csi-hostpath-snapclass
      source:
        persistentVolumeClaimName: pvc-raw
    EOF
    ```

    Wait for the snapshot to be ready:
    ```
    kubectl get vs raw-pvc-snap-1
    ```
    The output should look something like this when the snapshot is available:
    ```
    NAME             READYTOUSE   SOURCEPVC   SOURCESNAPSHOTCONTENT   RESTORESIZE   SNAPSHOTCLASS            SNAPSHOTCONTENT                                    CREATIONTIME   AGE
    raw-pvc-snap-1   true         pvc-raw                             10Mi          csi-hostpath-snapclass   snapcontent-e17ba543-b8be-4a8e-9b0f-d708d664a0ee   99s            100s
    ```

3. Now, inside the `csi-client` Pod that was created in previous steps, use the `snapshot-metadata-lister` tool to query
   for metadata on the allocated blocks:

    ```
    kubectl exec -n csi-client csi-client -c run-client -- /tools/snapshot-metadata-lister -n default -s raw-pvc-snap-1
    ```
    The command above requests metadata on the allocated blocks of the current snapshot (`-s` flag) of the
    VolumeSnapshot object in the specified Namespace (`-n` flag).
 
    The output should look something like this:
    ```
    Record#   VolCapBytes  BlockMetadataType   ByteOffset     SizeBytes
    ------- -------------- ----------------- -------------- --------------
         1       10485760      FIXED_LENGTH           4096           4096
         1       10485760      FIXED_LENGTH          12288           4096
         1       10485760      FIXED_LENGTH          20480           4096
         1       10485760      FIXED_LENGTH          28672           4096
         1       10485760      FIXED_LENGTH          36864           4096
    ```
    There should be 5 allocated blocks.

### Test the GetMetadataDelta RPC
In this test we make changes to the volume and then make another snapshot.
We then fetch metadata on the changes between the latest and the previous snapshots,
and then view the allocated blocks of the latest snapshot.

1. Now make some changes to the mounted volume:
    ```
    kubectl exec -it pod-raw -- dd if=/dev/urandom of=/dev/loop3 bs=4K count=1 seek=2 conv=notrunc
    kubectl exec -it pod-raw -- dd if=/dev/urandom of=/dev/loop3 bs=4K count=1 seek=4 conv=notrunc
    kubectl exec -it pod-raw -- dd if=/dev/urandom of=/dev/loop3 bs=4K count=1 seek=9 conv=notrunc
    ```
    The calls above add two new blocks at block addresses 2 and 4, and modifies the existing block at block address 9.

2. Snapshot the `pvc-raw` PVC again, creating VolumeSnapshot `raw-pvc-snap-2`:

    ```
    kubectl apply -f - <<EOF
    apiVersion: snapshot.storage.k8s.io/v1
    kind: VolumeSnapshot
    metadata:
      name: raw-pvc-snap-2
    spec:
      volumeSnapshotClassName: csi-hostpath-snapclass
      source:
        persistentVolumeClaimName: pvc-raw
    EOF
    ```

    Wait for the snapshot to be ready:
    ```
    kubectl get vs raw-pvc-snap-2
    ```
    The output should look something like this when the snapshot is available:
    ```
    NAME             READYTOUSE   SOURCEPVC   SOURCESNAPSHOTCONTENT   RESTORESIZE   SNAPSHOTCLASS            SNAPSHOTCONTENT                                    CREATIONTIME   AGE
    raw-pvc-snap-2   true         pvc-raw                             10Mi          csi-hostpath-snapclass   snapcontent-630fc6d8-6b42-48e2-8ec0-c215a1f65882   7s             7s
    ```

3. Once again use the `snapshot-metadata-lister` tool, but this time to view the changes between the two snapshots:
    ```
    kubectl exec -n csi-client csi-client -c run-client -- /tools/snapshot-metadata-lister -n default -s raw-pvc-snap-2 -p raw-pvc-snap-1
    ```
    The command requests metadata on the changed blocks between the current (`-s` flag) and the previous (`-p` flag)
    VolumeSnapshot objects in the specified Namespace (`-n` flag).

    The ouptut should look something like this:
    ```
    Record#   VolCapBytes  BlockMetadataType   ByteOffset     SizeBytes
    ------- -------------- ----------------- -------------- --------------
          1       10485760      FIXED_LENGTH           8192           4096
          1       10485760      FIXED_LENGTH          16384           4096
          1       10485760      FIXED_LENGTH          36864           4096
    ```

4. View the allocated blocks of the latest snapshot:
   ```
   kubectl exec -n csi-client csi-client -c run-client -- /tools/snapshot-metadata-lister -n default -s raw-pvc-snap-2
   ```
   The command above requests metadata on the allocated blocks of the latest snapshot (`-s` flag) of the
   VolumeSnapshot object in the specified Namespace (`-n` flag).

   The output should look something like this:
   ```
   Record#   VolCapBytes  BlockMetadataType   ByteOffset     SizeBytes
   ------- -------------- ----------------- -------------- --------------
         1       10485760      FIXED_LENGTH           4096           4096
         1       10485760      FIXED_LENGTH           8192           4096
         1       10485760      FIXED_LENGTH          12288           4096
         1       10485760      FIXED_LENGTH          16384           4096
         1       10485760      FIXED_LENGTH          20480           4096
         1       10485760      FIXED_LENGTH          28672           4096
         1       10485760      FIXED_LENGTH          36864           4096
   ```
   There should be 7 allocated blocks in total.
