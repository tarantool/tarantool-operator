package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/tarantool/tarantool-operator/pkg/utils"
)

var _ = Describe("cmp utils unit testing", func() {
	DescribeTable(
		"should return true if map is subset of another map",
		func(set, subset map[string]any) {
			Expect(utils.IsMapSubset(set, subset)).Should(BeTrue())
		},
		Entry("Empty and empty", map[string]any{}, map[string]any{}),
		Entry("Not empty empty", map[string]any{
			"key":     map[string]any{"val": 1},
			"key.yml": "{\"val\": 1}",
		}, map[string]any{}),
		Entry("Not empty and not empty", map[string]any{
			"key":     map[string]any{"val": 1},
			"key.yml": "{\"val\": 1}",
		}, map[string]any{
			"key": map[string]any{"val": 1},
		}),
	)

	DescribeTable(
		"should return false if map is NOT subset of another map",
		func(set, subset map[string]any) {
			Expect(utils.IsMapSubset(set, subset)).Should(BeFalse())
		},
		Entry("contains changed ket", map[string]any{
			"key":     map[string]any{"val": 1},
			"key.yml": "{\"val\": 1}",
		}, map[string]any{
			"new_key": map[string]any{"val": 1},
		}),
		Entry("contains changed value", map[string]any{
			"key":     map[string]any{"val": 1},
			"key.yml": "{\"val\": 1}",
		}, map[string]any{
			"key": map[string]any{"val": 2},
		}),

		Entry("empty contains new entry", map[string]any{
			"key":     map[string]any{"val": 1},
			"key.yml": "{\"val\": 1}",
		}, map[string]any{
			"key":     map[string]any{"val": 1},
			"new_key": map[string]any{"val": 2},
		}),
	)
})
