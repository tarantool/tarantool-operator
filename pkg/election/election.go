package election

import (
	"context"

	"github.com/pkg/errors"
	"github.com/tarantool/tarantool-operator/pkg/api"
	"github.com/tarantool/tarantool-operator/pkg/events"
	"github.com/tarantool/tarantool-operator/pkg/k8s"
	"github.com/tarantool/tarantool-operator/pkg/topology"
	"github.com/tarantool/tarantool-operator/pkg/utils"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type LeaderElection struct {
	client.Client
	k8s.ResourcesManager
	*events.Recorder

	Topology topology.CartridgeTopology
}

var (
	ErrNoAvailableLeader      = errors.New("no available leader pod")
	ErrLeaderElectionConflict = errors.New("race condition during leader election")
	ErrLeaderNotReady         = errors.New("leader pod is not ready")
	ErrLeaderWasNotElected    = errors.New("leader was not elected")
)

// GetLeaderInstance return first pod of first replicaset of first role.
// Its implements the following logic
// 1. Retrieve previously elected leader
// 1.1 When leader where is no previously elected leader, elects a new one
// 1.2 When previously elected leader is not available
// 1.2.1 and cluster was already bootstrapped, try to find new one instance
// 1.2.2 and cluster was NOT bootstrapped - it fails, because we can't safely change leader in this case,
//
//	previously elected leader may have some important config and topology data, which can be absent on new leader,
//	some empty instance may be elected in that case it will produce a split brain.
func (r *LeaderElection) GetLeaderInstance(ctx context.Context, cluster api.Cluster) (*v1.Pod, error) {
	leader, err := r.loadLeaderInstance(ctx, cluster)
	if err != nil {
		if errors.Is(err, ErrLeaderWasNotElected) {
			return r.ElectLeaderInstance(ctx, cluster)
		}

		if errors.Is(err, ErrLeaderNotReady) {
			if cluster.IsBootstrapped() {
				return r.ElectLeaderInstance(ctx, cluster)
			}

			return nil, err
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

	canBeLeader, err := r.CanBeLeader(ctx, cluster, leader)
	if err != nil {
		return nil, err
	}

	if !canBeLeader {
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
		stsName     string
		podName     string
		canBeLeader bool
		pod         *v1.Pod
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

				canBeLeader, err = r.CanBeLeader(ctx, cluster, pod)
				if err != nil {
					return nil, err
				}

				if canBeLeader {
					return pod, nil
				}
			}
		}
	}

	return nil, ErrNoAvailableLeader
}

func (r *LeaderElection) CanBeLeader(ctx context.Context, cluster api.Cluster, pod *v1.Pod) (bool, error) {
	if utils.IsPodDeleting(pod) {
		return false, nil
	}

	if !utils.IsPodRunning(pod) {
		return false, nil
	}

	if !cluster.IsBootstrapped() {
		started, err := r.Topology.IsCartridgeStarted(ctx, pod)
		if err != nil {
			return false, err
		}

		return started, nil
	}

	configured, err := r.Topology.IsCartridgeConfigured(ctx, pod)
	if err != nil {
		return false, err
	}

	return configured, nil
}

func (r *LeaderElection) IsLeader(cluster api.Cluster, pod *v1.Pod) bool {
	return pod.GetNamespace() == cluster.GetNamespace() && cluster.GetLeader() == pod.GetName()
}
