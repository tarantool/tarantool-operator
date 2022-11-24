package resources

const FakeDockerImageV1 = "fake/cartridge:1.0.0"
const FakeDockerImageV2 = "fake/cartridge:2.0.0"

const PodContainerName = "cartridge"

const DefaultCpuLimit = 1000
const DefaultMemoryLimit = 256 * 1024 * 1024
const DefaultStorageLimit = 1 * 1024 * 1024 * 1024

const DefaultDomain = "Cluster.local"
const DefaultListenPort = 3301
const ConstDefaultReplicasetWeight = 100

const RoleRouter = "router"
const RoleStorage = "storage"
