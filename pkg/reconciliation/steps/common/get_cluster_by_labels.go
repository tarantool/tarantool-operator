package common

import (
	"fmt"
	"time"

	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
)

func GetClusterByLabels[CtxType Context, CtrlType Controller]() *GetClusterByLabelsStep[CtxType, CtrlType] {
	return &GetClusterByLabelsStep[CtxType, CtrlType]{}
}

type GetClusterByLabelsStep[CtxType Context, CtrlType Controller] struct{}

func (r *GetClusterByLabelsStep[CtxType, CtrlType]) GetName() string {
	return "Get cluster by labels"
}

func (r *GetClusterByLabelsStep[CtxType, CtrlType]) Reconcile(ctx CtxType, ctrl CtrlType) (*Result, error) {
	if ctx.GetRelatedCluster() != nil {
		return NextStep()
	}

	object := ctx.GetRequestedObject()
	objectLabels := object.GetLabels()

	if objectLabels == nil {
		return Break(fmt.Errorf(" %s label is required", ctrl.GetLabelsManager().ClusterName()))
	}

	clusterName, ok := objectLabels[ctrl.GetLabelsManager().ClusterName()]
	if !ok || clusterName == "" {
		return Break(fmt.Errorf(" %s label is required", ctrl.GetLabelsManager().ClusterName()))
	}

	cluster, err := ctrl.GetResourcesManager().GetCluster(ctx, object.GetNamespace(), clusterName)
	if err != nil {
		return Error(err)
	}

	if err != nil {
		if apiErrors.IsNotFound(err) {
			return Requeue(10 * time.Second)
		}

		return Error(err)
	}

	ctx.SetRelatedCluster(cluster)

	return NextStep()
}
