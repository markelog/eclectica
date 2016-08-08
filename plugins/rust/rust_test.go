package rust_test

import (
  "regexp"
  "io/ioutil"

  "github.com/jarcoal/httpmock"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  ."github.com/markelog/eclectica/plugins/rust"
)

func Read(path string) string {
  bytes, _ := ioutil.ReadFile(path)

  return string(bytes)
}

var _ = Describe("rust", func() {
  BeforeEach(func() {
    content := Read("../../testdata/rust/dist.txt")

    httpmock.Activate()

    httpmock.RegisterResponder(
      "GET",
      "https://static.rust-lang.org/dist/index.txt",
      httpmock.NewStringResponder(200, content),
    )
  })

  AfterEach(func() {
    defer httpmock.DeactivateAndReset()
  })

  Describe("ListVersions", func() {
    var (
      remotes []string
      err error
    )

    BeforeEach(func() {
      remotes, err = ListVersions()
    })

    It("should not return an error", func() {
      Expect(err).To(BeNil())
    })

    It("gets list of versions", func() {
      Expect(remotes[0]).To(Equal("0.10"))
      Expect(remotes[2]).To(Equal("0.12.0"))
      Expect(remotes[8]).To(Equal("1.0.0-beta.4"))
    })

    It("should have correct version values", func() {
      rp := regexp.MustCompile("[[:digit:]]+\\.[[:digit:]]+")

      for _, element := range remotes {
        Expect(rp.MatchString(element)).To(Equal(true))
      }
    })
  })
})
