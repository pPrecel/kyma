module github.com/kyma-project/kyma/components/function-controller

go 1.13

require (
	contrib.go.opencensus.io/exporter/prometheus v0.1.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.12.8 // indirect
	github.com/Azure/azure-sdk-for-go v19.1.1+incompatible // indirect
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/containerd/containerd v1.3.0 // indirect
	github.com/docker/distribution v2.6.0-rc.1.0.20180327202408-83389a148052+incompatible // indirect
	github.com/docker/docker v1.4.2-0.20190924003213-a8608b5b67c7 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/gogo/protobuf v1.3.1
	github.com/google/go-containerregistry v0.0.0-20190320210540-8d4083db9aa0 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/mattbaird/jsonpatch v0.0.0-20171005235357-81af80346b1a
	github.com/onsi/ginkgo v1.11.0 // indirect
	github.com/onsi/gomega v1.8.1
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/openzipkin/zipkin-go v0.2.2 // indirect
	github.com/tektoncd/pipeline v0.7.0
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553
	golang.org/x/sys v0.0.0-20191010194322-b09406accb47 // indirect
	google.golang.org/grpc v1.24.0 // indirect
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.0.0-20191114101535-6c5935290e33
	k8s.io/kubernetes v1.11.10
	knative.dev/pkg v0.0.0-20191230183737-ead56ad1f3bd
	knative.dev/serving v0.8.1
	sigs.k8s.io/controller-runtime v0.4.0
)

// https://github.com/Azure/go-autorest/issues/414
// https://github.com/rancher/rio/pull/451
replace github.com/Azure/go-autorest => github.com/Azure/go-autorest v0.9.3
