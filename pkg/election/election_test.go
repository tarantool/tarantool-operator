package election_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	. "github.com/tarantool/tarantool-operator/internal/implementation"
	. "github.com/tarantool/tarantool-operator/pkg/election"
	"github.com/tarantool/tarantool-operator/pkg/events"
	"github.com/tarantool/tarantool-operator/pkg/k8s"
	"github.com/tarantool/tarantool-operator/test/mocks"
	"github.com/tarantool/tarantool-operator/test/resources"
	"github.com/tarantool/tarantool-operator/test/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func newTestElection(fakeClient client.Client, fakeTopologyService *mocks.FakeCartridgeTopology) *LeaderElection {
	resourcesManager := &ResourcesManager{
		LabelsManager: labelsManager,
		CommonResourcesManager: &k8s.CommonResourcesManager{
			Client: fakeClient,
			Scheme: scheme.Scheme,
		},
	}

	return &LeaderElection{
		Client:           fakeClient,
		Recorder:         events.NewRecorder(record.NewFakeRecorder(10)),
		ResourcesManager: resourcesManager,
		Topology:         fakeTopologyService,
	}
}

var _ = Describe("election unit testing", func() {
	var (
		namespace   = "default"
		clusterName string
		cartridge   *resources.FakeCartridge
	)

	BeforeEach(func() {
		clusterName = fmt.Sprintf("cluster-%s", utils.RandStringRunes(4))
		cartridge = resources.NewFakeCartridge(labelsManager).
			WithNamespace(namespace).
			WithClusterName(clusterName)
	})

	Context("Initial election", func() {
		It("should not elect foreign pods", func() {
			fakeClient := cartridge.NewFakeClientBuilder().WithObjects(
				&v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foreign-pod",
						Namespace: namespace,
					},
				},
			).Build()
			fakeTopologyService := new(mocks.FakeCartridgeTopology)

			election := newTestElection(fakeClient, fakeTopologyService)

			leader, err := election.GetLeaderInstance(ctx, cartridge.Cluster)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeEquivalentTo(ErrNoAvailableLeader))
			Expect(leader).To(BeNil(), "leader should not be elected")
		})

		It("should skip pod which is not running", func() {
			cartridge.WithRouterRole(2, 1)
			cartridge.WithStorageRole(2, 2)
			cartridge.WithRouterStatefulSetsCreated()
			cartridge.WithStorageStatefulSetsCreated()
			cartridge.WithRouterPodsCreated()
			cartridge.WithStoragePodsCreated()

			fakeClient := cartridge.BuildFakeClient()
			fakeTopologyService := new(mocks.FakeCartridgeTopology)

			election := newTestElection(fakeClient, fakeTopologyService)

			leader, err := election.GetLeaderInstance(ctx, cartridge.Cluster)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeEquivalentTo(ErrNoAvailableLeader))
			Expect(leader).To(BeNil(), "leader should not be elected")
		})

		It("should skip deleting pod", func() {
			cartridge.WithRouterRole(2, 1)
			cartridge.WithStorageRole(2, 2)
			cartridge.WithRouterStatefulSetsCreated()
			cartridge.WithStorageStatefulSetsCreated()
			cartridge.WithRouterPodsCreated()
			cartridge.WithStoragePodsCreated()
			cartridge.WithAllPodsDeleting()

			fakeClient := cartridge.BuildFakeClient()
			fakeTopologyService := new(mocks.FakeCartridgeTopology)

			election := newTestElection(fakeClient, fakeTopologyService)

			leader, err := election.GetLeaderInstance(ctx, cartridge.Cluster)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeEquivalentTo(ErrNoAvailableLeader))
			Expect(leader).To(BeNil(), "leader should not be elected")
		})

		It("should skip pods where cartridge not started", func() {
			cartridge.WithRouterRole(1, 1)
			cartridge.WithStorageRole(1, 1)
			cartridge.WithRouterStatefulSetsCreated()
			cartridge.WithStorageStatefulSetsCreated()
			cartridge.WithRouterPodsCreated()
			cartridge.WithStoragePodsCreated()
			cartridge.WithAllPodsRunning()

			fakeClient := cartridge.BuildFakeClient()
			fakeTopologyService := new(mocks.FakeCartridgeTopology)
			fakeTopologyService.
				On("IsCartridgeStarted", mock.Anything, mock.Anything).
				Return(false, nil)

			election := newTestElection(fakeClient, fakeTopologyService)

			leader, err := election.GetLeaderInstance(ctx, cartridge.Cluster)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeEquivalentTo(ErrNoAvailableLeader))
			Expect(leader).To(BeNil(), "leader should not be elected")
		})

		It("should elected any non-foreign running pod", func() {
			cartridge.WithRouterRole(1, 1)
			cartridge.WithStorageRole(1, 1)
			cartridge.WithRouterStatefulSetsCreated()
			cartridge.WithStorageStatefulSetsCreated()
			cartridge.WithRouterPodsCreated()
			cartridge.WithStoragePodsCreated()
			cartridge.WithAllPodsRunning()

			fakeClient := cartridge.NewFakeClientBuilder().WithObjects(
				&v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foreign-pod",
						Namespace: namespace,
					},
				},
			).Build()

			fakeTopologyService := new(mocks.FakeCartridgeTopology)
			fakeTopologyService.
				On("IsCartridgeStarted", mock.Anything, mock.Anything).
				Return(true, nil)

			election := newTestElection(fakeClient, fakeTopologyService)

			leader, err := election.GetLeaderInstance(ctx, cartridge.Cluster)
			Expect(err).NotTo(HaveOccurred())
			Expect(leader).NotTo(BeNil(), "leader should be elected")
			Expect(leader.GetName()).NotTo(BeEquivalentTo("foreign-pod"), "should not elect foreign pods")
		})
	})

	Context("Re-election on bootstrapped cluster", func() {
		It("should not change leader if it's alive, and has config", func() {
			cartridge.WithRouterRole(1, 1)
			cartridge.WithStorageRole(1, 1)
			cartridge.WithRouterStatefulSetsCreated()
			cartridge.WithStorageStatefulSetsCreated()
			cartridge.WithRouterPodsCreated()
			cartridge.WithStoragePodsCreated()
			cartridge.WithAllPodsRunning()
			cartridge.WithLeader("router-0-0")
			cartridge.Bootstrapped()

			fakeClient := cartridge.BuildFakeClient()
			fakeTopologyService := new(mocks.FakeCartridgeTopology)
			fakeTopologyService.
				On("IsCartridgeConfigured", mock.Anything, mock.Anything).
				Return(true, nil)

			election := newTestElection(fakeClient, fakeTopologyService)

			leader, err := election.GetLeaderInstance(ctx, cartridge.Cluster)
			Expect(err).NotTo(HaveOccurred())
			Expect(leader).NotTo(BeNil(), "leader should be elected")
			Expect(leader.GetName()).To(BeEquivalentTo("router-0-0"), "leader should not be changed")
		})

		Describe("should change leader if it's dead", func() {
			It("should skip pod where no config present", func() {
				cartridge.WithRouterRole(1, 1)
				cartridge.WithStorageRole(2, 1)
				cartridge.WithRouterStatefulSetsCreated()
				cartridge.WithStorageStatefulSetsCreated()
				cartridge.WithRouterPodsCreated()
				cartridge.WithStoragePodsCreated()
				cartridge.WithPodsRunning("storage-0-0", "storage-1-0")
				cartridge.WithLeader("router-0-0")
				cartridge.Bootstrapped()

				fakeClient := cartridge.BuildFakeClient()
				fakeTopologyService := new(mocks.FakeCartridgeTopology)
				fakeTopologyService.
					On("IsCartridgeConfigured", mock.Anything, mock.MatchedBy(func(pod *v1.Pod) bool {
						return pod.GetName() == "storage-0-0"
					})).
					Return(false, nil)

				fakeTopologyService.
					On("IsCartridgeConfigured", mock.Anything, mock.MatchedBy(func(pod *v1.Pod) bool {
						return pod.GetName() == "storage-1-0"
					})).
					Return(true, nil)

				election := newTestElection(fakeClient, fakeTopologyService)

				leader, err := election.GetLeaderInstance(ctx, cartridge.Cluster)
				Expect(err).NotTo(HaveOccurred())
				Expect(leader).NotTo(BeNil(), "leader should be changed")
				Expect(leader.GetName()).To(BeEquivalentTo("storage-1-0"), "new leader should be configured pod")
			})
		})
	})

	Context("Re-election on not bootstrapped cluster", func() {
		It("should elect same leader if it's running", func() {
			cartridge.WithRouterRole(1, 1)
			cartridge.WithStorageRole(2, 1)
			cartridge.WithRouterStatefulSetsCreated()
			cartridge.WithStorageStatefulSetsCreated()
			cartridge.WithRouterPodsCreated()
			cartridge.WithStoragePodsCreated()
			cartridge.WithAllPodsRunning()
			cartridge.WithLeader("router-0-0")

			fakeClient := cartridge.BuildFakeClient()
			fakeTopologyService := new(mocks.FakeCartridgeTopology)
			fakeTopologyService.
				On("IsCartridgeStarted", mock.Anything, mock.Anything).
				Return(true, nil)

			election := newTestElection(fakeClient, fakeTopologyService)

			leader, err := election.GetLeaderInstance(ctx, cartridge.Cluster)
			Expect(err).NotTo(HaveOccurred())
			Expect(leader).NotTo(BeNil(), "leader should be elected")
			Expect(leader.GetName()).To(BeEquivalentTo("router-0-0"), "leader should not be changed")
		})

		It("should error if leader is not running", func() {
			cartridge.WithRouterRole(1, 1)
			cartridge.WithStorageRole(2, 1)
			cartridge.WithRouterStatefulSetsCreated()
			cartridge.WithStorageStatefulSetsCreated()
			cartridge.WithRouterPodsCreated()
			cartridge.WithStoragePodsCreated()
			cartridge.WithPodsRunning("storage-0-0", "storage-1-0")
			cartridge.WithLeader("router-0-0")

			fakeClient := cartridge.BuildFakeClient()
			fakeTopologyService := new(mocks.FakeCartridgeTopology)
			fakeTopologyService.
				On("IsCartridgeStarted", mock.Anything, mock.MatchedBy(func(pod *v1.Pod) bool {
					return pod.GetName() == "router-0-0"
				})).
				Return(false, nil)

			fakeTopologyService.
				On("IsCartridgeStarted", mock.Anything, mock.MatchedBy(func(pod *v1.Pod) bool {
					return pod.GetName() != "router-0-0"
				})).
				Return(true, nil)

			election := newTestElection(fakeClient, fakeTopologyService)

			leader, err := election.GetLeaderInstance(ctx, cartridge.Cluster)
			Expect(err).To(HaveOccurred())
			Expect(leader).To(BeNil(), "leader should not be elected")
			Expect(err).To(BeEquivalentTo(ErrLeaderNotReady))
		})
	})
})
