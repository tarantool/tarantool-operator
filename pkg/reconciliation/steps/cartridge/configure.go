package cartridge

import (
	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
	"github.com/tarantool/tarantool-operator/pkg/utils"
	"gopkg.in/yaml.v3"
)

type ConfigureStep[ConfigType api.CartridgeConfig, CtxType CartridgeConfigContext[ConfigType], CtrlType CartridgeConfigController] struct{}

func (r *ConfigureStep[ConfigType, CtxType, CtrlType]) GetName() string {
	return "Configure cartridge"
}

func (r *ConfigureStep[ConfigType, CtxType, CtrlType]) Reconcile(ctx CtxType, ctrl CtrlType) (*Result, error) {
	actualConfig, err := ctrl.GetTopology().GetCartridgeConfig(ctx, ctx.GetLeader())
	if err != nil {
		return Error(err)
	}

	var desiredConfig map[string]interface{}

	err = yaml.Unmarshal(ctx.GetCartridgeConfig().GetData(), &desiredConfig)
	if err != nil {
		return Error(err)
	}

	if utils.IsMapSubset(actualConfig, desiredConfig) {
		ctx.GetLogger().Info("Nothing to change in config")

		return Complete()
	}

	err = ctrl.GetTopology().ApplyCartridgeConfig(ctx, ctx.GetLeader(), desiredConfig)
	if err != nil {
		ctx.GetLogger().Error(err, "Unable to apply cartridge config")

		return Error(err)
	}

	return NextStep()
}
