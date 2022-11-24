package resources

const (
	FakeDockerImageV1 = "fake/cartridge:1.0.0"
	FakeDockerImageV2 = "fake/cartridge:2.0.0"
)

const PodContainerName = "cartridge"

const (
	DefaultCpuLimit     = 1000
	DefaultMemoryLimit  = 256 * 1024 * 1024
	DefaultStorageLimit = 1 * 1024 * 1024 * 1024
)

const (
	DefaultDomain                = "Cluster.local"
	DefaultListenPort            = 3301
	ConstDefaultReplicasetWeight = 100
)

const (
	RoleRouter  = "router"
	RoleStorage = "storage"
)
