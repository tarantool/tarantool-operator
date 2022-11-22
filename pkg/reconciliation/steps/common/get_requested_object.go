package common

import (
	"github.com/pkg/errors"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetRequestedObject[CtxType Context, CtrlType Controller](target client.Object) *GetRequestedObjectStep[CtxType, CtrlType] {
	return &GetRequestedObjectStep[CtxType, CtrlType]{
		Target: target,
	}
}

type GetRequestedObjectStep[CtxType Context, CtrlType Controller] struct {
	Target client.Object
}

func (r *GetRequestedObjectStep[CtxType, CtrlType]) GetName() string {
	return "Get requested object"
}

func (r *GetRequestedObjectStep[CtxType, CtrlType]) Reconcile(ctx CtxType, controller CtrlType) (*Result, error) {
	err := controller.Get(ctx, ctx.GetRequest().NamespacedName, r.Target)
	if err != nil {
		if !apiErrors.IsNotFound(err) {
			return nil, errors.Wrap(err, "unable to retrieve object for reconcile")
		}

		return Complete()
	}

	err = ctx.SetRequestedObject(r.Target)
	if err != nil {
		return Error(err)
	}

	return NextStep()
}
