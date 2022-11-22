package role

import (
	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
)

type UpdateStatefulSetsStep[RoleType api.Role, CtxType RoleContext[RoleType], CtrlType RoleController[RoleType]] struct{}

func (r *UpdateStatefulSetsStep[RoleType, CtxType, CtrlType]) GetName() string {
	return "Update StatefulSets"
}

func (r *UpdateStatefulSetsStep[RoleType, CtxType, CtrlType]) Reconcile(ctx CtxType, ctrl CtrlType) (*Result, error) {
	role := ctx.GetRole()
	cluster := ctx.GetRelatedCluster()

	complete, err := ctrl.GetReplicasetsManger().UpdateStatefulSets(ctx, cluster, role)
	if err != nil || !complete {
		return Error(err)
	}

	return NextStep()
}
