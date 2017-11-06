package rvm_test

import (
	"github.com/markelog/release"

	"github.com/bouk/monkey"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/markelog/eclectica/plugins/ruby/rvm"
)

var _ = Describe("RVM related methods", func() {
	Describe("GetUrl", func() {
		It("should return same resulted url part", func() {
			monkey.Patch(release.All, func() (string, string, string) {
				return "osx", "x86_64", "10.12"
			})
			result := rvm.GetUrl("test")
			monkey.Unpatch(release.All)

			Expect(result).To(Equal("test/osx/10.12/x86_64"))
		})

		It("should return lower then returned version", func() {
			monkey.Patch(release.All, func() (string, string, string) {
				return "osx", "x86_64", "10.13"
			})
			result := rvm.GetUrl("test")
			monkey.Unpatch(release.All)

			Expect(result).To(Equal("test/osx/10.12/x86_64"))
		})
	})
})
