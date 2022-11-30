package cluster

import (
	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
)

type CheckDeletionStep[ClusterType api.Cluster, CtxType ClusterContext[ClusterType], CtrlType ClusterController] struct{}

func (r *CheckDeletionStep[ClusterType, CtxType, CtrlType]) GetName() string {
	return "Check cluster deletion"
}

func (r *CheckDeletionStep[ClusterType, CtxType, CtrlType]) Reconcile(ctx CtxType, ctrl CtrlType) (*Result, error) {
	if ctx.GetCluster().GetDeletionTimestamp() != nil {
		return Complete()
	}

	return NextStep()
}
