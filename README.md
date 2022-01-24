# CSI Hostpath Driver

This repository hosts the CSI Hostpath driver and all of its build and dependent configuration files to deploy the driver.

---
*WARNING: This driver is just a demo implementation and is used for CI testing. This has many fake implementations and other non-standard best practices, and should not be used as an example of how to write a real driver.
---

## Pre-requisite
- Kubernetes cluster
- Running version 1.17 or later
- Access to terminal with `kubectl` installed
- VolumeSnapshot CRDs and Snapshot Controller must be installed as part of the cluster deployment (see Kubernetes 1.17+ deployment instructions)

## Features

The driver can provide empty directories that are backed by the same filesystem as EmptyDir volumes. In addition, it can provide raw block volumes that are backed by a single file in that same filesystem and bound to a loop device.

[Various command line parameters](cmd/hostpathplugin/main.go) influence the behavior of the driver. This is relevant in particular for the end-to-end testing that this driver is used for in Kubernetes.

Usually, the driver implements all CSI operations itself. When deployed with the `-proxy-endpoint` parameter, it instead proxies all incoming connections for a CSI driver that is [embedded inside the Kubernetes E2E test suite](https://github.com/kubernetes/kubernetes/tree/master/test/e2e/storage/drivers/csi-test) and used for mocking a CSI driver [with callbacks provided by certain tests](https://github.com/kubernetes/kubernetes/blob/5ad79eae2dcbf33df3b35c48ec993d30fbda46dd/test/e2e/storage/csi_mock_volume.go#L110).

## Deployment
[Deployment for Kubernetes 1.17 and later](docs/deploy-1.17-and-later.md)

## Examples
The following examples assume that the CSI hostpath driver has been deployed and validated:
- [Volume snapshots](docs/example-snapshots-1.17-and-later.md)
- [Inline ephemeral volumes](docs/example-ephemeral.md)

## Building the binaries
If you want to build the driver yourself, you can do so with the following command from the root directory:

```shell
make
```

## Development

### Updating sidecar images
The `deploy/` directory contains manifests for deploying the CSI hostpath driver for different Kubernetes versions.

If you want to update the image versions used in these manifests, you can do so with the following command from the root directory:

```shell
hack/bump-image-versions.sh
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
