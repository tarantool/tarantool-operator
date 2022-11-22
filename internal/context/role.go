package context

import (
	"github.com/pkg/errors"
	"github.com/tarantool/tarantool-operator/apis/v1alpha2"
	"github.com/tarantool/tarantool-operator/pkg/reconciliation"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RoleContext struct {
	*reconciliation.CommonContext

	Role *v1alpha2.Role
}

func (r *RoleContext) SetRole(role *v1alpha2.Role) {
	r.Role = role
}

func (r *RoleContext) GetRole() *v1alpha2.Role {
	return r.Role
}

func (r *RoleContext) HasRequestedObject() bool {
	return r.GetRole() != nil
}

func (r *RoleContext) SetRequestedObject(obj client.Object) error {
	role, ok := obj.(*v1alpha2.Role)
	if !ok {
		return errors.New("RoleContext used with wrong k8s object")
	}

	r.SetRole(role)

	return nil
}

func (r *RoleContext) GetRequestedObject() client.Object {
	return r.GetRole()
}
