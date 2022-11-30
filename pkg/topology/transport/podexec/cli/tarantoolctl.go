package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type TarantoolCTLResult struct {
	Ok  bool   `json:"ok"`
	Hex string `json:"hex"`
}

type TarantoolCTL struct{}

func (*TarantoolCTL) CreateCommand(lua string, args ...any) (*Command, error) {
	var safeLua string

	tpl := `
		local json = require('json')
		local args = json.decode('%s')
		local func = function(...)
        	%s
		end
		local ok, result
		ok, result = pcall(func, unpack(args))
		return {ok=ok,hex=string.hex(json.encode(result)),}	
	`

	if args != nil {
		jsonArgs, err := json.Marshal(args)
		if err != nil {
			return nil, err
		}

		safeLua = fmt.Sprintf(tpl, jsonArgs, lua)
	} else {
		safeLua = fmt.Sprintf(tpl, "{}", lua)
	}

	safeLua = strings.ReplaceAll(safeLua, "\t", "")

	return &Command{
		Command: []string{
			"sh",
			"-c",
			"cat /dev/stdin | tarantoolctl connect `ls $CARTRIDGE_RUN_DIR/*.control` ",
		},
		StdIn: safeLua,
	}, nil
}

func (*TarantoolCTL) Unmarshal(result string, target any) error {
	var (
		leadingRe  = regexp.MustCompile(`(?m)^---\n-\s+?`)
		trailingRe = regexp.MustCompile(`(?m)\n.{3}\n?$`)
		resultRe   = regexp.MustCompile(`(?m)\n(\s*)hex:`)
	)

	result = leadingRe.ReplaceAllString(result, "")
	result = trailingRe.ReplaceAllString(result, "")
	result = resultRe.ReplaceAllString(result, "\nhex:")

	var hexRes *TarantoolCTLResult

	err := yaml.Unmarshal([]byte(result), &hexRes)
	if err != nil {
		return err
	}

	if hexRes == nil {
		return errors.New("can't parse tarantool output")
	}

	bytesRes, err := hex.DecodeString(hexRes.Hex)
	if err != nil {
		return err
	}

	if !hexRes.Ok {
		var errMsg string

		err = json.Unmarshal(bytesRes, &errMsg)
		if err != nil {
			return err
		}

		return errors.New(errMsg)
	}

	err = json.Unmarshal(bytesRes, &target)
	if err != nil {
		return err
	}

	return nil
}
