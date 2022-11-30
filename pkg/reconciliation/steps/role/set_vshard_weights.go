package role

import (
	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
)

type SetVShardWeightsStep[RoleType api.Role, CtxType RoleContext[RoleType], CtrlType RoleController[RoleType]] struct{}

func (r *SetVShardWeightsStep[RoleType, CtxType, CtrlType]) GetName() string {
	return "Set vshard weights"
}

func (r *SetVShardWeightsStep[RoleType, CtxType, CtrlType]) Reconcile(ctx CtxType, ctrl CtrlType) (*Result, error) {
	role := ctx.GetRole()
	config := role.GetVShardConfig()
	topologyClient := ctrl.GetTopology()
	weight := config.GetWeight()

	for i := int32(0); i < role.GetReplicasets(); i++ {
		stsList, err := ctrl.GetResourcesManager().ListStatefulSets(
			ctx,
			role.GetNamespace(),
			ctrl.GetLabelsManager().SelectorByReplicasetOrdinal(role, i),
		)
		if err != nil {
			return Error(err)
		}

		for _, sts := range stsList.Items {
			if sts.GetDeletionTimestamp() != nil {
				continue
			}

			stsLabels := sts.GetLabels()
			stsUUID := stsLabels[ctrl.GetLabelsManager().ReplicasetUUID()]

			err = topologyClient.SetWeight(ctx, ctx.GetLeader(), stsUUID, weight)
			if err != nil {
				return Error(err)
			}
		}
	}

	return NextStep()
}
