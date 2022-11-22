package topology

import "github.com/tarantool/tarantool-operator/pkg/api"

type JoinServerParams struct {
	URI  string `json:"uri"`
	UUID string `json:"uuid,omitempty"`
}

type EditReplicasetParams struct {
	UUID             string             `json:"uuid,omitempty"`
	Alias            string             `json:"alias,omitempty"`
	Roles            []string           `json:"roles,omitempty"`
	AllRw            bool               `json:"all_rw,omitempty"`
	Weight           int32              `json:"weight,omitempty"`
	FailoverPriority []string           `json:"failover_priority,omitempty"`
	VShardGroup      string             `json:"vshard_group,omitempty"`
	JoinServers      []JoinServerParams `json:"join_servers,omitempty"`
}

type EditServerParams struct {
	URI      string `json:"uri,omitempty"`
	UUID     string `json:"uuid"`
	Disabled bool   `json:"disabled,omitempty"`
	Expelled bool   `json:"expelled,omitempty"`
}

type EditTopologyParams struct {
	Servers     []EditServerParams     `json:"servers,omitempty"`
	Replicasets []EditReplicasetParams `json:"replicasets,omitempty"`
}

type GetReplicasetsQuery struct {
	UUID string `json:"uuid"`
}

// FailoverParams type struct for configure failover.
type FailoverParams struct {
	Mode             api.FailoverMode          `json:"mode"`
	Timeout          int32                     `json:"failover_timeout,omitempty"`
	StateProvider    api.FailoverStateProvider `json:"state_provider,omitempty"`
	FencingEnabled   bool                      `json:"fencing_enabled,omitempty"`
	FencingTimeout   int32                     `json:"fencing_timeout,omitempty"`
	FencingPause     int32                     `json:"fencing_pause,omitempty"`
	Etcd2Params      *Etcd2Params              `json:"etcd2_params,omitempty"`
	StateboardParams *StateboardParams         `json:"stateboard_params,omitempty"`
}

// Etcd2Params type struct for configure etcd2 failover state provider.
type Etcd2Params struct {
	Endpoints []string `json:"endpoints"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	LockDelay int32    `json:"lock_delay"`
	Prefix    string   `json:"prefix"`
}

// StateboardParams type struct for configure etcd2 failover state provider.
type StateboardParams struct {
	URI      string `json:"uri"`
	Password string `json:"password"`
}
