## Cluster setup
For Kubernetes 1.17+, some initial cluster setup is required to install the following:
- CSI VolumeSnapshot beta CRDs (custom resource definitions)
- Snapshot Controller

### Check if cluster components are already installed
Run the follow commands to ensure the VolumeSnapshot CRDs have been installed:
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
$ kubectl get pods --all-namespaces -o jsonpath="{range .items[*]}{range .spec.containers[*]}{.image}{'\n'}{end}{end}" | grep snapshot-controller
```

If no pods are running the snapshot-controller, follow the instructions below to create the snapshot-controller

__Note:__ The above command may not work for clusters running on managed k8s services. In this case, the presence of all VolumeSnapshot CRDs is an indicator that your cluster is ready for hostpath deployment.

### VolumeSnapshot CRDs and snapshot controller installation
Run the following commands to install these components: 
```shell
# Change to the latest supported snapshotter release branch
$ SNAPSHOTTER_BRANCH=release-6.3
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${SNAPSHOTTER_BRANCH}/client/config/crd/snapshot.storage.k8s.io_volumesnapshotclasses.yaml
customresourcedefinition.apiextensions.k8s.io/volumesnapshotclasses.snapshot.storage.k8s.io created
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${SNAPSHOTTER_BRANCH}/client/config/crd/snapshot.storage.k8s.io_volumesnapshotcontents.yaml
customresourcedefinition.apiextensions.k8s.io/volumesnapshotcontents.snapshot.storage.k8s.io created
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${SNAPSHOTTER_BRANCH}/client/config/crd/snapshot.storage.k8s.io_volumesnapshots.yaml
customresourcedefinition.apiextensions.k8s.io/volumesnapshots.snapshot.storage.k8s.io created

