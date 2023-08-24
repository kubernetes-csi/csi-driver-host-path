module github.com/kubernetes-csi/csi-driver-host-path

go 1.20

require (
	github.com/container-storage-interface/spec v1.8.0
	github.com/golang/glog v1.1.2
	github.com/golang/protobuf v1.5.3
	github.com/kubernetes-csi/csi-lib-utils v0.14.0
	github.com/pborman/uuid v1.2.1
	github.com/stretchr/testify v1.8.4
	golang.org/x/net v0.14.0
	google.golang.org/grpc v1.57.0
	k8s.io/apimachinery v0.27.0
	k8s.io/klog/v2 v2.100.1
	k8s.io/kubernetes v1.28.1
	k8s.io/utils v0.0.0-20230406110748-d93618cff8a2
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.16.0 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.10.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiextensions-apiserver v0.0.0 // indirect
	k8s.io/apiserver v0.27.0 // indirect
	k8s.io/component-base v0.27.0 // indirect
	k8s.io/mount-utils v0.26.0 // indirect
)

replace k8s.io/api => k8s.io/api v0.27.0

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.27.0

replace k8s.io/apimachinery => k8s.io/apimachinery v0.27.0

replace k8s.io/apiserver => k8s.io/apiserver v0.27.0

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.27.0

replace k8s.io/client-go => k8s.io/client-go v0.27.0

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.27.0

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.27.0

replace k8s.io/code-generator => k8s.io/code-generator v0.27.0

replace k8s.io/component-base => k8s.io/component-base v0.27.0

replace k8s.io/component-helpers => k8s.io/component-helpers v0.27.0

replace k8s.io/controller-manager => k8s.io/controller-manager v0.27.0

replace k8s.io/cri-api => k8s.io/cri-api v0.27.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.27.0

replace k8s.io/dynamic-resource-allocation => k8s.io/dynamic-resource-allocation v0.27.0

replace k8s.io/kms => k8s.io/kms v0.27.0

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.27.0

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.27.0

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.27.0

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.27.0

replace k8s.io/kubectl => k8s.io/kubectl v0.27.0

replace k8s.io/kubelet => k8s.io/kubelet v0.27.0

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.27.0

replace k8s.io/metrics => k8s.io/metrics v0.27.0

replace k8s.io/mount-utils => k8s.io/mount-utils v0.27.0

replace k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.27.0

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.27.0
