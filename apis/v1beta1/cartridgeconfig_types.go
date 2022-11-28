package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CartridgeConfigSpec defines the desired state of CartridgeConfig.
type CartridgeConfigSpec struct {
	// Data contains the configuration data.
	// +kubebuilder:validation:Required
	Data string `json:"data,omitempty"`
}

// CartridgeConfigPhase is a label for the condition of a CartridgeConfig at the current time.
// +enum.
type CartridgeConfigPhase string

const (
	CartridgeConfigWaitingForCluster CartridgeConfigPhase = "WaitingForCluster"
	CartridgeConfigWaitingForLeader  CartridgeConfigPhase = "WaitingForLeader"
	CartridgeConfigApplying          CartridgeConfigPhase = "Applying"
	CartridgeConfigReady             CartridgeConfigPhase = "Ready"
)

// CartridgeConfigStatus defines the observed state of CartridgeConfig.
type CartridgeConfigStatus struct {
	// Phase indicates current state of CartridgeConfig
	// +kubebuilder:default=Pending
	Phase CartridgeConfigPhase `json:"phase"`
}

// CartridgeConfig is the Schema for the cartridgeconfigs API
// More info: https://www.tarantool.io/doc/latest/book/cartridge/cartridge_api/modules/cartridge.clusterwide-config/
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",priority=0
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type CartridgeConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CartridgeConfigSpec   `json:"spec,omitempty"`
	Status CartridgeConfigStatus `json:"status,omitempty"`
}

func (in *CartridgeConfig) ResetStatus() {
	in.Status = CartridgeConfigStatus{}
}

func (in *CartridgeConfig) SetPhase(phase CartridgeConfigPhase) {
	in.Status.Phase = phase
}

func (in *CartridgeConfig) GetPhase() CartridgeConfigPhase {
	return in.Status.Phase
}

func (in *CartridgeConfig) GetData() []byte {
	return []byte(in.Spec.Data)
}

//+kubebuilder:object:root=true

// CartridgeConfigList contains a list of CartridgeConfig.
type CartridgeConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CartridgeConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CartridgeConfig{}, &CartridgeConfigList{})
}
