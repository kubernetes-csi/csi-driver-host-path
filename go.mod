module github.com/kubernetes-csi/csi-driver-host-path

go 1.20

require (
	github.com/container-storage-interface/spec v1.8.0
	github.com/golang/glog v1.1.2
	github.com/golang/protobuf v1.5.3
	github.com/kubernetes-csi/csi-lib-utils v0.15.0
	github.com/pborman/uuid v1.2.1
	github.com/stretchr/testify v1.8.4
	golang.org/x/net v0.17.0
	golang.org/x/sys v0.13.0
	google.golang.org/grpc v1.58.2
	k8s.io/apimachinery v0.28.2
	k8s.io/klog/v2 v2.100.1
	k8s.io/utils v0.0.0-20230406110748-d93618cff8a2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230711160842-782d3b101e98 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace k8s.io/api => k8s.io/api v0.28.2

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.28.2

replace k8s.io/apimachinery => k8s.io/apimachinery v0.28.2

replace k8s.io/apiserver => k8s.io/apiserver v0.28.2

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.28.2

replace k8s.io/client-go => k8s.io/client-go v0.28.2

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.28.2

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.28.2

replace k8s.io/code-generator => k8s.io/code-generator v0.28.2

replace k8s.io/component-base => k8s.io/component-base v0.28.2

replace k8s.io/component-helpers => k8s.io/component-helpers v0.28.2

replace k8s.io/controller-manager => k8s.io/controller-manager v0.28.2

replace k8s.io/cri-api => k8s.io/cri-api v0.28.2

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.28.2

replace k8s.io/dynamic-resource-allocation => k8s.io/dynamic-resource-allocation v0.28.2

replace k8s.io/kms => k8s.io/kms v0.28.2

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.28.2

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.28.2

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.28.2

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.28.2

replace k8s.io/kubectl => k8s.io/kubectl v0.28.2

replace k8s.io/kubelet => k8s.io/kubelet v0.28.2

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.28.2

replace k8s.io/metrics => k8s.io/metrics v0.28.2

replace k8s.io/mount-utils => k8s.io/mount-utils v0.28.2

replace k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.28.2

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.28.2