$ SNAPSHOTTER_VERSION=v6.3.3
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${SNAPSHOTTER_VERSION}/deploy/kubernetes/snapshot-controller/rbac-snapshot-controller.yaml
serviceaccount/snapshot-controller created
clusterrole.rbac.authorization.k8s.io/snapshot-controller-runner created
clusterrolebinding.rbac.authorization.k8s.io/snapshot-controller-role created
role.rbac.authorization.k8s.io/snapshot-controller-leaderelection created
rolebinding.rbac.authorization.k8s.io/snapshot-controller-leaderelection created
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${SNAPSHOTTER_VERSION}/deploy/kubernetes/snapshot-controller/setup-snapshot-controller.yaml
deployment.apps/snapshot-controller created
```

## Deployment
The simplest way to test the HostPath driver is by running the deploy.sh script corresponding to your cluster's Kubernetes version. 
For example, to deploy on the latest Kubernetes, use the following command:

```
# deploy hostpath driver
$ deploy/kubernetes-latest/deploy.sh
```

You should see an output similar to the following printed on the terminal showing the application of rbac rules and the
result of deploying the hostpath driver, external provisioner, external attacher and snapshotter components. 
Note that the following output is from Kubernetes 1.32.2:

```shell
csi-driver-host-path %  deploy/kubernetes-latest/deploy.sh
applying RBAC rules
curl https://raw.githubusercontent.com/kubernetes-csi/external-provisioner/v5.2.0/deploy/kubernetes/rbac.yaml --output /var/folders/42/l7fg3dk55xn7jm24ld4bpkyw0000gn/T/tmp.ZKWrXmZPJ9/rbac.yaml --silent --location
kubectl apply --kustomize /var/folders/42/l7fg3dk55xn7jm24ld4bpkyw0000gn/T/tmp.ZKWrXmZPJ9
serviceaccount/csi-provisioner created
role.rbac.authorization.k8s.io/external-provisioner-cfg created
clusterrole.rbac.authorization.k8s.io/external-provisioner-runner created
rolebinding.rbac.authorization.k8s.io/csi-provisioner-role-cfg created
clusterrolebinding.rbac.authorization.k8s.io/csi-provisioner-role created
curl https://raw.githubusercontent.com/kubernetes-csi/external-attacher/v4.8.0/deploy/kubernetes/rbac.yaml --output /var/folders/42/l7fg3dk55xn7jm24ld4bpkyw0000gn/T/tmp.ZKWrXmZPJ9/rbac.yaml --silent --location
kubectl apply --kustomize /var/folders/42/l7fg3dk55xn7jm24ld4bpkyw0000gn/T/tmp.ZKWrXmZPJ9
serviceaccount/csi-attacher created
role.rbac.authorization.k8s.io/external-attacher-cfg created
clusterrole.rbac.authorization.k8s.io/external-attacher-runner created
rolebinding.rbac.authorization.k8s.io/csi-attacher-role-cfg created
clusterrolebinding.rbac.authorization.k8s.io/csi-attacher-role created
curl https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/v8.2.0/deploy/kubernetes/csi-snapshotter/rbac-csi-snapshotter.yaml --output /var/folders/42/l7fg3dk55xn7jm24ld4bpkyw0000gn/T/tmp.ZKWrXmZPJ9/rbac.yaml --silent --location
kubectl apply --kustomize /var/folders/42/l7fg3dk55xn7jm24ld4bpkyw0000gn/T/tmp.ZKWrXmZPJ9
serviceaccount/csi-snapshotter created
role.rbac.authorization.k8s.io/external-snapshotter-leaderelection created
clusterrole.rbac.authorization.k8s.io/external-snapshotter-runner created
rolebinding.rbac.authorization.k8s.io/external-snapshotter-leaderelection created
clusterrolebinding.rbac.authorization.k8s.io/csi-snapshotter-role created
curl https://raw.githubusercontent.com/kubernetes-csi/external-resizer/v1.13.1/deploy/kubernetes/rbac.yaml --output /var/folders/42/l7fg3dk55xn7jm24ld4bpkyw0000gn/T/tmp.ZKWrXmZPJ9/rbac.yaml --silent --location
kubectl apply --kustomize /var/folders/42/l7fg3dk55xn7jm24ld4bpkyw0000gn/T/tmp.ZKWrXmZPJ9
serviceaccount/csi-resizer created
role.rbac.authorization.k8s.io/external-resizer-cfg created
clusterrole.rbac.authorization.k8s.io/external-resizer-runner created
rolebinding.rbac.authorization.k8s.io/csi-resizer-role-cfg created
clusterrolebinding.rbac.authorization.k8s.io/csi-resizer-role created
curl https://raw.githubusercontent.com/kubernetes-csi/external-health-monitor/v0.14.0/deploy/kubernetes/external-health-monitor-controller/rbac.yaml --output /var/folders/42/l7fg3dk55xn7jm24ld4bpkyw0000gn/T/tmp.ZKWrXmZPJ9/rbac.yaml --silent --location
kubectl apply --kustomize /var/folders/42/l7fg3dk55xn7jm24ld4bpkyw0000gn/T/tmp.ZKWrXmZPJ9
serviceaccount/csi-external-health-monitor-controller created
role.rbac.authorization.k8s.io/external-health-monitor-controller-cfg created
clusterrole.rbac.authorization.k8s.io/external-health-monitor-controller-runner created
rolebinding.rbac.authorization.k8s.io/csi-external-health-monitor-controller-role-cfg created
clusterrolebinding.rbac.authorization.k8s.io/csi-external-health-monitor-controller-role created
deploying hostpath components
   /Users/csi-driver-host-path/deploy/kubernetes-latest/hostpath/csi-hostpath-driverinfo.yaml
