package reconciliation

import (
	"fmt"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ErrorTimeout                    = 10 * time.Second
	SteppedReconcilerVerbosityLevel = 1
)

type SteppedReconciler[CtxType Context, CtrlType Controller] struct {
	client.Client

	Controller CtrlType
}

func (r *SteppedReconciler[CtxType, CtrlType]) Run(ctx CtxType, steps ...Step[CtxType, CtrlType]) (ctrl.Result, error) {
	var (
		res      *ctrl.Result
		err      error
		lastStep string
	)

	for _, step := range steps {
		lastStep = step.GetName()

		ctx.GetLogger().V(SteppedReconcilerVerbosityLevel).Info(fmt.Sprintf("Execute `%s` step", step.GetName()))
		res, err = step.Reconcile(ctx, r.Controller)

		if err != nil || res != nil {
			break
		}
	}

	updErr := r.UpdateStatus(ctx)
	if updErr != nil {
		return ctrl.Result{RequeueAfter: ErrorTimeout}, updErr
	}

	if err != nil && res == nil {
		res = &ctrl.Result{RequeueAfter: ErrorTimeout}
	}

	if res == nil {
		res = &ctrl.Result{}
	}

	if err != nil {
		ctx.GetLogger().V(SteppedReconcilerVerbosityLevel).Error(err, fmt.Sprintf("An error occurred during `%s` step", lastStep))
	}

	if res.Requeue || res.RequeueAfter > 0 {
		ctx.GetLogger().V(SteppedReconcilerVerbosityLevel).Info(fmt.Sprintf("Requeue at `%s` step", lastStep))
	} else {
		ctx.GetLogger().V(SteppedReconcilerVerbosityLevel).Info(fmt.Sprintf("Finish cycle at `%s` step", lastStep))
	}

	return *res, err
}

func (r *SteppedReconciler[CtxType, CtrlType]) UpdateStatus(ctx CtxType) error {
	if ctx.HasRequestedObject() {
		obj := ctx.GetRequestedObject()

		if obj.GetDeletionTimestamp() == nil || len(obj.GetFinalizers()) > 0 {
			return r.Status().Update(ctx, obj)
		}
	}

	return nil
}
