# CSI Hostpath Driver

This repository hosts the CSI Hostpath driver and all of its build and dependent configuration files to deploy the driver.

## Pre-requisite
- Kubernetes cluster
- Running verrsion 1.13 or later
- Access to terminal with `kubectl` installed

## Deployment
The easiest way to test the Hostpath driver is to run `deploy/deploy-hostpath.sh` scrip as show:

```shell
$ sh deploy/deploy-hostpath.sh
```

You should see an output similar to the following printed on the terminal showing the application of rbac rules and the result of deploying the hostpath driver, external privisioner and external attacher components:

```shell
applying RBAC rules
serviceaccount/csi-provisioner created
clusterrole.rbac.authorization.k8s.io/external-provisioner-runner created
clusterrolebinding.rbac.authorization.k8s.io/csi-provisioner-role created
role.rbac.authorization.k8s.io/external-provisioner-cfg created
rolebinding.rbac.authorization.k8s.io/csi-provisioner-role-cfg created
serviceaccount/csi-attacher created
clusterrole.rbac.authorization.k8s.io/external-attacher-runner created
clusterrolebinding.rbac.authorization.k8s.io/csi-attacher-role created
role.rbac.authorization.k8s.io/external-attacher-cfg created
rolebinding.rbac.authorization.k8s.io/csi-attacher-role-cfg created
deploying hostpath components
service/csi-hostpath-attacher created
statefulset.apps/csi-hostpath-attacher created
statefulset.apps/csi-hostpathplugin created
service/csi-hostpath-provisioner created
statefulset.apps/csi-hostpath-provisioner created
```

The script can also install CRDs that are needed for alpha features,
but as this is something that should be done by the cluster
provisioning tool it is disabled in the script by default. For this
and other customizations see the source code of the deploy script.

## Run example application and validate

Next, validate the deployment.  First, ensure all expected pods are running properly including the external attacher, provisioner, and the actual hostpath driver plugin:

```shell
$ kubectl get pods
NAME                         READY   STATUS    RESTARTS   AGE
csi-hostpath-attacher-0      1/1     Running   0          5m47s
csi-hostpath-provisioner-0   1/1     Running   0          5m47s
csi-hostpathplugin-0         2/2     Running   0          5m45s
```

From the root directory, deploy the application pods including a storage class, a PVC, and a pod which mounts a volume using the Hostpath driver found in directory `./examples`:

```shell
$ kubectl create -f ./examples
pod/my-csi-app created
persistentvolumeclaim/csi-pvc created
storageclass.storage.k8s.io/csi-hostpath-sc created
```

Let's validate the components are deployed:

```shell
$> kubectl get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM             STORAGECLASS      REASON   AGE
pvc-58d5ec38-03e5-11e9-be51-000c29e88ff1   1Gi        RWO            Delete           Bound    default/csi-pvc   csi-hostpath-sc            80s

$> kubectl get pvc
NAME      STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS      AGE
csi-pvc   Bound    pvc-58d5ec38-03e5-11e9-be51-000c29e88ff1   1Gi        RWO            csi-hostpath-sc   93s
```

Finally, inspect the application pod `my-csi-app`  which mounts a Hostpath volume:

```shell
$> kubectl describe pods/my-csi-app
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
$> kubectl describe volumeattachment
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
