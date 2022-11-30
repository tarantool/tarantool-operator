package topology

import (
	"fmt"
)

type LuaCallResult[TRes any] struct {
	Res TRes      `json:"res,omitempty"`
	Err *LuaError `json:"err,omitempty"`
}

type LuaError struct {
	Line      int64  `json:"line"`
	ClassName string `json:"class_name"`
	Err       string `json:"err"`
	File      string `json:"file"`
	Stack     string `json:"stack"`
}

func (r *LuaError) Error() string {
	return fmt.Sprintf("%s: %s", r.ClassName, r.Err)
}

type (
	BooleanResult = LuaCallResult[bool]
	Int64Result   = LuaCallResult[int64]
)

type CartridgeConfigData = map[string]interface{}
