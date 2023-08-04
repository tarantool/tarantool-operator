package role

import (
	"time"

	"github.com/pkg/errors"
	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
	"github.com/tarantool/tarantool-operator/pkg/topology"
	"github.com/tarantool/tarantool-operator/pkg/utils"
	v1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/utils/strings/slices"
)

type JoinInstancesStep[
	PhaseType comparable,
	RoleType api.RoleWithStatus[PhaseType],
	CtxType RoleContext[RoleType],
	CtrlType RoleController[RoleType],
] struct {
	ConfigErrorPhase PhaseType
}

func (r *JoinInstancesStep[PhaseType, RoleType, CtxType, CtrlType]) GetName() string {
	return "Join instances"
}

func (r *JoinInstancesStep[PhaseType, RoleType, CtxType, CtrlType]) Reconcile(ctx CtxType, ctrl CtrlType) (*Result, error) {
	role := ctx.GetRole()
	cluster := ctx.GetRelatedCluster()

	topologyClient := ctrl.GetTopology()
	allJoined := true

	var pod *v1.Pod

	for stsOrdinal := int32(0); stsOrdinal < role.GetReplicasets(); stsOrdinal++ {
		selector := ctrl.GetLabelsManager().SelectorByReplicasetOrdinal(role, stsOrdinal)

		stsList, err := ctrl.GetResourcesManager().ListStatefulSets(ctx, role.GetNamespace(), selector)
		if err != nil {
			return Error(err)
		}

		var advertiseURIs []string

		for key := range stsList.Items {
			sts := &stsList.Items[key]
			if sts.GetDeletionTimestamp() != nil {
				continue
			}

			for podOrdinal := int32(0); podOrdinal < role.GetReplicas(); podOrdinal++ {
				podName := utils.GetStatefulSetPodName(sts.GetName(), podOrdinal)

				pod, err = ctrl.GetResourcesManager().GetPod(ctx, role.GetNamespace(), podName)
				if err != nil {
					if !apiErrors.IsNotFound(err) {
						return Error(err)
					}

					allJoined = false
					advertiseURIs = nil

					break
				}

				if !utils.IsPodRunning(pod) || utils.IsPodDeleting(pod) {
					allJoined = false
					advertiseURIs = nil

					break
				}

				instanceUUID, uuidErr := topologyClient.GetInstanceUUID(ctx, pod)
				if uuidErr != nil {
					ctx.GetLogger().Error(uuidErr, "unable to get instance uuid")

					return Error(uuidErr)
				}

				if instanceUUID == "" {
					advertiseURIs = append(advertiseURIs, ctrl.GetReplicasetsManger().GetAdvertiseURI(cluster, pod))
				}
			}
		}

		if len(advertiseURIs) == 0 {
			continue
		}

		leader := ctx.GetLeader()

		if !slices.Contains(advertiseURIs, ctrl.GetReplicasetsManger().GetAdvertiseURI(cluster, leader)) {
			switch leaderUUID, uuidErr := topologyClient.GetInstanceUUID(ctx, leader); {
			case uuidErr != nil:
				ctx.GetLogger().Error(uuidErr, "unable to get leader instance uuid")

				return Error(uuidErr)
			case leaderUUID == "":
				allJoined = false

				continue
			}
		}

		vshard := role.GetVShardConfig()

		alias, err := role.GetReplicasetName(stsOrdinal)
		if err != nil {
			return Error(err)
		}

		joinErr := topologyClient.Join(
			ctx,
			leader,
			alias,
			ctrl.GetReplicasetsManger().GetReplicasetUUID(role, stsOrdinal),
			vshard.GetRoles(),
			vshard.GetWeight(),
			vshard.GetGroupName(),
			role.IsAllRw(),
			advertiseURIs...,
		)
		if joinErr != nil {
			ctx.GetLogger().Error(joinErr, "unable to join instance")

			var configErr *topology.UnknownRoleError
			if errors.As(joinErr, &configErr) {
				ctrl.GetEventsRecorder().Event(role, NewWrongVShardRolesEvent(configErr))
				role.SetPhase(r.ConfigErrorPhase)

				return Complete()
			}

			return Error(joinErr)
		}
	}

	if !allJoined {
		ctx.GetLogger().Info("Not all pods joined tarantool cluster")

		return Requeue(10 * time.Second)
	}

	return NextStep()
}
