module github.com/pennyscissors/image-exporter

go 1.16

// replace k8s.io/client-go => k8s.io/client-go v0.18.0

// require (
// 	github.com/rancher/wrangler v0.8.10
// 	github.com/sirupsen/logrus v1.4.2
// 	github.com/urfave/cli v1.22.2
// 	k8s.io/api v0.18.8
// 	k8s.io/apiextensions-apiserver v0.18.0
// 	k8s.io/apimachinery v0.18.8
// 	k8s.io/client-go v0.18.8
// 	k8s.io/code-generator v0.18.0
// )

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/rancher/rancher/pkg/apis => github.com/rancher/rancher/pkg/apis v0.0.0-20211218171736-3c1d5fac3248
	github.com/rancher/rancher/pkg/client => github.com/rancher/rancher/pkg/client v0.0.0-20211218171736-3c1d5fac3248
	helm.sh/helm/v3 => github.com/rancher/helm/v3 v3.5.4-rancher.1
	k8s.io/api => k8s.io/api v0.21.0
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.0
	k8s.io/apiserver => k8s.io/apiserver v0.21.0
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.21.0
	k8s.io/client-go => github.com/rancher/client-go v0.21.0-rancher.1
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.21.0
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.21.0
	k8s.io/code-generator => k8s.io/code-generator v0.21.0
	k8s.io/component-base => k8s.io/component-base v0.21.0
	k8s.io/component-helpers => k8s.io/component-helpers v0.21.0
	k8s.io/controller-manager => k8s.io/controller-manager v0.21.0
	k8s.io/cri-api => k8s.io/cri-api v0.21.0
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.21.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.21.0
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.21.0
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.21.0
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.21.0
	k8s.io/kubectl => k8s.io/kubectl v0.21.0
	k8s.io/kubelet => k8s.io/kubelet v0.21.0
	k8s.io/kubernetes => k8s.io/kubernetes v1.21.0
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.21.0
	k8s.io/metrics => k8s.io/metrics v0.21.0
	k8s.io/mount-utils => k8s.io/mount-utils v0.21.0
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.21.0
	launchpad.net/gocheck => github.com/go-check/check v0.0.0-20200227125254-8fa46927fb4f
)

require (
	github.com/Microsoft/go-winio v0.4.17 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/go-git/go-git/v5 v5.4.2
	github.com/pkg/errors v0.9.1
	github.com/rancher/rancher v0.0.0-20211218171736-3c1d5fac3248
	github.com/sirupsen/logrus v1.8.1
	github.com/urfave/cli/v2 v2.3.0
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
