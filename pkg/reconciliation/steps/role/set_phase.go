package role

import (
	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
)

type SetPhaseStep[PhaseType comparable, RoleType api.RoleWithStatus[PhaseType], CtxType RoleContext[RoleType], CtrlType RoleController[RoleType]] struct {
	Phase PhaseType
}

func (r *SetPhaseStep[PhaseType, RoleType, CtxType, CtrlType]) GetName() string {
	return "Set role phase"
}

func (r *SetPhaseStep[PhaseType, RoleType, CtxType, CtrlType]) Reconcile(ctx CtxType, _ CtrlType) (*Result, error) {
	ctx.GetRole().SetPhase(r.Phase)

	return NextStep()
}