csidriver.storage.k8s.io/hostpath.csi.k8s.io created
   /Users/csi-driver-host-path/deploy/kubernetes-latest/hostpath/csi-hostpath-plugin.yaml
        using           image: registry.k8s.io/sig-storage/hostpathplugin:v1.15.0
        using           image: registry.k8s.io/sig-storage/csi-external-health-monitor-controller:v0.14.0
        using           image: registry.k8s.io/sig-storage/csi-node-driver-registrar:v2.13.0
        using           image: registry.k8s.io/sig-storage/livenessprobe:v2.15.0
        using           image: registry.k8s.io/sig-storage/csi-attacher:v4.8.0
        using           image: registry.k8s.io/sig-storage/csi-provisioner:v5.2.0
        using           image: registry.k8s.io/sig-storage/csi-resizer:v1.13.1
        using           image: registry.k8s.io/sig-storage/csi-snapshotter:v8.2.0
serviceaccount/csi-hostpathplugin-sa created
clusterrolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-attacher-cluster-role created
clusterrolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-health-monitor-controller-cluster-role created
clusterrolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-provisioner-cluster-role created
clusterrolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-resizer-cluster-role created
clusterrolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-snapshotter-cluster-role created
clusterrolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-snapshot-metadata-cluster-role created
rolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-attacher-role created
rolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-health-monitor-controller-role created
rolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-provisioner-role created
rolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-resizer-role created
rolebinding.rbac.authorization.k8s.io/csi-hostpathplugin-snapshotter-role created
statefulset.apps/csi-hostpathplugin created
   /Users/csi-driver-host-path/deploy/kubernetes-latest/hostpath/csi-hostpath-snapshotclass.yaml
volumesnapshotclass.snapshot.storage.k8s.io/csi-hostpath-snapclass unchanged
   /Users/csi-driver-host-path/deploy/kubernetes-latest/hostpath/csi-hostpath-testing.yaml
        using           image: registry.k8s.io/sig-storage/hostpathplugin:v1.15.0
service/hostpath-service created
statefulset.apps/csi-hostpath-socat created
13:49:11 waiting for hostpath deployment to complete, attempt #0
```

The [livenessprobe side-container](https://github.com/kubernetes-csi/livenessprobe) provided by the CSI community is deployed with the CSI driver to provide the liveness checking of the CSI services.

## Modify Cluster Role

For example, if you want to modify external-resizer RBAC rules, you can do:
```
kubectl edit clusterrole external-resizer-runner 
```
Replace external-resizer-runner to the role you want to modify

## Run example application and validate

Next, validate the deployment.  
First, ensure all expected pods are running properly including the external attacher, provisioner, snapshotter and the actual hostpath driver plugin:

```shell
$ kubectl get pods
NAME                   READY   STATUS    RESTARTS   AGE
csi-hostpath-socat-0   1/1     Running   0          8m8s
csi-hostpathplugin-0   8/8     Running   0          8m9s
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
pvc-ad827273-8d08-430b-9d5a-e60e05a2bc3e   1Gi        RWO            Delete           Bound    default/csi-pvc   csi-hostpath-sc            45s

$ kubectl get pvc
NAME      STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS      AGE
csi-pvc   Bound    pvc-ad827273-8d08-430b-9d5a-e60e05a2bc3e   1Gi        RWO            csi-hostpath-sc   94s
```

Finally, inspect the application pod `my-csi-app`  which mounts a Hostpath volume:

```shell
kubectl describe pods/my-csi-app
Name:             my-csi-app
Namespace:        default
Priority:         0
Service Account:  default
Node:             kind-control-plane/172.19.0.2
Start Time:       Sat, 29 Mar 2025 13:59:51 -0700
Labels:           <none>
Annotations:      <none>
Status:           Running
IP:               10.244.0.22
IPs:
  IP:  10.244.0.22
