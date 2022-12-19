package api

import (
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	RolePending               string = "Pending"
	RoleWaitingForCluster     string = "WaitingForCluster"
	RoleWaitingForLeader      string = "WaitingForLeader"
	RoleWaitForCartridgeReady string = "WaitForCartridgeReady"
	RoleJoining               string = "Joining"
	RoleConfiguring           string = "Configuring"
	RoleWaitingForBootstrap   string = "WaitingForBootstrap"
	RoleConfiguringWeights    string = "ConfiguringWeights"
	RoleReady                 string = "Ready"
	RoleConfigError           string = "ConfigError"
)

type Role interface {
	client.Object

	IsAllRw() bool

	GetReplicasetName(ordinal int32) (string, error)

	GetReplicasets() int32
	GetReplicas() int32

	GetVolumeClaimTemplates() []v1.PersistentVolumeClaim
	GetVShardConfig() VShardConfig

	ResetStatus()

	SetReadyPodsCount(count int32)
}

type VShardConfig interface {
	GetGroupName() string
	GetRoles() []string
	GetWeight() int32
}

type RoleWithStatus[PhaseType comparable] interface {
	Role

	SetPhase(phase PhaseType)
	GetPhase() PhaseType
}
