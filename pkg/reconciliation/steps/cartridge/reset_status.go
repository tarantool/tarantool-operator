package cartridge

import (
	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
)

type ResetStatusStep[ConfigType api.CartridgeConfig, CtxType CartridgeConfigContext[ConfigType], CtrlType CartridgeConfigController] struct{}

func (r *ResetStatusStep[ConfigType, CtxType, CtrlType]) GetName() string {
	return "Reset cartridge config status"
}

func (r *ResetStatusStep[ConfigType, CtxType, CtrlType]) Reconcile(ctx CtxType, _ CtrlType) (*Result, error) {
	ctx.GetCartridgeConfig().ResetStatus()

	return NextStep()
}
