package resources

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *FakeCartridge) WithStatefulSetsCreated(roleName string) *FakeCartridge {
	cluster := r.Cluster

	role, ok := r.Roles[roleName]
	if !ok {
		panic(fmt.Errorf("role %s not added to fake cartridge", roleName))
	}

	r.StatefulSets[roleName] = map[string]*appsv1.StatefulSet{}

	for stsNum := int32(0); stsNum < *role.Spec.Replicasets; stsNum++ {
		stsName := fmt.Sprintf("%s-%d", role.Name, stsNum)
		r.StatefulSets[roleName][stsName] = &appsv1.StatefulSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      stsName,
				Namespace: cluster.Namespace,
				Labels: map[string]string{
					r.labelsManager.ClusterName(): cluster.Name,
				},
			},
			Spec: appsv1.StatefulSetSpec{
				Replicas:             role.Spec.ReplicasetTemplate.Replicas,
				Template:             role.Spec.ReplicasetTemplate.PodTemplate,
				VolumeClaimTemplates: role.Spec.ReplicasetTemplate.VolumeClaimTemplates,
			},
			Status: appsv1.StatefulSetStatus{},
		}
		r.object(r.StatefulSets[roleName][stsName])
	}

	return r
}

func (r *FakeCartridge) WithRouterStatefulSetsCreated() *FakeCartridge {
	return r.WithStatefulSetsCreated(RoleRouter)
}

func (r *FakeCartridge) WithStorageStatefulSetsCreated() *FakeCartridge {
	return r.WithStatefulSetsCreated(RoleStorage)
}
