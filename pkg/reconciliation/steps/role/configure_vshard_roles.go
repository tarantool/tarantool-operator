package role

import (
	"time"

	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
	"github.com/tarantool/tarantool-operator/pkg/utils"
)

type ConfigureVShardRolesStep[RoleType api.Role, CtxType RoleContext[RoleType], CtrlType RoleController[RoleType]] struct{}

func (r *ConfigureVShardRolesStep[RoleType, CtxType, CtrlType]) GetName() string {
	return "Check role orphan deletion"
}

func (r *ConfigureVShardRolesStep[RoleType, CtxType, CtrlType]) Reconcile(ctx CtxType, ctrl CtrlType) (*Result, error) {
	role := ctx.GetRole()
	topologyClient := ctrl.GetTopology()

	allRolesConfigured := true

	for stsOrdinal := int32(0); stsOrdinal < role.GetReplicasets(); stsOrdinal++ {
		stsList, err := ctrl.GetResourcesManager().ListStatefulSets(
			ctx,
			role.GetNamespace(),
			ctrl.GetLabelsManager().SelectorByReplicasetOrdinal(role, stsOrdinal),
		)
		if err != nil {
			allRolesConfigured = false

			continue
		}

		for _, sts := range stsList.Items {
			if sts.GetDeletionTimestamp() != nil {
				continue
			}

			hierarchy, err := topologyClient.GetRolesHierarchy(ctx, ctx.GetLeader())
			if err != nil {
				allRolesConfigured = false

				continue
			}

			replicasetUUID := ctrl.GetReplicasetsManger().GetReplicasetUUID(role, stsOrdinal)

			actualRoles, err := topologyClient.GetReplicasetRoles(ctx, ctx.GetLeader(), replicasetUUID)
			if err != nil {
				allRolesConfigured = false

				continue
			}

			vshardConfig := role.GetVShardConfig()
			desiredRoles := vshardConfig.GetRoles()

			if !utils.IsVShardRolesEquals(actualRoles, desiredRoles, hierarchy) {
				err = topologyClient.SetReplicasetRoles(ctx, ctx.GetLeader(), replicasetUUID, vshardConfig.GetRoles())
				if err != nil {
					allRolesConfigured = false

					continue
				}
			}
		}
	}

	if !allRolesConfigured {
		ctx.GetLogger().Info("Not all vshard roles configured")

		return Requeue(10 * time.Second)
	}

	return NextStep()
}
