package main

import (
	"os"

	"github.com/kosmos.io/simple-controller/cmd/operator/app"
	apiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/cli"
)

func main() {
	ctx := apiserver.SetupSignalContext()
	cmd := app.NewOperatorCommand(ctx)
	code := cli.Run(cmd)
	os.Exit(code)
}
