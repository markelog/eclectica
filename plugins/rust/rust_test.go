package rust_test

import (
  "regexp"
  "io/ioutil"
  "runtime"

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
  Describe("ListVersions", func() {
    var (
      remotes []string
      err error
    )

    Describe("fail", func() {
      BeforeEach(func() {
        httpmock.Activate()
      })

      AfterEach(func() {
        defer httpmock.DeactivateAndReset()
      })

      It("should return an error", func() {
        httpmock.RegisterResponder(
          "GET",
          "https://static.rust-lang.org/dist/index.txt",
          httpmock.NewStringResponder(500, ""),
        )

        remotes, err = ListVersions()

        Expect(err).Should(MatchError("Can't establish connection"))
      })
    })

    Describe("success", func() {
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

      // :/
      if runtime.GOOS == "darwin" {
        Expect(result["filename"]).To(Equal("node-v6.3.1-darwin-x64"))
        Expect(result["url"]).To(Equal("https://nodejs.org/dist/latest/node-v6.3.1-darwin-x64.tar.gz"))
      } else if runtime.GOOS == "linux" {
        Expect(result["filename"]).To(Equal("node-v6.3.1-linux-x64"))
        Expect(result["url"]).To(Equal("https://nodejs.org/dist/latest/node-v6.3.1-linux-x64.tar.gz"))
      }
    })

    It("should get info about lts version", func() {
      result, _ := Version("lts")

      Expect(result["name"]).To(Equal("node"))
      Expect(result["version"]).To(Equal("6.3.1"))

      // :/
      if runtime.GOOS == "darwin" {
        Expect(result["filename"]).To(Equal("node-v6.3.1-darwin-x64"))
        Expect(result["url"]).To(Equal("https://nodejs.org/dist/lts/node-v6.3.1-darwin-x64.tar.gz"))
      } else if runtime.GOOS == "linux" {
        Expect(result["filename"]).To(Equal("node-v6.3.1-linux-x64"))
        Expect(result["url"]).To(Equal("https://nodejs.org/dist/lts/node-v6.3.1-linux-x64.tar.gz"))
      }
    })

    FIt("should get info about 1.9.0 version", func() {
      result, _ := Version("1.9.0")

      Expect(result["name"]).To(Equal("rust"))
      Expect(result["version"]).To(Equal("1.9.0"))

      // :/
      if runtime.GOOS == "darwin" {
        Expect(result["filename"]).To(Equal("rust-1.9.0-x86_64-apple-darwin"))
        Expect(result["url"]).To(Equal("https://static.rust-lang.org/dist/rust-1.9.0-x86_64-apple-darwin.tar.gz"))
      } else if runtime.GOOS == "linux" {
        Expect(result["filename"]).To(Equal("rust-1.9.0-x86_64-unknown-linux-gnu"))
        Expect(result["url"]).To(Equal("https://static.rust-lang.org/dist/rust-1.9.0-x86_64-unknown-linux-gnu.tar.gz"))
      }
    })
  })
})
