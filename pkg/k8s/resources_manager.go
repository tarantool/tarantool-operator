package k8s

import (
	"context"

	"github.com/tarantool/tarantool-operator/pkg/api"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type ResourcesManager interface {
	CreateObject(ctx context.Context, obj client.Object) error
	UpdateObject(ctx context.Context, obj client.Object) error
	ControlObject(owner client.Object, target client.Object) (bool, error)
	GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error)
	GetService(ctx context.Context, namespace, name string) (*corev1.Service, error)
	ListStatefulSets(ctx context.Context, namespace string, selector labels.Selector) (*appsv1.StatefulSetList, error)
	ListPods(ctx context.Context, namespace string, selector labels.Selector) (*corev1.PodList, error)
	GetConfigMap(ctx context.Context, namespace, name string) (*corev1.ConfigMap, error)
	GetConfigMapValue(ctx context.Context, ns, name, key string) (string, error)
	GetSecret(ctx context.Context, namespace, name string) (*corev1.Secret, error)
	GetSecretValue(ctx context.Context, ns, name, key string) (string, error)

	GetCluster(ctx context.Context, ns, name string) (api.Cluster, error)
	GetClusterRoles(ctx context.Context, cluster api.Cluster) ([]api.Role, error)
}

type CommonResourcesManager struct {
	client.Client

	Scheme *runtime.Scheme
}

func (r *CommonResourcesManager) CreateObject(ctx context.Context, obj client.Object) error {
	err := r.Create(ctx, obj)
	if err != nil {
		return err
	}

	return nil
}

func (r *CommonResourcesManager) UpdateObject(ctx context.Context, obj client.Object) error {
	err := r.Update(ctx, obj)
	if err != nil {
		return err
	}

	return nil
}

func (r *CommonResourcesManager) ControlObject(owner client.Object, target client.Object) (bool, error) {
	if metav1.IsControlledBy(target, owner) {
		return false, nil
	}

	err := controllerutil.SetControllerReference(owner, target, r.Scheme)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *CommonResourcesManager) GetPod(ctx context.Context, namespace, name string) (*corev1.Pod, error) {
	pod := &corev1.Pod{}

	err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, pod)
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func (r *CommonResourcesManager) GetService(ctx context.Context, namespace, name string) (*corev1.Service, error) {
	svc := &corev1.Service{}

	err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, svc)
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func (r *CommonResourcesManager) ListStatefulSets(ctx context.Context, namespace string, selector labels.Selector) (*appsv1.StatefulSetList, error) {
	stsList := &appsv1.StatefulSetList{}

	err := r.List(ctx, stsList, &client.ListOptions{LabelSelector: selector, Namespace: namespace})
	if err != nil {
		return nil, err
	}

	return stsList, nil
}

func (r *CommonResourcesManager) ListPods(ctx context.Context, namespace string, selector labels.Selector) (*corev1.PodList, error) {
	podList := &corev1.PodList{}

	err := r.List(ctx, podList, &client.ListOptions{LabelSelector: selector, Namespace: namespace})
	if err != nil {
		return nil, err
	}

	return podList, nil
}

func (r *CommonResourcesManager) GetConfigMap(ctx context.Context, namespace, name string) (*corev1.ConfigMap, error) {
	cfgMap := &corev1.ConfigMap{}

	err := r.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, cfgMap)
	if err != nil {
		return nil, err
	}

	return cfgMap, nil
}

func (r *CommonResourcesManager) GetConfigMapValue(ctx context.Context, ns, name, key string) (string, error) {
	cfg, err := r.GetConfigMap(ctx, ns, name)
	if err != nil {
		return "", err
	}

	return cfg.Data[key], nil
}

func (r *CommonResourcesManager) GetSecret(ctx context.Context, ns, name string) (*corev1.Secret, error) {
	secret := &corev1.Secret{}

	err := r.Get(ctx, types.NamespacedName{Namespace: ns, Name: name}, secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (r *CommonResourcesManager) GetSecretValue(ctx context.Context, ns, name, key string) (string, error) {
	secret, err := r.GetSecret(ctx, ns, name)
	if err != nil {
		return "", err
	}

	return string(secret.Data[key]), nil
}
