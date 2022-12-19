package resources

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *FakeCartridge) WithPodsCreated(roleName string) *FakeCartridge {
	cluster := r.Cluster

	role, ok := r.Roles[roleName]
	if !ok {
		panic(fmt.Errorf("role %s not added to fake cartridge", roleName))
	}

	statefulSets, ok := r.StatefulSets[roleName]
	if !ok {
		panic(fmt.Errorf("stateful sets for role %s not added to fake cartridge", roleName))
	}

	for _, sts := range statefulSets {
		for podNum := int32(0); podNum < *sts.Spec.Replicas; podNum++ {
			pod := &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("%s-%d", sts.Name, podNum),
					Namespace: sts.Namespace,
					Labels: map[string]string{
						r.labelsManager.ClusterName(): cluster.Name,
						r.labelsManager.RoleName():    role.Name,
					},
				},
				Status: v1.PodStatus{
					Phase: v1.PodPending,
				},
			}
			r.Pods = append(r.Pods, pod)
			r.object(pod)
		}
	}

	return r
}

func (r *FakeCartridge) WithRouterPodsCreated() *FakeCartridge {
	return r.WithPodsCreated(RoleRouter)
}

func (r *FakeCartridge) WithStoragePodsCreated() *FakeCartridge {
	return r.WithPodsCreated(RoleStorage)
}

func (r *FakeCartridge) WithAllPodsRunning() *FakeCartridge {
	for _, pod := range r.Pods {
		r.setPodRunning(pod)
	}

	return r
}

func (r *FakeCartridge) WithPodsRunning(names ...string) *FakeCartridge {
	namesMap := make(map[string]bool, len(names))
	for _, name := range names {
		namesMap[name] = true
	}

	for _, pod := range r.Pods {
		_, ok := namesMap[pod.GetName()]
		if ok {
			r.setPodRunning(pod)
		}
	}

	return r
}

func (r *FakeCartridge) WithAllPodsDeleting() *FakeCartridge {
	for _, pod := range r.Pods {
		r.setPodDeleting(pod)
	}

	return r
}

func (r *FakeCartridge) setPodRunning(pod *v1.Pod) {
	pod.Status.Phase = v1.PodRunning
}

func (r *FakeCartridge) setPodDeleting(pod *v1.Pod) {
	now := metav1.Now()
	pod.DeletionTimestamp = &now
}
