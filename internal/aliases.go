package internal

import (
	"github.com/tarantool/tarantool-operator/internal/context"
	"github.com/tarantool/tarantool-operator/internal/controller"
)

type (
	ClusterControllerCE = controller.ClusterController
	ClusterContextCE    = context.ClusterContext
)

type (
	RoleControllerCE = controller.RoleController
	RoleContextCE    = context.RoleContext
)

type (
	CartridgeConfigControllerCE = controller.CartridgeConfigController
	CartridgeConfigContextCE    = context.CartridgeConfigContext
)
