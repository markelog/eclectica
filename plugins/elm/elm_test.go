package elm_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	eIO "github.com/markelog/eclectica/io"
	. "github.com/markelog/eclectica/plugins/elm"
	"github.com/markelog/eclectica/variables"
)

var _ = Describe("elm", func() {
	var (
		remotes []string
		err     error
	)

	elm := &Elm{}

	Describe("ListRemote", func() {
		old := VersionLink

		AfterEach(func() {
			VersionLink = old
		})

		Describe("success", func() {
			BeforeEach(func() {
				content := eIO.Read("../../testdata/plugins/elm/elm-platform.html")

				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					status := 200

					if _, ok := r.URL.Query()["status"]; ok {
						fmt.Sscanf(r.URL.Query().Get("status"), "%d", &status)
					}

					w.WriteHeader(status)
					io.WriteString(w, content)
				}))

				VersionLink = ts.URL

				remotes, err = elm.ListRemote()
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should have correct version values", func() {
				Expect(remotes[0]).To(Equal("0.15.1"))
			})
		})

		Describe("fail", func() {
			BeforeEach(func() {
				VersionLink = ""
				remotes, err = elm.ListRemote()
			})

			It("should return an error", func() {
				Expect(err).Should(MatchError(variables.ConnectionError))
			})
		})
	})

	Describe("Info", func() {
		It("should get info about 0.18.0 version", func() {
			result := (&Elm{Version: "0.18.0"}).Info()

			Expect(result["archive-folder"]).Should(ContainSubstring("elm-archive-0.18.0/"))

			// :/
			if runtime.GOOS == "darwin" {
				Expect(result["unarchive-filename"]).To(Equal(""))
				Expect(result["filename"]).To(Equal("darwin-x64"))
				Expect(result["url"]).To(Equal("https://dl.bintray.com/elmlang/elm-platform/0.18.0/darwin-x64.tar.gz"))
			} else if runtime.GOOS == "linux" {
				Expect(result["unarchive-filename"]).To(Equal("dist_binaries"))
				Expect(result["filename"]).To(Equal("linux-x64"))
				Expect(result["url"]).To(Equal("https://dl.bintray.com/elmlang/elm-platform/0.18.0/linux-x64.tar.gz"))
			}
		})

		It("should get info about 0.17.1 version", func() {
			result := (&Elm{Version: "0.17.1"}).Info()

			Expect(result["archive-folder"]).Should(ContainSubstring("elm-archive-0.17.1/"))

			// :/
			if runtime.GOOS == "darwin" {
				Expect(result["unarchive-filename"]).To(Equal(""))
				Expect(result["filename"]).To(Equal("darwin-x64"))
				Expect(result["url"]).To(Equal("https://dl.bintray.com/elmlang/elm-platform/0.17.1/darwin-x64.tar.gz"))
			} else if runtime.GOOS == "linux" {
				Expect(result["unarchive-filename"]).To(Equal("dist_binaries"))
				Expect(result["filename"]).To(Equal("linux-x64"))
				Expect(result["url"]).To(Equal("https://dl.bintray.com/elmlang/elm-platform/0.17.1/linux-x64.tar.gz"))
			}
		})

		It("should get info about 0.17.0 version", func() {
			result := (&Elm{Version: "0.17.0"}).Info()

			// :/
			if runtime.GOOS == "darwin" {
				Expect(result["unarchive-filename"]).To(Equal("dist_binaries"))
				Expect(result["url"]).To(Equal("https://dl.bintray.com/elmlang/elm-platform/0.17.0/darwin-x64.tar.gz"))
			} else if runtime.GOOS == "linux" {
				Expect(result["unarchive-filename"]).To(Equal("dist_binaries"))
				Expect(result["url"]).To(Equal("https://dl.bintray.com/elmlang/elm-platform/0.17.0/linux-x64.tar.gz"))
			}
		})

		It("should get info about 0.15.1 version", func() {
			result := (&Elm{Version: "0.15.1"}).Info()

			// :/
			if runtime.GOOS == "darwin" {
				Expect(result["unarchive-filename"]).To(Equal("osx"))
				Expect(result["url"]).To(Equal("https://dl.bintray.com/elmlang/elm-platform/0.15.1/darwin-x64.tar.gz"))
			} else if runtime.GOOS == "linux" {
				Expect(result["unarchive-filename"]).To(Equal("linux64"))
				Expect(result["url"]).To(Equal("https://dl.bintray.com/elmlang/elm-platform/0.15.1/linux-x64.tar.gz"))
			}
		})
	})
})
