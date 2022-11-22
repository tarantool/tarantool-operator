package v1alpha2

import (
	"fmt"

	"github.com/tarantool/tarantool-operator/pkg/api"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RoleSpec defines the desired state of Role.
type RoleSpec struct {
	// Replicasets is a number of StatefulSets (Tarantool replicasets) created under this Role, defaults to 1
	// +optional
	// +kubebuilder:default=1
	Replicasets *int32 `json:"replicasets"`

	// ReplicasetTemplate is the object that describes the StatefulSet that will be created if
	// insufficient Replicasets are detected. Each StatefulSet stamped out by the Role
	// will fulfill this ReplicasetTemplate, but have a unique identity from the rest
	// of the Role.
	ReplicasetTemplate *ReplicasetTemplate `json:"replicasetTemplate"`

	// AllRw is a flag indicating that all servers in the replicaset should be read-write.
	// +optional
	AllRw bool `json:"allRw"`

	// VShard defines config for vshard
	// See more: https://www.tarantool.io/doc/latest/reference/reference_rock/vshard/
	// +kubebuilder:validation:Required
	VShard RoleVShardConfig `json:"vshard"`
}

// ReplicasetTemplate is StatefulSet.Spec but with some fields omit
// +k8s:openapi-gen=true
type ReplicasetTemplate struct {
	// Replicas is the desired number of replicas of the given PodTemplate.
	// These are replicas in the sense that they are instantiations of the
	// same PodTemplate, but individual replicas also have a consistent identity.
	// If unspecified, defaults to 1.
	// +optional
	// +kubebuilder:default=1
	Replicas *int32 `json:"replicas,omitempty"`

	// PodTemplate is the object that describes the pod that will be created if
	// insufficient replicas are detected. Each pod stamped out by the StatefulSet
	// will fulfill this PodTemplate, but have a unique identity from the rest
	// of the StatefulSet.
	// +kubebuilder:validation:Required
	PodTemplate v1.PodTemplateSpec `json:"podTemplate,omitempty"`

	// VolumeClaimTemplates is a list of claims that pods are allowed to reference.
	// The StatefulSet controller is responsible for mapping network identities to
	// claims in a way that maintains the identity of a pod. Every claim in
	// this list must have at least one matching (by name) volumeMount in one
	// container in the template. A claim in this list takes precedence over
	// any volumes in the template, with the same name.
	// +optional
	VolumeClaimTemplates []v1.PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty"`

	// MinReadySeconds is a minimum number of seconds for which a newly created pod should be ready
	// without any of its container crashing for it to be considered available.
	// Defaults to 0 (pod will be considered available as soon as it is ready)
	// This is an alpha field and requires enabling StatefulSetMinReadySeconds feature gate.
	// +optional
	// +kubebuilder:default=0
	MinReadySeconds int32 `json:"minReadySeconds,omitempty"`
}

// RoleVShardConfig defines config for vshard
// See more: https://www.tarantool.io/doc/latest/reference/reference_rock/vshard/
// +k8s:openapi-gen=true
type RoleVShardConfig struct {
	// VshardGroupName defines group for ReplicaSet with `vshard-storage` role, defaults to: "default"
	// Seem more: https://www.tarantool.io/doc/latest/book/cartridge/cartridge_dev/#using-multiple-vshard-storage-groups
	// +optional
	// +kubebuilder:default=default
	VshardGroupName string `json:"vshardGroupName,omitempty"`

	// roles https://www.tarantool.io/ru/doc/latest/book/cartridge/cartridge_dev/#cluster-roles
	// +kubebuilder:validation:Required
	ClusterRoles []string `json:"clusterRoles"`

	// Weight is a number which determinate vshard weight of each replicaset in role, defaults to 100
	// See more: https://www.tarantool.io/ru/doc/latest/reference/reference_rock/vshard/vshard_admin/#replica-weights
	// +optional
	// +kubebuilder:default=100
	Weight int32 `json:"weight,omitempty"`
}

func (in *RoleVShardConfig) GetGroupName() string {
	return in.VshardGroupName
}

func (in *RoleVShardConfig) GetRoles() []string {
	return in.ClusterRoles
}

func (in *RoleVShardConfig) GetWeight() int32 {
	return in.Weight
}

// RolePhase is a label for the condition of a role at the current time.
// +enum.
type RolePhase string

// These are the valid statuses of Role.
const (
	RolePending               RolePhase = "Pending"
	RoleWaitingForCluster     RolePhase = "WaitingForCluster"
	RoleWaitingForLeader      RolePhase = "WaitingForLeader"
	RoleWaitForCartridgeReady RolePhase = "WaitForCartridgeReady"
	RoleJoining               RolePhase = "Joining"
	RoleConfiguring           RolePhase = "Configuring"
	RoleWaitingForBootstrap   RolePhase = "WaitingForBootstrap"
	RoleConfiguringWeights    RolePhase = "ConfiguringWeights"
	RoleReady                 RolePhase = "Ready"
	RoleConfigError           RolePhase = "ConfigError"
)

// RoleStatus defines the observed state of Role
// +k8s:openapi-gen=true
type RoleStatus struct {
	// Phase of roles
	// +kubebuilder:default=Pending
	Phase RolePhase `json:"phase"`
}

// Role is the Schema for the roles API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",priority=0
// +kubebuilder:printcolumn:name="Replicasets",type="number",JSONPath=".status.replicasets",priority=1
// +kubebuilder:printcolumn:name="Ready replicasets",type="number",JSONPath=".status.readyReplicasets",priority=1
// +kubebuilder:printcolumn:name="Replicas",type="number",JSONPath=".status.replicas",priority=1
// +kubebuilder:printcolumn:name="Ready replicas",type="number",JSONPath=".status.readyReplicas",priority=1
// +kubebuilder:printcolumn:name="Weight",type="number",JSONPath=".status.weight",priority=0
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",priority=0
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Role struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   RoleSpec   `json:"spec,omitempty"`
	Status RoleStatus `json:"status,omitempty"`
}

func (in *Role) IsAllRw() bool {
	return in.Spec.AllRw
}

func (in *Role) GetReplicasetName(ordinal int32) (string, error) {
	return fmt.Sprintf("%s-%d", in.Name, ordinal), nil
}

func (in *Role) GetReplicasets() int32 {
	return *in.Spec.Replicasets
}

func (in *Role) GetReplicas() int32 {
	return *in.Spec.ReplicasetTemplate.Replicas
}

func (in *Role) GetVolumeClaimTemplates() []v1.PersistentVolumeClaim {
	return in.Spec.ReplicasetTemplate.VolumeClaimTemplates
}

func (in *Role) ResetStatus() {
	in.Status = RoleStatus{}
}

func (in *Role) SetPhase(phase RolePhase) {
	in.Status.Phase = phase
}

func (in *Role) GetPhase() RolePhase {
	return in.Status.Phase
}

func (in *Role) GetVShardConfig() api.VShardConfig {
	return &in.Spec.VShard
}

// RoleList contains a list of Role
// +kubebuilder:object:root=true
type RoleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Role `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Role{}, &RoleList{})
}
