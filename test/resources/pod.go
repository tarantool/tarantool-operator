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
		panic(fmt.Errorf("role %s not added to fake cartrdige", roleName))
	}

	statefulSets, ok := r.StatefulSets[roleName]
	if !ok {
		panic(fmt.Errorf("stateful sets for role %s not added to fake cartrdige", roleName))
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

func (r *FakeCartridge) WithAllPodsReady() *FakeCartridge {
	for _, pod := range r.Pods {
		r.setPodReady(pod)
		r.setPodContainerReady(pod, PodContainerName)
	}

	return r
}

func (r *FakeCartridge) setPodReady(pod *v1.Pod) {
	if pod.Status.Conditions == nil {
		pod.Status.Conditions = []v1.PodCondition{}
	}

	pod.Status.Phase = v1.PodRunning
	pod.Status.Conditions = append(pod.Status.Conditions, v1.PodCondition{
		Type:    v1.PodReady,
		Status:  v1.ConditionTrue,
		Reason:  "Ready",
		Message: "Ready",
	})
}

func (r *FakeCartridge) setPodContainerReady(pod *v1.Pod, containerName string) {
	if pod.Status.ContainerStatuses == nil {
		pod.Status.ContainerStatuses = []v1.ContainerStatus{}
	}

	started := true
	pod.Status.ContainerStatuses = append(
		pod.Status.ContainerStatuses,
		v1.ContainerStatus{
			Name:    containerName,
			Ready:   true,
			Started: &started,
		},
	)
}
