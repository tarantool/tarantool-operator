package reconciliation

import (
	"github.com/tarantool/tarantool-operator/pkg/api"
	"github.com/tarantool/tarantool-operator/pkg/k8s"
)

type RoleController[RoleType api.Role] interface {
	Controller

	GetReplicasetsManger() k8s.ReplicasetsManger[RoleType]
}

type CommonRoleController struct {
	*CommonController
}
