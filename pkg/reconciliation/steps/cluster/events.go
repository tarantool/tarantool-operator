package cluster

import (
	"github.com/tarantool/tarantool-operator/pkg/events"
	corev1 "k8s.io/api/core/v1"
)

const (
	EventUnableToBootstrap = "UnableToBootstrap"
	EventBootstrapped      = "Bootstrapped"
)

func NewUnableToBootstrapEvent(err error) *events.Event {
	return &events.Event{
		EventType: corev1.EventTypeWarning,
		Reason:    EventUnableToBootstrap,
		Message:   err.Error(),
	}
}

func NewBootstrappedEvent() *events.Event {
	return &events.Event{
		EventType: corev1.EventTypeNormal,
		Reason:    EventBootstrapped,
		Message:   "Bootstrapped successfully.",
	}
}
