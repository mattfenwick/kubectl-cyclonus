package main

import (
	"github.com/mattfenwick/kubectl-cyclonus/pkg/cli"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // required for GKE
)

func main() {
	cli.RunRootCommand()
}
