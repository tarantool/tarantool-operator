package topology

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/tarantool/tarantool-operator/pkg/topology/transport"
	v1 "k8s.io/api/core/v1"
)

// CommonCartridgeTopology .
type CommonCartridgeTopology struct {
	Transport transport.Transport
}

func (r *CommonCartridgeTopology) Exec(ctx context.Context, instance *v1.Pod, res any, lua string, args ...any) error {
	return r.Transport.Exec(ctx, instance, &res, lua, args...)
}

// Join pod to cluster.
func (r *CommonCartridgeTopology) Join(
	ctx context.Context,
	leader *v1.Pod,
	replicasetAlias string,
	replicasetUUID string,
	replicasetRoles []string,
	replicasetWeight int32,
	replicasetVshardGroup string,
	replicasetIsAllRw bool,
	advertiseURIs ...string,
) error {
	joinServers := make([]JoinServerParams, len(advertiseURIs))

	for i, uri := range advertiseURIs {
		joinServers[i].URI = uri
	}

	editTopology := EditTopologyParams{
		Replicasets: []EditReplicasetParams{
			{
				UUID:        replicasetUUID,
				Alias:       replicasetAlias,
				Roles:       replicasetRoles,
				Weight:      replicasetWeight,
				VShardGroup: replicasetVshardGroup,
				AllRw:       replicasetIsAllRw,
				JoinServers: joinServers,
			},
		},
	}

	success, err := r.adminEditTopology(ctx, leader, editTopology)
	if err != nil {
		var luaErr *LuaError
		if errors.As(err, &luaErr) {
			if strings.Contains(luaErr.Error(), "is already joined") {
				return ErrAlreadyJoined
			}

			if strings.Contains(luaErr.Error(), "can not enable unknown role") {
				return NewUnknownRoleError(luaErr)
			}
		}

		return err
	}

	if !success {
		return ErrTopologyIsDown
	}

	return nil
}

func (r *CommonCartridgeTopology) GetInstanceUUID(ctx context.Context, pod *v1.Pod) (string, error) {
	// language=lua
	lua := `
		if type(box.cfg) == 'function' then
			return ""
		end

		return box.info().uuid
	`

	var res string

	err := r.Exec(ctx, pod, &res, lua)
	if err != nil {
		return "", err
	}

	return res, nil
}

// SetWeight sets weight of a replicaset.
func (r *CommonCartridgeTopology) SetWeight(ctx context.Context, leader *v1.Pod, replicasetUUID string, replicaWeight int32) error {
	var res BooleanResult

	// language=lua
	lua := `
		local box = require('box')
		local cartridge = require('cartridge')
		local uuid, weight = ...

		local replicaset = cartridge.admin_get_replicasets(uuid)[1]
		local actualWight = 0

		if replicaset ~= nil and replicaset.weight ~= nil and actualWight ~= box.NULL then
			actualWight = replicaset.weight
		end

    	if replicaset == nil or actualWight ~= weight then
			local topology, err =  cartridge.admin_edit_topology({
				replicasets = {
					{ uuid = uuid , weight = weight}
				}
			})
			if err ~= nil then
				return { res=false, err=err}
			end
		end

		return { res=true ~= nil, err=nil}
	`

	err := r.Exec(ctx, leader, &res, lua, replicasetUUID, replicaWeight)
	if err != nil {
		return err
	}

	if res.Err != nil {
		if strings.Contains(res.Err.Error(), "not in config") {
			return ErrNotInConfig
		}

		if strings.Contains(res.Err.Error(), "At least one vshard-storage (default) must have weight > 0") {
			return ErrLastStorageWeight
		}

		return res.Err
	}

	return nil
}

