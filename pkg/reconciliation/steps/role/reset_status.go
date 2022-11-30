package role

import (
	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ResetStatusStep[RoleType api.Role, CtxType RoleContext[RoleType], CtrlType RoleController[RoleType]] struct{}

func (r *ResetStatusStep[RoleType, CtxType, CtrlType]) GetName() string {
	return "Reset role status"
}

func (r *ResetStatusStep[RoleType, CtxType, CtrlType]) Reconcile(ctx CtxType, _ CtrlType) (*ctrl.Result, error) {
	ctx.GetRole().ResetStatus()

	return NextStep()
}
