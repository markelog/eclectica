package rust_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRust(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rust Suite")
}
