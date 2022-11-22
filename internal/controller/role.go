package controller

import (
	"github.com/tarantool/tarantool-operator/apis/v1alpha2"
	"github.com/tarantool/tarantool-operator/internal/implementation"
	"github.com/tarantool/tarantool-operator/pkg/k8s"
	"github.com/tarantool/tarantool-operator/pkg/reconciliation"
)

type RoleController struct {
	*reconciliation.CommonRoleController

	ReplicasetsManger *implementation.ReplicasetsManger
}

func (r *RoleController) GetReplicasetsManger() k8s.ReplicasetsManger[*v1alpha2.Role] {
	return r.ReplicasetsManger
}
