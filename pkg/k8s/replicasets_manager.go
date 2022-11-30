package k8s

import (
	"context"

	"github.com/tarantool/tarantool-operator/pkg/api"
	v1 "k8s.io/api/core/v1"
)

type ReplicasetsManger[RoleType api.Role] interface {
	// GetReplicasetUUID must return stable UUIDv5 as string for each replicaset
	GetReplicasetUUID(role RoleType, ordinal int32) string

	// GetAdvertiseURI must return stable URI as string for each replicaset
	GetAdvertiseURI(cluster api.Cluster, pod *v1.Pod) string

	CreateStatefulSets(ctx context.Context, cluster api.Cluster, role RoleType) error
	UpdateStatefulSets(ctx context.Context, cluster api.Cluster, role RoleType) (complete bool, err error)
}
