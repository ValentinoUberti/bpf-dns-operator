/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BpfdnsSpec defines the desired state of Bpfdns

type DnsNameStruct struct {
	DnsName string `json:"dnsName,omitempty"`
}

type BpfdnsSpec struct {
	/* Bpfdns yaml
	spec:
	  blockDns:
	    - dnsName: wwww.google.com
		- dnsName: www.example.com

	*/
	// Foo is an example field of Bpfdns. Edit Bpfdns_types.go to remove/update
	BlockDns []DnsNameStruct `json:"blockDns,omitempty"`
}

// BpfdnsStatus defines the observed state of Bpfdns
type BpfdnsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes []string `json:"nodes"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Bpfdns is the Schema for the bpfdns API
type Bpfdns struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BpfdnsSpec   `json:"spec,omitempty"`
	Status BpfdnsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BpfdnsList contains a list of Bpfdns
type BpfdnsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Bpfdns `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Bpfdns{}, &BpfdnsList{})
}
