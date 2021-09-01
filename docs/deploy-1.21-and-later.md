## Cluster setup
For Kubernetes 1.21+, some initial cluster setup is required to install the following:
- CSI VolumeSnapshot beta CRDs (custom resource definitions)
- Snapshot Controller

### Check if cluster components are already installed
Run the following commands to ensure the VolumeSnapshot CRDs have been installed:
```
$ kubectl get volumesnapshotclasses.snapshot.storage.k8s.io 
$ kubectl get volumesnapshots.snapshot.storage.k8s.io 
$ kubectl get volumesnapshotcontents.snapshot.storage.k8s.io
```
If any of these commands return the following error message, you must install the corresponding CRD:
```
error: the server doesn't have a resource type "volumesnapshotclasses"
```

Next, check if any pods are running the snapshot-controller image:
```
$ kubectl get pods --all-namespaces -o=jsonpath='{range .items[*]}{"\n"}{range .spec.containers[*]}{.image}{", "}{end}{end}' | grep snapshot-controller
k8s.gcr.io/sig-storage/snapshot-controller:v4.2.0, 
```

If no pods are running the snapshot-controller, follow the instructions below to create the snapshot-controller

__Note:__ The above command may not work for clusters running on managed k8s services. In this case, the presence of all VolumeSnapshot CRDs is an indicator that your cluster is ready for hostpath deployment.

### VolumeSnapshot CRDs and snapshot controller installation
Run the following commands to install these components: 
```shell
# Change to the latest supported snapshotter version
$ SNAPSHOTTER_VERSION=v4.0.1

# Apply VolumeSnapshot CRDs
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${SNAPSHOTTER_VERSION}/client/config/crd/snapshot.storage.k8s.io_volumesnapshotclasses.yaml
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${SNAPSHOTTER_VERSION}/client/config/crd/snapshot.storage.k8s.io_volumesnapshotcontents.yaml
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${SNAPSHOTTER_VERSION}/client/config/crd/snapshot.storage.k8s.io_volumesnapshots.yaml

# Create snapshot controller
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${SNAPSHOTTER_VERSION}/deploy/kubernetes/snapshot-controller/rbac-snapshot-controller.yaml
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${SNAPSHOTTER_VERSION}/deploy/kubernetes/snapshot-controller/setup-snapshot-controller.yaml
```

## Deployment
The easiest way to test the Hostpath driver is to run the `deploy.sh` script for the Kubernetes version used by
the cluster as shown below for Kubernetes 1.20. This creates the deployment that is maintained specifically for that
release of Kubernetes. However, other deployments may also work.

```
# deploy hostpath driver
$ deploy/kubernetes-latest/deploy.sh
```

You should see an output similar to the following printed on the terminal showing the application of rbac rules 
Note that the following output is from Kubernetes 1.20:

