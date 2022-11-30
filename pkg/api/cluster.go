package api

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Cluster interface {
	client.Object

	GetDomain() string
	GetListenPort() int32

	GetFailoverConfig() FailoverConfig

	SetLeader(leader string)
	GetLeader() string

	MarkBootstrapped()
	IsBootstrapped() bool

	ResetStatus()
}

type ClusterWithStatus[PhaseType comparable] interface {
	Cluster

	SetPhase(phase PhaseType)
	GetPhase() PhaseType
}
