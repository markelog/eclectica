package compile_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	eio "github.com/markelog/eclectica/io"
	. "github.com/markelog/eclectica/plugins/ruby/compile"
	"github.com/markelog/eclectica/variables"
)

var _ = Describe("compile ruby", func() {
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
				content := eio.Read("../../../testdata/plugins/ruby/compile-dist.html")

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
				Expect(remotes).To(ContainElement("2.3.3"))
			})

			It("should not contain preview version (2.5.0 in this case)", func() {
				Expect(remotes).ToNot(ContainElement("2.5.0"))
			})
		})

		Describe("fail", func() {
			BeforeEach(func() {
				VersionLink = ""
				remotes, err = ruby.ListRemote()
			})

			It("should return an error", func() {
				Expect(err).Should(MatchError(variables.ConnectionError))
			})
		})
	})

	Describe("Info", func() {
		It("should get info about 2.2.3 version", func() {
			result := (&Ruby{Version: "2.2.3"}).Info()

			Expect(result["filename"]).To(Equal("ruby-2.2.3"))
			Expect(result["url"]).To(Equal("https://cache.ruby-lang.org/pub/ruby/ruby-2.2.3.tar.gz"))
		})
	})
})
