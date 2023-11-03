package app

import (
	"context"
	"fmt"
	apisv1 "github.com/kosmos.io/simple-controller/pkg/apis/v1"
	"github.com/kosmos.io/simple-controller/pkg/controller"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	setupLog = ctrl.Log.WithName("setup")
)

type Options struct {
	KubeConfig string
}

// NewOperatorCommand creates a *cobra.Command object with default parameters
func NewOperatorCommand(ctx context.Context) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:  "operator",
		Long: `starting operator`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := Run(ctx, opts); err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&opts.KubeConfig, "kubeconfig", "", "Absolute path to the kubeconfig file.")

	return cmd
}

func Run(ctx context.Context, opts *Options) error {
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", opts.KubeConfig)
	if err != nil {
		return fmt.Errorf("error building kubeconfig: %s", err.Error())
	}

	ctrl.SetLogger(zap.New())
	setupLog.Info("starting manager")

	mgr, err := ctrl.NewManager(kubeconfig, ctrl.Options{})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// in a real controller, we'd create a new scheme for this
	if err = apisv1.AddToScheme(mgr.GetScheme()); err != nil {
		setupLog.Error(err, "unable to add scheme")
		os.Exit(1)
	}

	if err = (&controller.AppServiceReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "AppService")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctx); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
	return nil
}
