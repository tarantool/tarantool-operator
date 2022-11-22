package steps

import (
	. "github.com/tarantool/tarantool-operator/apis/v1alpha2"
	. "github.com/tarantool/tarantool-operator/internal/context"
	. "github.com/tarantool/tarantool-operator/internal/controller"
	"github.com/tarantool/tarantool-operator/pkg/reconciliation/steps/cluster"
)

func CheckClusterDeletion() *cluster.CheckDeletionStep[*Cluster, *ClusterContext, *ClusterController] {
	return &cluster.CheckDeletionStep[*Cluster, *ClusterContext, *ClusterController]{}
}

func ResetClusterStatus() *cluster.ResetStatusStep[*Cluster, *ClusterContext, *ClusterController] {
	return &cluster.ResetStatusStep[*Cluster, *ClusterContext, *ClusterController]{}
}

func SetClusterPhase(phase ClusterPhase) *cluster.SetPhaseStep[ClusterPhase, *Cluster, *ClusterContext, *ClusterController] {
	return &cluster.SetPhaseStep[ClusterPhase, *Cluster, *ClusterContext, *ClusterController]{
		Phase: phase,
	}
}

func SyncClusterWideService() *cluster.SyncClusterWideServiceStep[*Cluster, *ClusterContext, *ClusterController] {
	return &cluster.SyncClusterWideServiceStep[*Cluster, *ClusterContext, *ClusterController]{}
}

func WaitForRolesPhases(expectedPhases ...RolePhase) *cluster.WaitForRolesPhaseStep[RolePhase, *Cluster, *ClusterContext, *ClusterController] {
	return &cluster.WaitForRolesPhaseStep[RolePhase, *Cluster, *ClusterContext, *ClusterController]{
		ExpectedPhases: expectedPhases,
	}
}

type BootstrapParams struct {
	OnError ClusterPhase
}

func Bootstrap(params BootstrapParams) *cluster.BootstrapStep[ClusterPhase, *Cluster, *ClusterContext, *ClusterController] {
	return &cluster.BootstrapStep[ClusterPhase, *Cluster, *ClusterContext, *ClusterController]{
		ErrorPhase: params.OnError,
	}
}

func ConfigureFailover() *cluster.ConfigureFailoverStep[*Cluster, *ClusterContext, *ClusterController] {
	return &cluster.ConfigureFailoverStep[*Cluster, *ClusterContext, *ClusterController]{}
}
