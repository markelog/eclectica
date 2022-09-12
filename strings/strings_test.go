package strings_test

import (
	"github.com/markelog/monkey"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"golang.org/x/crypto/ssh/terminal"

	. "github.com/markelog/eclectica/strings"
)

var _ = Describe("Strings", func() {
	Describe("Elipsis", func() {
		It("Should truncate the string", func() {
			result := Elipsis("test", 1)
			Expect(result).To(Equal("t..."))
		})

		It("Should not truncate the string if number is just weird", func() {
			result := Elipsis("test", -1)
			Expect(result).To(Equal("test"))
		})

		It("Should not truncate the string if number is longer then string", func() {
			println(Elipsis("test", 50))
			result := Elipsis("test", 50)
			Expect(result).To(Equal("test"))
		})
	})

	Describe("ElipsisForTerminal", func() {
		BeforeEach(func() {
			monkey.Patch(terminal.GetSize, func(num int) (int, int, error) {
				return 51, 1, nil
			})
		})

		AfterEach(func() {
			monkey.Unpatch(terminal.GetSize)
		})

		It("Should truncate the string", func() {
			result := ElipsisForTerminal("test")
			Expect(result).To(Equal("t..."))
		})
	})
})
