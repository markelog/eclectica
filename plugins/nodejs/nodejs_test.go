package nodejs_test

import (
  "regexp"

  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  ."github.com/markelog/eclectica/plugins/nodejs"
)

var _ = Describe("nodejs", func() {
  var (
    remotes []string
    err error
  )

  Describe("ListVersions", func() {
    BeforeEach(func() {
      remotes, err = ListVersions()
    })

    It("should not return an error", func() {
      Expect(err).To(BeNil())
    })

    It("should have correct version values", func() {
      rp := regexp.MustCompile("[[:digit:]]+\\.[[:digit:]]+\\.[[:digit:]]+$")

      for _, element := range remotes {
        Expect(rp.MatchString(element)).To(Equal(true))
      }
    })
  })
})
