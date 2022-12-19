package implementation

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/tarantool/tarantool-operator/apis/v1beta1"
	"github.com/tarantool/tarantool-operator/pkg/api"
	"github.com/tarantool/tarantool-operator/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ReplicasetsManger struct {
	*ResourcesManager

	UUIDSpace uuid.UUID
}

func (r *ReplicasetsManger) GetReplicasetUUID(role *v1beta1.Role, ordinal int32) string {
	replicasetUUID := uuid.NewSHA1(
		r.UUIDSpace,
		[]byte(fmt.Sprintf("%s-%d", role.GetName(), ordinal)),
	)

	return replicasetUUID.String()
}

func (r *ReplicasetsManger) GetAdvertiseURI(cluster api.Cluster, pod *v1.Pod) string {
	return fmt.Sprintf(
		"%s.%s.%s.svc.%s:%d",
		pod.GetObjectMeta().GetName(),      // Instance name
		cluster.GetName(),                  // Cartridge cluster name
		pod.GetObjectMeta().GetNamespace(), // Namespace
		cluster.GetDomain(),                // Cluster domain name
		cluster.GetListenPort(),
	)
}

func (r *ReplicasetsManger) CreateStatefulSets(ctx context.Context, cluster api.Cluster, role *v1beta1.Role) error {
	for ordinal := int32(0); ordinal < role.GetReplicasets(); ordinal++ {
		selector := r.LabelsManager.SelectorByReplicasetOrdinal(role, ordinal)

		stsList, err := r.ListStatefulSets(ctx, role.GetNamespace(), selector)
		if err == nil && len(stsList.Items) > 0 {
			continue
		}

		if err != nil && !apierrors.IsNotFound(err) {
			return err
		}

		err = r.createStatefulSet(ctx, cluster, role, ordinal)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *ReplicasetsManger) UpdateStatefulSets(ctx context.Context, cluster api.Cluster, role *v1beta1.Role) (complete bool, err error) {
	var (
		updated bool
		ordinal int64
		done    = true
		stsList *appsv1.StatefulSetList
	)

	stsList, err = r.ListStatefulSets(ctx, role.GetNamespace(), r.LabelsManager.SelectorByRoleName(role))
	if err != nil {
		return
	}

	for key := range stsList.Items {
		sts := &stsList.Items[key]

		ordinal, err = strconv.ParseInt(sts.GetLabels()[r.LabelsManager.ReplicasetOrdinal()], 10, 32)
		if err != nil {
			return false, err
		}

		updated, err = r.updateStatefulSet(ctx, cluster, role, sts, int32(ordinal))
		if err != nil {
			return false, err
		}

		if !updated {
			done = false
		}
	}

	return done, nil
}

func (r *ReplicasetsManger) createStatefulSet(ctx context.Context, cluster api.Cluster, role *v1beta1.Role, ordinal int32) error {
	stsName, err := role.GetReplicasetName(ordinal)
	if err != nil {
		return err
	}

	replicasetUUID := r.GetReplicasetUUID(role, ordinal)
	revisionHistoryLimit := int32(1)

	selectorLabels := map[string]string{
		r.LabelsManager.ClusterName():    cluster.GetName(),
		r.LabelsManager.RoleName():       role.GetName(),
		r.LabelsManager.ReplicasetName(): stsName,
		r.LabelsManager.ReplicasetUUID(): replicasetUUID,
	}

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      stsName,
			Namespace: role.GetNamespace(),
			Labels:    map[string]string{},
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: role.Spec.ReplicasetTemplate.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: selectorLabels,
			},
			VolumeClaimTemplates: role.Spec.ReplicasetTemplate.VolumeClaimTemplates,
			ServiceName:          cluster.GetName(),
			PodManagementPolicy:  appsv1.ParallelPodManagement,
			RevisionHistoryLimit: &revisionHistoryLimit,
			PersistentVolumeClaimRetentionPolicy: &appsv1.StatefulSetPersistentVolumeClaimRetentionPolicy{
				WhenDeleted: appsv1.RetainPersistentVolumeClaimRetentionPolicyType,
				WhenScaled:  appsv1.RetainPersistentVolumeClaimRetentionPolicyType,
			},
		},
	}

	_, err = r.syncStatefulSet(cluster, role, sts, ordinal, replicasetUUID)
	if err != nil {
		return err
	}

	return r.CreateObject(ctx, sts)
}

