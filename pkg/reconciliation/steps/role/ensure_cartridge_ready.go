package role

import (
	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
	"github.com/tarantool/tarantool-operator/pkg/utils"
)

type EnsureCartridgeReadyStep[RoleType api.Role, CtxType RoleContext[RoleType], CtrlType RoleController[RoleType]] struct{}

func (r *EnsureCartridgeReadyStep[RoleType, CtxType, CtrlType]) GetName() string {
	return "Ensure cartridge ready"
}

func (r *EnsureCartridgeReadyStep[RoleType, CtxType, CtrlType]) Reconcile(ctx CtxType, ctrl CtrlType) (*Result, error) {
	allPodsRunning, err := r.IsCartridgeReady(ctx, ctrl)
	if err != nil || !allPodsRunning {
		ctx.GetLogger().Info("Not all pods of role are running. Wait for it.")

		return Error(err)
	}

	return NextStep()
}

func (r *EnsureCartridgeReadyStep[RoleType, CtxType, CtrlType]) IsCartridgeReady(ctx CtxType, ctrl CtrlType) (bool, error) {
	cluster := ctx.GetRelatedCluster()
	role := ctx.GetRole()

	selector := ctrl.GetLabelsManager().SelectorByRoleName(role)

	pods, err := ctrl.GetResourcesManager().ListPods(ctx, cluster.GetNamespace(), selector)
	if err != nil {
		return false, err
	}

	expectedCount := role.GetReplicasets() * role.GetReplicas()
	if expectedCount > int32(len(pods.Items)) {
		return false, nil
	}

	readyCount := int32(0)

	for key := range pods.Items {
		if utils.IsPodDeleting(&pods.Items[key]) {
			continue
		}

		if !utils.IsPodDefaultContainerReady(&pods.Items[key]) {
			continue
		}

		readyCount++
	}

	return readyCount >= expectedCount, nil
}
