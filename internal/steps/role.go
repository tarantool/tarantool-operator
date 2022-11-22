package steps

import (
	. "github.com/tarantool/tarantool-operator/apis/v1alpha2"
	. "github.com/tarantool/tarantool-operator/internal"
	"github.com/tarantool/tarantool-operator/pkg/reconciliation/steps/role"
)

func ResetRoleStatus() *role.ResetStatusStep[*Role, *RoleContextCE, *RoleControllerCE] {
	return &role.ResetStatusStep[*Role, *RoleContextCE, *RoleControllerCE]{}
}

func SetRolePhase(phase RolePhase) *role.SetPhaseStep[RolePhase, *Role, *RoleContextCE, *RoleControllerCE] {
	return &role.SetPhaseStep[RolePhase, *Role, *RoleContextCE, *RoleControllerCE]{
		Phase: phase,
	}
}

func CreateStatefulSets() *role.CreateStatefulSetsStep[*Role, *RoleContextCE, *RoleControllerCE] {
	return &role.CreateStatefulSetsStep[*Role, *RoleContextCE, *RoleControllerCE]{}
}

func UpdateStatefulSets() *role.UpdateStatefulSetsStep[*Role, *RoleContextCE, *RoleControllerCE] {
	return &role.UpdateStatefulSetsStep[*Role, *RoleContextCE, *RoleControllerCE]{}
}

func EnsureCartridgeReady() *role.EnsureCartridgeReadyStep[*Role, *RoleContextCE, *RoleControllerCE] {
	return &role.EnsureCartridgeReadyStep[*Role, *RoleContextCE, *RoleControllerCE]{}
}

type JoinInstancesParams struct {
	ConfigErrorPhase RolePhase
}

func JoinInstances(params JoinInstancesParams) *role.JoinInstancesStep[RolePhase, *Role, *RoleContextCE, *RoleControllerCE] {
	return &role.JoinInstancesStep[
		RolePhase,
		*Role,
		*RoleContextCE,
		*RoleControllerCE,
	]{
		ConfigErrorPhase: params.ConfigErrorPhase,
	}
}

func ConfigureVShardRoles() *role.ConfigureVShardRolesStep[*Role, *RoleContextCE, *RoleControllerCE] {
	return &role.ConfigureVShardRolesStep[*Role, *RoleContextCE, *RoleControllerCE]{}
}

func SetVShardWeights() *role.SetVShardWeightsStep[*Role, *RoleContextCE, *RoleControllerCE] {
	return &role.SetVShardWeightsStep[*Role, *RoleContextCE, *RoleControllerCE]{}
}
