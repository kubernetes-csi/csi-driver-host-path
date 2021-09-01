## Snapshot support

Volume snapshot was introduced as an alpha feature in Kubernetes v1.12 and it was promoted to beta in Kubernetes v1.17. Volume snapshot is GA in Kubernetes 1.20.

### Pre-requisites 
- For Kubernetes 1.21+, you must install the [VolumeSnapshot beta CRDs and the Snapshot Controller](deploy-1.21-and-later.md)

### Creating and restoring volume snapshots
Ensure your volumesnapshotclass was created during hostpath deployment:


$ kubectl get volumesnapshotclass
 ```
 NAME                     DRIVER                DELETIONPOLICY   AGE
  csi-hostpath-snapclass   hostpath.csi.k8s.io   Delete           149m
```

$ kubectl describe volumesnapshotclass
```
  Name:             csi-hostpath-snapclass
Namespace:        
Labels:           app.kubernetes.io/component=volumesnapshotclass
                  app.kubernetes.io/instance=hostpath.csi.k8s.io
                  app.kubernetes.io/name=csi-hostpath-snapclass
                  app.kubernetes.io/part-of=csi-driver-host-path
Annotations:      <none>
API Version:      snapshot.storage.k8s.io/v1
Deletion Policy:  Delete
Driver:           hostpath.csi.k8s.io
Kind:             VolumeSnapshotClass
Metadata:
  Creation Timestamp:  2022-07-18T04:10:47Z
  Generation:          1
  Managed Fields:
    API Version:  snapshot.storage.k8s.io/v1
    Fields Type:  FieldsV1
    fieldsV1:
      f:deletionPolicy:
      f:driver:
      f:metadata:
        f:annotations:
          .:
          f:kubectl.kubernetes.io/last-applied-configuration:
        f:labels:
          .:
          f:app.kubernetes.io/component:
          f:app.kubernetes.io/instance:
          f:app.kubernetes.io/name:
          f:app.kubernetes.io/part-of:
    Manager:         kubectl-client-side-apply
    Operation:       Update
    Time:            2022-07-18T04:10:47Z
  Resource Version:  1112
  UID:               c82cc073-4ef0-4659-9522-7b2a65df9ff5
Events:              <none>
 ```

After having created the `csi-pvc` as described in the deployment validation,
use the volume snapshot class to dynamically create a volume snapshot:
  - `$ kubectl apply -f examples/csi-snapshot-v1.yaml`
 ```
 volumesnapshot.snapshot.storage.k8s.io/new-snapshot-demo created
 ```


 $ kubectl get volumesnapshot
 ```
NAME                   READYTOUSE   SOURCEPVC   SOURCESNAPSHOTCONTENT   RESTORESIZE   SNAPSHOTCLASS            SNAPSHOTCONTENT                                    CREATIONTIME   AGE
new-snapshot-v1-demo   true         csi-pvc                             1Gi           csi-hostpath-snapclass   snapcontent-aae1446f-e594-46b9-806f-d6775296f804   23s            46s
```

 $ kubectl get volumesnapshotcontent
```
NAME                                               READYTOUSE   RESTORESIZE   DELETIONPOLICY   DRIVER                VOLUMESNAPSHOTCLASS      VOLUMESNAPSHOT         AGE
snapcontent-aae1446f-e594-46b9-806f-d6775296f804   true         1073741824    Delete           hostpath.csi.k8s.io   csi-hostpath-snapclass   new-snapshot-v1-demo   80s
 ```

