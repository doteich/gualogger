/*
Copyright 2025.

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

package v1alpha

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// GeistConnector is the Schema for the geistconnectors API
type GeistConnector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GeistConnectorSpec   `json:"spec,omitempty"`
	Status GeistConnectorStatus `json:"status,omitempty"`
}

// GeistConnectorSpec defines the desired state of GeistConnector
type GeistConnectorSpec struct {
	// OpcuaSpec contains the configuration for the OPC UA application.
	// +kubebuilder:validation:Required
	ConnectorSpec ConnectorSpec `json:"connectorSpec"`

	// DeploymentSpec contains Kubernetes-specific configuration for the application's deployment.
	// +kubebuilder:validation:Required
	DeploymentSpec DeploymentSpec `json:"deploymentSpec"`
}

// DeploymentSpec holds Kubernetes deployment-specific configuration.
type DeploymentSpec struct {
	// The container image repository (e.g., myregistry/collector).
	// +kubebuilder:default=doteich/geist-connector
	// +kubebuilder:validation:Required
	ImageRepo string `json:"imageRepo"`

	// The container image version/tag (e.g., v1.0.0).
	// +kubebuilder:default=latest
	// +kubebuilder:validation:Required
	ImageVersion string `json:"imageVersion"`

	// The image pull policy.
	// +kubebuilder:validation:Enum=Always;IfNotPresent;Never
	// +kubebuilder:default=IfNotPresent
	// +kubebuilder:validation:Required
	PullPolicy string `json:"pullPolicy"`

	// Map of custom annotations to add to the Deployment/Pod.
	// +optional
	CustomAnnotations map[string]string `json:"customAnnotations,omitempty"`

	// Map of custom labels to add to the Deployment/Pod.
	// +optional
	CustomLabels map[string]string `json:"customLabels,omitempty"`
}

// OpcuaSpec holds the application configuration.
// It now contains both Opcua and Redpanda configs as siblings.
type ConnectorSpec struct {
	// +kubebuilder:validation:Required
	OPCUA OpcuaConfig `json:"opcua"`

	// +kubebuilder:validation:Required
	Redpanda RedpandaConfig `json:"redpanda"`
}

// GeistConnectorStatus defines the observed state of GeistConnector
type GeistConnectorStatus struct {
	// Conditions represents the latest available observations of the resource's state.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true

// GeistConnectorList contains a list of GeistConnector
type GeistConnectorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GeistConnector `json:"items"`
}

// OpcuaConfig reflects the 'opcua' block
type OpcuaConfig struct {
	Connection   ConnectionConfig   `json:"connection"`
	Subscription SubscriptionConfig `json:"subscription"`
}

// ConnectionConfig reflects the 'connection' block
type ConnectionConfig struct {
	Endpoint string `json:"endpoint"`
	Port     int32  `json:"port"`

	// +kubebuilder:validation:Enum=None;Sign;SignAndEncrypt
	Mode string `json:"mode"`

	// +kubebuilder:validation:Enum=None;Basic256;Basic256Sha256;Aes256Sha256RsaPss;Aes128Sha256RsaOaep
	Policy         string               `json:"policy"`
	Authentication AuthenticationConfig `json:"authentication"`
	Certificate    CertificateConfig    `json:"certificate"`
	RetryCount     int32                `json:"retry_count"`
}

// AuthenticationConfig reflects the 'authentication' block
type AuthenticationConfig struct {
	// +kubebuilder:validation:Enum=None;User&Password;Certificate
	Type        string           `json:"type"`
	Credentials CredentialConfig `json:"credentials,omitempty"`
	Certificate AuthCertConfig   `json:"certificate,omitempty"`
}

// CredentialConfig reflects the 'credentials' block
type CredentialConfig struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// AuthCertConfig reflects the 'authentication.certificate' block
type AuthCertConfig struct {
	CertificatePath string `json:"certificate_path,omitempty"`
	PrivateKeyPath  string `json:"private_key_path,omitempty"`
}

// CertificateConfig reflects the 'connection.certificate' block
type CertificateConfig struct {
	AutoCreate      bool   `json:"auto_create,omitempty"`
	CertificatePath string `json:"certificate_path,omitempty"`
	PrivateKeyPath  string `json:"private_key_path,omitempty"`
}

// SubscriptionConfig reflects the 'subscription' block
type SubscriptionConfig struct {
	SubInterval int32        `json:"sub_interval"`
	NodeIDs     []NodeIDInfo `json:"nodeids"`
}

// NodeIDInfo reflects an element in the 'nodeids' list
type NodeIDInfo struct {
	ID string `json:"id"`
}

// RedpandaConfig reflects the 'redpanda' block
type RedpandaConfig struct {
	Brokers []string           `json:"brokers"`
	Topic   string             `json:"topic"`
	Auth    RedpandaAuthConfig `json:"auth"`
	TLS     RedpandaTLSConfig  `json:"tls,omitempty"`
}

// RedpandaAuthConfig reflects the 'auth' block
type RedpandaAuthConfig struct {
	SASL SASLConfig `json:"sasl"`
}

// SASLConfig reflects the 'sasl' block
type SASLConfig struct {
	// +kubebuilder:validation:Enum=plain;scram-sha-256;scram-sha-512
	Type     string `json:"type"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// RedpandaTLSConfig reflects the 'tls' block
type RedpandaTLSConfig struct {
	InsecureSkipVerify bool `json:"insecure_skip_verify,omitempty"`
}

func init() {
	SchemeBuilder.Register(&GeistConnector{}, &GeistConnectorList{})
}
