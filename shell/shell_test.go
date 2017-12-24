package shell_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/markelog/eclectica/shell"
)

var _ = Describe("shell", func() {
	Describe("Compose", func() {
		It("puts main bin folder first", func() {
			plugins := []string{}
			result := strings.Split(Compose(plugins), ":")

			Expect(result[1]).To(ContainSubstring(".eclectica/bin"))
		})

		It("has only node if only node plugin was past", func() {
			plugins := []string{
				"node",
			}
			result := strings.Split(Compose(plugins), ":")

			Expect(result[2]).To(ContainSubstring(".eclectica/versions/node/current/bin"))
		})

		It("has only node if only node plugin was past", func() {
			plugins := []string{"node"}
			result := strings.Split(Compose(plugins), ":")

			Expect(result[2]).To(ContainSubstring(".eclectica/versions/node/current/bin"))

			for _, elem := range result {
				Expect(elem).NotTo(ContainSubstring(".eclectica/versions/rust/current/bin"))
			}
		})
	})
})
