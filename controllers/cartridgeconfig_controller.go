package controllers

import (
	"context"
	"net/http"

	"github.com/go-logr/logr"
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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

//+kubebuilder:rbac:groups=tarantool.io,resources=clusters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tarantool.io,resources=clusters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tarantool.io,resources=cartridgeconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tarantool.io,resources=cartridgeconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tarantool.io,resources=cartridgeconfigs/finalizers,verbs=update

func NewCartridgeConfigReconciler(mgr Manager) *CartridgeConfigReconciler {
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
		&http.Client{},
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
	eventsRecorder := events.NewRecorder(mgr.GetEventRecorderFor("cartridge-config-controller"))
	luaTopology := &topology.CommonCartridgeTopology{
		Transport: &podexec.PodExec{
			RestClient:    restClient,
			RestConfig:    k8sConfig,
			RuntimeScheme: k8sScheme,
			CLI:           &cli.TarantoolCTL{},
		},
	}

	return &CartridgeConfigReconciler{
		SteppedReconciler: &SteppedReconciler[*CartridgeConfigContextCE, *CartridgeConfigControllerCE]{
			Client: k8sClient,
			Controller: &CartridgeConfigControllerCE{
				CommonCartridgeConfigController: &CommonCartridgeConfigController{
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

// CartridgeConfigReconciler reconciles a Config object.
type CartridgeConfigReconciler struct {
	*SteppedReconciler[*CartridgeConfigContextCE, *CartridgeConfigControllerCE]
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *CartridgeConfigReconciler) Reconcile(ctx context.Context, req Request) (Result, error) {
	return r.Run(
		&CartridgeConfigContextCE{
			CommonContext: &CommonContext{
				Context: ctx,
				Request: req,
				Logger:  logr.FromContextOrDiscard(ctx),
			},
		},
		Info[*CartridgeConfigContextCE, *CartridgeConfigControllerCE]("Reconcile Config"),
		GetRequestedObject[*CartridgeConfigContextCE, *CartridgeConfigControllerCE](&CartridgeConfig{}),
		ResetCartridgeConfigStatus(),
		SetCartridgeConfigPhase(CartridgeConfigWaitingForCluster),
		GetClusterByLabels[*CartridgeConfigContextCE, *CartridgeConfigControllerCE](),
		WaitForClusterBootstrapped[*CartridgeConfigContextCE, *CartridgeConfigControllerCE](),
		SetCartridgeConfigPhase(CartridgeConfigWaitingForLeader),
		GetLeader[*CartridgeConfigContextCE, *CartridgeConfigControllerCE](),
		SetCartridgeConfigPhase(CartridgeConfigApplying),
		ConfigureCartridge(),
		SetCartridgeConfigPhase(CartridgeConfigReady),
		Info[*CartridgeConfigContextCE, *CartridgeConfigControllerCE]("CartridgeConfig ready"),
	)
}

// SetupWithManager sets up the controller with the Manager.
func (r *CartridgeConfigReconciler) SetupWithManager(mgr Manager) error {
	return NewControllerManagedBy(mgr).
		For(&CartridgeConfig{}).
		Complete(r)
}
