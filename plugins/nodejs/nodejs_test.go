package nodejs_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/markelog/eclectica/plugins/nodejs"
)

func Read(path string) string {
	bytes, _ := ioutil.ReadFile(path)

	return string(bytes)
}

var _ = Describe("nodejs", func() {
	var (
		remotes []string
		err     error
	)

	node := &Node{}

	Describe("ListRemote", func() {
		old := VersionsLink

		AfterEach(func() {
			VersionsLink = old
		})

		Describe("success", func() {
			BeforeEach(func() {
				content := Read("../../testdata/plugins/nodejs/dist.html")

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

				VersionsLink = ts.URL

				remotes, err = node.ListRemote()
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should have correct version values", func() {
				Expect(remotes[0]).To(Equal("6.4.0"))
			})
		})

		Describe("fail", func() {
			BeforeEach(func() {
				VersionsLink = ""
				remotes, err = node.ListRemote()
			})

			It("should return an error", func() {
				Expect(err).Should(MatchError("Can't establish connection"))
			})
		})
	})

	Describe("Info", func() {
		BeforeEach(func() {
			content := Read("../../testdata/plugins/nodejs/latest.txt")

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
			Skip("Waiting on #10")
			result, _ := (&Node{Version: "latest"}).Info()

			// :/
			if runtime.GOOS == "darwin" {
				Expect(result["filename"]).To(Equal("node-v6.3.1-darwin-x64"))
				Expect(result["url"]).To(Equal("https://nodejs.org/dist/latest/node-v6.3.1-darwin-x64.tar.gz"))
			} else if runtime.GOOS == "linux" {
				Expect(result["filename"]).To(Equal("node-v6.3.1-linux-x64"))
				Expect(result["url"]).To(Equal("https://nodejs.org/dist/latest/node-v6.3.1-linux-x64.tar.gz"))
			}
		})

		It("should get info about 6.3.1 version", func() {
			result, _ := (&Node{Version: "6.3.1"}).Info()

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
