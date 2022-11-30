package election

import (
	"context"

	"github.com/pkg/errors"
	"github.com/tarantool/tarantool-operator/pkg/api"
	"github.com/tarantool/tarantool-operator/pkg/events"
	"github.com/tarantool/tarantool-operator/pkg/k8s"
	"github.com/tarantool/tarantool-operator/pkg/utils"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type LeaderElection struct {
	client.Client
	k8s.ResourcesManager
	*events.Recorder
}

var (
	ErrNoAvailableLeader      = errors.New("no available leader pod")
	ErrLeaderElectionConflict = errors.New("race condition during leader election")
	ErrLeaderNotReady         = errors.New("leader pod is not ready")
	ErrLeaderWasNotElected    = errors.New("leader was not elected")
)

// GetLeaderInstance return first pod of first replicaset of first role.
func (r *LeaderElection) GetLeaderInstance(ctx context.Context, cluster api.Cluster) (*v1.Pod, error) {
	leader, err := r.loadLeaderInstance(ctx, cluster)
	if err != nil {
		if errors.Is(err, ErrLeaderWasNotElected) || errors.Is(err, ErrLeaderNotReady) {
			return r.ElectLeaderInstance(ctx, cluster)
		}
	}

	return leader, nil
}

func (r *LeaderElection) ElectLeaderInstance(ctx context.Context, cluster api.Cluster) (*v1.Pod, error) {
	leader, err := r.findNewLeaderInstance(ctx, cluster)
	if err != nil {
		return nil, err
	}

	cluster.SetLeader(leader.GetName())

	err = r.Status().Update(ctx, cluster)
	if err != nil {
		if apierrors.IsConflict(err) {
			return nil, ErrLeaderElectionConflict
		}

		return nil, err
	}

	r.Event(cluster, NewTopologyLeaderElectedEvent(leader))

	return leader, nil
}

func (r *LeaderElection) loadLeaderInstance(ctx context.Context, cluster api.Cluster) (*v1.Pod, error) {
	if cluster.GetLeader() == "" {
		return nil, ErrLeaderWasNotElected
	}

	leader, err := r.GetPod(ctx, cluster.GetNamespace(), cluster.GetLeader())
	if err != nil {
		return nil, err
	}

	if !r.CanBeLeader(cluster, leader) {
		return nil, ErrLeaderNotReady
	}

	return leader, nil
}

// findLeaderInstance
// Leader is a running pod
// Leader should not be the same as previous leader
// cartridge container should be ready
// sidecar container should be ready if vshard was bootstrapped
// Order by:
// - RoleCE ordinal ascending
// - StatefulSet ordinal ascending
// - Pod ordinal ascending.
func (r *LeaderElection) findNewLeaderInstance(ctx context.Context, cluster api.Cluster) (*v1.Pod, error) {
	var (
		err            error
		roles          []api.Role
		maxReplicasets = int32(0)
		maxReplicas    = int32(0)
	)

	roles, err = r.GetClusterRoles(ctx, cluster)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if maxReplicasets < role.GetReplicasets() {
			maxReplicasets = role.GetReplicasets()
		}

		if maxReplicas < role.GetReplicas() {
			maxReplicas = role.GetReplicas()
		}
	}

	var (
		stsName string
		podName string
		pod     *v1.Pod
	)

	for podOrdinal := int32(0); podOrdinal < maxReplicas; podOrdinal++ {
		for replicasetOrdinal := int32(0); replicasetOrdinal < maxReplicasets; replicasetOrdinal++ {
			for _, role := range roles {
				if role.GetReplicasets() <= replicasetOrdinal {
					continue
				}

				if role.GetReplicas() <= podOrdinal {
					continue
				}

				stsName, err = role.GetReplicasetName(replicasetOrdinal)
				if err != nil {
					return nil, err
				}

				podName = utils.GetStatefulSetPodName(stsName, podOrdinal)
				if podName == cluster.GetLeader() {
					continue
				}

				pod, err = r.GetPod(ctx, cluster.GetNamespace(), podName)
				if err != nil {
					if apierrors.IsNotFound(err) {
						continue
					}

					return nil, err
				}

				if r.CanBeLeader(cluster, pod) {
					return pod, nil
				}
			}
		}
	}

	return nil, ErrNoAvailableLeader
}

func (r *LeaderElection) CanBeLeader(cluster api.Cluster, pod *v1.Pod) bool {
	if pod.GetDeletionTimestamp() != nil {
		return false
	}

	if !utils.IsPodRunning(pod) {
		return false
	}

	if !cluster.IsBootstrapped() {
		return utils.IsPodDefaultContainerReady(pod)
	}

	return utils.IsPodReady(pod)
}

func (r *LeaderElection) IsLeader(cluster api.Cluster, pod *v1.Pod) bool {
	return pod.GetNamespace() == cluster.GetNamespace() && cluster.GetLeader() == pod.GetName()
}
