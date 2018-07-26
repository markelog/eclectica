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
	Describe("ExtractVersion", func() {
		It("gets version with v1.2.3 format", func() {
			result, err := ExtractVersion("v1.2.3")

			Expect(err).To(BeNil())
			Expect(result).To(Equal("1.2.3"))
		})

		It("gets version with 1.2.3 format", func() {
			result, err := ExtractVersion("1.2.3")

			Expect(err).To(BeNil())
			Expect(result).To(Equal("1.2.3"))
		})

		It("gets version with 1.2 format", func() {
			result, err := ExtractVersion("1.2")

			Expect(err).To(BeNil())
			Expect(result).To(Equal("1.2.0"))
		})

		It("gets version with one digit format", func() {
			result, err := ExtractVersion("1")

			Expect(err).To(BeNil())
			Expect(result).To(Equal("1.0.0"))
		})

		It("returns an error if there is no version", func() {
			result, err := ExtractVersion("test")

			Expect(err).To(HaveOccurred())
			Expect(result).To(Equal(""))
		})

		It("gets version with v8.11.2 format", func() {
			result, err := ExtractVersion("v8.11.2")

			Expect(err).To(BeNil())
			Expect(result).To(Equal("8.11.2"))
		})
	})

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

		It("Should get version for node from .nvmrc file with v-string", func() {
			dots := plugins.New(&plugins.Args{
				Language: "node",
			}).Dots()
			path, _ := filepath.Abs("../testdata/io/node-with-nvm-v-string/")
			result, dotPath, _ := GetVersion(dots, path)

			Expect(dotPath).To(ContainSubstring("io/node-with-nvm-v-string/.nvmrc"))
			Expect(result).To(Equal("6.8.0"))
		})
	})
})
