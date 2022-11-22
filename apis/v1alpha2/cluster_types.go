package v1alpha2

import (
	"github.com/tarantool/tarantool-operator/pkg/api"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterSpec defines the desired state of Cluster
// +k8s:openapi-gen=true
// +kubebuilder:printcolumn:name="Failover",type="string",JSONPath=".failover.mode",priority=0
type ClusterSpec struct {
	// Domain is kubernetes cluster domain, defaults to: "cluster.local".
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=cluster.local
	Domain string `json:"domain,omitempty"`

	// ListenPort is port where tarantool iproto starts
	// +optional
	// +kubebuilder:default=3301
	ListenPort int32 `json:"listenPort"`

	// Failover defines cartridge failover params
	// More info: https://www.tarantool.io/doc/latest/book/cartridge/cartridge_api/modules/cartridge/#failoverparams
	Failover FailoverConfig `json:"failover"`

	// Failover defines foreign cartridge instance as topology leader
	ForeignLeader string `json:"foreignLeader,omitempty"`
}

// FailoverConfig defines cartridge failover params
// More info: https://www.tarantool.io/doc/latest/book/cartridge/cartridge_api/modules/cartridge/#failoverparams
// +k8s:openapi-gen=true
type FailoverConfig struct {
	// Mode is a string value which determinate Cartridge failover mode. Allowed value is one of disabled, eventual, stateful.
	// More info: https://www.tarantool.io/doc/latest/book/cartridge/cartridge_api/modules/cartridge/#automatic-failover-management
	// +kubebuilder:validation:Enum=disabled;eventual;stateful;raft
	// +kubebuilder:default=disabled
	// +optional
	Mode api.FailoverMode `json:"mode"`

	// Timeout (in seconds), used by membership tomark suspect members as dead (default: 20)
	// +kubebuilder:default=20
	// +optional
	Timeout int32 `json:"timeout"`

	// StateProvider - storage for failover state. Only «etcd2» supported at the moment.
	// +kubebuilder:validation:Enum=etcd2;stateboard
	// +optional
	StateProvider api.FailoverStateProvider `json:"stateProvider"`

	// Fencing - Abandon leadership when both the state provider quorum and atleast one replica are lost (suitable in stateful mode only,default: false)
	// +kubebuilder:default=false
	// +optional
	Fencing bool `json:"fencing"`

	// FencingTimeout - Time (in seconds) to actuate fencing after the check fails(default: 10)
	// +kubebuilder:default=10
	// +optional
	FencingTimeout int32 `json:"fencingTimeout"`

	// FencingPause - The period (in seconds) of performing the check(default: 2)
	// +kubebuilder:default=2
	// +optional
	FencingPause int32 `json:"fencingPause"`

	// Etcd2 is a config for "etcd2" failover state provider
	// +optional
	Etcd2 *FailoverEtcd2 `json:"etcd2"`

	// Stateboard is a config for "stateboard" failover state provider
	// +optional
	Stateboard *FailoverStateboard `json:"stateboard"`
}

func (in *FailoverConfig) GetMode() api.FailoverMode {
	return in.Mode
}

func (in *FailoverConfig) GetTimeout() int32 {
	return in.Timeout
}

func (in *FailoverConfig) GetStateProvider() api.FailoverStateProvider {
	return in.StateProvider
}

func (in *FailoverConfig) GetFencing() bool {
	return in.Fencing
}

func (in *FailoverConfig) GetFencingTimeout() int32 {
	return in.FencingTimeout
}

func (in *FailoverConfig) GetFencingPause() int32 {
	return in.FencingPause
}

func (in *FailoverConfig) GetETCD2Config() api.FailoverETCD2Config {
	return in.Etcd2
}

func (in *FailoverConfig) GetStateboardConfig() api.FailoverStateboardConfig {
	return in.Stateboard
}

// FailoverEtcd2 is a config for "etcd2" failover state provider
// +k8s:openapi-gen=true
type FailoverEtcd2 struct {
	// Defines "etcd2" failover state provider params
	// More info: https://www.tarantool.io/doc/latest/book/cartridge/cartridge_api/modules/cartridge/#failoverparams

	// Endpoints - URIs that are used to discover and to access etcd cluster instances.
	// +optional
	Endpoints []string `json:"endpoints"`

	// Username for etcd2 connection
	// +kubebuilder:validation:Required
	Username string `json:"username"`

	// Password for etcd2 connection
	// +optional
	Password corev1.SecretReference `json:"password"`

	// LockDelay - Timeout (in seconds), determines lock’s time-to-live (default: 10)
	// +kubebuilder:default=10
	// +optional
	LockDelay int32 `json:"lockDelay"`

	// Prefix used for etcd keys: <prefix>/lock and`<prefix>/leaders`, defaults to empty string
	// +optional
	Prefix string `json:"prefix"`
}

func (in *FailoverEtcd2) GetEndpoints() []string {
	return in.Endpoints
}

func (in *FailoverEtcd2) GetUsername() string {
	return in.Username
}

func (in *FailoverEtcd2) GetPassword() corev1.SecretReference {
	return in.Password
}

func (in *FailoverEtcd2) GetLockDelay() int32 {
	return in.LockDelay
}

func (in *FailoverEtcd2) GetPrefix() string {
	return in.Prefix
}

// FailoverStateboard is a config for "stateboard" failover state provider
// +k8s:openapi-gen=true
type FailoverStateboard struct {
	// URI of stateboard
	// +kubebuilder:validation:Required
	URI string `json:"uril"`

	// Password for etcd2 connection
	// +optional
	Password corev1.SecretReference `json:"password"`
}

func (in *FailoverStateboard) GetURI() string {
	return in.URI
}

func (in *FailoverStateboard) GetPassword() corev1.SecretReference {
	return in.Password
}

// ClusterPhase is a label for the condition of a Cluster at the current time.
// +enum.
type ClusterPhase string

const (
	ClusterPending             ClusterPhase = "Pending"
	ClusterSyncingService      ClusterPhase = "SyncingService"
	ClusterWaitingForRoles     ClusterPhase = "WaitingForRoles"
	ClusterWaitingForLeader    ClusterPhase = "WaitingForLeader"
	ClusterReady               ClusterPhase = "Ready"
	ClusterUnableToBootstrap   ClusterPhase = "UnableToBootstrap"
	ClusterFailoverConfiguring ClusterPhase = "FailoverConfiguring"
)

// ClusterStatus defines the observed state of Cluster
// +k8s:openapi-gen=true
type ClusterStatus struct {
	// Phase indicates current state of cluster
	// +kubebuilder:default=Pending
	Phase ClusterPhase `json:"phase"`

	// Bootstrapped indicates that vshard was bootstrapped
	// +kubebuilder:default=false
	Bootstrapped bool `json:"bootstrapped"`

	// Leader indicates name of pod which use to control topology
	// +optional
	Leader string `json:"leader"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cluster is the Schema for the clusters API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",priority=0
// +kubebuilder:printcolumn:name="Leader",type="string",JSONPath=".status.leader",priority=0
// +kubebuilder:printcolumn:name="Bootstrapped",type="boolean",JSONPath=".status.bootstrapped",priority=0
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

func (in *Cluster) GetDomain() string {
	return in.Spec.Domain
}

func (in *Cluster) GetListenPort() int32 {
	return in.Spec.ListenPort
}

func (in *Cluster) GetFailoverConfig() api.FailoverConfig {
	return &in.Spec.Failover
}

func (in *Cluster) SetLeader(leader string) {
	in.Status.Leader = leader
}

func (in *Cluster) GetLeader() string {
	if in.Spec.ForeignLeader != "" {
		return in.Spec.ForeignLeader
	}

	return in.Status.Leader
}

func (in *Cluster) MarkBootstrapped() {
	in.Status.Bootstrapped = true
}

func (in *Cluster) IsBootstrapped() bool {
	return in.Status.Bootstrapped
}

func (in *Cluster) ResetStatus() {
	in.Status = ClusterStatus{
		Phase:        "",
		Bootstrapped: in.Status.Bootstrapped,
		Leader:       in.Status.Leader,
	}
}

func (in *Cluster) SetPhase(phase ClusterPhase) {
	in.Status.Phase = phase
}

func (in *Cluster) GetPhase() ClusterPhase {
	return in.Status.Phase
}

//+kubebuilder:object:root=true

// ClusterList contains a list of Cluster.
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
