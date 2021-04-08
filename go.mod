module github.com/libopenstorage/csi-test

go 1.15

require (
	github.com/cloudfoundry/gosigar v1.1.0 // indirect
	github.com/container-storage-interface/spec v1.3.0
	github.com/gobuffalo/packr v1.30.1 // indirect
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2 // indirect
	github.com/kubernetes-csi/csi-test v2.2.0+incompatible
	github.com/kubernetes-csi/csi-test/v4 v4.1.0
	github.com/libopenstorage/openstorage v8.0.1-0.20190926212733-daaed777713e+incompatible
	github.com/libopenstorage/systemutils v0.0.0-20180202003325-dc0e35340804 // indirect
	github.com/onsi/ginkgo v1.15.0
	github.com/onsi/gomega v1.10.5
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/portworx/kvdb v0.0.0-20210316224406-4c29d8da7240
	github.com/prometheus/client_golang v1.10.0 // indirect
	github.com/robertkrimen/otto v0.0.0-20200922221731-ef014fd054ac
	github.com/rs/cors v1.7.0 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.7.0
	github.com/vbatts/tar-split v0.11.1 // indirect
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b
	google.golang.org/grpc v1.34.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/klog/v2 v2.4.0
)

replace (
	github.com/docker/docker => github.com/moby/moby v20.10.3-0.20210324213045-797b974cb90e+incompatible
	github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1
	github.com/kubernetes-incubator/external-storage => github.com/libopenstorage/external-storage v0.20.4-openstorage-rc3
	github.com/libopenstorage/openstorage => github.com/libopenstorage/openstorage v0.0.0-k8s120-rc1
	github.com/libopenstorage/openstorage v8.0.1-0.20190926212733-daaed777713e+incompatible => github.com/libopenstorage/openstorage v0.0.0-k8s120-rc1
	github.com/satori/go.uuid => github.com/satori/go.uuid v0.0.0-20160324112244-f9ab0dce87d8
	google.golang.org/grpc => google.golang.org/grpc v1.29.1

	k8s.io/api => k8s.io/api v0.20.4
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.20.4
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.4
	k8s.io/apiserver => k8s.io/apiserver v0.20.4
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.20.4
	k8s.io/client-go => k8s.io/client-go v0.20.4
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.20.4
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.20.4
	k8s.io/code-generator => k8s.io/code-generator v0.20.4
	k8s.io/component-base => k8s.io/component-base v0.20.4
	k8s.io/component-helpers => k8s.io/component-helpers v0.20.4
	k8s.io/controller-manager => k8s.io/controller-manager v0.20.4
	k8s.io/cri-api => k8s.io/cri-api v0.20.4
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.20.4
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.20.4
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.20.4
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.20.4
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.20.4
	k8s.io/kubectl => k8s.io/kubectl v0.20.4
	k8s.io/kubelet => k8s.io/kubelet v0.20.4
	k8s.io/kubernetes => k8s.io/kubernetes v1.20.4
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.20.4
	k8s.io/metrics => k8s.io/metrics v0.20.4
	k8s.io/mount-utils => k8s.io/mount-utils v0.20.4
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.20.4
)
