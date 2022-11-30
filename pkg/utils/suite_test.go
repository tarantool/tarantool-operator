package utils_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestControllerUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pkg utils suite")
}
