package api

import v1 "k8s.io/api/core/v1"

type FailoverMode string

const (
	FailoverModeDisabled FailoverMode = "disabled"
	FailoverModeEventual FailoverMode = "eventual"
	FailoverModeStateful FailoverMode = "stateful"
	FailoverModeRaft     FailoverMode = "raft"
)

type FailoverStateProvider string

const (
	FailoverStateProviderETCD2      FailoverStateProvider = "etcd2"
	FailoverStateProviderStateboard FailoverStateProvider = "stateboard"
)

type FailoverConfig interface {
	GetMode() FailoverMode
	GetTimeout() int32
	GetStateProvider() FailoverStateProvider
	GetFencing() bool
	GetFencingTimeout() int32
	GetFencingPause() int32

	GetETCD2Config() FailoverETCD2Config
	GetStateboardConfig() FailoverStateboardConfig
}

type FailoverETCD2Config interface {
	GetEndpoints() []string
	GetUsername() string
	GetPassword() v1.SecretReference
	GetLockDelay() int32
	GetPrefix() string
}

type FailoverStateboardConfig interface {
	GetURI() string
	GetPassword() v1.SecretReference
}
