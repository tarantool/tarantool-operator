package resources

import (
	"github.com/tarantool/tarantool-operator/apis/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *FakeCartridge) WithRole(name string, replicasets, replicas int32) *FakeCartridge {
	cluster := r.Cluster
	r.Roles[name] = &v1beta1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: cluster.GetNamespace(),
			Labels: map[string]string{
				r.labelsManager.ClusterName(): cluster.GetName(),
				r.labelsManager.RoleName():    name,
			},
		},
		Spec: v1beta1.RoleSpec{
			Replicasets: &replicasets,
			VShard: v1beta1.RoleVShardConfig{
				Weight: ConstDefaultReplicasetWeight,
			},
			ReplicasetTemplate: &v1beta1.ReplicasetTemplate{
				Replicas:             &replicas,
				VolumeClaimTemplates: []v1.PersistentVolumeClaim{},
				PodTemplate: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							r.labelsManager.ClusterName(): cluster.GetName(),
						},
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							r.NewCartridgeContainer(),
						},
					},
				},
			},
		},
	}
	r.object(r.Roles[name])

	return r
}

func (r *FakeCartridge) WithRouterRole(replicasets, replicas int32) *FakeCartridge {
	return r.WithRole(RoleRouter, replicasets, replicas)
}

func (r *FakeCartridge) WithStorageRole(replicasets, replicas int32) *FakeCartridge {
	return r.WithRole(RoleStorage, replicasets, replicas)
}

func (r *FakeCartridge) WithAllRolesInPhase(phase v1beta1.RolePhase) *FakeCartridge {
	for _, role := range r.Roles {
		role.Status.Phase = phase
	}

	return r
}
