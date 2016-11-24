package ruby_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	eio "github.com/markelog/eclectica/io"
	. "github.com/markelog/eclectica/plugins/ruby"
)

var _ = Describe("ruby", func() {
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
				content := eio.Read("../../testdata/plugins/ruby/dist.xml")

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

			It("should have correct version values (first and last)", func() {
				if runtime.GOOS == "darwin" {
					Expect(remotes[0]).To(Equal("2.1.5"))
				}

				if runtime.GOOS == "linux" {
					Expect(remotes[0]).To(Equal("1.8.7"))
				}

				last := len(remotes) - 1
				Expect(remotes[last]).To(Equal("2.3.3"))
			})
		})

		Describe("fail", func() {
			BeforeEach(func() {
				VersionLink = ""
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

			if runtime.GOOS == "darwin" {
				Expect(result["url"]).To(Equal("https://s3.amazonaws.com/travis-rubies/binaries/osx/10.12/x86_64/ruby-2.2.3.tar.bz2"))
			}

			if runtime.GOOS == "linux" {
				Expect(result["url"]).To(Equal("https://s3.amazonaws.com/travis-rubies/binaries/ubuntu/14.04/x86_64/ruby-2.2.3.tar.bz2"))
			}
		})
	})
})
