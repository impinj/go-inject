package inject_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestInject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Inject Suite")
}
