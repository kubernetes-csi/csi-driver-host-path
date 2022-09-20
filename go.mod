module github.com/kubernetes-csi/csi-driver-host-path

go 1.19

require (
	github.com/container-storage-interface/spec v1.6.0
	github.com/golang/glog v1.0.0
	github.com/golang/protobuf v1.5.2
	github.com/kubernetes-csi/csi-lib-utils v0.11.0
	github.com/pborman/uuid v1.2.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20220921203646-d300de134e69
	google.golang.org/grpc v1.49.0
	k8s.io/apimachinery v0.25.2
	k8s.io/klog/v2 v2.80.1
	k8s.io/kubernetes v1.25.2
	k8s.io/utils v0.0.0-20220922104903-7796b5f52b7e
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.13.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20220919091848-fb04ddd9f9c8 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220921223823-23cae91e6737 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiserver v0.25.2 // indirect
	k8s.io/component-base v0.25.2 // indirect
	k8s.io/mount-utils v0.25.2 // indirect
)

replace k8s.io/api => k8s.io/api v0.25.0

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.25.0

replace k8s.io/apimachinery => k8s.io/apimachinery v0.26.0-alpha.0

replace k8s.io/apiserver => k8s.io/apiserver v0.25.0

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.25.0

replace k8s.io/client-go => k8s.io/client-go v0.25.0

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.25.0

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.25.0

replace k8s.io/code-generator => k8s.io/code-generator v0.25.1-rc.0

replace k8s.io/component-base => k8s.io/component-base v0.25.0

replace k8s.io/component-helpers => k8s.io/component-helpers v0.25.0

replace k8s.io/controller-manager => k8s.io/controller-manager v0.25.0

replace k8s.io/cri-api => k8s.io/cri-api v0.25.1-rc.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.25.0

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.25.0

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.25.0

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.25.0

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.25.0

replace k8s.io/kubectl => k8s.io/kubectl v0.25.0

replace k8s.io/kubelet => k8s.io/kubelet v0.25.0

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.25.0

replace k8s.io/metrics => k8s.io/metrics v0.25.0

replace k8s.io/mount-utils => k8s.io/mount-utils v0.25.2-rc.0

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.25.0

replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.25.0

replace k8s.io/sample-controller => k8s.io/sample-controller v0.25.0

replace k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.25.0
