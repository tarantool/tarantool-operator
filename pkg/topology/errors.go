package topology

import (
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrTopologyIsDown    = errors.New("topology service is down")
	ErrAlreadyJoined     = errors.New("already joined")
	ErrNotInConfig       = errors.New("not in config")
	ErrLastStorageWeight = errors.New("at least one vshard-storage (default) must have weight > 0")
)

type UnknownRoleError struct {
	*LuaError
}

func NewUnknownRoleError(err *LuaError) *UnknownRoleError {
	return &UnknownRoleError{
		LuaError: err,
	}
}

func isAlreadyBootstrapped(err *LuaError) bool {
	return err.ClassName == "Bootstrapping vshard failed" &&
		strings.Contains(err.Err, "already bootstrapped")
}
