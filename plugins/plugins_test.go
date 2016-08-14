package plugins_test

import (
  "os"
  "io/ioutil"
  "path/filepath"

  "github.com/jarcoal/httpmock"
  "github.com/bouk/monkey"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "github.com/markelog/eclectica/variables"
  ."github.com/markelog/eclectica/plugins"
)

func Read(path string) string {
  bytes, _ := ioutil.ReadFile(path)

  return string(bytes)
}

var _ = Describe("plugins", func() {
  var (
    name string
    path string
    version string
    archivePath string
    destFolder string
    versionsFolder string
    url string
    filename string
    info map[string]string
  )

  Describe("Activate", func() {
    It("should call Extract then Install methods", func() {
      result := ""

      monkey.Patch(Extract, func(info map[string]string) error {
        result += "Extract"
        return nil
      })

      monkey.Patch(Install, func(info map[string]string) error {
        result += "Install"
        return nil
      })

      Activate(map[string]string{})

      Expect(result).To(Equal("ExtractInstall"))

      monkey.Unpatch(Extract)
      monkey.Unpatch(Install)
    })
  })

  Describe("Extract", func() {
    BeforeEach(func() {
      path, _ = filepath.Abs("../testdata/plugins")
      versionsFolder, _ = filepath.Abs("../testdata/plugins/versions")
      name = "node"
      version = "1.0.0"
      filename = "node-arch"
      destFolder, _ = filepath.Abs("../testdata/plugins/versions/" + name + "/" + version)
      archivePath = path + "/" + filename + ".tar.gz"

      info = map[string]string{
        "name": name,
        "version": version,
        "archive-path": archivePath,
        "destination-folder": destFolder,
        "filename": filename,
      }

      monkey.Patch(variables.Home, func() string {
        return versionsFolder
      })
    })

    AfterEach(func() {
      monkey.Unpatch(variables.Home)
      os.RemoveAll(versionsFolder + "/" + name)
    })

    It("should extract langauge", func() {
      Extract(info)

      _, err := os.Stat(destFolder + "/test.txt");
      Expect(err).To(BeNil())
    })

    It("should extract even if previous archive was downloaded, but not extracted", func() {
      failedAttempt := versionsFolder + "/" + name + "/" + filename

      os.MkdirAll(failedAttempt, 0700)

      Extract(info)

      _, err := os.Stat(destFolder + "/test.txt");
      Expect(err).To(BeNil())

      _, err = os.Stat(failedAttempt);
      Expect(err).ShouldNot(BeNil())
    })
  })

  Describe("Download", func() {
    BeforeEach(func() {
      path, _ = filepath.Abs("../testdata/plugins")
      filename = "node-v6.3.0-darwin-x64.tar.gz"
      archivePath = path + "/" + filename
      destFolder = path + "/test"
      url = "https://example.com/" + filename

      info = map[string]string{
        "name": "node",
        "version": "6.3.0",
        "destination-folder": destFolder,
        "archive-folder": path,
        "archive-path": archivePath,
        "url": url,
      }
    })

    AfterEach(func() {
      defer httpmock.DeactivateAndReset()
      os.RemoveAll(archivePath)
      os.RemoveAll(destFolder)
    })

    Describe("200 response", func() {
      BeforeEach(func() {
        content := Read("../testdata/plugins/download.txt")

        httpmock.Activate()

        httpmock.RegisterResponder(
          "GET", url, httpmock.NewStringResponder(200, content),
        )
      })

      It("should download tar", func() {
        Download(info)

        _, err := os.Stat(archivePath);
        Expect(err).To(BeNil())
      })

      It("should not download anything if file already exist", func() {
        os.MkdirAll(destFolder, 0777)

        response, _ := Download(info)
        Expect(response).To(BeNil())
      })
    })

    Describe("404 response", func() {
      BeforeEach(func() {
        httpmock.Activate()

        httpmock.RegisterResponder(
          "GET", "", httpmock.NewStringResponder(404, ""),
        )
      })

      It("should return error", func() {
        _, err := Download(info)

        Expect(err).Should(MatchError("Incorrect version 6.3.0"))
      })
    })
  })

  Describe("ComposeVersions", func() {
    It("should compose versions", func() {
      compose := ComposeVersions([]string{"0.8.2", "4.4.7", "6.3.0", "6.4.2"})

      Expect(compose["0.x"]).To(Equal([]string{"0.8.2"}))
      Expect(compose["4.x"]).To(Equal([]string{"4.4.7"}))
      Expect(compose["6.x"]).To(Equal([]string{"6.3.0", "6.4.2"}))
    })
  })

  Describe("GetKeys", func() {
    It("should get version keys", func() {
      list := map[string][]string{"4.x": []string{}, "0.x": []string{"0.8.2"}}
      keys := GetKeys(list)

      Expect(keys[0]).To(Equal("0.x"))
      Expect(keys[1]).To(Equal("4.x"))
    })
  })

  Describe("GetElements", func() {
    It("should get version elements", func() {
      list := ComposeVersions([]string{"0.8.2", "4.4.7", "6.3.0", "6.4.2"})
      elements := GetElements("6.x", list)

      Expect(elements[0]).To(Equal("6.3.0"))
      Expect(elements[1]).To(Equal("6.4.2"))
    })
  })
})
