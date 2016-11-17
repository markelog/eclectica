package python_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	eio "github.com/markelog/eclectica/io"
	. "github.com/markelog/eclectica/plugins/python"
)

var _ = Describe("python", func() {
	var (
		remotes []string
		err     error
	)

	python := &Python{}

	Describe("ListRemote", func() {
		old := VersionsLink

		AfterEach(func() {
			VersionsLink = old
		})

		Describe("success", func() {
			BeforeEach(func() {
				content := eio.Read("../../testdata/plugins/python/index.html")

				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					status := 200

					if _, ok := r.URL.Query()["status"]; ok {
						fmt.Sscanf(r.URL.Query().Get("status"), "%d", &status)
					}

					w.WriteHeader(status)
					io.WriteString(w, content)
				}))

				VersionsLink = ts.URL

				remotes, err = python.ListRemote()
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should exclude everything less then 2.6 versions", func() {
				Expect(remotes[0]).To(Equal("2.6"))
			})
		})

		Describe("fail", func() {
			BeforeEach(func() {
				VersionsLink = ""
				remotes, err = python.ListRemote()
			})

			It("should return an error", func() {
				Expect(err).Should(MatchError("Can't establish connection"))
			})
		})
	})

	Describe("Info", func() {
		BeforeEach(func() {
			content := eio.Read("../../testdata/plugins/python/index.html")

			httpmock.Activate()

			httpmock.RegisterResponder(
				"GET",
				"https://pythonjs.org/dist/latest/SHASUMS256.txt",
				httpmock.NewStringResponder(200, content),
			)

			httpmock.RegisterResponder(
				"GET",
				"https://pythonjs.org/dist/lts/SHASUMS256.txt",
				httpmock.NewStringResponder(200, content),
			)
		})

		AfterEach(func() {
			defer httpmock.DeactivateAndReset()
		})

		It("should get info about latest version", func() {
			Skip("Waiting on #10")
			// result, _ := (&Python{Version: "latest"}).Info()

			// TODO
		})

		It("should get info about 2.0 version", func() {
			result, _ := (&Python{Version: "3.0.0"}).Info()

			Expect(result["version"]).To(Equal("3.0"))
			Expect(result["filename"]).To(Equal("Python-3.0"))
			Expect(result["url"]).To(Equal("https://www.python.org/ftp/python/3.0/Python-3.0.tgz"))
		})

		It("should get info about 3.0 version", func() {
			result, _ := (&Python{Version: "3.0.0"}).Info()

			Expect(result["version"]).To(Equal("3.0"))
			Expect(result["filename"]).To(Equal("Python-3.0"))
			Expect(result["url"]).To(Equal("https://www.python.org/ftp/python/3.0/Python-3.0.tgz"))
		})

		It("up the not ante for 3.2", func() {
			result, _ := (&Python{Version: "3.2.0"}).Info()

			Expect(result["version"]).To(Equal("3.2"))
			Expect(result["filename"]).To(Equal("Python-3.2"))
			Expect(result["url"]).To(Equal("https://www.python.org/ftp/python/3.2/Python-3.2.tgz"))
		})

		It("up the ante for 3.3", func() {
			result, _ := (&Python{Version: "3.3.0"}).Info()

			Expect(result["version"]).To(Equal("3.3.0"))
			Expect(result["filename"]).To(Equal("Python-3.3.0"))
			Expect(result["url"]).To(Equal("https://www.python.org/ftp/python/3.3.0/Python-3.3.0.tgz"))
		})
	})
})
