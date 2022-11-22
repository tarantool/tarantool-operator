package context

import (
	"github.com/pkg/errors"
	"github.com/tarantool/tarantool-operator/apis/v1alpha2"
	"github.com/tarantool/tarantool-operator/pkg/reconciliation"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ClusterContext struct {
	*reconciliation.CommonContext

	Cluster *v1alpha2.Cluster
}

func (r *ClusterContext) SetCluster(cluster *v1alpha2.Cluster) {
	r.Cluster = cluster
	r.CommonContext.SetRelatedCluster(cluster)
}

func (r *ClusterContext) GetCluster() *v1alpha2.Cluster {
	return r.Cluster
}

func (r *ClusterContext) HasRequestedObject() bool {
	return r.GetCluster() != nil
}

func (r *ClusterContext) SetRequestedObject(obj client.Object) error {
	cluster, ok := obj.(*v1alpha2.Cluster)
	if !ok {
		return errors.New("ClusterContext used with wrong k8s object")
	}

	r.SetCluster(cluster)

	return nil
}

func (r *ClusterContext) GetRequestedObject() client.Object {
	return r.GetCluster()
}
