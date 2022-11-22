package topology

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

// CartridgeTopology .
type CartridgeTopology interface {
	Exec(ctx context.Context, instance *v1.Pod, res interface{}, lua string, args ...interface{}) error

	Join(
		ctx context.Context,
		leader *v1.Pod,
		replicasetAlias string,
		replicasetUUID string,
		replicasetRoles []string,
		replicasetWeight int32,
		replicasetVshardGroup string,
		replicasetIsAllRw bool,
		advertiseURI string,
	) error

	GetInstanceUUID(ctx context.Context, pod *v1.Pod) (string, error)

	SetWeight(ctx context.Context, leader *v1.Pod, replicasetUUID string, replicaWeight int32) error
	GetReplicasetRoles(ctx context.Context, leader *v1.Pod, replicasetUUID string) ([]string, error)
	SetReplicasetRoles(ctx context.Context, leader *v1.Pod, replicasetUUID string, roles []string) error

	BootstrapVshard(ctx context.Context, leader *v1.Pod) error

	GetRolesHierarchy(ctx context.Context, leader *v1.Pod) (map[string][]string, error)

	SetFailoverParams(ctx context.Context, leader *v1.Pod, params *FailoverParams) error
	GetFailoverParams(ctx context.Context, leader *v1.Pod) (*FailoverParams, error)

	GetCartridgeConfig(ctx context.Context, leader *v1.Pod) (CartridgeConfigData, error)
	ApplyCartridgeConfig(ctx context.Context, leader *v1.Pod, config CartridgeConfigData) error
}
