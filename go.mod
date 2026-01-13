module github.com/kubernetes-csi/csi-driver-host-path

go 1.25.5

require (
	github.com/container-storage-interface/spec v1.12.0
	github.com/kubernetes-csi/csi-lib-utils v0.23.1
	github.com/pborman/uuid v1.2.1
	github.com/stretchr/testify v1.11.1
	golang.org/x/net v0.49.0
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11
	k8s.io/apimachinery v0.35.0
	k8s.io/klog/v2 v2.130.1
	k8s.io/kubernetes v1.35.0
	k8s.io/utils v0.0.0-20260108192941-914a6e750570
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fxamacker/cbor/v2 v2.9.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/moby/sys/mountinfo v0.7.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_golang v1.23.2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.5 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.opentelemetry.io/otel v1.39.0 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	go.uber.org/automaxprocs v1.6.0 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260112192933-99fd39fd28a9 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiextensions-apiserver v0.35.0 // indirect
	k8s.io/apiserver v0.35.0 // indirect
	k8s.io/client-go v1.5.2 // indirect
	k8s.io/component-base v0.35.0 // indirect
	k8s.io/controller-manager v0.35.0 // indirect
	k8s.io/mount-utils v0.35.0 // indirect
	sigs.k8s.io/json v0.0.0-20250730193827-2d320260d730 // indirect
)

replace k8s.io/api => k8s.io/api v0.35.0

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.35.0

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.33.0

replace k8s.io/apiserver => k8s.io/apiserver v0.35.0

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.35.0

replace k8s.io/client-go => k8s.io/client-go v0.35.0

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.35.0

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.35.0

replace k8s.io/code-generator => k8s.io/code-generator v0.35.0

replace k8s.io/component-base => k8s.io/component-base v0.35.0

replace k8s.io/component-helpers => k8s.io/component-helpers v0.35.0

replace k8s.io/controller-manager => k8s.io/controller-manager v0.35.0

replace k8s.io/cri-api => k8s.io/cri-api v0.35.0

replace k8s.io/cri-client => k8s.io/cri-client v0.35.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.35.0

replace k8s.io/dynamic-resource-allocation => k8s.io/dynamic-resource-allocation v0.35.0

replace k8s.io/endpointslice => k8s.io/endpointslice v0.35.0

replace k8s.io/externaljwt => k8s.io/externaljwt v0.35.0

replace k8s.io/kms => k8s.io/kms v0.35.0

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.35.0

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.35.0

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.35.0

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.35.0

replace k8s.io/kubectl => k8s.io/kubectl v0.35.0

replace k8s.io/kubelet => k8s.io/kubelet v0.35.0

replace k8s.io/metrics => k8s.io/metrics v0.35.0

replace k8s.io/mount-utils => k8s.io/mount-utils v0.35.0

replace k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.35.0

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.35.0