func (r *ReplicasetsManger) updateStatefulSet(ctx context.Context, cluster api.Cluster, role *v1beta1.Role, sts *appsv1.StatefulSet, ordinal int32) (bool, error) {
	replicasetUUID := r.GetReplicasetUUID(role, ordinal)

	changed, err := r.syncStatefulSet(cluster, role, sts, ordinal, replicasetUUID)
	if err != nil {
		return false, err
	}

	if changed {
		err = r.UpdateObject(ctx, sts)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (r *ReplicasetsManger) syncStatefulSet(
	cluster api.Cluster,
	role *v1beta1.Role,
	sts *appsv1.StatefulSet,
	ordinal int32,
	replicasetUUID string,
) (bool, error) {
	var (
		changed bool
		err     error
	)

	// Prepare revision hash
	rsPodTemplateHash, err := utils.HashObject(&role.Spec.ReplicasetTemplate.PodTemplate)
	if err != nil {
		return changed, err
	}

	// Prepare labels for all
	clusterLabels := map[string]string{
		r.LabelsManager.ClusterName(): cluster.GetName(),
	}
	roleLabels := role.GetLabels()
	stsLabels := utils.MergeMaps(clusterLabels, roleLabels, map[string]string{
		r.LabelsManager.ClusterName():               cluster.GetName(),
		r.LabelsManager.RoleName():                  role.GetName(),
		r.LabelsManager.ReplicasetName():            sts.GetName(),
		r.LabelsManager.ReplicasetUUID():            replicasetUUID,
		r.LabelsManager.ReplicasetOrdinal():         fmt.Sprintf("%d", ordinal),
		r.LabelsManager.ReplicasetPodTemplateHash(): rsPodTemplateHash,
	})

	podTemplateLabels := utils.MergeMaps(
		role.Spec.ReplicasetTemplate.PodTemplate.GetLabels(),
		stsLabels,
	)

	//// If StatefulSet still not controlled by role we need to take control
	changed, err = r.ControlObject(role, sts)
	if err != nil {
		return changed, err
	}

	if *role.Spec.ReplicasetTemplate.Replicas > *sts.Spec.Replicas {
		sts.Spec.Replicas = role.Spec.ReplicasetTemplate.Replicas

		changed = true
	}

	if role.Spec.ReplicasetTemplate.MinReadySeconds != sts.Spec.MinReadySeconds {
		sts.Spec.MinReadySeconds = role.Spec.ReplicasetTemplate.MinReadySeconds

		changed = true
	}

	if role.Spec.ReplicasetTemplate.UpdateStrategy.Type != sts.Spec.UpdateStrategy.Type {
		sts.Spec.UpdateStrategy.Type = role.Spec.ReplicasetTemplate.UpdateStrategy.Type

		changed = true
	}

	if sts.Spec.UpdateStrategy.Type == appsv1.OnDeleteStatefulSetStrategyType {
		if sts.Spec.UpdateStrategy.RollingUpdate != nil {
			sts.Spec.UpdateStrategy.RollingUpdate = nil

			changed = true
		}
	}

	if sts.Spec.UpdateStrategy.Type == appsv1.RollingUpdateStatefulSetStrategyType {
		if !cmp.Equal(role.Spec.ReplicasetTemplate.UpdateStrategy.RollingUpdate, sts.Spec.UpdateStrategy.RollingUpdate) {
			sts.Spec.UpdateStrategy.RollingUpdate = role.Spec.ReplicasetTemplate.UpdateStrategy.RollingUpdate

			changed = true
		}
	}

	// If current hash of pod template stored in StatefulSet labels does not match
	// calculated hash of role's pod template, when it means something changed.
	// So we need to update StatefulSet
	if rsPodTemplateHash != sts.Labels[r.LabelsManager.ReplicasetPodTemplateHash()] {
		role.Spec.ReplicasetTemplate.PodTemplate.DeepCopyInto(&sts.Spec.Template)

		changed = true
	}

	// Labels of StatefulSet should match
	// When actual StatefulSet labels does not match that rule we need to update StatefulSet.
	if !cmp.Equal(stsLabels, sts.GetLabels()) {
		sts.SetLabels(stsLabels)

		changed = true
	}

	// Labels of PodTemplate should be sum of required labels and labels defined in ReplicasetTemplate.PodTemplate.Labels
	// When actual StatefulSet labels does not match that rule we need to update StatefulSet.
	if !cmp.Equal(podTemplateLabels, sts.Spec.Template.GetLabels()) {
		sts.Spec.Template.SetLabels(podTemplateLabels)

		changed = true
	}

	// Annotations of StatefulSet is a copy role annotations
	if !cmp.Equal(role.GetAnnotations(), sts.GetAnnotations()) {
		sts.SetAnnotations(role.GetAnnotations())

		changed = true
	}

	return changed, nil
}
