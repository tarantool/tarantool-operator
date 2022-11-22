package context

import (
	"github.com/pkg/errors"
	"github.com/tarantool/tarantool-operator/apis/v1alpha2"
	"github.com/tarantool/tarantool-operator/pkg/reconciliation"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CartridgeConfigContext struct {
	*reconciliation.CommonContext

	CartridgeConfig *v1alpha2.CartridgeConfig
}

func (r *CartridgeConfigContext) SetCartridgeConfig(config *v1alpha2.CartridgeConfig) {
	r.CartridgeConfig = config
}

func (r *CartridgeConfigContext) GetCartridgeConfig() *v1alpha2.CartridgeConfig {
	return r.CartridgeConfig
}

func (r *CartridgeConfigContext) HasRequestedObject() bool {
	return r.CartridgeConfig != nil
}

func (r *CartridgeConfigContext) SetRequestedObject(obj client.Object) error {
	cartridgeConfig, ok := obj.(*v1alpha2.CartridgeConfig)
	if !ok {
		return errors.New("CartridgeConfigContext used with wrong k8s object")
	}

	r.CartridgeConfig = cartridgeConfig

	return nil
}

func (r *CartridgeConfigContext) GetRequestedObject() client.Object {
	return r.CartridgeConfig
}