// SetReplicasetRoles set roles list of replicaset in the Tarantool service.
func (r *CommonCartridgeTopology) SetReplicasetRoles(ctx context.Context, leader *v1.Pod, replicasetUUID string, roles []string) error {
	editTopology := EditTopologyParams{
		Replicasets: []EditReplicasetParams{
			{
				UUID:  replicasetUUID,
				Roles: roles,
			},
		},
	}

	success, err := r.adminEditTopology(ctx, leader, editTopology)
	if err != nil {
		return err
	}

	if !success {
		return ErrTopologyIsDown
	}

	return nil
}

// GetReplicasetRoles get roles list of replicaset from the Tarantool service.
func (r *CommonCartridgeTopology) GetReplicasetRoles(ctx context.Context, leader *v1.Pod, replicasetUUID string) ([]string, error) {
	var res []string

	// language=lua
	lua := `
		local args = ...
		local cartridge = require('cartridge')
		local res = cartridge.admin_get_replicasets(args.uuid)
		if res[1] ~= nil then
			return res[1].roles
		end
		return nil
	`

	err := r.Exec(ctx, leader, &res, lua,
		GetReplicasetsQuery{
			UUID: replicasetUUID,
		})
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve cartridge replicasets params")
	}

	if res == nil {
		return []string{}, nil
	}

	return res, nil
}

// BootstrapVshard enable the vshard service on the cluster.
func (r *CommonCartridgeTopology) BootstrapVshard(ctx context.Context, leader *v1.Pod) error {
	// language=lua
	lua := `
		local cartridge = require('cartridge')
		local res, err = cartridge.admin_bootstrap_vshard()
		return { res = res, err = err }
	`

	var res BooleanResult

	err := r.Exec(ctx, leader, &res, lua)
	if err != nil {
		return errors.Wrap(err, "unable to bootstrap cluster")
	}

	if res.Err != nil && !isAlreadyBootstrapped(res.Err) {
		return errors.Wrap(res.Err, "unable to bootstrap cluster")
	}

	return nil
}

func (r *CommonCartridgeTopology) GetRolesHierarchy(ctx context.Context, leader *v1.Pod) (map[string][]string, error) {
	// language=lua
	lua := `
		local roles = require('cartridge.roles')
		local hierarchy = {}
	
		for _, roleName in pairs(roles.get_all_roles()) do
			hierarchy[roleName] = roles.get_role_dependencies(roleName)
		end
	
		return hierarchy
	`

	var hierarchy map[string][]string

	err := r.Exec(ctx, leader, &hierarchy, lua)
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve roles hierarchy")
	}

	return hierarchy, nil
}

func (r *CommonCartridgeTopology) adminEditTopology(ctx context.Context, leader *v1.Pod, topologyObject EditTopologyParams) (bool, error) {
	// language=lua
	lua := `
		local args = ...
		local cartridge = require('cartridge')
		local topology, err = cartridge.admin_edit_topology(args)
		return { res = topology ~= nil, err=err }
	`

	var res BooleanResult

	err := r.Exec(ctx, leader, &res, lua, topologyObject)
	if err != nil {
		return false, errors.Wrap(err, "unable to edit topology")
	}

	if res.Err != nil {
		return res.Res, errors.Wrap(res.Err, "unable to edit topology")
	}

	return res.Res, nil
}

// SetFailoverParams configures cluster failover.
func (r *CommonCartridgeTopology) SetFailoverParams(ctx context.Context, leader *v1.Pod, params *FailoverParams) error {
	var res *BooleanResult

	// language=lua
	lua := `
		local cartridge = require('cartridge')
		local res, err = cartridge.failover_set_params(...)

		return { res = res, err = err }
	`

	err := r.Exec(ctx, leader, &res, lua, params)
	if err != nil {
		return err
	}

	if res.Err != nil {
		return res.Err
	}

	return nil
}

// GetFailoverParams retrieves cluster failover params.
func (r *CommonCartridgeTopology) GetFailoverParams(ctx context.Context, leader *v1.Pod) (*FailoverParams, error) {
	var res FailoverParams

	// language=lua
	lua := `
		local cartridge = require('cartridge')
		return cartridge.failover_get_params()
	`

	err := r.Exec(ctx, leader, &res, lua)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve cartridge failover params")
	}

	return &res, err
}

