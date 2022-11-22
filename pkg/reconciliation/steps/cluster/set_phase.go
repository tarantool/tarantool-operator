package cluster

import (
	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
)

type SetPhaseStep[PhaseType comparable, ClusterType api.ClusterWithStatus[PhaseType], CtxType ClusterContext[ClusterType], CtrlType ClusterController] struct {
	Phase PhaseType
}

func (r *SetPhaseStep[PhaseType, ClusterType, CtxType, CtrlType]) GetName() string {
	return "Set cluster phase"
}

func (r *SetPhaseStep[PhaseType, ClusterType, CtxType, CtrlType]) Reconcile(ctx CtxType, _ CtrlType) (*Result, error) {
	ctx.GetCluster().SetPhase(r.Phase)

	return NextStep()
}