Containers:
  my-frontend:
    Container ID:  containerd://6ec737ab0ef8510a2d8c4fcbaa869a6e58785fe7bc53e8fd83740aa0244a969a
    Image:         busybox
    Image ID:      docker.io/library/busybox@sha256:37f7b378a29ceb4c551b1b5582e27747b855bbfaa73fa11914fe0df028dc581f
    Port:          <none>
    Host Port:     <none>
    Command:
      sleep
      1000000
    State:          Running
      Started:      Sat, 29 Mar 2025 14:00:02 -0700
    Ready:          True
    Restart Count:  0
    Environment:    <none>
    Mounts:
      /data from my-csi-volume (rw)
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-kwlwh (ro)
Conditions:
  Type                        Status
  PodReadyToStartContainers   True 
  Initialized                 True 
  Ready                       True 
  ContainersReady             True 
  PodScheduled                True 
Volumes:
  my-csi-volume:
    Type:       PersistentVolumeClaim (a reference to a PersistentVolumeClaim in the same namespace)
    ClaimName:  csi-pvc
    ReadOnly:   false
  kube-api-access-kwlwh:
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
  Normal  Scheduled               67s   default-scheduler        Successfully assigned default/my-csi-app to kind-control-plane
  Normal  SuccessfulAttachVolume  66s   attachdetach-controller  AttachVolume.Attach succeeded for volume "pvc-80c31c4e-27d1-45ef-b302-8b29704f3415"
  Normal  Pulling                 57s   kubelet                  Pulling image "busybox"
  Normal  Pulled                  56s   kubelet                  Successfully pulled image "busybox" in 807ms (807ms including waiting). Image size: 1855985 bytes.
  Normal  Created                 56s   kubelet                  Created container: my-frontend
  Normal  Started                 56s   kubelet                  Started container my-frontend
```

## Confirm Hostpath driver works
The Hostpath driver is configured to create new volumes under `/csi-data-dir` inside the hostpath container that is specified in the plugin StatefulSet found [here](../deploy/kubernetes-1.31-test/hostpath/csi-hostpath-plugin.yaml).  
This path persist as long as the StatefulSet pod is up and running.

A file written in a properly mounted Hostpath volume inside an application should show up inside the Hostpath container.  
The following steps confirms that Hostpath is working properly.  First, create a file from the application pod as shown:

```shell
$ kubectl exec -it my-csi-app /bin/sh
/ # touch /data/hello-world
/ # exit
```

Next, ssh into the Hostpath container and verify that the file shows up there:
```shell
$ kubectl exec -it $(kubectl get pods --selector app.kubernetes.io/name=csi-hostpathplugin -o jsonpath='{.items[*].metadata.name}') -c hostpath /bin/sh

```
Then, use the following command to locate the file. If everything works OK you should get a result similar to the following:

```shell
/ # find / -name hello-world
/var/lib/kubelet/pods/907ee44d-582f-401a-bf87-8c7d42de619d/volumes/kubernetes.io~csi/pvc-80c31c4e-27d1-45ef-b302-8b29704f3415/mount/hello-world
/csi-data-dir/5f8cc66b-0c52-11f0-ae3c-12a0ddb447ec/hello-world
/ # exit
```

## Confirm the creation of the VolumeAttachment object
An additional way to ensure the driver is working properly is by inspecting the VolumeAttachment API object created that represents the attached volume:

```shell
$  kubectl describe volumeattachment
Name:         csi-76020859ca347da4de55748c73810c3b1f9bbb9721651fabfacee8992a903aeb
Namespace:    
Labels:       <none>
Annotations:  <none>
API Version:  storage.k8s.io/v1
Kind:         VolumeAttachment
Metadata:
  Creation Timestamp:  2025-03-29T20:59:51Z
  Resource Version:    131288
  UID:                 464a73bc-b296-4d6f-8324-ec2cde6bfc41
Spec:
  Attacher:   hostpath.csi.k8s.io
  Node Name:  kind-control-plane
  Source:
    Persistent Volume Name:  pvc-80c31c4e-27d1-45ef-b302-8b29704f3415
Status:
  Attached:  true
