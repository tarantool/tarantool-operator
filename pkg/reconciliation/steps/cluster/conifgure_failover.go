package cluster

import (
	"github.com/google/go-cmp/cmp"
	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
	"github.com/tarantool/tarantool-operator/pkg/topology"
)

type ConfigureFailoverStep[ClusterType api.Cluster, CtxType ClusterContext[ClusterType], CtrlType ClusterController] struct{}

func (r *ConfigureFailoverStep[ClusterType, CtxType, CtrlType]) GetName() string {
	return "Configure failover"
}

func (r *ConfigureFailoverStep[ClusterType, CtxType, CtrlType]) Reconcile(ctx CtxType, ctrl CtrlType) (*Result, error) {
	config := ctx.GetCluster().GetFailoverConfig()
	if config == nil {
		return NextStep()
	}

	params, err := r.LoadFailoverParams(ctx, ctrl)
	if err != nil {
		ctx.GetLogger().Error(err, "failed to enable cluster failover, unable to retrieve params")

		return Error(err)
	}

	currentParams, err := ctrl.GetTopology().GetFailoverParams(ctx, ctx.GetLeader())
	if err != nil {
		ctx.GetLogger().Error(
			err,
			"failed to enable cluster failover, unable to retrieve current params",
		)

		return Error(err)
	}

	if cmp.Equal(params, currentParams) {
		ctx.GetLogger().Info("failover has same configuration, not retrying")

		return NextStep()
	}

	if err = ctrl.GetTopology().SetFailoverParams(ctx, ctx.GetLeader(), params); err != nil {
		ctx.GetLogger().Error(err, "failed to enable cluster failover")

		return Error(err)
	}

	ctx.GetLogger().Info("failover enabled")

	return NextStep()
}

func (r *ConfigureFailoverStep[ClusterType, CtxType, CtrlType]) LoadFailoverParams(ctx CtxType, ctrl CtrlType) (*topology.FailoverParams, error) {
	config := ctx.GetCluster().GetFailoverConfig()
	params := &topology.FailoverParams{
		Mode:           config.GetMode(),
		Timeout:        config.GetTimeout(),
		StateProvider:  config.GetStateProvider(),
		FencingEnabled: config.GetFencing(),
		FencingTimeout: config.GetFencingTimeout(),
		FencingPause:   config.GetFencingPause(),
	}

	if config.GetMode() == api.FailoverModeStateful {
		switch config.GetStateProvider() {
		case api.FailoverStateProviderETCD2:
			var (
				etcd2Password string
				err           error
			)

			etcd2config := config.GetETCD2Config()
			etcd2PasswordRef := etcd2config.GetPassword()

			if etcd2PasswordRef.Name != "" {
				etcd2Password, err = ctrl.GetResourcesManager().GetSecretValue(
					ctx,
					etcd2PasswordRef.Namespace,
					etcd2PasswordRef.Name,
					"etcd2-password", // fixme: make not hardcoded
				)
				if err != nil {
					return nil, err
				}
			}

			params.Etcd2Params = &topology.Etcd2Params{
				Endpoints: etcd2config.GetEndpoints(),
				Username:  etcd2config.GetUsername(),
				Password:  etcd2Password,
				LockDelay: etcd2config.GetLockDelay(),
				Prefix:    etcd2config.GetPrefix(),
			}
		case api.FailoverStateProviderStateboard:
			var (
				stateboardPassword string
				err                error
			)

			stateboardConfig := config.GetStateboardConfig()
			stateboardPasswordRef := stateboardConfig.GetPassword()

			if stateboardPasswordRef.Name != "" {
				stateboardPassword, err = ctrl.GetResourcesManager().GetSecretValue(
					ctx,
					stateboardPasswordRef.Namespace,
					stateboardPasswordRef.Name,
					"stateboard-password", // fixme: make not hardcoded
				)
				if err != nil {
					return nil, err
				}
			}

			params.StateboardParams = &topology.StateboardParams{
				URI:      stateboardConfig.GetURI(),
				Password: stateboardPassword,
			}
		}
	}

	return params, nil
}
