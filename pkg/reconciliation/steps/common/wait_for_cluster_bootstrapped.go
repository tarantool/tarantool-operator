package common

import (
	"time"

	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
)

func WaitForClusterBootstrapped[CtxType Context, CtrlType Controller]() *WaitForClusterBootstrappedStep[CtxType, CtrlType] {
	return &WaitForClusterBootstrappedStep[CtxType, CtrlType]{}
}

type WaitForClusterBootstrappedStep[CtxType Context, CtrlType Controller] struct{}

func (r *WaitForClusterBootstrappedStep[CtxType, CtrlType]) GetName() string {
	return "Wait cluster to be bootstrapped"
}

func (r *WaitForClusterBootstrappedStep[CtxType, CtrlType]) Reconcile(ctx CtxType, _ CtrlType) (*Result, error) {
	cluster := ctx.GetRelatedCluster()
	if cluster == nil || !cluster.IsBootstrapped() {
		return Requeue(10 * time.Second)
	}

	return NextStep()
}
