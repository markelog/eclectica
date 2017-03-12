package info_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/markelog/eclectica/cmd/info"
)

var _ = Describe("info", func() {
	Describe("GetLanguage", func() {
		It("should get language", func() {
			language, version := info.GetLanguage([]string{"-r", "rust"})

			Expect(language).To(Equal("rust"))
			Expect(version).To(Equal(""))
		})

		It("should get language in different sequence", func() {
			language, version := info.GetLanguage([]string{"rust", "-r"})

			Expect(language).To(Equal("rust"))
			Expect(version).To(Equal(""))
		})

		It("should not get non-existing language with `-r`", func() {
			language, _ := info.GetLanguage([]string{"-r", "rustc@1.2.3"})

			Expect(language).To(Equal(""))
		})

		It("should not get non-existing language", func() {
			language, _ := info.GetLanguage([]string{"rustc@1.2.3"})

			Expect(language).To(Equal(""))
		})

		It("should get language without additional data", func() {
			language, version := info.GetLanguage([]string{"rust@1.2.3"})

			Expect(language).To(Equal("rust"))
			Expect(version).To(Equal("1.2.3"))
		})

		It("should not get non-existing language without version number", func() {
			language, _ := info.GetLanguage([]string{"rustc"})

			Expect(language).To(Equal(""))
		})
	})

	Describe("NonInstallCommand", func() {
		It("should detect command", func() {
			hasCommand := info.NonInstallCommand([]string{"-r", "ls"})

			Expect(hasCommand).To(Equal(true))
		})

		It("should detect command in different sequence", func() {
			hasCommand := info.NonInstallCommand([]string{"ls", "-r"})

			Expect(hasCommand).To(Equal(true))
		})

		It("should not detect command", func() {
			hasCommand := info.NonInstallCommand([]string{"-r", "rustc@1.2.3"})

			Expect(hasCommand).To(Equal(false))
		})
	})
})
