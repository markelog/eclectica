package variables_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestEclectica(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Eclectica Suite")
}
