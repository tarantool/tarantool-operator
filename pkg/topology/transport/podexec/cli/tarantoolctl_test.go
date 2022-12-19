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
  res: WyJmYWlsb3Zlci1jb29yZGluYXRvciIsInZzaGFyZC1yb3V0ZXIiLCJhcHAucm9sZXMucm91dGVyIl0=
...`,
			expectedResult: []string{"failover-coordinator", "vshard-router", "app.roles.router"},
			err:            "",
		},
		{
			input: `---
- ok: true
  res: WyJmYWlsb3Zlci1jb29yZGluYXRvciJd
...`,
			expectedResult: []string{"failover-coordinator"},
			err:            "",
		},
		{
			input: `---
- ok: true
  res: W10=
...`,
			expectedResult: []string{},
			err:            "",
		},
		{
			input: `---
- ok: true
  res: ImVycm9yIg==
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
