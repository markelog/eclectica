package variables_test

import (
	"errors"
	"os"

	"github.com/bouk/monkey"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/markelog/eclectica/variables"
)

var _ = Describe("variables", func() {
	Describe("Index", func() {
		It("doesn't provide index", func() {
			exec := Index(DefaultBins, "/home/markelog/bin")

			Expect(exec).To(Equal(-1))
		})

		It("doesn provide index", func() {
			exec := Index(DefaultBins, "/usr/bin")

			Expect(exec).To(Equal(0))
		})
	})

	Describe("InLocalBin", func() {
		Describe("local path", func() {
			var firstPath string
			var secondPath string

			BeforeEach(func() {
				monkey.Patch(os.Stat, func(path string) (os.FileInfo, error) {
					if firstPath == path {
						return nil, nil
					}

					if secondPath == path {
						return nil, nil
					}

					return nil, errors.New("test")
				})
			})

			AfterEach(func() {
				firstPath = ""
				secondPath = ""
				monkey.Unpatch(os.Stat)
			})

			It("should be in local path since $PATH doesn't contain it", func() {
				firstPath = "/first/go"
				secondPath = "/second/go"

				result := InLocalBin("/first:/second", "/third", "go")

				Expect(result).To(Equal(false))
			})

			It("should be in local path since it has lower index", func() {
				firstPath = "/first/go"
				secondPath = "/second/go"

				result := InLocalBin("/first:/second", "/first", "go")

				Expect(result).To(Equal(true))
			})

			It("shouldn't be in local path since it has higher index", func() {
				firstPath = "/first/go"
				secondPath = "/second/go"

				result := InLocalBin("/first:/second", "/second", "go")

				Expect(result).To(Equal(false))
			})
		})
	})
})
