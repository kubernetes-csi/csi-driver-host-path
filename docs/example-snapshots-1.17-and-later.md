## Snapshot support

Volume snapshot was introduced as an alpha feature in Kubernetes v1.12 and it was promoted to beta in Kubernetes v1.17.

### Pre-requisites 
- For Kubernetes 1.17+, you must install the [VolumeSnapshot beta CRDs and the Snapshot Controller](deploy-1.17-and-later.md)

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
> Name:             csi-hostpath-snapclass
> Namespace:        
> Labels:           <none>
> Annotations:      kubectl.kubernetes.io/last-applied-configuration:
>                     {"apiVersion":"snapshot.storage.k8s.io/v1beta1","deletionPolicy":"Delete","driver":"hostpath.csi.k8s.io","kind":"VolumeSnapshotClass","met...
> API Version:      snapshot.storage.k8s.io/v1beta1
> Deletion Policy:  Delete
> Driver:           hostpath.csi.k8s.io
> Kind:             VolumeSnapshotClass
> Metadata:
>   Creation Timestamp:  2020-03-09T20:53:32Z
>   Generation:          1
>   Resource Version:    938
>   Self Link:           /apis/snapshot.storage.k8s.io/v1beta1/volumesnapshotclasses/csi-hostpath-snapclass
>   UID:                 8d2320cb-85fc-4908-9895-5ff8867169e2
> Events:                <none>
> ```

After having created the `csi-pvc` as described in the deployment validation,
use the volume snapshot class to dynamically create a volume snapshot:
>  - `$ kubectl apply -f examples/csi-snapshot-v1beta1.yaml`
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
> snapcontent-1b461d4e-6279-4f1d-9910-61d35d80c888   14s
> ```
>
> $ kubectl describe volumesnapshot
> ```
> Name:         new-snapshot-demo
> Namespace:    default
> Labels:       <none>
> Annotations:  kubectl.kubernetes.io/last-applied-configuration:
>                 {"apiVersion":"snapshot.storage.k8s.io/v1beta1","kind":"VolumeSnapshot","metadata":{"annotations":{},"name":"new-snapshot-demo","namespace...
> API Version:  snapshot.storage.k8s.io/v1beta1
> Kind:         VolumeSnapshot
> Metadata:
>   Creation Timestamp:  2020-03-09T21:45:04Z
>   Finalizers:
>     snapshot.storage.kubernetes.io/volumesnapshot-as-source-protection
>     snapshot.storage.kubernetes.io/volumesnapshot-bound-protection
>   Generation:        1
>   Resource Version:  11146
>   Self Link:         /apis/snapshot.storage.k8s.io/v1beta1/namespaces/default/volumesnapshots/new-snapshot-demo
>   UID:               1b461d4e-6279-4f1d-9910-61d35d80c888
> Spec:
>   Source:
>     Persistent Volume Claim Name:  csi-pvc
>   Volume Snapshot Class Name:      csi-hostpath-snapclass
> Status:
>   Bound Volume Snapshot Content Name:  snapcontent-1b461d4e-6279-4f1d-9910-61d35d80c888
>   Creation Time:                       2020-03-09T21:45:04Z
>   Ready To Use:                        true
>   Restore Size:                        1Gi
> Events:                                <none>
> ```
>
> 
> $ kubectl describe volumesnapshotcontent
> ```
> Name:         snapcontent-1b461d4e-6279-4f1d-9910-61d35d80c888
> Namespace:    
> Labels:       <none>
> Annotations:  <none>
> API Version:  snapshot.storage.k8s.io/v1beta1
> Kind:         VolumeSnapshotContent
> Metadata:
>   Creation Timestamp:  2020-03-09T21:45:04Z
>   Finalizers:
>     snapshot.storage.kubernetes.io/volumesnapshotcontent-bound-protection
>   Generation:        1
>   Resource Version:  11145
>   Self Link:         /apis/snapshot.storage.k8s.io/v1beta1/volumesnapshotcontents/snapcontent-1b461d4e-6279-4f1d-9910-61d35d80c888
>   UID:               665657cd-4461-476c-9cdb-5c0490c58945
> Spec:
>   Deletion Policy:  Delete
>   Driver:           hostpath.csi.k8s.io
>   Source:
>     Volume Handle:             42bdc1e0-624e-11ea-beee-42d40678b2d1
>   Volume Snapshot Class Name:  csi-hostpath-snapclass
>   Volume Snapshot Ref:
>     API Version:       snapshot.storage.k8s.io/v1beta1
>     Kind:              VolumeSnapshot
>     Name:              new-snapshot-demo
>     Namespace:         default
>     Resource Version:  11136
>     UID:               1b461d4e-6279-4f1d-9910-61d35d80c888
> Status:
>   Creation Time:    1583790304342000422
>   Ready To Use:     true
>   Restore Size:     1073741824
>   Snapshot Handle:  3c651edc-624f-11ea-beee-42d40678b2d1
> Events:             <none>
> ```

## Restore volume from snapshot support

Follow the following example to create a volume from a volume snapshot:

> $ kubectl apply -f examples/csi-restore.yaml
> `persistentvolumeclaim/hpvc-restore created`
>
> $ kubectl get pvc
> ```
> NAME           STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS      AGE
> csi-pvc        Bound    pvc-ad827273-8d08-430b-9d5a-e60e05a2bc3e   1Gi        RWO            csi-hostpath-sc   31m
> hpvc-restore   Bound    pvc-6d79a775-09f0-4bd9-968d-05c38d189bc4   1Gi        RWO            csi-hostpath-sc   23s
> ```
>
> $ kubectl get pv
> ```
> NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM                  STORAGECLASS      REASON   AGE
> pvc-6d79a775-09f0-4bd9-968d-05c38d189bc4   1Gi        RWO            Delete           Bound    default/hpvc-restore   csi-hostpath-sc            55s
> pvc-ad827273-8d08-430b-9d5a-e60e05a2bc3e   1Gi        RWO            Delete           Bound    default/csi-pvc        csi-hostpath-sc            31m
> ```

