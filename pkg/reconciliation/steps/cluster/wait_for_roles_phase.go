package cluster

import (
	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
)

type WaitForRolesPhaseStep[
	RolePhaseType comparable,
	ClusterType api.Cluster,
	CtxType ClusterContext[ClusterType],
	CtrlType ClusterWithRolesController[RolePhaseType],
] struct {
	ExpectedPhases []RolePhaseType
}

func (r *WaitForRolesPhaseStep[RolePhaseType, ClusterType, CtxType, CtrlType]) GetName() string {
	return "Wait roles to be in phase"
}

func (r *WaitForRolesPhaseStep[RolePhaseType, ClusterType, CtxType, CtrlType]) Reconcile(ctx CtxType, ctrl CtrlType) (*Result, error) {
	for _, expectedPhase := range r.ExpectedPhases {
		allRolesAtCorrectPhase, err := ctrl.IsAllRolesAtPhase(ctx, ctx.GetCluster(), expectedPhase)
		if err != nil {
			return Error(err)
		}

		if allRolesAtCorrectPhase {
			return NextStep()
		}
	}

	return Complete()
}
