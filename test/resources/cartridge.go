package resources

import (
	"github.com/tarantool/tarantool-operator/apis/v1beta1"
	"github.com/tarantool/tarantool-operator/pkg/k8s"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func NewFakeCartridge(labelsManager k8s.LabelsManager) *FakeCartridge {
	cluster := &v1beta1.Cluster{
		ObjectMeta: metav1.ObjectMeta{},
		Spec: v1beta1.ClusterSpec{
			Domain:     DefaultDomain,
			ListenPort: DefaultListenPort,
		},
		Status: v1beta1.ClusterStatus{
			Phase:        v1beta1.ClusterPending,
			Leader:       "",
			Bootstrapped: false,
		},
	}

	builder := &FakeCartridge{
		labelsManager: labelsManager,
		Cluster:       cluster,
		Roles:         map[string]*v1beta1.Role{},
		StatefulSets:  map[string]map[string]*appsv1.StatefulSet{},
		objects: []client.Object{
			cluster,
		},
	}

	return builder
}

type FakeCartridge struct {
	labelsManager k8s.LabelsManager

	Cluster      *v1beta1.Cluster
	Roles        map[string]*v1beta1.Role
	StatefulSets map[string]map[string]*appsv1.StatefulSet
	Pods         []*v1.Pod

	objects []client.Object
}

func (r *FakeCartridge) object(obj client.Object) {
	r.objects = append(r.objects, obj)
}

func (r *FakeCartridge) WithNamespace(namespace string) *FakeCartridge {
	r.Cluster.Namespace = namespace

	return r
}

func (r *FakeCartridge) BuildFakeClient() client.WithWatch {
	fakeClientBuilder := fake.NewClientBuilder()

	return fakeClientBuilder.
		WithObjects(r.objects...).
		Build()
}
