package controllers_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/tarantool/tarantool-operator/apis/v1alpha2"
	. "github.com/tarantool/tarantool-operator/controllers"
	. "github.com/tarantool/tarantool-operator/internal"
	. "github.com/tarantool/tarantool-operator/internal/implementation"
	"github.com/tarantool/tarantool-operator/pkg/election"
	"github.com/tarantool/tarantool-operator/pkg/events"
	"github.com/tarantool/tarantool-operator/pkg/k8s"
	"github.com/tarantool/tarantool-operator/pkg/reconciliation"
	"github.com/tarantool/tarantool-operator/pkg/topology"
	"github.com/tarantool/tarantool-operator/test/mocks"
	"github.com/tarantool/tarantool-operator/test/resources"
	"github.com/tarantool/tarantool-operator/test/utils"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("cluster_controller unit testing", func() {
	var (
		ctx         = context.Background()
		namespace   = "default"
		clusterName string
	)

	BeforeEach(func() {
		clusterName = fmt.Sprintf("cluster-%s", utils.RandStringRunes(4))
	})

	Context("Common logic", func() {
		Describe("cluster controller should reconcile deletion", func() {
			It("must accept and process request with cluster which does not exists", func() {
				fakeClient := fake.NewClientBuilder().Build()
				fakeTopologyService := new(mocks.FakeCartridgeTopology)

				labelsManager := &k8s.NamespacedLabelsManager{
					Namespace: "tarantool.io",
				}

				resourcesManager := &ResourcesManager{
					LabelsManager: labelsManager,
					CommonResourcesManager: &k8s.CommonResourcesManager{
						Client: fakeClient,
						Scheme: scheme.Scheme,
					},
				}

				eventsRecorder := events.NewRecorder(record.NewFakeRecorder(10))

				clusterReconciler := &ClusterReconciler{
					LabelsManager: labelsManager,
					SteppedReconciler: &reconciliation.SteppedReconciler[*ClusterContextCE, *ClusterControllerCE]{
						Client: fakeClient,
						Controller: &ClusterControllerCE{
							CommonClusterController: &reconciliation.CommonClusterController{
								CommonController: &reconciliation.CommonController{
									Client: fakeClient,
									Schema: scheme.Scheme,
									LeaderElection: &election.LeaderElection{
										Client:           fakeClient,
										Recorder:         eventsRecorder,
										ResourcesManager: resourcesManager,
									},
									ResourcesManager: resourcesManager,
									LabelsManager:    labelsManager,
									EventsRecorder:   eventsRecorder,
									Topology:         fakeTopologyService,
								},
							},
						},
					},
				}

				reconcile, err := clusterReconciler.Reconcile(ctx, utils.ReconcileRequest(namespace, clusterName))
				if err != nil {
					return
				}

				Expect(err).NotTo(HaveOccurred(), "an error during reconcile")
				Expect(reconcile.Requeue).To(BeFalse(), "should not be re-queued")
				Expect(reconcile.RequeueAfter).To(BeZero(), "should not be re-queued")
			})
		})
	})

	Context("cluster bootstrapping", func() {
		Describe("cluster should control vshard bootstrapping", func() {
			var cartridge *resources.FakeCartridge
			var fakeTopologyService *mocks.FakeCartridgeTopology
			labelsManager := &k8s.NamespacedLabelsManager{
				Namespace: "tarantool.io",
			}

			BeforeEach(func() {
				cartridge = resources.NewFakeCartridge(labelsManager).
					WithNamespace(namespace).
					WithClusterName(clusterName).
					WithRouterRole(2, 1).
					WithStorageRole(2, 3).
					WithRouterStatefulSetsCreated().
					WithStorageStatefulSetsCreated().
					WithRouterPodsCreated().
					WithStoragePodsCreated()

				fakeTopologyService = new(mocks.FakeCartridgeTopology)

				fakeTopologyService.
					On("BootstrapVshard", mock.Anything).
					Return(nil)
			})

			It("Must bootstrap cluster if all roles ready", func() {
				cartridge.
					WithAllRolesInPhase(v1alpha2.RoleWaitingForBootstrap).
					WithAllPodsReady()

				fakeTopologyService.
					On("BootstrapVshard", mock.Anything, mock.Anything).
					Return(nil).
					Once()

				fakeClient := cartridge.BuildFakeClient()

				resourcesManager := &ResourcesManager{
					LabelsManager: labelsManager,
					CommonResourcesManager: &k8s.CommonResourcesManager{
						Client: fakeClient,
						Scheme: scheme.Scheme,
					},
				}

				fakeTopologyService.
					On("GetFailoverParams", mock.Anything, mock.Anything).
					Return(&topology.FailoverParams{}, nil)

				eventsRecorder := events.NewRecorder(record.NewFakeRecorder(10))

				clusterReconciler := &ClusterReconciler{
					LabelsManager: labelsManager,
					SteppedReconciler: &reconciliation.SteppedReconciler[*ClusterContextCE, *ClusterControllerCE]{
						Client: fakeClient,
						Controller: &ClusterControllerCE{
							CommonClusterController: &reconciliation.CommonClusterController{
								CommonController: &reconciliation.CommonController{
									Client: fakeClient,
									Schema: scheme.Scheme,
									LeaderElection: &election.LeaderElection{
										Client:           fakeClient,
										Recorder:         eventsRecorder,
										ResourcesManager: resourcesManager,
									},
									ResourcesManager: resourcesManager,
									LabelsManager:    labelsManager,
									EventsRecorder:   eventsRecorder,
									Topology:         fakeTopologyService,
								},
							},
						},
					},
				}
				_, err := clusterReconciler.Reconcile(ctx, ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: namespace,
						Name:      clusterName,
					},
				})
				Expect(err).NotTo(HaveOccurred(), "an error during reconcile")

				err = fakeClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: clusterName}, cartridge.Cluster)
				Expect(err).NotTo(HaveOccurred(), "cluster gone")

				Expect(cartridge.Cluster.Status.Bootstrapped).To(BeTrue(), "cluster not bootstrapped")
			})

			It("Must NOT bootstrap cluster if NOT all roles ready", func() {
				fakeClient := cartridge.BuildFakeClient()

				resourcesManager := &ResourcesManager{
					LabelsManager: labelsManager,
					CommonResourcesManager: &k8s.CommonResourcesManager{
						Client: fakeClient,
						Scheme: scheme.Scheme,
					},
				}

				eventsRecorder := events.NewRecorder(record.NewFakeRecorder(10))

				clusterReconciler := &ClusterReconciler{
					LabelsManager: labelsManager,
					SteppedReconciler: &reconciliation.SteppedReconciler[*ClusterContextCE, *ClusterControllerCE]{
						Client: fakeClient,
						Controller: &ClusterControllerCE{
							CommonClusterController: &reconciliation.CommonClusterController{
								CommonController: &reconciliation.CommonController{
									Client: fakeClient,
									Schema: scheme.Scheme,
									LeaderElection: &election.LeaderElection{
										Client:           fakeClient,
										Recorder:         eventsRecorder,
										ResourcesManager: resourcesManager,
									},
									ResourcesManager: resourcesManager,
									LabelsManager:    labelsManager,
									EventsRecorder:   eventsRecorder,
									Topology:         fakeTopologyService,
								},
							},
						},
					},
				}

				_, err := clusterReconciler.Reconcile(ctx, ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: namespace,
						Name:      clusterName,
					},
				})
				Expect(err).NotTo(HaveOccurred(), "an error during reconcile")

				err = fakeClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: clusterName}, cartridge.Cluster)
				Expect(err).NotTo(HaveOccurred())

				Expect(cartridge.Cluster.Status.Bootstrapped).To(BeFalse(), "cluster unexpectedly bootstrapped")
			})

			It("Must NOT bootstrap cluster if it is already bootstrapped", func() {
				cartridge.Bootstrapped()
				fakeClient := cartridge.BuildFakeClient()

				resourcesManager := &ResourcesManager{
					LabelsManager: labelsManager,
					CommonResourcesManager: &k8s.CommonResourcesManager{
						Client: fakeClient,
						Scheme: scheme.Scheme,
					},
				}

				eventsRecorder := events.NewRecorder(record.NewFakeRecorder(10))

				clusterReconciler := &ClusterReconciler{
					LabelsManager: labelsManager,
					SteppedReconciler: &reconciliation.SteppedReconciler[*ClusterContextCE, *ClusterControllerCE]{
						Client: fakeClient,
						Controller: &ClusterControllerCE{
							CommonClusterController: &reconciliation.CommonClusterController{
								CommonController: &reconciliation.CommonController{
									Client: fakeClient,
									Schema: scheme.Scheme,
									LeaderElection: &election.LeaderElection{
										Client:           fakeClient,
										Recorder:         eventsRecorder,
										ResourcesManager: resourcesManager,
									},
									ResourcesManager: resourcesManager,
									LabelsManager:    labelsManager,
									EventsRecorder:   eventsRecorder,
									Topology:         fakeTopologyService,
								},
							},
						},
					},
				}

				_, err := clusterReconciler.Reconcile(ctx, ctrl.Request{
					NamespacedName: types.NamespacedName{
						Namespace: namespace,
						Name:      clusterName,
					},
				})
				Expect(err).NotTo(HaveOccurred(), "an error during reconcile")

				err = fakeClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: clusterName}, cartridge.Cluster)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
