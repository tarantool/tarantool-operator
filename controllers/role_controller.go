package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	. "github.com/tarantool/tarantool-operator/apis/v1beta1"
	. "github.com/tarantool/tarantool-operator/internal"
	"github.com/tarantool/tarantool-operator/internal/implementation"
	. "github.com/tarantool/tarantool-operator/internal/steps"
	"github.com/tarantool/tarantool-operator/pkg/election"
	"github.com/tarantool/tarantool-operator/pkg/events"
	"github.com/tarantool/tarantool-operator/pkg/k8s"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation/steps/common"
	"github.com/tarantool/tarantool-operator/pkg/topology"
	"github.com/tarantool/tarantool-operator/pkg/topology/transport/podexec"
	"github.com/tarantool/tarantool-operator/pkg/topology/transport/podexec/cli"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

//+kubebuilder:rbac:groups=tarantool.io,resources=roles,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tarantool.io,resources=roles/status,verbs=get;update;patch

func NewRoleReconciler(mgr Manager) *RoleReconciler {
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
	resourcesManager := &implementation.ResourcesManager{
		LabelsManager: labelsManager,
		CommonResourcesManager: &k8s.CommonResourcesManager{
			Client: k8sClient,
			Scheme: k8sScheme,
		},
	}
	eventsRecorder := events.NewRecorder(mgr.GetEventRecorderFor("role-controller"))
	luaTopology := &topology.CommonCartridgeTopology{
		Transport: &podexec.PodExec{
			RestClient:    restClient,
			RestConfig:    k8sConfig,
			RuntimeScheme: k8sScheme,
			CLI:           &cli.TarantoolCTL{},
		},
	}

	return &RoleReconciler{
		SteppedReconciler: &SteppedReconciler[*RoleContextCE, *RoleControllerCE]{
			Client: k8sClient,
			Controller: &RoleControllerCE{
				ReplicasetsManger: &implementation.ReplicasetsManger{
					ResourcesManager: resourcesManager,
					// UUIDSpace is a randomly generate uuid v4 used as base for generation v5 uuid of replicasets.
					UUIDSpace: uuid.MustParse("8601615b-c39f-4fd9-a88b-9df688656cbd"),
				},
				CommonRoleController: &CommonRoleController{
					CommonController: &CommonController{
						Client: k8sClient,
						Schema: k8sScheme,
						LeaderElection: &election.LeaderElection{
							Client:           k8sClient,
							Recorder:         eventsRecorder,
							ResourcesManager: resourcesManager,
							Topology:         luaTopology,
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

// RoleReconciler reconciles a Role object.
type RoleReconciler struct {
	*SteppedReconciler[*RoleContextCE, *RoleControllerCE]
}

func (r *RoleReconciler) Reconcile(ctx context.Context, req Request) (Result, error) {
	return r.Run(
		&RoleContextCE{
			CommonContext: &CommonContext{
				Context: ctx,
				Request: req,
				Logger:  logr.FromContextOrDiscard(ctx),
			},
		},
		Info[*RoleContextCE, *RoleControllerCE]("Reconcile role"),
		GetRequestedObject[*RoleContextCE, *RoleControllerCE](&Role{}),
		ResetRoleStatus(),

		SetRolePhase(RoleWaitingForCluster),
		GetClusterByLabels[*RoleContextCE, *RoleControllerCE](),

		SetRolePhase(RolePending),
		CreateStatefulSets(),
		UpdateStatefulSets(),

		SetRolePhase(RoleWaitingForLeader),
		GetLeader[*RoleContextCE, *RoleControllerCE](),

		SetRolePhase(RoleWaitForCartridgeReady),
		EnsureCartridgeReady(),

		SetRolePhase(RoleJoining),
		JoinInstances(JoinInstancesParams{
			ConfigErrorPhase: RoleConfigError,
		}),

		SetRolePhase(RoleConfiguring),
		ConfigureVShardRoles(),

		SetRolePhase(RoleWaitingForBootstrap),
		WaitForClusterBootstrapped[*RoleContextCE, *RoleControllerCE](),

		SetRolePhase(RoleConfiguringWeights),
		SetVShardWeights(),

		SetRolePhase(RoleReady),
		Info[*RoleContextCE, *RoleControllerCE]("Role ready"),
	)
}

// SetupWithManager sets up the controller with the Manager.
func (r *RoleReconciler) SetupWithManager(mgr Manager) error {
	return NewControllerManagedBy(mgr).
		For(&Role{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&v1.Pod{}).
		Complete(r)
}
