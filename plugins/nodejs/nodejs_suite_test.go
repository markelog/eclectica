package nodejs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCprf(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nodejs Suite")
}
