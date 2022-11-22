package reconciliation

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/tarantool/tarantool-operator/pkg/api"
	v1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type CommonContext struct {
	context.Context

	Request ctrl.Request
	Logger  logr.Logger

	cluster api.Cluster
	leader  *v1.Pod
}

func (r *CommonContext) GetRequest() ctrl.Request {
	return r.Request
}

func (r *CommonContext) GetLogger() logr.Logger {
	return r.Logger
}

func (r *CommonContext) SetRelatedCluster(cluster api.Cluster) {
	r.cluster = cluster
}

func (r *CommonContext) GetRelatedCluster() api.Cluster {
	return r.cluster
}

func (r *CommonContext) SetLeader(leader *v1.Pod) {
	r.leader = leader
}

func (r *CommonContext) GetLeader() *v1.Pod {
	return r.leader
}
