package controller

import (
	"context"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apisv1 "github.com/kosmos.io/simple-controller/pkg/apis/v1"
	"github.com/kosmos.io/simple-controller/pkg/controller/resources"
)

// AppServiceReconciler reconciles a AppService object
type AppServiceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

var (
	oldSpecAnnotation = "old/spec"
)

// Reconcile is the core logical part of your controller
// For more details, you can refer to here
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *AppServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Get the AppService instance
	var appService apisv1.AppService
	err := r.Get(ctx, req.NamespacedName, &appService)
	if err != nil {
		// Ignore when the AppService is deleted
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	logger.Info("fetch appservice objects: ", "appservice", req.NamespacedName, "yaml", appService)

	// If no, create an associated resource
	// If yes, determine whether to update it
	deploy := &appsv1.Deployment{}
	if err := r.Get(ctx, req.NamespacedName, deploy); err != nil && errors.IsNotFound(err) {
		// 1. Related Annotations
		data, _ := json.Marshal(appService.Spec)
		if appService.Annotations != nil {
			appService.Annotations[oldSpecAnnotation] = string(data)
		} else {
			appService.Annotations = map[string]string{oldSpecAnnotation: string(data)}
		}
		if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			return r.Client.Update(ctx, &appService)
		}); err != nil {
			return ctrl.Result{}, err
		}
		// Creating Associated Resources
		// 2. Create Deployment
		deploy := resources.NewDeploy(&appService)
		if err := r.Client.Create(ctx, deploy); err != nil {
			return ctrl.Result{}, err
		}
		// 3. Create Service
		service := resources.NewService(&appService)
		if err := r.Create(ctx, service); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	// Gets the spec of the original AppService object
	oldspec := apisv1.AppServiceSpec{}
	if err := json.Unmarshal([]byte(appService.Annotations[oldSpecAnnotation]), &oldspec); err != nil {
		return ctrl.Result{}, err
	}
	// The current spec is inconsistent with the old object and needs to be updated
	// Otherwise return
	if !reflect.DeepEqual(appService.Spec, oldspec) {
		// Update associated resources
		newDeploy := resources.NewDeploy(&appService)
		oldDeploy := &appsv1.Deployment{}
		if err := r.Get(ctx, req.NamespacedName, oldDeploy); err != nil {
			return ctrl.Result{}, err
		}
		oldDeploy.Spec = newDeploy.Spec
		if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			return r.Client.Update(ctx, oldDeploy)
		}); err != nil {
			return ctrl.Result{}, err
		}

		newService := resources.NewService(&appService)
		oldService := &corev1.Service{}
		if err := r.Get(ctx, req.NamespacedName, oldService); err != nil {
			return ctrl.Result{}, err
		}
		// You need to specify the ClusterIP to the previous one; otherwise, an error will be reported during the update
		newService.Spec.ClusterIP = oldService.Spec.ClusterIP
		oldService.Spec = newService.Spec
		if err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			return r.Client.Update(ctx, oldService)
		}); err != nil {
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AppServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apisv1.AppService{}).
		Complete(r)
}
