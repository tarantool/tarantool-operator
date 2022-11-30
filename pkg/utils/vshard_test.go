package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/tarantool/tarantool-operator/pkg/utils"
)

var _ = Describe("role utils unit testing", func() {
	Describe(`IsVShardRolesEquals function must compare two arrays of roles for equality, respecting the hierarchy and regardless of the order of the elements`,
		func() {
			Context("positive cases (equal arrays)", func() {
				DescribeTable(
					"should regard order if roles are same",
					func(a, b []string) {
						hierarchy := map[string][]string{
							"A": {},
							"B": {},
							"C": {},
						}
						Expect(utils.IsVShardRolesEquals(a, b, hierarchy)).Should(BeTrue())
						Expect(utils.IsVShardRolesEquals(b, a, hierarchy)).Should(BeTrue())
					},
					Entry("Empty and empty", []string{}, []string{}),
					Entry("A and A", []string{"A"}, []string{"A"}),
					Entry("A-B-C and A-B-C", []string{"A", "B", "C"}, []string{"A", "B", "C"}),
					Entry("A-B-C and A-C-B", []string{"A", "B", "C"}, []string{"A", "C", "B"}),
					Entry("A-B-C and B-A-C", []string{"A", "B", "C"}, []string{"B", "A", "C"}),
					Entry("A-B-C and B-C-A", []string{"A", "B", "C"}, []string{"B", "C", "A"}),
					Entry("A-B-C and C-A-B", []string{"A", "B", "C"}, []string{"C", "A", "B"}),
					Entry("A-B-C and C-B-A", []string{"A", "B", "C"}, []string{"C", "B", "A"}),
					Entry("A-C-B and A-B-C", []string{"A", "C", "B"}, []string{"A", "B", "C"}),
					Entry("A-C-B and A-C-B", []string{"A", "C", "B"}, []string{"A", "C", "B"}),
					Entry("A-C-B and B-A-C", []string{"A", "C", "B"}, []string{"B", "A", "C"}),
					Entry("A-C-B and B-C-A", []string{"A", "C", "B"}, []string{"B", "C", "A"}),
					Entry("A-C-B and C-A-B", []string{"A", "C", "B"}, []string{"C", "A", "B"}),
					Entry("A-C-B and C-B-A", []string{"A", "C", "B"}, []string{"C", "B", "A"}),
					Entry("B-A-C and A-B-C", []string{"B", "A", "C"}, []string{"A", "B", "C"}),
					Entry("B-A-C and A-C-B", []string{"B", "A", "C"}, []string{"A", "C", "B"}),
					Entry("B-A-C and B-A-C", []string{"B", "A", "C"}, []string{"B", "A", "C"}),
					Entry("B-A-C and B-C-A", []string{"B", "A", "C"}, []string{"B", "C", "A"}),
					Entry("B-A-C and C-A-B", []string{"B", "A", "C"}, []string{"C", "A", "B"}),
					Entry("B-A-C and C-B-A", []string{"B", "A", "C"}, []string{"C", "B", "A"}),
					Entry("B-C-A and A-B-C", []string{"B", "C", "A"}, []string{"A", "B", "C"}),
					Entry("B-C-A and A-C-B", []string{"B", "C", "A"}, []string{"A", "C", "B"}),
					Entry("B-C-A and B-A-C", []string{"B", "C", "A"}, []string{"B", "A", "C"}),
					Entry("B-C-A and B-C-A", []string{"B", "C", "A"}, []string{"B", "C", "A"}),
					Entry("B-C-A and C-A-B", []string{"B", "C", "A"}, []string{"C", "A", "B"}),
					Entry("B-C-A and C-B-A", []string{"B", "C", "A"}, []string{"C", "B", "A"}),
					Entry("C-A-B and A-B-C", []string{"C", "A", "B"}, []string{"A", "B", "C"}),
					Entry("C-A-B and A-C-B", []string{"C", "A", "B"}, []string{"A", "C", "B"}),
					Entry("C-A-B and B-A-C", []string{"C", "A", "B"}, []string{"B", "A", "C"}),
					Entry("C-A-B and B-C-A", []string{"C", "A", "B"}, []string{"B", "C", "A"}),
					Entry("C-A-B and C-A-B", []string{"C", "A", "B"}, []string{"C", "A", "B"}),
					Entry("C-A-B and C-B-A", []string{"C", "A", "B"}, []string{"C", "B", "A"}),
					Entry("C-B-A and A-B-C", []string{"C", "B", "A"}, []string{"A", "B", "C"}),
					Entry("C-B-A and A-C-B", []string{"C", "B", "A"}, []string{"A", "C", "B"}),
					Entry("C-B-A and B-A-C", []string{"C", "B", "A"}, []string{"B", "A", "C"}),
					Entry("C-B-A and B-C-A", []string{"C", "B", "A"}, []string{"B", "C", "A"}),
					Entry("C-B-A and C-A-B", []string{"C", "B", "A"}, []string{"C", "A", "B"}),
					Entry("C-B-A and C-B-A", []string{"C", "B", "A"}, []string{"C", "B", "A"}),
				)
				DescribeTable(
					"Should respect roles hierarchy, roles can be skipped in array if parent role in array",
					func(a, b []string) {
						hierarchy := map[string][]string{
							"A": {"A1"},
							"B": {"B1"},
							"C": {"C1", "C2"},
						}
						Expect(utils.IsVShardRolesEquals(a, b, hierarchy)).Should(BeTrue())
						Expect(utils.IsVShardRolesEquals(b, a, hierarchy)).Should(BeTrue())
					},
					Entry("Simple example A", []string{"A", "A1"}, []string{"A"}),
					Entry("Simple example B", []string{"B", "B1"}, []string{"B"}),
					Entry("Simple example C", []string{"C", "C1", "C2"}, []string{"C"}),
					Entry("A and B", []string{"A", "B"}, []string{"A", "A1", "B", "B1"}),
					Entry("B and C", []string{"B", "C", "C1", "C2"}, []string{"B", "B1", "C"}),
					Entry("A and B rand order", []string{"B", "A"}, []string{"B1", "A", "B", "A1"}),
				)
			})

			Context("negative cases (unequal arrays)", func() {
				DescribeTable("should return false array doe not match",
					func(a, b []string) {
						hierarchy := map[string][]string{
							"A": {},
							"B": {},
							"C": {},
						}

						Expect(utils.IsVShardRolesEquals(a, b, hierarchy)).Should(BeFalse())
						Expect(utils.IsVShardRolesEquals(b, a, hierarchy)).Should(BeFalse())
					},
					Entry("A and empty", []string{"A"}, []string{}),
					Entry("A and B", []string{"A"}, []string{"B"}),
					Entry("A-B and B", []string{"A", "B"}, []string{"B"}),
					Entry("A-B and C-A", []string{"A", "B"}, []string{"C", "A"}),
				)

				DescribeTable("all child roles != parent role",
					func(a, b []string) {
						hierarchy := map[string][]string{
							"A": {"A1"},
							"B": {"B1"},
							"C": {"C1", "C2"},
						}

						Expect(utils.IsVShardRolesEquals(a, b, hierarchy)).Should(BeFalse())
						Expect(utils.IsVShardRolesEquals(b, a, hierarchy)).Should(BeFalse())
					},
					Entry("A-B-B1 and B-B1-A1", []string{"B", "B1"}, []string{"B", "B1", "A"}),
					Entry("C1-C2 and C", []string{"C1", "C2"}, []string{"C"}),
				)
			})
		})
})
