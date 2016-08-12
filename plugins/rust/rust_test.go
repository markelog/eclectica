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
    AfterEach(func() {
      defer httpmock.DeactivateAndReset()
    })

    It("should get info about nightly version", func() {
      result, _ := Version("nightly")

      Expect(result["name"]).To(Equal("rust"))
      Expect(result["version"]).To(Equal("nightly"))

      // :/
      if runtime.GOOS == "darwin" {
        Expect(result["filename"]).To(Equal("rust-nightly-x86_64-apple-darwin"))
        Expect(result["url"]).To(Equal("https://static.rust-lang.org/dist/rust-nightly-x86_64-apple-darwin.tar.gz"))
      } else if runtime.GOOS == "linux" {
        Expect(result["filename"]).To(Equal("rust-nightly-x86_64-unknown-linux-gnu"))
        Expect(result["url"]).To(Equal("https://static.rust-lang.org/dist/rust-nightly-x86_64-unknown-linux-gnu.tar.gz"))
      }
    })

    It("should get info about lts version", func() {
      result, _ := Version("beta")

      Expect(result["name"]).To(Equal("rust"))
      Expect(result["version"]).To(Equal("beta"))

      // :/
      if runtime.GOOS == "darwin" {
        Expect(result["filename"]).To(Equal("rust-beta-x86_64-apple-darwin"))
        Expect(result["url"]).To(Equal("https://static.rust-lang.org/dist/rust-beta-x86_64-apple-darwin.tar.gz"))
      } else if runtime.GOOS == "linux" {
        Expect(result["filename"]).To(Equal("rust-beta-x86_64-unknown-linux-gnu"))
        Expect(result["url"]).To(Equal("https://static.rust-lang.org/dist/rust-beta-x86_64-unknown-linux-gnu.tar.gz"))
      }
    })

    It("should get info about 1.9.0 version", func() {
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
