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

	Describe("PossibleLanguage", func() {
		It("should get language", func() {
			language := info.PossibleLanguage([]string{"-r", "rust"})

			Expect(language).To(Equal("rust"))
		})

		It("should get language in different sequence", func() {
			language := info.PossibleLanguage([]string{"rust", "-r"})

			Expect(language).To(Equal("rust"))
		})

		It("should get non-existing language with `-r`", func() {
			language := info.PossibleLanguage([]string{"-r", "boom@1.2.3"})

			Expect(language).To(Equal("boom"))
		})
	})
})
