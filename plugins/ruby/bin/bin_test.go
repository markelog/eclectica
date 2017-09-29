package bin_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	eio "github.com/markelog/eclectica/io"
	. "github.com/markelog/eclectica/plugins/ruby/bin"
)

var _ = Describe("bin ruby", func() {
	var (
		remotes []string
		err     error
	)

	ruby := &Ruby{}

	Describe("ListRemote", func() {
		old := VersionLink

		AfterEach(func() {
			VersionLink = old
		})

		Describe("success", func() {
			BeforeEach(func() {
				content := eio.Read("../../../testdata/plugins/ruby/bin-dist.html")

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

				remotes, err = ruby.ListRemote()
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should contain 2.3.3 version", func() {
				Expect(remotes).To(ContainElement("2.3.3"))
			})
		})

		Describe("fail", func() {
			BeforeEach(func() {
				VersionLink = ""
				remotes, err = ruby.ListRemote()
			})

			It("should return an error", func() {
				Expect(err).Should(MatchError("Connection cannot be established"))
			})
		})
	})

	Describe("Info", func() {
		It("should get info about 2.2.3 version", func() {
			result := (&Ruby{Version: "2.2.3"}).Info()

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
