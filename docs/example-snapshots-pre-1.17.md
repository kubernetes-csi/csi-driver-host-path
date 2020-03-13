## Snapshot support

Volume snapshot was introduced as an alpha feature in Kubernetes v1.12 and it was promoted to beta in Kubernetes v1.17.

### Pre-requisites 
- Enable feature gate called `VolumeSnapshotDataSource` in Kubernetes

### Creating and restoring volume snapshots
Ensure your volumesnapshotclass was created during hostpath deployment:
>
> $ kubectl get volumesnapshotclass
> ```
> NAME                     AGE
> csi-hostpath-snapclass   11s
> ```
>
> $ kubectl describe volumesnapshotclass
> ```
> Name:         csi-hostpath-snapclass
> Namespace:
> Labels:       <none>
> Annotations:  <none>
> API Version:  snapshot.storage.k8s.io/v1alpha1
> Kind:         VolumeSnapshotClass
> Metadata:
>   Creation Timestamp:  2018-10-03T14:15:30Z
>   Generation:          1
>   Resource Version:    2418
>   Self Link:           /apis/snapshot.storage.k8s.io/v1alpha1/volumesnapshotclasses/csi-hostpath-snapclass
>   UID:                 c8f5bc47-c716-11e8-8911-000c2967769a
> Snapshotter:           hostpath.csi.k8s.io
> Events:                <none>
> ```

After having created the `csi-pvc` as described in the example above,
use the volume snapshot class to dynamically create a volume snapshot:
>  - `$ kubectl apply -f examples/csi-snapshot.yaml`
> ```
> volumesnapshot.snapshot.storage.k8s.io/new-snapshot-demo created
> ```

>
> $ kubectl get volumesnapshot
> ```
> NAME                AGE
> new-snapshot-demo   12s
> ```
>
> $ kubectl get volumesnapshotcontent
>```
> NAME                                               AGE
> snapcontent-f55db632-c716-11e8-8911-000c2967769a   14s
> ```
>
> $ kubectl describe volumesnapshot
> ```
> Name:         new-snapshot-demo
> Namespace:    default
> Labels:       <none>
> Annotations:  <none>
> API Version:  snapshot.storage.k8s.io/v1alpha1
> Kind:         VolumeSnapshot
> Metadata:
>   Creation Timestamp:  2018-10-03T14:16:45Z
>   Generation:          1
>   Resource Version:    2476
>   Self Link:           /apis/snapshot.storage.k8s.io/v1alpha1/namespaces/default/volumesnapshots/new-snapshot-demo
>   UID:                 f55db632-c716-11e8-8911-000c2967769a
> Spec:
>   Snapshot Class Name:    csi-hostpath-snapclass
>   Snapshot Content Name:  snapcontent-f55db632-c716-11e8-8911-000c2967769a
>   Source:
>     API Group:  <nil>
>     Kind:  PersistentVolumeClaim
>     Name:  csi-pvc
> Status:
>   Creation Time:  2018-10-03T14:16:45Z
>   Ready:          true
>   Restore Size:   1Gi
> Events:           <none>
> ```
>
> $ kubectl describe volumesnapshotcontent
> ```
> Name:         snapcontent-f55db632-c716-11e8-8911-000c2967769a
> Namespace:
> Labels:       <none>
> Annotations:  <none>
> API Version:  snapshot.storage.k8s.io/v1alpha1
> Kind:         VolumeSnapshotContent
> Metadata:
>   Creation Timestamp:  2018-10-03T14:16:45Z
>   Generation:          1
>   Resource Version:    2474
>   Self Link:           /apis/snapshot.storage.k8s.io/v1alpha1/volumesnapshotcontents/snapcontent-f55db632-c716-11e8-8911-000c2967769a
>   UID:                 f561411f-c716-11e8-8911-000c2967769a
> Spec:
>   Csi Volume Snapshot Source:
>     Creation Time:    1538576205471577525
>     Driver:           hostpath.csi.k8s.io
>     Restore Size:     1073741824
>     Snapshot Handle:  f55ff979-c716-11e8-bb16-000c2967769a
>   Deletion Policy:    Delete
>   Persistent Volume Ref:
>     API Version:        v1
>     Kind:               PersistentVolume
>     Name:               pvc-0571cc14-c714-11e8-8911-000c2967769a
>     Resource Version:   1573
>     UID:                0575b966-c714-11e8-8911-000c2967769a
>   Snapshot Class Name:  csi-hostpath-snapclass
>   Volume Snapshot Ref:
>     API Version:       snapshot.storage.k8s.io/v1alpha1
>     Kind:              VolumeSnapshot
>     Name:              new-snapshot-demo
>     Namespace:         default
>     Resource Version:  2472
>     UID:               f55db632-c716-11e8-8911-000c2967769a
> Events:                <none>
> ```

## Restore volume from snapshot support

Follow the following example to create a volume from a volume snapshot:

> $ kubectl apply -f examples/csi-restore.yaml
> `persistentvolumeclaim/hpvc-restore created`
>
> $ kubectl get pvc
> ```
> NAME           STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS      AGE
> csi-pvc        Bound    pvc-0571cc14-c714-11e8-8911-000c2967769a   1Gi        RWO            csi-hostpath-sc   24m
> hpvc-restore   Bound    pvc-77324684-c717-11e8-8911-000c2967769a   1Gi        RWO            csi-hostpath-sc   6s
> ```
>
> $ kubectl get pv
> ```
> NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                  STORAGECLASS      REASON   AGE
> pvc-0571cc14-c714-11e8-8911-000c2967769a   1Gi        RWO            Delete           Bound    default/csi-pvc        csi-hostpath-sc            25m
> pvc-77324684-c717-11e8-8911-000c2967769a   1Gi        RWO            Delete           Bound    default/hpvc-restore   csi-hostpath-sc            33s
> ```

