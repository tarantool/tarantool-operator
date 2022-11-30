package reconciliation

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/tarantool/tarantool-operator/pkg/api"
	v1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Context interface {
	context.Context

	GetLogger() logr.Logger

	GetRequest() ctrl.Request

	SetRelatedCluster(cluster api.Cluster)
	GetRelatedCluster() api.Cluster

	SetLeader(leader *v1.Pod)
	GetLeader() *v1.Pod

	HasRequestedObject() bool
	SetRequestedObject(obj client.Object) error
	GetRequestedObject() client.Object
}
