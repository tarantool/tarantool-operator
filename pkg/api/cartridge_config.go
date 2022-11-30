package api

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CartridgeConfig interface {
	client.Object

	GetData() []byte

	ResetStatus()
}

type CartridgeConfigWithStatus[PhaseType comparable] interface {
	CartridgeConfig

	SetPhase(phase PhaseType)
	GetPhase() PhaseType
}
