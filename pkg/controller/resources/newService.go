package resources

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	apisv1 "github.com/kosmos.io/simple-controller/pkg/apis/v1"
)

func NewService(app *apisv1.AppService) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, schema.GroupVersionKind{
					Group:   app.GroupVersionKind().Group,
					Version: app.GroupVersionKind().Version,
					Kind:    app.GroupVersionKind().Kind,
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Type:  corev1.ServiceTypeNodePort,
			Ports: app.Spec.Ports,
			Selector: map[string]string{
				"app": app.Name,
			},
		},
	}
}
