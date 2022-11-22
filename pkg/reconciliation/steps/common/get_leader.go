package common

import (
	"time"

	"github.com/pkg/errors"
	"github.com/tarantool/tarantool-operator/pkg/election"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
)

func GetLeader[CtxType Context, CtrlType Controller]() *GetLeaderStep[CtxType, CtrlType] {
	return &GetLeaderStep[CtxType, CtrlType]{}
}

type GetLeaderStep[CtxType Context, CtrlType Controller] struct {
	Message string
}

func (r *GetLeaderStep[CtxType, CtrlType]) GetName() string {
	return "Get leader"
}

func (r *GetLeaderStep[CtxType, CtrlType]) Reconcile(ctx CtxType, ctrl CtrlType) (*Result, error) {
	leader, err := ctrl.GetLeaderElection().GetLeaderInstance(ctx, ctx.GetRelatedCluster())
	if err != nil {
		if errors.Is(err, election.ErrLeaderElectionConflict) {
			return Requeue(time.Second * 10)
		}

		if errors.Is(err, election.ErrNoAvailableLeader) && !ctx.GetRelatedCluster().IsBootstrapped() {
			return Requeue(time.Second * 10)
		}

		return Error(err)
	}

	ctx.SetLeader(leader)

	return NextStep()
}
