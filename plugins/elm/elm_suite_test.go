package elm_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestElm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Elm Suite")
}
