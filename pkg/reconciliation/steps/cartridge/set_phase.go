package cartridge

import (
	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
)

type SetPhaseStep[PhaseType comparable, ConfigType api.CartridgeConfigWithStatus[PhaseType], CtxType CartridgeConfigContext[ConfigType], CtrlType CartridgeConfigController] struct {
	Phase PhaseType
}

func (r *SetPhaseStep[PhaseType, ConfigType, CtxType, CtrlType]) GetName() string {
	return "Set cartridge config phase"
}

func (r *SetPhaseStep[PhaseType, ConfigType, CtxType, CtrlType]) Reconcile(ctx CtxType, _ CtrlType) (*Result, error) {
	ctx.GetCartridgeConfig().SetPhase(r.Phase)

	return NextStep()
}
