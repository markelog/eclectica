package info_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/markelog/eclectica/cmd/info"
)

var _ = Describe("info", func() {
	Describe("GetFullVersion", func() {
		It("support for partial major version", func() {
			version := "6"
			versions := []string{
				"6.1.0", "5.2.0", "6.2.0", "6.8.3",
			}

			test, _ := info.GetFullVersion(version, versions)

			Expect(test).To(Equal("6.8.3"))
		})

		It("support for partial minor version", func() {
			version := "6.4"
			versions := []string{
				"6.1.0", "5.2.0", "6.2.0", "6.8.3", "6.4.2", "6.4.0",
			}

			test, _ := info.GetFullVersion(version, versions)

			Expect(test).To(Equal("6.4.2"))
		})

		It("shouldn't do anything for full version", func() {
			version := "6.1.1"
			versions := []string{
				"6.1.0", "5.2.0", "6.2.0", "6.8.3", "6.4.2", "6.4.0",
			}

			test, err := info.GetFullVersion(version, versions)

			Expect(err).To(BeNil())
			Expect(test).To(Equal("6.1.1"))
		})
	})

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