```shell
applying RBAC rules
curl https://raw.githubusercontent.com/kubernetes-csi/external-provisioner/v3.1.0/deploy/kubernetes/rbac.yaml --output /tmp/tmp.qXBMCAiF8m/rbac.yaml --silent --location
kubectl apply --kustomize /tmp/tmp.qXBMCAiF8m
serviceaccount/csi-provisioner unchanged
role.rbac.authorization.k8s.io/external-provisioner-cfg unchanged
clusterrole.rbac.authorization.k8s.io/external-provisioner-runner unchanged
rolebinding.rbac.authorization.k8s.io/csi-provisioner-role-cfg unchanged
clusterrolebinding.rbac.authorization.k8s.io/csi-provisioner-role unchanged
curl https://raw.githubusercontent.com/kubernetes-csi/external-attacher/v3.4.0/deploy/kubernetes/rbac.yaml --output /tmp/tmp.qXBMCAiF8m/rbac.yaml --silent --location
kubectl apply --kustomize /tmp/tmp.qXBMCAiF8m
serviceaccount/csi-attacher unchanged
role.rbac.authorization.k8s.io/external-attacher-cfg created
clusterrole.rbac.authorization.k8s.io/external-attacher-runner created
rolebinding.rbac.authorization.k8s.io/csi-attacher-role-cfg created
clusterrolebinding.rbac.authorization.k8s.io/csi-attacher-role created
curl https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/v5.0.1/deploy/kubernetes/csi-snapshotter/rbac-csi-snapshotter.yaml --output /tmp/tmp.qXBMCAiF8m/rbac.yaml --silent --location
kubectl apply --kustomize /tmp/tmp.qXBMCAiF8m
serviceaccount/csi-snapshotter created
role.rbac.authorization.k8s.io/external-snapshotter-leaderelection created
clusterrole.rbac.authorization.k8s.io/external-snapshotter-runner created
rolebinding.rbac.authorization.k8s.io/external-snapshotter-leaderelection created
clusterrolebinding.rbac.authorization.k8s.io/csi-snapshotter-role created
curl https://raw.githubusercontent.com/kubernetes-csi/external-resizer/v1.4.0/deploy/kubernetes/rbac.yaml --output /tmp/tmp.qXBMCAiF8m/rbac.yaml --silent --location
kubectl apply --kustomize /tmp/tmp.qXBMCAiF8m
serviceaccount/csi-resizer created
role.rbac.authorization.k8s.io/external-resizer-cfg created
clusterrole.rbac.authorization.k8s.io/external-resizer-runner created
rolebinding.rbac.authorization.k8s.io/csi-resizer-role-cfg created
clusterrolebinding.rbac.authorization.k8s.io/csi-resizer-role created
curl https://raw.githubusercontent.com/kubernetes-csi/external-health-monitor/v0.4.0/deploy/kubernetes/external-health-monitor-controller/rbac.yaml --output /tmp/tmp.qXBMCAiF8m/rbac.yaml --silent --location
kubectl apply --kustomize /tmp/tmp.qXBMCAiF8m
serviceaccount/csi-external-health-monitor-controller created
role.rbac.authorization.k8s.io/external-health-monitor-controller-cfg created
clusterrole.rbac.authorization.k8s.io/external-health-monitor-controller-runner created
rolebinding.rbac.authorization.k8s.io/csi-external-health-monitor-controller-role-cfg created
clusterrolebinding.rbac.authorization.k8s.io/csi-external-health-monitor-controller-role created
deploying hostpath components
   /go/src/github.com/kubernetes-csi/csi-driver-host-path/deploy/kubernetes-latest/hostpath/csi-hostpath-driverinfo.yaml
csidriver.storage.k8s.io/hostpath.csi.k8s.io created
   /go/src/github.com/kubernetes-csi/csi-driver-host-path/deploy/kubernetes-latest/hostpath/csi-hostpath-plugin.yaml
        using           image: k8s.gcr.io/sig-storage/hostpathplugin:v1.7.3
        using           image: k8s.gcr.io/sig-storage/csi-external-health-monitor-controller:v0.4.0
        using           image: k8s.gcr.io/sig-storage/csi-node-driver-registrar:v2.5.0
        using           image: k8s.gcr.io/sig-storage/livenessprobe:v2.6.0
        using           image: k8s.gcr.io/sig-storage/csi-attacher:v3.4.0
        using           image: k8s.gcr.io/sig-storage/csi-provisioner:v3.1.0
        using           image: k8s.gcr.io/sig-storage/csi-resizer:v1.4.0
        using           image: k8s.gcr.io/sig-storage/csi-snapshotter:v5.0.1
serviceaccount/csi-hostpathplugin-sa created
clusterrolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-attacher-cluster-role created
clusterrolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-health-monitor-controller-cluster-role created
clusterrolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-provisioner-cluster-role created
clusterrolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-resizer-cluster-role created
clusterrolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-snapshotter-cluster-role created
rolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-attacher-role created
rolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-health-monitor-controller-role created
rolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-provisioner-role created
rolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-resizer-role created
rolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-snapshotter-role created
statefulset.apps/csi-hostpathplugin created
   /go/src/github.com/kubernetes-csi/csi-driver-host-path/deploy/kubernetes-latest/hostpath/csi-hostpath-snapshotclass.yaml
volumesnapshotclass.snapshot.storage.k8s.io/csi-hostpath-snapclass created
   /go/src/github.com/kubernetes-csi/csi-driver-host-path/deploy/kubernetes-latest/hostpath/csi-hostpath-testing.yaml
        using           image: alpine/socat:1.0.3
service/hostpath-service created
statefulset.apps/csi-hostpath-socat created
09:40:49 waiting for hostpath deployment to complete, attempt #0
09:40:59 waiting for hostpath deployment to complete, attempt #1
09:41:13 waiting for hostpath deployment to complete, attempt #2
09:41:23 waiting for hostpath deployment to complete, attempt #3
09:41:33 waiting for hostpath deployment to complete, attempt #4
```

