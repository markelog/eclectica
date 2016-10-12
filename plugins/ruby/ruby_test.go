package ruby_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/markelog/eclectica/plugins/ruby"
)

func Read(path string) string {
	bytes, _ := ioutil.ReadFile(path)

	return string(bytes)
}

var _ = Describe("ruby", func() {
	var (
		remotes []string
		err     error
	)

	ruby := &Ruby{}

	Describe("ListRemote", func() {
		old := VersionsLink

		AfterEach(func() {
			VersionsLink = old
		})

		Describe("success", func() {
			BeforeEach(func() {
				content := Read("../../testdata/ruby/dist.html")

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

				remotes, err = ruby.ListRemote()
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should have correct version values", func() {
				Expect(remotes[0]).To(Equal("2.0.0-p451"))
			})
		})

		Describe("fail", func() {
			BeforeEach(func() {
				VersionsLink = ""
				remotes, err = ruby.ListRemote()
			})

			It("should return an error", func() {
				Expect(err).Should(MatchError("Can't establish connection"))
			})
		})
	})

	Describe("Info", func() {
		It("should get info about 2.2.3 version", func() {
			result, _ := (&Ruby{Version: "2.2.3"}).Info()

			Expect(result["filename"]).To(Equal("ruby-2.2.3"))

			Expect(result["url"]).Should(ContainSubstring("https://rvm.io/binaries"))
			Expect(result["url"]).Should(ContainSubstring("x86_64"))
			Expect(result["url"]).Should(ContainSubstring("ruby-2.2.3.tar.bz2"))

			if runtime.GOOS == "darwin" {
				Expect(result["url"]).Should(ContainSubstring("osx"))
			}

			if runtime.GOOS == "linux" {
				Expect(result["url"]).Should(ContainSubstring("ubuntu"))
			}
		})
	})
})
