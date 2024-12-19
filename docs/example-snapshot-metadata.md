## Snapshot Changed Block Metadata Support

The CSI HostPath driver now includes support for the CSI [SnapshotMetadata](https://github.com/container-storage-interface/spec/blob/master/csi.proto#L130) service. This service provides APIs to retrieve metadata about the allocated blocks of a CSI VolumeSnapshot or the changed blocks between any two CSI VolumeSnapshot objects of the same PersistentVolume.

This document outlines the steps to test this feature on a Kubernetes cluster.

### Deploying CSI Hostpath driver with SnapshotMetadata service

Setting up CSI Hostpath driver with SnapshotMetadata service requires provisioning TLS certificates, creating TLS secrets, SnapshotMetadata custom resource, patching up csi-hostpathplugin deployments, etc. These steps are automated in `deploy.sh` script used to deploy CSI Hostpath driver.

Follow the steps below to deploy CSI Hostpath driver with SnapshotMetadata service:

  a. Create `SnapshotMetadata` CRD

  ```
  $ kubectl create -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshot-metadata/main/client/config/crd/cbt.storage.k8s.io_snapshotmetadataservices.yaml
  ```

  b. Execute deploy script to setup hostpath plugin driver with external-snapshot-metadata change

  ```
  $ SNAPSHOT_METADATA_TESTS=true ./deploy/kubernetes-1.27/deploy.sh
  ```

### Setup SnapshotMetadata client

The `SnapshotMetadata` service implements gRPC APIs. A gRPC client can query these APIs to retrieve metadata about the allocated blocks of a CSI VolumeSnapshot or the changed blocks between any two CSI VolumeSnapshot objects.

For our testing, we will be using a sample client implementation in Go provided as a example in [external-snapshot-metadata](https://github.com/kubernetes-csi/external-snapshot-metadata/tree/main/examples/snapshot-metadata-lister) repo.

Follow the following steps to setup client with all the required permissions:

1. Setup RBAC

   a. Create `ClusterRole` containing all the required permissions for the client

   ```bash
   $ kubectl create -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshot-metadata/main/deploy/snapshot-metadata-client-cluster-role.yaml
   ```

   b. Create a namespace to deploy client

   ```
   $ kubectl create namespace csi-client
   ```

   c. Create service account

   ```
   $ kubectl create serviceaccount csi-client-sa -n csi-client
   ```

   d. Bind the clusterrole to the service account

   ```
   $ kubectl create clusterrolebinding csi-client-cluster-role-binding --clusterrole=external-snapshot-metadata-client-runner --serviceaccount=csi-client:csi-client-sa
   ```

2. Deploy sample client pod which contains [snapshot-metadata-lister](https://github.com/kubernetes-csi/external-snapshot-metadata) tool which can be used as a client to call SnapshotMetadata APIs

    ```
    $ kubectl create -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshot-metadata/main/examples/snapshot-metadata-lister/deploy/snapshot-medata-lister-pod.yaml -n csi-client
    ```

This client performs following actions:
1. Find Driver name for the snapshot.
2. Discover `SnapshotMetadataService` resource for the driver which contains endpoint, audience and CA cert.
3. Create SA Token with expected audience and permissions.
4. Make gRPC call `GetMetadataAllocated` and `GetMetadataDelta` with appropriate params from `SnapshotMetadataService` resource, generated SA token.
5. Stream response and print on console.


### Test GetMetadataAllocated

1. Create CSI Hostpath storageclass

    ```
    $ kubectl create -f examples/csi-storageclass.yaml
    ```

2. Create a volume with Block mode access

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

3. Mount the PVC to a pod

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

4. Snapshot PVC `pvc-raw`

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

    Wait for snapshot to be ready

    ```
    $ kg vs snapshot-1
    NAME         READYTOUSE   SOURCEPVC   SOURCESNAPSHOTCONTENT   RESTORESIZE   SNAPSHOTCLASS            SNAPSHOTCONTENT                                    CREATIONTIME   AGE
    snapshot-1   true         pvc-raw                             10Mi          csi-hostpath-snapclass   snapcontent-70b40b27-80d4-448b-bd9a-a87079c1a248   28s            29s
    ```

5. Now, inside `csi-client` pod which is created in previous steps, use `snapshot-metadata-lister` tool query allocated blocks metadata

    ```
    $ kubectl exec -n csi-client csi-client -c run-client -- /tools/snapshot-metadata-lister -n default -s raw-pvc-snap-1

    Record#   VolCapBytes  BlockMetadataType   ByteOffset     SizeBytes   
    ------- -------------- ----------------- -------------- --------------
          1       10485760      FIXED_LENGTH              0           4096
          1       10485760      FIXED_LENGTH           4096           4096
          1       10485760      FIXED_LENGTH           8192           4096
          1       10485760      FIXED_LENGTH          12288           4096
          1       10485760      FIXED_LENGTH          16384           4096
          1       10485760      FIXED_LENGTH          20480           4096
          1       10485760      FIXED_LENGTH          24576           4096
          1       10485760      FIXED_LENGTH          28672           4096
          1       10485760      FIXED_LENGTH          32768           4096
          1       10485760      FIXED_LENGTH          36864           4096
          1       10485760      FIXED_LENGTH          40960           4096
    .
    .
    .
    .
         10       10485760      FIXED_LENGTH       10452992           4096
         10       10485760      FIXED_LENGTH       10457088           4096
         10       10485760      FIXED_LENGTH       10461184           4096
         10       10485760      FIXED_LENGTH       10465280           4096
         10       10485760      FIXED_LENGTH       10469376           4096
         10       10485760      FIXED_LENGTH       10473472           4096
         10       10485760      FIXED_LENGTH       10477568           4096
         10       10485760      FIXED_LENGTH       10481664           4096
    ```


### Test GetMetadataDelta

1. Change couple of blocks in the mounted device file in `pod-raw` Pod

    ```
    $ kubectl exec -ti pod-raw -- sh

    ### change blocks 12, 13, 15 and 20
    / # dd if=/dev/urandom of=/dev/loop3 bs=4K count=1 seek=12 conv=notrunc
    1+0 records in
    1+0 records out
    / # dd if=/dev/urandom of=/dev/loop3 bs=4K count=1 seek=13 conv=notrunc
    1+0 records in
    1+0 records out
    / # dd if=/dev/urandom of=/dev/loop3 bs=4K count=1 seek=15 conv=notrunc
    1+0 records in
    1+0 records out
    / # dd if=/dev/urandom of=/dev/loop3 bs=4K count=1 seek=20 conv=notrunc
    1+0 records in
    1+0 records out

    ```

2. Snapshot `pvc-raw` again

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

    Wait for snapshot to be ready

    ```
    $ kubectl get vs
    NAME             READYTOUSE   SOURCEPVC   SOURCESNAPSHOTCONTENT   RESTORESIZE   SNAPSHOTCLASS            SNAPSHOTCONTENT                                    CREATIONTIME   AGE
    raw-pvc-snap-1   true         pvc-raw                             10Mi          csi-hostpath-snapclass   snapcontent-ef10f725-4261-4e80-af37-906708796700   7m40s          7m40s
    raw-pvc-snap-2   true         pvc-raw                             10Mi          csi-hostpath-snapclass   snapcontent-188562cb-03b3-4b70-b12d-28900527bca8   23s            23s
    ```

3. Using `external-snapshot-metadata-client` which uses `GetMetadataDelta` gRPC to allocated blocks metadata

    ```
    $ kubectl exec -n csi-client csi-client -c run-client -- /tools/snapshot-metadata-lister -n default -s raw-pvc-snap-1 -p raw-pvc-snap-2


    Record#   VolCapBytes  BlockMetadataType   ByteOffset     SizeBytes
    ------- -------------- ----------------- -------------- --------------
          1       10485760      FIXED_LENGTH          49152           4096
          1       10485760      FIXED_LENGTH          53248           4096
          1       10485760      FIXED_LENGTH          61440           4096
          1       10485760      FIXED_LENGTH          81920           4096
    ```