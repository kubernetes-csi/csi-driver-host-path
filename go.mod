module github.com/kubernetes-csi/csi-driver-host-path

go 1.12

require (
	github.com/container-storage-interface/spec v1.4.0
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.5.1
	github.com/google/uuid v1.2.0 // indirect
	github.com/kubernetes-csi/csi-lib-utils v0.9.1
	github.com/pborman/uuid v1.2.1
	github.com/stretchr/testify v1.6.1
	golang.org/x/net v0.0.0-20210316092652-d523dce5a7f4
	golang.org/x/sys v0.0.0-20210317225723-c4fcb01b228e // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210317182105-75c7a8546eb9 // indirect
	google.golang.org/grpc v1.36.0
	k8s.io/apimachinery v0.21.0-beta.1
	k8s.io/apiserver v0.20.4 // indirect
	k8s.io/component-base v0.20.4 // indirect
	k8s.io/klog/v2 v2.8.0 // indirect
	k8s.io/kubernetes v1.20.4
	k8s.io/mount-utils v0.20.4 // indirect
	k8s.io/utils v0.0.0-20210305010621-2afb4311ab10
)

replace k8s.io/api => k8s.io/api v0.20.0

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.20.0

replace k8s.io/apimachinery => k8s.io/apimachinery v0.21.0-alpha.0

replace k8s.io/apiserver => k8s.io/apiserver v0.20.0

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.20.0

replace k8s.io/client-go => k8s.io/client-go v0.20.0

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.20.0

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.20.0

replace k8s.io/code-generator => k8s.io/code-generator v0.20.1-rc.0

replace k8s.io/component-base => k8s.io/component-base v0.20.0

replace k8s.io/component-helpers => k8s.io/component-helpers v0.20.0

replace k8s.io/controller-manager => k8s.io/controller-manager v0.20.0

replace k8s.io/cri-api => k8s.io/cri-api v0.20.1-rc.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.20.0

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.20.0

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.20.0

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.20.0

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.20.0

replace k8s.io/kubectl => k8s.io/kubectl v0.20.0

replace k8s.io/kubelet => k8s.io/kubelet v0.20.0

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.20.0

replace k8s.io/metrics => k8s.io/metrics v0.20.0

replace k8s.io/mount-utils => k8s.io/mount-utils v0.20.1-rc.0

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.20.0

replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.20.0

replace k8s.io/sample-controller => k8s.io/sample-controller v0.20.0
