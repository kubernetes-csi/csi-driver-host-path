# CSI Hostpath Driver

This repository hosts the CSI Hostpath driver and all of its build and dependent configuration files to deploy the driver.

## Pre-requisite
- Kubernetes cluster
- Running version 1.13 or later
- Access to terminal with `kubectl` installed

## Deployment
The easiest way to test the Hostpath driver is to run the `deploy-hostpath.sh` script for the Kubernetes version used by
the cluster as shown below for Kubernetes 1.13. This creates the deployment that is maintained specifically for that
release of Kubernetes. However, other deployments may also work. For details see the individual READMEs.

```shell
$ deploy/kubernetes-1.13/deploy-hostpath.sh
```

You should see an output similar to the following printed on the terminal showing the application of rbac rules and the result of deploying the hostpath driver, external provisioner, external attacher and snapshotter components:

```shell
applying RBAC rules
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-provisioner/v1.0.1/deploy/kubernetes/rbac.yaml
serviceaccount/csi-provisioner created
clusterrole.rbac.authorization.k8s.io/external-provisioner-runner created
clusterrolebinding.rbac.authorization.k8s.io/csi-provisioner-role created
role.rbac.authorization.k8s.io/external-provisioner-cfg created
rolebinding.rbac.authorization.k8s.io/csi-provisioner-role-cfg created
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-attacher/v1.0.1/deploy/kubernetes/rbac.yaml
serviceaccount/csi-attacher created
clusterrole.rbac.authorization.k8s.io/external-attacher-runner created
clusterrolebinding.rbac.authorization.k8s.io/csi-attacher-role created
role.rbac.authorization.k8s.io/external-attacher-cfg created
rolebinding.rbac.authorization.k8s.io/csi-attacher-role-cfg created
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/v1.0.1/deploy/kubernetes/rbac.yaml
serviceaccount/csi-snapshotter created
clusterrole.rbac.authorization.k8s.io/external-snapshotter-runner created
clusterrolebinding.rbac.authorization.k8s.io/csi-snapshotter-role created
deploying hostpath components
   deploy/kubernetes-1.13/hostpath/csi-hostpath-attacher.yaml
        using           image: quay.io/k8scsi/csi-attacher:v1.0.1
service/csi-hostpath-attacher created
statefulset.apps/csi-hostpath-attacher created
   deploy/kubernetes-1.13/hostpath/csi-hostpath-plugin.yaml
        using           image: quay.io/k8scsi/csi-node-driver-registrar:v1.0.2
        using           image: quay.io/k8scsi/hostpathplugin:v1.0.1
        using           image: quay.io/k8scsi/livenessprobe:v1.0.2
service/csi-hostpathplugin created
statefulset.apps/csi-hostpathplugin created
   deploy/kubernetes-1.13/hostpath/csi-hostpath-provisioner.yaml
        using           image: quay.io/k8scsi/csi-provisioner:v1.0.1
service/csi-hostpath-provisioner created
statefulset.apps/csi-hostpath-provisioner created
   deploy/kubernetes-1.13/hostpath/csi-hostpath-snapshotter.yaml
        using           image: quay.io/k8scsi/csi-snapshotter:v1.0.1
service/csi-hostpath-snapshotter created
statefulset.apps/csi-hostpath-snapshotter created
   deploy/kubernetes-1.13/hostpath/csi-hostpath-testing.yaml
        using           image: alpine/socat:1.0.3
service/hostpath-service created
statefulset.apps/csi-hostpath-socat created
23:16:10 waiting for hostpath deployment to complete, attempt #0
deploying snapshotclass
volumesnapshotclass.snapshot.storage.k8s.io/csi-hostpath-snapclass created
```

