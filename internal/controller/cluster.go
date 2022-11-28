package controller

import (
	"context"

	"github.com/tarantool/tarantool-operator/apis/v1beta1"
	"github.com/tarantool/tarantool-operator/pkg/api"
	"github.com/tarantool/tarantool-operator/pkg/reconciliation"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ClusterController struct {
	*reconciliation.CommonClusterController
}

func (r *ClusterController) IsAllRolesAtPhase(ctx context.Context, cluster api.Cluster, phase v1beta1.RolePhase) (bool, error) {
	selector := r.LabelsManager.SelectorByClusterName(cluster)

	roleList := &v1beta1.RoleList{}

	err := r.List(ctx, roleList, &client.ListOptions{LabelSelector: selector, Namespace: cluster.GetNamespace()})
	if err != nil {
		return false, err
	}

	if len(roleList.Items) == 0 {
		return false, nil
	}

	for _, role := range roleList.Items {
		if role.GetPhase() != phase {
			return false, nil
		}
	}

	return true, nil
}