$ kubectl describe volumesnapshot
 ```
Name:         new-snapshot-v1-demo
Namespace:    default
Labels:       <none>
Annotations:  <none>
API Version:  snapshot.storage.k8s.io/v1
Kind:         VolumeSnapshot
Metadata:
  Creation Timestamp:  2022-07-18T07:40:55Z
  Finalizers:
    snapshot.storage.kubernetes.io/volumesnapshot-as-source-protection
    snapshot.storage.kubernetes.io/volumesnapshot-bound-protection
  Generation:  1
  Managed Fields:
    API Version:  snapshot.storage.k8s.io/v1
    Fields Type:  FieldsV1
    fieldsV1:
      f:metadata:
        f:annotations:
          .:
          f:kubectl.kubernetes.io/last-applied-configuration:
      f:spec:
        .:
        f:source:
          .:
          f:persistentVolumeClaimName:
        f:volumeSnapshotClassName:
    Manager:      kubectl-client-side-apply
    Operation:    Update
    Time:         2022-07-18T07:40:55Z
    API Version:  snapshot.storage.k8s.io/v1
    Fields Type:  FieldsV1
    fieldsV1:
      f:metadata:
        f:finalizers:
          v:"snapshot.storage.kubernetes.io/volumesnapshot-bound-protection":
      f:status:
        .:
        f:boundVolumeSnapshotContentName:
    Manager:      snapshot-controller
    Operation:    Update
    Time:         2022-07-18T07:41:21Z
    API Version:  snapshot.storage.k8s.io/v1
    Fields Type:  FieldsV1
    fieldsV1:
      f:status:
        f:creationTime:
        f:readyToUse:
        f:restoreSize:
    Manager:         snapshot-controller
    Operation:       Update
    Subresource:     status
    Time:            2022-07-18T07:41:25Z
  Resource Version:  8705
  UID:               aae1446f-e594-46b9-806f-d6775296f804
Spec:
  Source:
    Persistent Volume Claim Name:  csi-pvc
  Volume Snapshot Class Name:      csi-hostpath-snapclass
Status:
  Bound Volume Snapshot Content Name:  snapcontent-aae1446f-e594-46b9-806f-d6775296f804
  Creation Time:                       2022-07-18T07:41:18Z
  Ready To Use:                        true
  Restore Size:                        1Gi
Events:             <none>
 ```

 
 $ kubectl describe volumesnapshotcontent
 ```
 Name:         snapcontent-aae1446f-e594-46b9-806f-d6775296f804
Namespace:    
Labels:       <none>
Annotations:  <none>
API Version:  snapshot.storage.k8s.io/v1
Kind:         VolumeSnapshotContent
Metadata:
  Creation Timestamp:  2022-07-18T07:41:10Z
  Finalizers:
    snapshot.storage.kubernetes.io/volumesnapshotcontent-bound-protection
  Generation:  1
  Managed Fields:
    API Version:  snapshot.storage.k8s.io/v1
    Fields Type:  FieldsV1
    fieldsV1:
      f:metadata:
        f:finalizers:
          .:
          v:"snapshot.storage.kubernetes.io/volumesnapshotcontent-bound-protection":
      f:spec:
        .:
        f:deletionPolicy:
        f:driver:
        f:source:
          .:
          f:volumeHandle:
        f:volumeSnapshotClassName:
        f:volumeSnapshotRef:
          .:
          f:apiVersion:
          f:kind:
          f:name:
          f:namespace:
          f:resourceVersion:
          f:uid:
    Manager:      snapshot-controller
    Operation:    Update
    Time:         2022-07-18T07:41:11Z
    API Version:  snapshot.storage.k8s.io/v1
    Fields Type:  FieldsV1
    fieldsV1:
      f:status:
        .:
        f:creationTime:
        f:readyToUse:
        f:restoreSize:
        f:snapshotHandle:
    Manager:         csi-snapshotter
    Operation:       Update
    Subresource:     status
    Time:            2022-07-18T07:41:45Z
  Resource Version:  8720
  UID:               55d8b9c7-b263-4440-948c-86e18b6074a7
Spec:
  Deletion Policy:  Delete
  Driver:           hostpath.csi.k8s.io
  Source:
    Volume Handle:             b1dbb381-065b-11ed-8de2-826a39489514
  Volume Snapshot Class Name:  csi-hostpath-snapclass
  Volume Snapshot Ref:
    API Version:       snapshot.storage.k8s.io/v1
    Kind:              VolumeSnapshot
    Name:              new-snapshot-v1-demo
    Namespace:         default
    Resource Version:  8673
    UID:               aae1446f-e594-46b9-806f-d6775296f804
Status:
  Creation Time:    1658130078827905238
  Ready To Use:     true
  Restore Size:     1073741824
  Snapshot Handle:  02e776f3-066d-11ed-8de2-826a39489514
Events:             <none>
> ```

## Restore volume from snapshot support

Follow the following example to create a volume from a volume snapshot:
Note that as the PVC size goes larger, the restore can be slower.

$ kubectl apply -f examples/csi-restore.yaml
`persistentvolumeclaim/hpvc-restore created`

$ kubectl get pvc
```
NAME           STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS      AGE
csi-pvc        Bound    pvc-f2d87c5d-23de-433a-acb2-96649cbb75af   1Gi        RWO            csi-hostpath-sc   150m
hpvc-restore   Bound    pvc-31938eff-a107-46ed-9d1d-c40634ac4435   1Gi        RWO            csi-hostpath-sc   6s
 ```

$ kubectl get pv
 ```
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                  STORAGECLASS      REASON   AGE
pvc-31938eff-a107-46ed-9d1d-c40634ac4435   1Gi        RWO            Delete           Bound    default/hpvc-restore   csi-hostpath-sc            38s
pvc-f2d87c5d-23de-433a-acb2-96649cbb75af   1Gi        RWO            Delete           Bound    default/csi-pvc        csi-hostpath-sc            150m
 ```

