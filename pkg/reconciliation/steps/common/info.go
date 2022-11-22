package common

import (
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
)

func Info[CtxType Context, CtrlType Controller](message string) *InfoStep[CtxType, CtrlType] {
	return &InfoStep[CtxType, CtrlType]{
		Message: message,
	}
}

type InfoStep[CtxType Context, CtrlType Controller] struct {
	Message string
}

func (r *InfoStep[CtxType, CtrlType]) GetName() string {
	return "Print info"
}

func (r *InfoStep[CtxType, CtrlType]) Reconcile(ctx CtxType, _ CtrlType) (*Result, error) {
	ctx.GetLogger().Info(r.Message)

	return NextStep()
}