The [livenessprobe side-container](https://github.com/kubernetes-csi/livenessprobe) provided by the CSI community is deployed with the CSI driver to provide the liveness checking of the CSI services.

## Run example application and validate

Next, validate the deployment.  First, ensure all expected pods are running properly including the external attacher, provisioner, snapshotter and the actual hostpath driver plugin:

```shell
$ kubectl get pods
NAME                         READY   STATUS    RESTARTS   AGE
csi-hostpath-attacher-0      1/1     Running   0          5m47s
csi-hostpath-provisioner-0   1/1     Running   0          5m47s
csi-hostpath-snapshotter-0   1/1     Running   0          5m47s
csi-hostpathplugin-0         2/2     Running   0          5m45s
```

From the root directory, deploy the application pods including a storage class, a PVC, and a pod which mounts a volume using the Hostpath driver found in directory `./examples`:

```shell
$ for i in ./examples/csi-storageclass.yaml ./examples/csi-pvc.yaml ./examples/csi-app.yaml; do kubectl apply -f $i; done
pod/my-csi-app created
persistentvolumeclaim/csi-pvc created
storageclass.storage.k8s.io/csi-hostpath-sc created
```

Let's validate the components are deployed:

```shell
$ kubectl get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM             STORAGECLASS      REASON   AGE
pvc-58d5ec38-03e5-11e9-be51-000c29e88ff1   1Gi        RWO            Delete           Bound    default/csi-pvc   csi-hostpath-sc            80s

$ kubectl get pvc
NAME      STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS      AGE
csi-pvc   Bound    pvc-58d5ec38-03e5-11e9-be51-000c29e88ff1   1Gi        RWO            csi-hostpath-sc   93s
```

Finally, inspect the application pod `my-csi-app`  which mounts a Hostpath volume:

```shell
$ kubectl describe pods/my-csi-app
Name:               my-csi-app
Namespace:          default
Priority:           0
PriorityClassName:  <none>
Node:               127.0.0.1/127.0.0.1
Start Time:         Wed, 19 Dec 2018 18:25:29 -0500
Labels:             <none>
Annotations:        <none>
Status:             Running
IP:                 172.17.0.5
Containers:
  my-frontend:
    Container ID:  docker://927dc537fd14704794e1167b75a5aa040eb86eff76e155672be65c5cf9bda798
    Image:         busybox
    Image ID:      docker-pullable://busybox@sha256:2a03a6059f21e150ae84b0973863609494aad70f0a80eaeb64bddd8d92465812
    Port:          <none>
    Host Port:     <none>
    Command:
      sleep
      1000000
    State:          Running
      Started:      Wed, 19 Dec 2018 18:25:33 -0500
    Ready:          True
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /data from my-csi-volume (rw)
      /var/run/secrets/kubernetes.io/serviceaccount from default-token-wm562 (ro)
Conditions:
  Type              Status
  Initialized       True
  Ready             True
  ContainersReady   True
  PodScheduled      True
Volumes:
  my-csi-volume:
    Type:       PersistentVolumeClaim (a reference to a PersistentVolumeClaim in the same namespace)
    ClaimName:  csi-pvc
    ReadOnly:   false
  default-token-wm562:
    Type:        Secret (a volume populated by a Secret)
    SecretName:  default-token-wm562
    Optional:    false
QoS Class:       BestEffort
Node-Selectors:  <none>
Tolerations:     node.kubernetes.io/not-ready:NoExecute for 300s
                 node.kubernetes.io/unreachable:NoExecute for 300s
Events:          <none>
```

## Confirm Hostpath driver works
The Hostpath driver is configured to create new volumes under `/tmp` inside the hostpath container that is specified in the plugin DaemonSet found [here](./deploy/hostpath/csi-hostpath-plugin.yaml).  This path persist as long as the DaemonSet pod is up and running. 

A file written in a properly mounted Hostpath volume inside an application should show up inside the Hostpath container.  The following steps confirms that Hostpath is working properly.  First, create a file from the application pod as shown:

```shell
$ kubectl exec -it my-csi-app /bin/sh
/ # touch /data/hello-world
/ # exit
```

Next, ssh into the Hostpath container and verify that the file shows up there:
```shell
$ kubectl exec -it $(kubectl get pods --selector app=csi-hostpathplugin -o jsonpath='{.items[*].metadata.name}') -c hostpath /bin/sh

```
Then, use the following command to locate the file. If everything works OK you should get a result similar to the following:

```shell
/ # find / -name hello-world
/tmp/057485ab-c714-11e8-bb16-000c2967769a/hello-world
/ # exit
```

## Confirm the creation of the VolumeAttachment object
An additional way to ensure the driver is working properly is by inspecting the VolumeAttachment API object created that represents the attached volume:

```shell
$ kubectl describe volumeattachment
Name:         csi-a7515d53b30a1193fd70b822b18181cff1d16422fd922692bce5ea234cb191e9
Namespace:
Labels:       <none>
Annotations:  <none>
API Version:  storage.k8s.io/v1
Kind:         VolumeAttachment
Metadata:
  Creation Timestamp:  2018-12-19T23:25:29Z
  Resource Version:    533
  Self Link:           /apis/storage.k8s.io/v1/volumeattachments/csi-a7515d53b30a1193fd70b822b18181cff1d16422fd922692bce5ea234cb191e9
  UID:                 5fb4874f-03e5-11e9-be51-000c29e88ff1
Spec:
  Attacher:   csi-hostpath
  Node Name:  127.0.0.1
  Source:
    Persistent Volume Name:  pvc-58d5ec38-03e5-11e9-be51-000c29e88ff1
Status:
  Attached:  true
Events:      <none>
```


## Snapshot support

Since volume snapshot is an alpha feature starting in Kubernetes v1.12, you need to enable feature gate called `VolumeSnapshotDataSource` in the Kubernetes.

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
> Snapshotter:           csi-hostpath
> Events:                <none>
> ```

After having created the `csi-pvc` as described in the example above,
use the volume snapshot class to dynamically create a volume snapshot:

> $ kubectl apply -f examples/csi-snapshot.yaml
> ```
> volumesnapshot.snapshot.storage.k8s.io/new-snapshot-demo created
> ```
>
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
>     Driver:           csi-hostpath
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

## Inline ephemeral support
As of version 1.15 of Kubernetes, the CSI Hostpath driver (starting with version 1.0.1) now includes support for inline ephemeral volume.  This means that a volume can be specified directly inside a pod spec without the need to use a persistent volume object.
Find out how to enable or create a CSI inline driver [here](https://kubernetes-csi.github.io/docs/ephemeral-local-volumes.html)

To test this feature, redeploy the CSI Hostpath plugin YAML by updating the `hostpath` container to use  the inline ephemeral mode by setting the `ephemeral` flag, of the driver binary, to true as shown in the following setup:

```yaml
kind: DaemonSet
apiVersion: apps/v1
metadata:
  name: csi-hostpathplugin
spec:
...
  template:
    spec:
      containers:
        - name: hostpath
          image: image: quay.io/k8scsi/hostpathplugin:v1.2.0
          args:
            - "--v=5"
            - "--endpoint=$(CSI_ENDPOINT)"
            - "--nodeid=$(KUBE_NODE_NAME)"
            - "--ephemeral=true"      
...

```
Notice the addition of the `ephemeral=true` flag used in the `args:` block in the previous snippet.

Once the driver plugin has been deployed, it can be tested by deploying a simple pod which has an inline volume specified in the spec:

```yaml
kind: Pod
apiVersion: v1
metadata:
  name: my-csi-app
spec:
  containers:
    - name: my-frontend
      image: busybox
      volumeMounts:
      - mountPath: "/data"
        name: my-csi-volume
      command: [ "sleep", "1000000" ]
  volumes:
    - name: my-csi-volume
      csi:
        driver: csi-hostpath
``` 

> See sample YAML file [here](./examples/csi-app-inline.yaml).

Notice the CSI driver is now specified directly in the container spec inside the `volumes:` block.  You can use the [same steps as above][Confirm Hostpath driver works] 
to verify that the volume has been created and deleted (when the pod is removed).


## Building the binaries
If you want to build the driver yourself, you can do so with the following command from the root directory:

```shell
make hostpath
```


## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

You can reach the maintainers of this project at:

- [Slack](http://slack.k8s.io/)
- [Mailing List](https://groups.google.com/forum/#!forum/kubernetes-dev)

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).

[owners]: https://git.k8s.io/community/contributors/guide/owners.md
[Creative Commons 4.0]: https://git.k8s.io/website/LICENSE
