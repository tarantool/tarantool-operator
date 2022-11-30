package k8s

import (
	"fmt"

	"github.com/tarantool/tarantool-operator/pkg/api"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type LabelsManager interface {
	ClusterName() string
	RoleName() string
	ReplicasetName() string
	ReplicasetUUID() string
	ReplicasetOrdinal() string
	ReplicasetPodTemplateHash() string

	SelectorByClusterName(cluster api.Cluster) labels.Selector
	SelectorByRoleName(role api.Role) labels.Selector
	SelectorByReplicasetOrdinal(role api.Role, ordinal int32) labels.Selector
	SelectorByReplicasetName(role api.Role, name string) labels.Selector
}

type NamespacedLabelsManager struct {
	Namespace string
}

func (r *NamespacedLabelsManager) namespacedLabel(label string) string {
	return fmt.Sprintf("%s/%s", r.Namespace, label)
}

func (r *NamespacedLabelsManager) getClusterNameFromLabels(obj client.Object) string {
	objLabels := obj.GetLabels()

	clusterName, ok := objLabels[r.ClusterName()]
	if !ok {
		return ""
	}

	return clusterName
}

func (r *NamespacedLabelsManager) ClusterName() string {
	return r.namespacedLabel("cluster-name")
}

func (r *NamespacedLabelsManager) RoleName() string {
	return r.namespacedLabel("role-name")
}

func (r *NamespacedLabelsManager) ReplicasetName() string {
	return r.namespacedLabel("replicaset-name")
}

func (r *NamespacedLabelsManager) ReplicasetUUID() string {
	return r.namespacedLabel("replicaset-uuid")
}

func (r *NamespacedLabelsManager) ReplicasetOrdinal() string {
	return r.namespacedLabel("replicaset-ordinal")
}

func (r *NamespacedLabelsManager) ReplicasetPodTemplateHash() string {
	return r.namespacedLabel("replicaset-pod-template-hash")
}

func (r *NamespacedLabelsManager) SelectorByClusterName(cluster api.Cluster) labels.Selector {
	return labels.SelectorFromSet(map[string]string{
		r.ClusterName(): cluster.GetName(),
	})
}

func (r *NamespacedLabelsManager) SelectorByRoleName(role api.Role) labels.Selector {
	return labels.SelectorFromSet(map[string]string{
		r.ClusterName(): r.getClusterNameFromLabels(role),
		r.RoleName():    role.GetName(),
	})
}

func (r *NamespacedLabelsManager) SelectorByReplicasetOrdinal(role api.Role, ordinal int32) labels.Selector {
	return labels.SelectorFromSet(map[string]string{
		r.ClusterName():       r.getClusterNameFromLabels(role),
		r.RoleName():          role.GetName(),
		r.ReplicasetOrdinal(): fmt.Sprintf("%d", ordinal),
	})
}

func (r *NamespacedLabelsManager) SelectorByReplicasetName(role api.Role, name string) labels.Selector {
	return labels.SelectorFromSet(map[string]string{
		r.ClusterName():    r.getClusterNameFromLabels(role),
		r.RoleName():       role.GetName(),
		r.ReplicasetName(): name,
	})
}
