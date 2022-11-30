package cluster

import (
	"github.com/google/go-cmp/cmp"
	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SyncClusterWideServiceStep[ClusterType api.Cluster, CtxType ClusterContext[ClusterType], CtrlType ClusterController] struct {
	Phase string
}

func (r *SyncClusterWideServiceStep[ClusterType, CtxType, CtrlType]) GetName() string {
	return "Sync cluster-wide service"
}

func (r *SyncClusterWideServiceStep[ClusterType, CtxType, CtrlType]) Reconcile(ctx CtxType, ctrl CtrlType) (*Result, error) {
	ctx.GetLogger().Info("SyncClusterWideServiceStep")

	cluster := ctx.GetCluster()

	saveFunc := ctrl.GetResourcesManager().UpdateObject

	svc, err := ctrl.GetResourcesManager().GetService(ctx, cluster.GetNamespace(), cluster.GetName())
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return Error(err)
		}

		saveFunc = ctrl.GetResourcesManager().CreateObject
		svc = &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cluster.GetName(),
				Namespace: cluster.GetNamespace(),
			},
		}
	}

	changed, err := ctrl.GetResourcesManager().ControlObject(cluster, svc)
	if err != nil {
		return Error(err)
	}

	selector := map[string]string{
		ctrl.GetLabelsManager().ClusterName(): cluster.GetName(),
	}
	if !cmp.Equal(svc.Spec.Selector, selector) {
		svc.Spec.Selector = selector
		changed = true
	}

	if !svc.Spec.PublishNotReadyAddresses {
		svc.Spec.PublishNotReadyAddresses = true
		changed = true
	}

	if svc.Spec.ClusterIP != "None" {
		svc.Spec.ClusterIP = "None"
		changed = true
	}

	ports := []corev1.ServicePort{
		{
			Name:     "app",
			Port:     cluster.GetListenPort(),
			Protocol: "TCP",
		},
		{
			Name:     "app-udp",
			Port:     cluster.GetListenPort(),
			Protocol: "UDP",
		},
	}

	if !cmp.Equal(svc.Spec.Ports, ports) {
		svc.Spec.Ports = ports
		changed = true
	}

	if changed {
		err = saveFunc(ctx, svc)
		if err != nil {
			return Error(err)
		}
	}

	return NextStep()
}
