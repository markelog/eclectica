package nodejs_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	eio "github.com/markelog/eclectica/io"
	. "github.com/markelog/eclectica/plugins/nodejs"
	"github.com/markelog/eclectica/variables"
)

var _ = Describe("nodejs", func() {
	var (
		remotes []string
		err     error
	)

	node := &Node{}

	Describe("ListRemote", func() {
		old := VersionLink

		AfterEach(func() {
			VersionLink = old
		})

		Describe("success", func() {
			BeforeEach(func() {
				content := eio.Read("../../testdata/plugins/nodejs/dist.html")

				// httpmock is not incompatible with goquery :/.
				// See https://github.com/jarcoal/httpmock/issues/18
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					status := 200

					if _, ok := r.URL.Query()["status"]; ok {
						fmt.Sscanf(r.URL.Query().Get("status"), "%d", &status)
					}

					w.WriteHeader(status)
					io.WriteString(w, content)
				}))

				VersionLink = ts.URL

				remotes, err = node.ListRemote()
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should have have 6.x versions", func() {
				Expect(remotes).To(ContainElement("6.4.0"))
			})

			It("should contain 0.10.x versions", func() {
				Expect(remotes).To(ContainElement("0.10.13"))
			})

			It("should contain 0.12.x versions", func() {
				Expect(remotes).To(ContainElement("0.12.6"))
			})

			It("should not contain 0.8.x versions", func() {
				Expect(remotes).NotTo(ContainElement("0.8.12"))
			})

			It("should not contain 0.1.x versions", func() {
				Expect(remotes).NotTo(ContainElement("0.1.14"))
			})
		})

		Describe("fail", func() {
			BeforeEach(func() {
				VersionLink = ""
				remotes, err = node.ListRemote()
			})

			It("should return an error", func() {
				Expect(err).Should(MatchError(variables.ConnectionError))
			})
		})
	})

	Describe("Info", func() {
		BeforeEach(func() {
			content := eio.Read("../../testdata/plugins/nodejs/latest.txt")

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

		It("should get info about 6.3.1 version", func() {
			result := (&Node{Version: "6.3.1"}).Info()

			// :/
			if runtime.GOOS == "darwin" {
				Expect(result["filename"]).To(Equal("node-v6.3.1-darwin-x64"))
				Expect(result["url"]).To(Equal("https://nodejs.org/dist/v6.3.1/node-v6.3.1-darwin-x64.tar.gz"))
			} else if runtime.GOOS == "linux" {
				Expect(result["filename"]).To(Equal("node-v6.3.1-linux-x64"))
				Expect(result["url"]).To(Equal("https://nodejs.org/dist/v6.3.1/node-v6.3.1-linux-x64.tar.gz"))
			}
		})
	})
})
