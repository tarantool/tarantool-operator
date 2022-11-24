package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/tarantool/tarantool-operator/pkg/topology"
	"k8s.io/api/core/v1"
)

type FakeCartridgeTopology struct {
	mock.Mock
}

func (f *FakeCartridgeTopology) Exec(ctx context.Context, instance *v1.Pod, res interface{}, lua string, args ...interface{}) error {
	callArgs := []interface{}{ctx, instance, res, lua}
	callArgs = append(callArgs, args...)
	called := f.Called(callArgs...)

	return called.Error(0)
}

func (f *FakeCartridgeTopology) Join(
	ctx context.Context,
	leader *v1.Pod,
	replicasetAlias string,
	replicasetUUID string,
	replicasetRoles []string,
	replicasetWeight int32,
	replicasetVshardGroup string,
	replicasetIsAllRw bool,
	advertiseURI string,
) error {
	args := f.Called(
		ctx,
		leader,
		replicasetAlias,
		replicasetUUID,
		replicasetRoles,
		replicasetWeight,
		replicasetVshardGroup,
		replicasetIsAllRw,
		advertiseURI,
	)
	return args.Error(0)
}

func (f *FakeCartridgeTopology) GetInstanceUUID(ctx context.Context, pod *v1.Pod) (string, error) {
	args := f.Called(ctx, pod)

	return args.String(0), args.Error(1)
}

func (f *FakeCartridgeTopology) SetWeight(ctx context.Context, leader *v1.Pod, replicasetUUID string, replicaWeight int32) error {
	args := f.Called(ctx, leader, replicasetUUID, replicaWeight)

	return args.Error(0)
}

func (f *FakeCartridgeTopology) GetReplicasetRoles(ctx context.Context, leader *v1.Pod, replicasetUUID string) ([]string, error) {
	args := f.Called(ctx, leader, replicasetUUID)

	return args.Get(0).([]string), args.Error(1)
}

func (f *FakeCartridgeTopology) SetReplicasetRoles(ctx context.Context, leader *v1.Pod, replicasetUUID string, roles []string) error {
	args := f.Called(ctx, leader, replicasetUUID, roles)

	return args.Error(0)
}

func (f *FakeCartridgeTopology) BootstrapVshard(ctx context.Context, leader *v1.Pod) error {
	args := f.Called(ctx, leader)

	return args.Error(0)
}

func (f *FakeCartridgeTopology) GetRolesHierarchy(ctx context.Context, leader *v1.Pod) (map[string][]string, error) {
	args := f.Called(ctx, leader)

	return args.Get(0).(map[string][]string), args.Error(1)
}

func (f *FakeCartridgeTopology) SetFailoverParams(ctx context.Context, leader *v1.Pod, params *topology.FailoverParams) error {
	args := f.Called(ctx, leader, params)

	return args.Error(0)
}

func (f *FakeCartridgeTopology) GetFailoverParams(ctx context.Context, leader *v1.Pod) (*topology.FailoverParams, error) {
	args := f.Called(ctx, leader)

	return args.Get(0).(*topology.FailoverParams), args.Error(1)
}

func (f *FakeCartridgeTopology) GetCartridgeConfig(ctx context.Context, leader *v1.Pod) (topology.CartridgeConfigData, error) {
	args := f.Called(ctx, leader)

	return args.Get(0).(topology.CartridgeConfigData), args.Error(1)
}

func (f *FakeCartridgeTopology) ApplyCartridgeConfig(ctx context.Context, leader *v1.Pod, config topology.CartridgeConfigData) error {
	args := f.Called(ctx, leader, config)

	return args.Error(0)
}
