package role

import (
	"github.com/tarantool/tarantool-operator/pkg/events"
	"github.com/tarantool/tarantool-operator/pkg/topology"
	corev1 "k8s.io/api/core/v1"
)

const (
	EventTypeWrongVShardRoles = "WrongVShardRoles"
)

func NewWrongVShardRolesEvent(err *topology.UnknownRoleError) *events.Event {
	return &events.Event{
		EventType: corev1.EventTypeWarning,
		Reason:    EventTypeWrongVShardRoles,
		Message:   err.Error(),
	}
}
