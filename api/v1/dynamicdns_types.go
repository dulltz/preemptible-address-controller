package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DynamicDNSSpec defines the desired state of DynamicDNS
type DynamicDNSSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// DynamicDNSStatus defines the observed state of DynamicDNS
type DynamicDNSStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// DynamicDNS is the Schema for the dynamicdns API
type DynamicDNS struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DynamicDNSSpec   `json:"spec,omitempty"`
	Status DynamicDNSStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DynamicDNSList contains a list of DynamicDNS
type DynamicDNSList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DynamicDNS `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DynamicDNS{}, &DynamicDNSList{})
}
