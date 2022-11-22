package utils

import v1 "k8s.io/api/core/v1"

func IsPodDeleting(pod *v1.Pod) bool {
	return pod.DeletionTimestamp != nil
}

func IsPodRunning(pod *v1.Pod) bool {
	return pod.Status.Phase == v1.PodRunning
}

func IsPodDefaultContainerReady(pod *v1.Pod) bool {
	if !IsPodRunning(pod) {
		return false
	}

	if pod.Status.ContainerStatuses == nil || len(pod.Status.ContainerStatuses) == 0 {
		return false
	}

	return pod.Status.ContainerStatuses[0].Ready
}

// IsPodReady returns true if a pod is ready; false otherwise.
func IsPodReady(pod *v1.Pod) bool {
	return IsPodReadyConditionTrue(pod.Status)
}

// IsPodReadyConditionTrue returns true if a pod is ready; false otherwise.
func IsPodReadyConditionTrue(status v1.PodStatus) bool {
	condition := GetPodReadyCondition(status)

	return condition != nil && condition.Status == v1.ConditionTrue
}

// GetPodReadyCondition extracts the pod ready condition from the given status and returns that.
// Returns nil if the condition is not present.
func GetPodReadyCondition(status v1.PodStatus) *v1.PodCondition {
	_, condition := GetPodCondition(&status, v1.PodReady)

	return condition
}

// GetPodCondition extracts the provided condition from the given status and returns that.
// Returns nil and -1 if the condition is not present, and the index of the located condition.
func GetPodCondition(status *v1.PodStatus, conditionType v1.PodConditionType) (int, *v1.PodCondition) {
	if status == nil {
		return -1, nil
	}

	return GetPodConditionFromList(status.Conditions, conditionType)
}

// GetPodConditionFromList extracts the provided condition from the given list of condition and
// returns the index of the condition and the condition. Returns -1 and nil if the condition is not present.
func GetPodConditionFromList(
	conditions []v1.PodCondition,
	conditionType v1.PodConditionType,
) (int, *v1.PodCondition) {
	if conditions == nil {
		return -1, nil
	}

	for i := range conditions {
		if conditions[i].Type == conditionType {
			return i, &conditions[i]
		}
	}

	return -1, nil
}