Events:      <none>
```
## 
The simplest way to Destroy the HostPath driver is by running the destroy.sh script.
For example, to destroy on Kubernetes 1.32.2, use the following command:
```shell
csi-driver-host-path % deploy/kubernetes-latest/destroy.sh
pod "csi-hostpath-socat-0" deleted
pod "csi-hostpathplugin-0" deleted
service "hostpath-service" deleted
statefulset.apps "csi-hostpath-socat" deleted
statefulset.apps "csi-hostpathplugin" deleted
role.rbac.authorization.k8s.io "external-attacher-cfg" deleted
role.rbac.authorization.k8s.io "external-health-monitor-controller-cfg" deleted
role.rbac.authorization.k8s.io "external-provisioner-cfg" deleted
role.rbac.authorization.k8s.io "external-resizer-cfg" deleted
role.rbac.authorization.k8s.io "external-snapshotter-leaderelection" deleted
clusterrole.rbac.authorization.k8s.io "external-attacher-runner" deleted
clusterrole.rbac.authorization.k8s.io "external-health-monitor-controller-runner" deleted
clusterrole.rbac.authorization.k8s.io "external-provisioner-runner" deleted
clusterrole.rbac.authorization.k8s.io "external-resizer-runner" deleted
clusterrole.rbac.authorization.k8s.io "external-snapshotter-runner" deleted
rolebinding.rbac.authorization.k8s.io "csi-attacher-role-cfg" deleted
rolebinding.rbac.authorization.k8s.io "csi-external-health-monitor-controller-role-cfg" deleted
rolebinding.rbac.authorization.k8s.io "csi-hostpathplugin-attacher-role" deleted
rolebinding.rbac.authorization.k8s.io "csi-hostpathplugin-health-monitor-controller-role" deleted
rolebinding.rbac.authorization.k8s.io "csi-hostpathplugin-provisioner-role" deleted
rolebinding.rbac.authorization.k8s.io "csi-hostpathplugin-resizer-role" deleted
rolebinding.rbac.authorization.k8s.io "csi-hostpathplugin-snapshotter-role" deleted
rolebinding.rbac.authorization.k8s.io "csi-provisioner-role-cfg" deleted
rolebinding.rbac.authorization.k8s.io "csi-resizer-role-cfg" deleted
rolebinding.rbac.authorization.k8s.io "external-snapshotter-leaderelection" deleted
clusterrolebinding.rbac.authorization.k8s.io "csi-attacher-role" deleted
clusterrolebinding.rbac.authorization.k8s.io "csi-external-health-monitor-controller-role" deleted
clusterrolebinding.rbac.authorization.k8s.io "csi-hostpathplugin-attacher-cluster-role" deleted
clusterrolebinding.rbac.authorization.k8s.io "csi-hostpathplugin-health-monitor-controller-cluster-role" deleted
clusterrolebinding.rbac.authorization.k8s.io "csi-hostpathplugin-provisioner-cluster-role" deleted
clusterrolebinding.rbac.authorization.k8s.io "csi-hostpathplugin-resizer-cluster-role" deleted
clusterrolebinding.rbac.authorization.k8s.io "csi-hostpathplugin-snapshot-metadata-cluster-role" deleted
clusterrolebinding.rbac.authorization.k8s.io "csi-hostpathplugin-snapshotter-cluster-role" deleted
clusterrolebinding.rbac.authorization.k8s.io "csi-provisioner-role" deleted
clusterrolebinding.rbac.authorization.k8s.io "csi-resizer-role" deleted
clusterrolebinding.rbac.authorization.k8s.io "csi-snapshotter-role" deleted
serviceaccount "csi-attacher" deleted
serviceaccount "csi-external-health-monitor-controller" deleted
serviceaccount "csi-hostpathplugin-sa" deleted
serviceaccount "csi-provisioner" deleted
serviceaccount "csi-resizer" deleted
serviceaccount "csi-snapshotter" deleted
csidriver.storage.k8s.io "hostpath.csi.k8s.io" deleted
```