package io_test

import (
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/plugins"
)

var _ = Describe("io", func() {
	Describe("FindDotFile", func() {
		It("Should find .nvmrc file for nodejs", func() {
			dots := plugins.New(&plugins.Args{
				Language: "node",
			}).Dots()
			path, _ := filepath.Abs("../testdata/io/node-with-nvm/")
			result, _ := FindDotFile(dots, path)

			Expect(strings.Contains(result, ".nvmrc")).To(Equal(true))
		})
	})

	Describe("GetVersion", func() {
		It("Should get version for node from .nvmrc file", func() {
			dots := plugins.New(&plugins.Args{
				Language: "node",
			}).Dots()
			path, _ := filepath.Abs("../testdata/io/node-with-nvm/")
			result, dotPath, _ := GetVersion(dots, path)

			Expect(dotPath).To(ContainSubstring("io/node-with-nvm/.nvmrc"))
			Expect(result).To(Equal("6.8.0"))
		})
	})
})
