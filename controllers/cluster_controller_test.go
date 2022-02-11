package controllers

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	helpers "github.com/tarantool/tarantool-operator/test/helpers"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	tarantooliov1alpha1 "github.com/tarantool/tarantool-operator/api/v1alpha1"
	"github.com/tarantool/tarantool-operator/controllers/utils"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type MockClient struct {
	Items []client.Object
}

func (c MockClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object) error {
	if obj == nil {
		return fmt.Errorf("required target object")
	}
	for _, o := range c.Items {
		if key.Name == o.GetName() && key.Namespace == o.GetNamespace() {
			obj = o
		}
	}
	return nil
}

func (c MockClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return nil
}

var _ = Describe("cluster_controller unit testing", func() {
	var (
		namespace = "default"
		ctx       = context.TODO()

		roleName       = "" // setup for every spec in hook
		rsTemplateName = ""

		clusterName = "test"
		clusterId   = clusterName

		defaultRolesToAssign = "[\"A\",\"B\"]"
		podName              = ""
		serviceName          = clusterName
	)

	Describe("cluster_controller manage cluster resources", func() {
		BeforeEach(func() {
			// setup variables for each spec
			roleName = fmt.Sprintf("test-role-%s", RandStringRunes(4))
			rsTemplateName = fmt.Sprintf("test-rs-%s", RandStringRunes(4))

			By("create new Role " + roleName)
			role := helpers.NewRole(helpers.RoleParams{
				Name:           roleName,
				Namespace:      namespace,
				RolesToAssign:  defaultRolesToAssign,
				RsNum:          int32(1),
				RsTemplateName: rsTemplateName,
				ClusterId:      clusterId,
			})
			// mock owner reference
			role.SetOwnerReferences([]metav1.OwnerReference{
				{
					APIVersion: "v0",
					Kind:       "mockRef",
					Name:       "mockRef",
					UID:        "-",
				},
			})
			Expect(k8sClient.Create(ctx, &role)).NotTo(HaveOccurred(), "failed to create Role")

			By("create new Cluster " + clusterName)
			cluster := helpers.NewCluster(helpers.ClusterParams{
				Name:      clusterName,
				Namespace: namespace,
				Id:        clusterId,
			})
			Expect(k8sClient.Create(ctx, &cluster)).NotTo(HaveOccurred(), "failed to create Cluster")
		})

		AfterEach(func() {
			By("remove role object " + roleName)
			role := &tarantooliov1alpha1.Role{}
			Expect(
				k8sClient.Get(ctx, client.ObjectKey{Name: roleName, Namespace: namespace}, role),
			).NotTo(HaveOccurred(), "failed to get Role")

			Expect(k8sClient.Delete(ctx, role)).NotTo(HaveOccurred(), "failed to delete Role")

			By("remove Cluster object " + clusterName)
			cluster := &tarantooliov1alpha1.Cluster{}
			Expect(
				k8sClient.Get(ctx, client.ObjectKey{Name: clusterName, Namespace: namespace}, cluster),
			).NotTo(HaveOccurred(), "failed to get Cluster")

			Expect(k8sClient.Delete(ctx, cluster)).NotTo(HaveOccurred(), "failed to delete Cluster")
		})

		Context("manage cluster leader", func() {
			BeforeEach(func() {
				// By ReplicasetTemplate object Role controller creating StatefulSets
				rsTemplate := helpers.NewReplicasetTemplate(helpers.ReplicasetTemplateParams{
					Name:           rsTemplateName,
					Namespace:      namespace,
					RoleName:       roleName,
					RolesToAssign:  defaultRolesToAssign,
					ContainerImage: "image:0.0.0",
					ContainerName:  "test",
					ServiceName:    serviceName,
				})
				Expect(k8sClient.Create(ctx, &rsTemplate)).NotTo(HaveOccurred(), "failed to create ReplicasetTemplate")

				// But, if tests running in envtest k8s cluster
				// StatefulSet controller not active and we need to run
				// required pod handly
				podName = roleName + "-0-0"
				pod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      podName,
						Namespace: namespace,
						Labels: map[string]string{
							"tarantool.io/cluster-id": clusterId,
						},
					},
					Spec: rsTemplate.Spec.Template.Spec,
				}
				Expect(k8sClient.Create(ctx, pod)).NotTo(HaveOccurred(), "failed to create Pod")
			})

			AfterEach(func() {
				rsTemplate := &tarantooliov1alpha1.ReplicasetTemplate{}
				Expect(
					k8sClient.Get(ctx, client.ObjectKey{Name: rsTemplateName, Namespace: namespace}, rsTemplate),
				).NotTo(HaveOccurred(), "failed to get ReplicasetTemplate")

				Expect(
					k8sClient.Delete(ctx, rsTemplate),
				).NotTo(HaveOccurred(), "failed to delete ReplicasetTemplate")

				pod := &corev1.Pod{}
				Expect(
					k8sClient.Get(ctx, client.ObjectKey{Name: podName, Namespace: namespace}, pod),
				).NotTo(HaveOccurred(), "failed to get Pod")

				Expect(
					k8sClient.Delete(ctx, pod),
				).NotTo(HaveOccurred(), "failed to delete Pod")
			})

			It("change leader if it is not defined", func() {
				expectedLeader := utils.MakeStaticPodAddr(podName, serviceName, namespace, "", 8081)
				Eventually(
					func() bool {
						cluster := tarantooliov1alpha1.Cluster{}
						err := k8sClient.Get(ctx, client.ObjectKey{Name: clusterId, Namespace: namespace}, &cluster)
						if err != nil {
							return false
						}

						leader := cluster.GetAnnotations()["tarantool.io/topology-manage-leader"]
						return leader == expectedLeader
					},
					time.Second*10, time.Millisecond*500,
				).Should(BeTrue())
			})
		})
	})

	Describe("cluster_controller.IsTopologyManageLeaderExists test", func() {
		Describe("the function must return true if the leader is exist and false if not exist", func() {
			Context("positive cases (leader exist)", func() {
				It("return True if leader assigned and exist", func() {
					replicasNum := int32(1)

					sts := appsv1.StatefulSet{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "sts-0",
							Namespace: "ns",
						},
						Spec: appsv1.StatefulSetSpec{
							ServiceName: "svc-0",
							Replicas:    &replicasNum,
						},
					}

					stsList := appsv1.StatefulSetList{
						Items: []appsv1.StatefulSet{sts},
					}

					podName := sts.Name + "-0"
					leader := utils.MakeStaticPodAddr(
						podName, sts.Spec.ServiceName, sts.Namespace, "", 8081)

					c := MockClient{
						Items: []client.Object{&stsList.Items[0]},
					}

					Expect(IsTopologyManageLeaderExists(c, &stsList, leader)).To(BeTrue())
				})
			})

			Context("negative cases (leader does not exist)", func() {
				It("return False if stsList empty", func() {
					c := MockClient{}
					stsList := appsv1.StatefulSetList{}

					leader := utils.MakeStaticPodAddr("pod", "svc", "ns", "", 8081)
					Expect(IsTopologyManageLeaderExists(c, &stsList, leader)).To(BeFalse())
				})

				It("return False if leader not in the StatefulSet list", func() {
					replicasNum := int32(1)

					sts := appsv1.StatefulSet{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "sts-0",
							Namespace: "ns",
						},
						Spec: appsv1.StatefulSetSpec{
							ServiceName: "svc-0",
							Replicas:    &replicasNum,
						},
					}

					stsList := appsv1.StatefulSetList{
						Items: []appsv1.StatefulSet{sts},
					}

					leader := utils.MakeStaticPodAddr("unexistent", "svc", "ns", "", 8081)

					c := MockClient{
						Items: []client.Object{&stsList.Items[0]},
					}

					Expect(IsTopologyManageLeaderExists(c, &stsList, leader)).To(BeFalse())
				})
			})
		})
	})
})
