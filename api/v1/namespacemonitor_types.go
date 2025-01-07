/*
Copyright 2024.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NamespaceMonitorSpec defines the desired state of NamespaceMonitor.
type NamespaceMonitorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of NamespaceMonitor. Edit namespacemonitor_types.go to remove/update
	// Foo string `json:"foo,omitempty"`
	Namespace      string `json:"namespace"`
	UpdateInterval string `json:"updateInterval"`
}

type ContainerMetrics struct {
	ContainerName string `json:"containerName"`
	CPUUsage      string `json:"cpuUsage"`
	MemoryUsage   string `json:"memoryUsage"`
}

type PodMetrics struct {
	PodName          string             `json:"podName"`
	ContainerMetrics []ContainerMetrics `json:"containerMetrics"`
}

// NamespaceMonitorStatus defines the observed state of NamespaceMonitor.
type NamespaceMonitorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	LastUpdated metav1.Time  `json:"lastUpdated,omitempty"`
	PodMetrics  []PodMetrics `json:"podMetrics,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// NamespaceMonitor is the Schema for the namespacemonitors API.
type NamespaceMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NamespaceMonitorSpec   `json:"spec,omitempty"`
	Status NamespaceMonitorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NamespaceMonitorList contains a list of NamespaceMonitor.
type NamespaceMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamespaceMonitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NamespaceMonitor{}, &NamespaceMonitorList{})
}
