package reconciliation

import ctrl "sigs.k8s.io/controller-runtime"

type Step[CtxType Context, CtrlType Controller] interface {
	GetName() string
	Reconcile(ctx CtxType, ctrl CtrlType) (*ctrl.Result, error)
}
