package controllers

import (
	"context"

	"github.com/go-logr/logr"
	. "github.com/tarantool/tarantool-operator/apis/v1alpha2"
	. "github.com/tarantool/tarantool-operator/internal"
	. "github.com/tarantool/tarantool-operator/internal/implementation"
	. "github.com/tarantool/tarantool-operator/internal/steps"
	"github.com/tarantool/tarantool-operator/pkg/election"
	"github.com/tarantool/tarantool-operator/pkg/events"
	"github.com/tarantool/tarantool-operator/pkg/k8s"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation/steps/common"
	"github.com/tarantool/tarantool-operator/pkg/topology"
	"github.com/tarantool/tarantool-operator/pkg/topology/transport/podexec"
	"github.com/tarantool/tarantool-operator/pkg/topology/transport/podexec/cli"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

//+kubebuilder:rbac:groups=tarantool.io,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tarantool.io,resources=clusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;create;update;watch;list;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;create;update;watch;list;patch;delete
//+kubebuilder:rbac:groups=core,resources=events,verbs=get;create;update;watch;list;patch;delete
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;create;update;watch;list;patch;delete
//+kubebuilder:rbac:groups="",resources=pods/exec,verbs=get;create;update;watch;list;patch;delete
//+kubebuilder:rbac:groups="",resources=pods/logs,verbs=get;watch;list;
//+kubebuilder:rbac:groups="",resources=services,verbs=get;create;update;watch;list;patch;delete
//+kubebuilder:rbac:groups="",resources=endpoints,verbs=get;create;update;watch;list;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumes,verbs=get;create;update;watch;list;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;create;update;watch;list;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;watch;list

func NewClusterReconciler(mgr Manager) *ClusterReconciler {
	k8sConfig := mgr.GetConfig()
	k8sClient := mgr.GetClient()
	k8sScheme := mgr.GetScheme()
	restClient, _ := apiutil.RESTClientForGVK(
		schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    "Pod",
		},
		false,
		k8sConfig,
		serializer.NewCodecFactory(k8sScheme),
	)

	labelsManager := &k8s.NamespacedLabelsManager{
		Namespace: "tarantool.io",
	}
	resourcesManager := &ResourcesManager{
		LabelsManager: labelsManager,
		CommonResourcesManager: &k8s.CommonResourcesManager{
			Client: k8sClient,
			Scheme: k8sScheme,
		},
	}
	eventsRecorder := events.NewRecorder(mgr.GetEventRecorderFor("cluster-controller"))
	luaTopology := &topology.CommonCartridgeTopology{
		Transport: &podexec.PodExec{
			RestClient:    restClient,
			RestConfig:    k8sConfig,
			RuntimeScheme: k8sScheme,
			CLI:           &cli.TarantoolCTL{},
		},
	}

	return &ClusterReconciler{
		LabelsManager: labelsManager,
		SteppedReconciler: &SteppedReconciler[*ClusterContextCE, *ClusterControllerCE]{
			Client: k8sClient,
			Controller: &ClusterControllerCE{
				CommonClusterController: &CommonClusterController{
					CommonController: &CommonController{
						Client: k8sClient,
						Schema: k8sScheme,
						LeaderElection: &election.LeaderElection{
							Client:           k8sClient,
							Recorder:         eventsRecorder,
							ResourcesManager: resourcesManager,
						},
						ResourcesManager: resourcesManager,
						EventsRecorder:   eventsRecorder,
						Topology:         luaTopology,
						LabelsManager:    labelsManager,
					},
				},
			},
		},
	}
}

// ClusterReconciler reconciles a ClusterCE object.
type ClusterReconciler struct {
	*SteppedReconciler[*ClusterContextCE, *ClusterControllerCE]

	LabelsManager k8s.LabelsManager
}

func (r *ClusterReconciler) Reconcile(ctx context.Context, req Request) (Result, error) {
	return r.Run(
		&ClusterContextCE{
			CommonContext: &CommonContext{
				Context: ctx,
				Request: req,
				Logger:  logr.FromContextOrDiscard(ctx),
			},
		},
		Info[*ClusterContextCE, *ClusterControllerCE]("Reconcile cluster"),
		GetRequestedObject[*ClusterContextCE, *ClusterControllerCE](&Cluster{}),
		SetClusterPhase(ClusterPending),
		CheckClusterDeletion(),
		ResetClusterStatus(),
		SetClusterPhase(ClusterSyncingService),
		SyncClusterWideService(),
		SetClusterPhase(ClusterWaitingForRoles),
		WaitForRolesPhases(RoleWaitingForBootstrap, RoleReady),
		Info[*ClusterContextCE, *ClusterControllerCE]("All roles ready, we are going to bootstrap cluster"),
		SetClusterPhase(ClusterWaitingForLeader),
		GetLeader[*ClusterContextCE, *ClusterControllerCE](),
		Bootstrap(BootstrapParams{
			OnError: ClusterUnableToBootstrap,
		}),
		SetClusterPhase(ClusterFailoverConfiguring),
		ConfigureFailover(),
		SetClusterPhase(ClusterReady),
		Info[*ClusterContextCE, *ClusterControllerCE]("Community cluster ready"),
	)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterReconciler) SetupWithManager(mgr Manager) error {
	return NewControllerManagedBy(mgr).
		For(&Cluster{}).
		Owns(&v1.Service{}).
		Watches(&source.Kind{Type: &Role{}}, handler.EnqueueRequestsFromMapFunc(func(obj client.Object) []reconcile.Request {
			role, ok := obj.(*Role)
			if !ok {
				return []Request{}
			}

			if labels := role.GetLabels(); role.Status.Phase == RoleWaitingForBootstrap && labels != nil {
				return []Request{
					{
						NamespacedName: types.NamespacedName{
							Namespace: role.GetNamespace(),
							Name:      labels[r.LabelsManager.ClusterName()],
						},
					},
				}
			}

			return []Request{}
		})).
		Complete(r)
}
