package election

import (
	"fmt"

	"github.com/tarantool/tarantool-operator/pkg/events"
	corev1 "k8s.io/api/core/v1"
)

const (
	EventNewTopologyLeaderElected = "NewTopologyLeaderElected"
)

func NewTopologyLeaderElectedEvent(leader *corev1.Pod) *events.Event {
	return &events.Event{
		EventType: corev1.EventTypeNormal,
		Reason:    EventNewTopologyLeaderElected,
		Message:   fmt.Sprintf("Pod %s is a new topology leader", leader.GetName()),
	}
}
