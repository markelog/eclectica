package nodejs_test

import (
  "regexp"
  "io/ioutil"

  "github.com/jarcoal/httpmock"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  ."github.com/markelog/eclectica/plugins/nodejs"
)

func Read(path string) string {
  bytes, _ := ioutil.ReadFile(path)

  return string(bytes)
}

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

  Describe("Keyword", func() {
    var (
      result map[string]string
      err error
    )

    BeforeEach(func() {
      content := Read("../../testdata/nodejs/latest.txt")

      httpmock.Activate()

      httpmock.RegisterResponder(
        "GET",
        "https://nodejs.org/dist/latest/SHASUMS256.txt",
        httpmock.NewStringResponder(200, content),
      )
    })

    AfterEach(func() {
      defer httpmock.DeactivateAndReset()
    })

    BeforeEach(func() {
      result, err = Keyword("latest")
    })

    It("should not return an error", func() {
      Expect(err).To(BeNil())
    })

    It("should get info about latest version", func() {
      Expect(result["name"]).To(Equal("node"))
      Expect(result["version"]).To(Equal("6.3.1"))
      Expect(result["filename"]).To(Equal("node-v6.3.1-darwin-x64"))
      Expect(result["url"]).To(Equal("https://nodejs.org/dist/latest/node-v6.3.1-darwin-x64.tar.gz"))
    })
  })

  Describe("Version", func() {
    BeforeEach(func() {
      content := Read("../../testdata/nodejs/latest.txt")

      httpmock.Activate()

      httpmock.RegisterResponder(
        "GET",
        "https://nodejs.org/dist/latest/SHASUMS256.txt",
        httpmock.NewStringResponder(200, content),
      )

      httpmock.RegisterResponder(
        "GET",
        "https://nodejs.org/dist/lts/SHASUMS256.txt",
        httpmock.NewStringResponder(200, content),
      )
    })

    AfterEach(func() {
      defer httpmock.DeactivateAndReset()
    })

    It("should get info about latest version", func() {
      result, _ := Version("latest")

      Expect(result["name"]).To(Equal("node"))
      Expect(result["version"]).To(Equal("6.3.1"))
      Expect(result["filename"]).To(Equal("node-v6.3.1-darwin-x64"))
      Expect(result["url"]).To(Equal("https://nodejs.org/dist/latest/node-v6.3.1-darwin-x64.tar.gz"))
    })

    It("should get info about lts version", func() {
      result, _ := Version("lts")

      Expect(result["name"]).To(Equal("node"))
      Expect(result["version"]).To(Equal("6.3.1"))
      Expect(result["filename"]).To(Equal("node-v6.3.1-darwin-x64"))
      Expect(result["url"]).To(Equal("https://nodejs.org/dist/lts/node-v6.3.1-darwin-x64.tar.gz"))
    })

    It("should get info about 6.3.1 version", func() {
      result, _ := Version("6.3.1")

      Expect(result["name"]).To(Equal("node"))
      Expect(result["version"]).To(Equal("6.3.1"))
      Expect(result["filename"]).To(Equal("node-v6.3.1-darwin-x64"))
      Expect(result["url"]).To(Equal("https://nodejs.org/dist/v6.3.1/node-v6.3.1-darwin-x64.tar.gz"))
    })
  })
})
