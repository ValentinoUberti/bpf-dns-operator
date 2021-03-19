module github.com/ValentinoUberti/bpf-dns-operator

go 1.15

require (
	github.com/Wifx/gonetworkmanager v0.2.0
	github.com/apex/log v1.9.0
	github.com/dropbox/goebpf v0.0.0-20210223223402-d54e462ac389
	github.com/go-logr/logr v0.3.0
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	golang.org/x/sys v0.0.0-20200722175500-76b94024e4b6
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/controller-runtime v0.7.0
)
