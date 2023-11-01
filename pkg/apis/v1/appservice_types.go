package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// You need to edit this file according to your crd configuration
// Note: json tags are required. '+k8s:deepcopy-gen' annotation is required to generate deepcopy files
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppService is the Schema for the appservices API
type AppService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppServiceSpec   `json:"spec"`
	Status AppServiceStatus `json:"status"`
}

// AppServiceSpec is what we need for our spec in crd. Neither it nor AppServiceStatus need to generate runtime.Object
type AppServiceSpec struct {
	Size  *int32               `json:"size"`
	Image string               `json:"image"`
	Ports []corev1.ServicePort `json:"ports,omitempty"`
}

// AppServiceStatus defines the observed state of AppService
type AppServiceStatus struct {
	appsv1.DeploymentStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AppServiceList contains a list of AppService
type AppServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AppService{}, &AppServiceList{})
}
