package cluster

import (
	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
)

type ResetStatusStep[ClusterType api.Cluster, CtxType ClusterContext[ClusterType], CtrlType ClusterController] struct{}

func (r *ResetStatusStep[ClusterType, CtxType, CtrlType]) GetName() string {
	return "Reset cluster status"
}

func (r *ResetStatusStep[ClusterType, CtxType, CtrlType]) Reconcile(ctx CtxType, _ CtrlType) (*Result, error) {
	ctx.GetCluster().ResetStatus()

	return NextStep()
}
