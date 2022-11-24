package resources

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func (r *FakeCartridge) NewCartridgeContainer() v1.Container {
	return v1.Container{
		Name:  PodContainerName,
		Image: FakeDockerImageV1,
		VolumeMounts: []v1.VolumeMount{
			{
				Name:      "data",
				MountPath: "/var/lib/tarantool",
			},
		},
		Resources: v1.ResourceRequirements{
			Limits: v1.ResourceList{
				"cpu":    *resource.NewMilliQuantity(DefaultCpuLimit, resource.DecimalSI),
				"memory": *resource.NewQuantity(DefaultMemoryLimit, resource.BinarySI),
			},
		},
		Ports: []v1.ContainerPort{
			{
				Name:          "app",
				Protocol:      v1.ProtocolTCP,
				ContainerPort: r.Cluster.GetListenPort(),
			},
			{
				Name:          "app-udp",
				Protocol:      v1.ProtocolUDP,
				ContainerPort: r.Cluster.GetListenPort(),
			},
		},
	}
}
