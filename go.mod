module github.com/pennyscissors/image-exporter

go 1.16

replace k8s.io/client-go => k8s.io/client-go v0.18.0

require (
	github.com/rancher/wrangler v0.8.10
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.22.2
	k8s.io/api v0.18.8
	k8s.io/apiextensions-apiserver v0.18.0
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v0.18.8
	k8s.io/code-generator v0.18.0
)
