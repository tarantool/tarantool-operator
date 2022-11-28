package steps

import (
	"github.com/tarantool/tarantool-operator/apis/v1beta1"
	. "github.com/tarantool/tarantool-operator/internal/context"
	. "github.com/tarantool/tarantool-operator/internal/controller"
	"github.com/tarantool/tarantool-operator/pkg/reconciliation/steps/cartridge"
)

func ResetCartridgeConfigStatus() *cartridge.ResetStatusStep[*v1beta1.CartridgeConfig, *CartridgeConfigContext, *CartridgeConfigController] {
	return &cartridge.ResetStatusStep[*v1beta1.CartridgeConfig, *CartridgeConfigContext, *CartridgeConfigController]{}
}

func SetCartridgeConfigPhase(phase v1beta1.CartridgeConfigPhase) *cartridge.SetPhaseStep[v1beta1.CartridgeConfigPhase, *v1beta1.CartridgeConfig, *CartridgeConfigContext, *CartridgeConfigController] {
	return &cartridge.SetPhaseStep[v1beta1.CartridgeConfigPhase, *v1beta1.CartridgeConfig, *CartridgeConfigContext, *CartridgeConfigController]{
		Phase: phase,
	}
}

func ConfigureCartridge() *cartridge.ConfigureStep[*v1beta1.CartridgeConfig, *CartridgeConfigContext, *CartridgeConfigController] {
	return &cartridge.ConfigureStep[*v1beta1.CartridgeConfig, *CartridgeConfigContext, *CartridgeConfigController]{}
}
