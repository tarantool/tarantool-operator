package role

import (
	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
)

type CreateStatefulSetsStep[RoleType api.Role, CtxType RoleContext[RoleType], CtrlType RoleController[RoleType]] struct{}

func (r *CreateStatefulSetsStep[RoleType, CtxType, CtrlType]) GetName() string {
	return "Create StatefulSets"
}

func (r *CreateStatefulSetsStep[RoleType, CtxType, CtrlType]) Reconcile(ctx CtxType, ctrl CtrlType) (*Result, error) {
	role := ctx.GetRole()
	cluster := ctx.GetRelatedCluster()

	if ctx.GetRole().GetDeletionTimestamp() != nil {
		return NextStep()
	}

	err := ctrl.GetReplicasetsManger().CreateStatefulSets(ctx, cluster, role)
	if err != nil {
		return Error(err)
	}

	return NextStep()
}