The [livenessprobe side-container](https://github.com/kubernetes-csi/livenessprobe) provided by the CSI community is deployed with the CSI driver to provide the liveness checking of the CSI services.

```
# Check object status for snapshotclass based on snapshotter version
$  kubectl get volumesnapshotclass
NAME                     DRIVER                DELETIONPOLICY   AGE
csi-hostpath-snapclass   hostpath.csi.k8s.io   Delete           149m
```

## Deploying the hostpath driver, external provisioner, external attacher and snapshotter components

From the root directory, Deploy the hostpath driver, external provisioner, external attacher and snapshotter components

```shell
for i in ./deploy/kubernetes-1.21-test/hostpath/csi-hostpath-attacher.yaml ./deploy/kubernetes-1.21-test/hostpath/csi-hostpath-provisioner.yaml ./deploy/kubernetes-1.21-test/hostpath/csi-hostpath-resizer.yaml ./deploy/kubernetes-1.21-test/hostpath/csi-hostpath-snapshotter.yaml; do kubectl apply -f $i; done

statefulset.apps/csi-hostpath-attacher created
statefulset.apps/csi-hostpath-provisioner created
statefulset.apps/csi-hostpath-resizer created
statefulset.apps/csi-hostpath-snapshotter created
```
Next, validate the deployment.  First, ensure all expected pods are running properly including the external attacher, provisioner, snapshotter and the actual hostpath driver plugin:

```shell
$ kubectl get pods
NAME                         READY   STATUS    RESTARTS   AGE
NAME                         READY   STATUS    RESTARTS       AGE
csi-hostpath-attacher-0      1/1     Running   0              2m59s
csi-hostpath-provisioner-0   1/1     Running   0              2m59s
csi-hostpath-resizer-0       1/1     Running   0              2m59s
csi-hostpath-snapshotter-0   1/1     Running   0              2m59s
csi-hostpath-socat-0         1/1     Running   0              127m
csi-hostpathplugin-0         8/8     Running   19 (10m ago)   127m
snapshot-controller-0        1/1     Running   0              140m
```

From the root directory, deploy the application pods including a storage class, a PVC, and a pod which mounts a volume using the Hostpath driver found in directory `./examples`:

```shell
$ for i in ./examples/csi-storageclass.yaml ./examples/csi-pvc.yaml ./examples/csi-app.yaml; do kubectl apply -f $i; done
storageclass.storage.k8s.io/csi-hostpath-sc created
persistentvolumeclaim/csi-pvc created
pod/my-csi-app created
```

Let's validate the components are deployed:

```shell
$ kubectl get pv
NAME                                       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM             STORAGECLASS      REASON   AGE
pvc-f2d87c5d-23de-433a-acb2-96649cbb75af   1Gi        RWO            Delete           Bound    default/csi-pvc   csi-hostpath-sc            51m

$ kubectl get pvc
NAME      STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS      AGE
csi-pvc   Bound    pvc-f2d87c5d-23de-433a-acb2-96649cbb75af   1Gi        RWO            csi-hostpath-sc   52m
```

Finally, inspect the application pod `my-csi-app`  which mounts a Hostpath volume:

```shell
$ kubectl describe pods/my-csi-app
Name:         my-csi-app
Namespace:    default
Priority:     0
Node:         kind-control-plane/172.19.0.2
Start Time:   Mon, 18 Jul 2022 11:07:22 +0530
Labels:       <none>
Annotations:  <none>
Status:       Running
IP:           10.244.0.8
IPs:
  IP:  10.244.0.8
Containers:
  my-frontend:
    Container ID:  containerd://73397b2cc3fb276058066e59d116508d0adb04c9b8ce2cf2bb17f897e5bb96ac
    Image:         busybox
    Image ID:      docker.io/library/busybox@sha256:3614ca5eacf0a3a1bcc361c939202a974b4902b9334ff36eb29ffe9011aaad83
    Port:          <none>
    Host Port:     <none>
    Command:
      sleep
      1000000
    State:          Running
      Started:      Mon, 18 Jul 2022 11:07:54 +0530
    Ready:          True
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /data from my-csi-volume (rw)
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-6w2db (ro)
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
  kube-api-access-6w2db:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   BestEffort
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type    Reason                  Age   From                     Message
  ----    ------                  ----  ----                     -------
  Normal  Scheduled               54m   default-scheduler        Successfully assigned default/my-csi-app to kind-control-plane
  Normal  SuccessfulAttachVolume  54m   attachdetach-controller  AttachVolume.Attach succeeded for volume "pvc-f2d87c5d-23de-433a-acb2-96649cbb75af"
  Normal  Pulling                 53m   kubelet                  Pulling image "busybox"
  Normal  Pulled                  53m   kubelet                  Successfully pulled image "busybox" in 11.97964778s
  Normal  Created                 53m   kubelet                  Created container my-frontend
  Normal  Started                 53m   kubelet                  Started container my-frontend
```

## Confirm Hostpath driver works
The Hostpath driver is configured to create new volumes under `/csi-data-dir` inside the hostpath container that is specified in the plugin StatefulSet found [here](../deploy/kubernetes-1.21-test/hostpath/csi-hostpath-plugin.yaml).  This path persist as long as the StatefulSet pod is up and running.

A file written in a properly mounted Hostpath volume inside an application should show up inside the Hostpath container.  The following steps confirms that Hostpath is working properly.  First, create a file from the application pod as shown:

```shell
$ kubectl exec -it my-csi-app /bin/sh
/ # touch /data/hello-world
/ # exit
```

Next, ssh into the Hostpath container and verify that the file shows up there:
```shell
$  kubectl exec -it csi-hostpathplugin-0 -c hostpath -- /bin/sh
```
Then, use the following command to locate the file. If everything works OK you should get a result similar to the following:

```shell
/ # find / -name hello-world
/var/lib/kubelet/pods/745eeaf7-7e5a-42ec-a4b8-96ba1877d6cb/volumes/kubernetes.io~csi/pvc-f2d87c5d-23de-433a-acb2-96649cbb75af/mount/hello-world
/csi-data-dir/b1dbb381-065b-11ed-8de2-826a39489514/hello-world
/ # exit
```

## Confirm the creation of the VolumeAttachment object
An additional way to ensure the driver is working properly is by inspecting the VolumeAttachment API object created that represents the attached volume:

```shell
$ kubectl describe volumeattachment
Name:         csi-c364ec7dcdaa3c6c68f77b6b267a856e70eb8bf05064c287b5fc7b87469a422a
Namespace:    
Labels:       <none>
Annotations:  <none>
API Version:  storage.k8s.io/v1
Kind:         VolumeAttachment
Metadata:
  Creation Timestamp:  2022-07-18T05:37:22Z
  Managed Fields:
    API Version:  storage.k8s.io/v1
    Fields Type:  FieldsV1
    fieldsV1:
      f:spec:
        f:attacher:
        f:nodeName:
        f:source:
          f:persistentVolumeName:
    Manager:      kube-controller-manager
    Operation:    Update
    Time:         2022-07-18T05:37:22Z
    API Version:  storage.k8s.io/v1
    Fields Type:  FieldsV1
    fieldsV1:
      f:status:
        f:attached:
    Manager:         Go-http-client
    Operation:       Update
    Subresource:     status
    Time:            2022-07-18T05:37:23Z
  Resource Version:  3249
  UID:               3a50e9a9-4c23-465d-bd56-b4b91780990c
Spec:
  Attacher:   hostpath.csi.k8s.io
  Node Name:  kind-control-plane
  Source:
    Persistent Volume Name:  pvc-f2d87c5d-23de-433a-acb2-96649cbb75af
Status:
  Attached:  true
Events:      <none>
```

