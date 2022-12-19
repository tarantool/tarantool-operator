package cluster

import (
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/tarantool/tarantool-operator/pkg/api"
	. "github.com/tarantool/tarantool-operator/pkg/reconciliation"
	"github.com/tarantool/tarantool-operator/pkg/topology"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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
		var (
			secret   *corev1.Secret
			password []byte
			err      error
			ok       bool
		)

		switch config.GetStateProvider() {
		case api.FailoverStateProviderETCD2:
			etcd2config := config.GetETCD2Config()
			etcd2PasswordRef := etcd2config.GetPassword()

			if etcd2PasswordRef.GetName() != "" {
				secret, err = ctrl.GetResourcesManager().GetSecret(ctx, etcd2PasswordRef.GetNamespace(), etcd2PasswordRef.GetName())
				if err != nil {
					return nil, err
				}

				etcd2secretKey := etcd2PasswordRef.GetKey()
				if etcd2secretKey == "" {
					etcd2secretKey = "etcd2-password"
				}

				password, ok = secret.Data[etcd2secretKey]
				if !ok {
					return nil, errors.New("no such key")
				}
			}

			params.Etcd2Params = &topology.Etcd2Params{
				Endpoints: etcd2config.GetEndpoints(),
				Username:  etcd2config.GetUsername(),
				Password:  string(password),
				LockDelay: etcd2config.GetLockDelay(),
				Prefix:    etcd2config.GetPrefix(),
			}
		case api.FailoverStateProviderStateboard:
			stateboardConfig := config.GetStateboardConfig()
			stateboardPasswordRef := stateboardConfig.GetPassword()

			if stateboardPasswordRef.GetName() != "" {
				secret, err = ctrl.GetResourcesManager().GetSecret(
					ctx,
					stateboardPasswordRef.GetNamespace(),
					stateboardPasswordRef.GetName(),
				)
				if err != nil {
					return nil, err
				}

				stateboardSecretKey := stateboardPasswordRef.GetKey()
				if stateboardSecretKey == "" {
					stateboardSecretKey = "stateboard-password"
				}

				password, ok = secret.Data[stateboardSecretKey]
				if !ok {
					return nil, errors.New("no such key")
				}
			}

			params.StateboardParams = &topology.StateboardParams{
				URI:      stateboardConfig.GetURI(),
				Password: string(password),
			}
		}

		if secret != nil {
			err = controllerutil.SetOwnerReference(ctx.GetCluster(), secret, ctrl.GetScheme())
			if err != nil {
				return nil, err
			}

			err = ctrl.Update(ctx, secret)
			if err != nil {
				return nil, err
			}
		}
	}

	return params, nil
}
