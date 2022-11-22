package reconciliation

import (
	"github.com/tarantool/tarantool-operator/pkg/api"
)

type RoleContext[RoleType api.Role] interface {
	Context

	SetRole(role RoleType)
	GetRole() RoleType
}
