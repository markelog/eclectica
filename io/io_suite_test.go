package io_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRequest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IO Suite")
}
