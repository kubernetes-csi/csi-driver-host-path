## Deployment
The easiest way to test the Hostpath driver is to run the `deploy-hostpath.sh` script for the Kubernetes version used by
the cluster as shown below for Kubernetes 1.16. This creates the deployment that is maintained specifically for that
release of Kubernetes. However, other deployments may also work.

```shell
$ deploy/kubernetes-1.16/deploy-hostpath.sh
```

You should see an output similar to the following printed on the terminal showing the application of rbac rules and the result of deploying the hostpath driver, external provisioner, external attacher and snapshotter components:

```shell
applying RBAC rules
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-provisioner/v1.4.0/deploy/kubernetes/rbac.yaml
serviceaccount/csi-provisioner created
clusterrole.rbac.authorization.k8s.io/external-provisioner-runner created
clusterrolebinding.rbac.authorization.k8s.io/csi-provisioner-role created
role.rbac.authorization.k8s.io/external-provisioner-cfg created
rolebinding.rbac.authorization.k8s.io/csi-provisioner-role-cfg created
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-attacher/v2.0.0/deploy/kubernetes/rbac.yaml
serviceaccount/csi-attacher created
clusterrole.rbac.authorization.k8s.io/external-attacher-runner created
clusterrolebinding.rbac.authorization.k8s.io/csi-attacher-role created
role.rbac.authorization.k8s.io/external-attacher-cfg created
rolebinding.rbac.authorization.k8s.io/csi-attacher-role-cfg created
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/v1.2.0/deploy/kubernetes/rbac.yaml
serviceaccount/csi-snapshotter created
clusterrole.rbac.authorization.k8s.io/external-snapshotter-runner created
clusterrolebinding.rbac.authorization.k8s.io/csi-snapshotter-role created
role.rbac.authorization.k8s.io/external-snapshotter-leaderelection created
rolebinding.rbac.authorization.k8s.io/external-snapshotter-leaderelection created
kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-resizer/v0.3.0/deploy/kubernetes/rbac.yaml
serviceaccount/csi-resizer created
clusterrole.rbac.authorization.k8s.io/external-resizer-runner created
clusterrolebinding.rbac.authorization.k8s.io/csi-resizer-role created
role.rbac.authorization.k8s.io/external-resizer-cfg created
rolebinding.rbac.authorization.k8s.io/csi-resizer-role-cfg created
deploying hostpath components
   ./hostpath/csi-hostpath-attacher.yaml
        using           image: quay.io/k8scsi/csi-attacher:v2.0.0
service/csi-hostpath-attacher created
statefulset.apps/csi-hostpath-attacher created
   ./hostpath/csi-hostpath-driverinfo.yaml
csidriver.storage.k8s.io/hostpath.csi.k8s.io created
   ./hostpath/csi-hostpath-plugin.yaml
        using           image: quay.io/k8scsi/csi-node-driver-registrar:v1.2.0
        using           image: quay.io/k8scsi/hostpathplugin:v1.3.0
        using           image: quay.io/k8scsi/livenessprobe:v1.1.0
service/csi-hostpathplugin created
statefulset.apps/csi-hostpathplugin created
   ./hostpath/csi-hostpath-provisioner.yaml
        using           image: quay.io/k8scsi/csi-provisioner:v1.4.0
service/csi-hostpath-provisioner created
statefulset.apps/csi-hostpath-provisioner created
   ./hostpath/csi-hostpath-resizer.yaml
        using           image: quay.io/k8scsi/csi-resizer:v0.3.0
service/csi-hostpath-resizer created
statefulset.apps/csi-hostpath-resizer created
   ./hostpath/csi-hostpath-snapshotter.yaml
        using           image: quay.io/k8scsi/csi-snapshotter:v1.2.0
service/csi-hostpath-snapshotter created
statefulset.apps/csi-hostpath-snapshotter created
   ./hostpath/csi-hostpath-testing.yaml
        using           image: alpine/socat:1.0.3
service/hostpath-service created
statefulset.apps/csi-hostpath-socat created
12:08:01 waiting for hostpath deployment to complete, attempt #0
12:08:11 waiting for hostpath deployment to complete, attempt #1
deploying snapshotclass based on snapshotter version
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
The Hostpath driver is configured to create new volumes under `/csi-data-dir` inside the hostpath container that is specified in the plugin StatefulSet found [here](../deploy/kubernetes-latest/hostpath/csi-hostpath-plugin.yaml).  This path persist as long as the StatefulSet pod is up and running.

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
