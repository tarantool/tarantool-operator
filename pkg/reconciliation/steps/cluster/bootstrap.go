package cluster

import (
	"strings"

	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
)

type BootstrapStep[
	PhaseType comparable,
	ClusterType api.ClusterWithStatus[PhaseType],
	CtxType ClusterContext[ClusterType],
	CtrlType ClusterController,
] struct {
	ErrorPhase PhaseType
}

func (r *BootstrapStep[PhaseType, ClusterType, CtxType, CtrlType]) GetName() string {
	return "Bootstrap vshard"
}

func (r *BootstrapStep[PhaseType, ClusterType, CtxType, CtrlType]) Reconcile(ctx CtxType, ctrl CtrlType) (*Result, error) {
	if ctx.GetCluster().IsBootstrapped() {
		return NextStep()
	}

	err := ctrl.GetTopology().BootstrapVshard(ctx, ctx.GetLeader())
	if err != nil {
		if strings.Contains(err.Error(), "No remotes with role \"vshard-router\" available") ||
			strings.Contains(err.Error(), "No remotes with role \"vshard-storage\" available") {
			ctrl.GetEventsRecorder().Event(ctx.GetCluster(), NewUnableToBootstrapEvent(err))
			ctx.GetCluster().SetPhase(r.ErrorPhase)

			return Complete()
		}

		return Error(err)
	}

	ctx.GetCluster().MarkBootstrapped()
	ctrl.GetEventsRecorder().Event(ctx.GetCluster(), NewBootstrappedEvent())

	return NextStep()
}
