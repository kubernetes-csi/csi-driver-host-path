## Inline ephemeral support
The CSI Hostpath driver (starting with version 1.2.0) now includes support for inline ephemeral volumes. This means that a volume can be specified directly inside a pod spec without the need to use a persistent volume object.
Inline ephemeral support was introduced as an alpha feature in Kubernetes 1.15 and was promoted to beta in 1.16.
Find out how to enable or create a CSI driver with support for such volumes [here](https://kubernetes-csi.github.io/docs/ephemeral-local-volumes.html)

To test this feature on Kubernetes 1.15 (and only with that release), redeploy the CSI Hostpath plugin YAML by updating the `hostpath` container to use  the inline ephemeral mode by setting the `ephemeral` flag, of the driver binary, to true as shown in the following setup:

```yaml
kind: StatefulSet
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
This is an intermediate solution for Kubernetes 1.15. With Kubernetes 1.16 and later, the normal
deployment supports both inline ephemeral volumes and persistent volumes.

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
        driver: hostpath.csi.k8s.io
``` 

> See sample YAML file [here](../examples/csi-app-inline.yaml).

Notice the CSI driver is now specified directly in the container spec inside the `volumes:` block.  You can use the [same steps as above][Confirm Hostpath driver works] 
to verify that the volume has been created and deleted (when the pod is removed).
