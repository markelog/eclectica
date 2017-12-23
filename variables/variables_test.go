package variables_test

import (
	"os"
	"os/user"

	"github.com/bouk/monkey"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/markelog/eclectica/variables"
)

var _ = Describe("variables", func() {
	Describe("Base", func() {
		BeforeEach(func() {
			monkey.Patch(user.Current, func() (*user.User, error) {
				user := &user.User{
					Name:    "root",
					HomeDir: "/test",
				}
				return user, nil
			})

			monkey.Patch(os.Getenv, func(name string) string {
				return "test"
			})
		})

		AfterEach(func() {
			monkey.Unpatch(user.Current)
			monkey.Unpatch(os.Getenv)
		})

		It("should return correct path to eclectica meta folder with sudo use", func() {
			result := variables.Base()

			Expect(result).To(Equal("/test/.eclectica"))
		})
	})
})
