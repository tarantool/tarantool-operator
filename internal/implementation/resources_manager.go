package implementation

import (
	"context"

	"github.com/tarantool/tarantool-operator/apis/v1alpha2"
	"github.com/tarantool/tarantool-operator/pkg/api"
	"github.com/tarantool/tarantool-operator/pkg/k8s"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ResourcesManager struct {
	*k8s.CommonResourcesManager

	LabelsManager k8s.LabelsManager
}

func (r *ResourcesManager) GetCluster(ctx context.Context, ns, name string) (api.Cluster, error) {
	cluster := &v1alpha2.Cluster{}
	selector := types.NamespacedName{
		Namespace: ns,
		Name:      name,
	}

	err := r.Get(ctx, selector, cluster)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func (r *ResourcesManager) GetClusterRoles(ctx context.Context, cluster api.Cluster) ([]api.Role, error) {
	selector := r.LabelsManager.SelectorByClusterName(cluster)

	roleList := &v1alpha2.RoleList{}

	err := r.List(ctx, roleList, &client.ListOptions{LabelSelector: selector, Namespace: cluster.GetNamespace()})
	if err != nil {
		return nil, err
	}

	result := make([]api.Role, len(roleList.Items))
	for k := range roleList.Items {
		result[k] = &roleList.Items[k]
	}

	return result, nil
}