func (r *CommonCartridgeTopology) GetCartridgeConfig(ctx context.Context, leader *v1.Pod) (CartridgeConfigData, error) {
	// language=lua
	lua := `
		local cartridge = require('cartridge')
		local cfg = cartridge.config_get_readonly()
		local blacklist = {
			['auth'] = true,
			['auth.yml'] = true,
			['topology'] = true,
			['topology.yml'] = true,
			['users_acl'] = true,
			['users_acl.yml'] = true,
			['vshard'] = true,
			['vshard.yml'] = true,
			['vshard_groups'] = true,
			['vshard_groups.yml'] = true,
			['schema.yml'] = true,
		}
	
		local ret = {}
		for section, data in pairs(cfg) do
			if not blacklist[section] then
				ret[section] = data
			end
		end
		return ret
	`

	var (
		res     map[string]interface{}
		jsonErr *json.UnmarshalTypeError
	)

	err := r.Exec(ctx, leader, &res, lua)

	if errors.As(err, &jsonErr) {
		return CartridgeConfigData{}, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to download cartridge config")
	}

	return res, nil
}

func (r *CommonCartridgeTopology) ApplyCartridgeConfig(ctx context.Context, leader *v1.Pod, config CartridgeConfigData) error {
	// language=lua
	lua := `
	local cartridge = require('cartridge')
	local blacklist = {
		['auth'] = true,
		['auth.yml'] = true,
		['topology'] = true,
		['topology.yml'] = true,
		['users_acl'] = true,
		['users_acl.yml'] = true,
		['vshard'] = true,
		['vshard.yml'] = true,
		['vshard_groups'] = true,
		['vshard_groups.yml'] = true,
		['schema.yml'] = true,
	}

	local desiredConfig = ...
	local safeConfig = {}
	for key, value in pairs(desiredConfig) do
		if not blacklist[key] then
			safeConfig[key] = value
		end
	end

	return cartridge.config_patch_clusterwide(safeConfig)
	`

	var res bool

	err := r.Exec(ctx, leader, &res, lua, config)
	if err != nil {
		return errors.Wrap(err, "failed to upload cartridge config")
	}

	if !res {
		return fmt.Errorf("failed to upload cartridge config")
	}

	return nil
}

func (r *CommonCartridgeTopology) IsCartridgeStarted(ctx context.Context, pod *v1.Pod) (bool, error) {
	// language=lua
	lua := `
		local confapplier = require('cartridge.confapplier')
		local state = confapplier.get_state()

		if state == '' then
			return { res = false, err=nil }
		end

		if state == 'InitError' or state == 'BootError' or state == 'OperationError' or state == 'ReloadError' then
			return { res = false, err=nil }
		end

		return { res = true, err=nil }
	`

	var res *BooleanResult

	err := r.Exec(ctx, pod, &res, lua)
	if err != nil {
		return false, errors.Wrap(err, "unable to retrieve instance state")
	}

	if res.Err != nil {
		return res.Res, errors.Wrap(res.Err, "unable to retrieve instance state")
	}

	return res.Res, nil
}

func (r *CommonCartridgeTopology) IsCartridgeConfigured(ctx context.Context, pod *v1.Pod) (bool, error) {
	// language=lua
	lua := `
		local confapplier = require('cartridge.confapplier')
		local state = confapplier.get_state()

		if state ~= 'RolesConfigured' and state ~= 'OperationError' then
			return { res = false, err=nil }
		end

		return { res = true, err=nil }
	`

	var res *BooleanResult

	err := r.Exec(ctx, pod, &res, lua)
	if err != nil {
		return false, errors.Wrap(err, "unable to retrieve instance state")
	}

	if res.Err != nil {
		return res.Res, errors.Wrap(res.Err, "unable to retrieve instance state")
	}

	return res.Res, nil
}
