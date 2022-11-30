package reconciliation

import (
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
)

type (
	Manager = ctrl.Manager
	Result  = ctrl.Result
	Request = ctrl.Request
)

var NewControllerManagedBy = ctrl.NewControllerManagedBy

func Requeue(requeueTime time.Duration) (*Result, error) {
	return &ctrl.Result{RequeueAfter: requeueTime}, nil
}

func Error(err error) (*Result, error) {
	return nil, err
}

func Break(err error) (*Result, error) {
	return &Result{}, err
}

func NextStep() (*Result, error) {
	return nil, nil
}

func Complete() (*Result, error) {
	return &Result{}, nil
}
