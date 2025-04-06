// Package kubeconfig is a minimal implementation of the Kubernetes kubeconfig
// configuration file format. It aims to be simple, fast, and use minimal
// dependencies.
//
// The package makes a best-effort attempt to follow the same [schema] and
// [merging logic] published on the Kubernetes documentation. If perfect
// compatibility with Kubernetes is needed, the official Kubernetes
// [client-go/tools/clientcmd] package should be used instead.
//
// Types have the same identifiers and YAML/JSON names as those used specified
// by Kubernetes source code for the [kubeconfig structs], with the exception
// of changes that make it possible to differentiate between nil and empty
// basic types.
//
// [schema]: https://kubernetes.io/docs/reference/config-api/kubeconfig.v1/#Config
// [merging logic]: https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#merging-kubeconfig-files
// [client-go/tools/clientcmd]: https://pkg.go.dev/k8s.io/client-go@v0.32.3/tools/clientcmd
// [kubeconfig structs]: https://github.com/kubernetes/client-go/blob/2086688a727d00268de695a9c701e52458939701/tools/clientcmd/api/v1/types.go
package kubeconfig
