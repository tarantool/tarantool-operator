package cli

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type TarantoolCTLResult struct {
	Ok  bool   `json:"ok"`
	Res string `json:"res"`
}

type TarantoolCTL struct{}

func (*TarantoolCTL) CreateCommand(lua string, args ...any) (*Command, error) {
	var safeLua string

	tpl := `
		local digest = require('digest')
		local json = require('json')
		local args = json.decode(%s)
		local func = function(...)
        	%s
		end
		local ok, result
		ok, result = pcall(func, unpack(args))
		return {ok=ok,res=digest.base64_encode(json.encode(result)),}	
	`

	if args != nil {
		jsonArgs, err := json.Marshal(args)
		if err != nil {
			return nil, err
		}

		quoted := strconv.Quote(string(jsonArgs))
		safeLua = fmt.Sprintf(tpl, quoted, lua)
	} else {
		safeLua = fmt.Sprintf(tpl, "'{}'", lua)
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
		resultRe   = regexp.MustCompile(`(?m)\n(\s*)res:`)
	)

	result = leadingRe.ReplaceAllString(result, "")
	result = trailingRe.ReplaceAllString(result, "")
	result = resultRe.ReplaceAllString(result, "\nres:")

	var encodedRes *TarantoolCTLResult

	err := yaml.Unmarshal([]byte(result), &encodedRes)
	if err != nil {
		return err
	}

	if encodedRes == nil {
		return errors.New("can't parse tarantool output")
	}

	bytesRes, err := base64.StdEncoding.DecodeString(encodedRes.Res)
	if err != nil {
		return err
	}

	if !encodedRes.Ok {
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
