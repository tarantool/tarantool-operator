package cli_test

import (
	"strings"
	"testing"

	"github.com/tarantool/tarantool-operator/pkg/topology/transport/podexec/cli"
	"github.com/tarantool/tarantool-operator/pkg/utils"
)

type parseStringTestCase struct {
	input          string
	expectedResult interface{}
	err            string
}

func TestParseTarantoolOutputParseRoles(t *testing.T) {
	tarantoolctl := &cli.TarantoolCTL{}

	cases := []parseStringTestCase{
		{
			input: `---
- ok: true
  hex: 5b226661696c6f7665722d636f6f7264696e61746f72222c227673686172642d726f75746572222c226170702e726f6c65732e726f75746572225d
...`,
			expectedResult: []string{"failover-coordinator", "vshard-router", "app.roles.router"},
			err:            "",
		},
		{
			input: `---
- ok: true
  hex: 5B226661696C6F7665722D636F6F7264696E61746F72225D
...`,
			expectedResult: []string{"failover-coordinator"},
			err:            "",
		},
		{
			input: `---
- ok: true
  hex: 5B5D
...`,
			expectedResult: []string{},
			err:            "",
		},
		{
			input: `---
- ok: true
  hex: 226572726F7222
...`,
			expectedResult: []string{},
			err:            "json: cannot unmarshal string into Go value of type []string",
		},
	}
	for i, c := range cases {
		var result []string
		err := tarantoolctl.Unmarshal(c.input, &result)

		if (err != nil || c.err != "") && !strings.Contains(err.Error(), c.err) {
			t.Fatalf("%d: expected %s err, got: %s", i, c.err, err)
		}

		expectedResult, ok := c.expectedResult.([]string)
		if !ok {
			t.Fatalf("Wrong expectedResult type")
		}

		for _, v := range result {
			if utils.SliceContains(expectedResult, v) == false {
				t.Fatalf("%d: roles must not contains %s", i, v)
			}
		}

		for _, v := range expectedResult {
			if utils.SliceContains(result, v) == false {
				t.Fatalf("%d: roles must contains %s", i, v)
			}
		}
	}
}
