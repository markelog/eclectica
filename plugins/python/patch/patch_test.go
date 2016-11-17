package patch_test

import (
	"io"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	eio "github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/plugins/python/patch"
)

var _ = Describe("patch", func() {
	var (
		ts      *httptest.Server
		oldLink = patch.Link
	)

	BeforeEach(func() {
		content := eio.Read("./testdata/patches.html")

		ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)

			io.WriteString(w, content)
		}))

		patch.Link = ts.URL
	})

	AfterEach(func() {
		ts.Close()
		ts = nil

		patch.Link = oldLink
	})

	Describe("Urls", func() {
		It("should contains url list for specific version", func() {
			urls, err := patch.Urls("2.6.9")

			Expect(err).To(BeNil())

			Expect(urls[0]).Should(ContainSubstring("/000_patch-setup.py.diff"))
			Expect(urls[0]).Should(ContainSubstring("2.6.9/Python-2.6.9"))
			Expect(urls[0]).Should(ContainSubstring(patch.RawLink))

			Expect(urls[3]).Should(ContainSubstring("/003_tk86.patch"))
			Expect(urls[3]).Should(ContainSubstring("2.6.9/Python-2.6.9"))
			Expect(urls[3]).Should(ContainSubstring(patch.RawLink))
		})

		It("should contains url list for specific version", func() {
			urls, err := patch.Urls("2.7.0")

			Expect(err).To(BeNil())

			Expect(urls[0]).Should(ContainSubstring("/000_patch-setup.py.diff"))
			Expect(urls[0]).Should(ContainSubstring("2.7/Python-2.7"))
			Expect(urls[0]).Should(ContainSubstring(patch.RawLink))

			Expect(urls[3]).Should(ContainSubstring("/003_tk86.patch"))
			Expect(urls[3]).Should(ContainSubstring("2.7/Python-2.7"))
			Expect(urls[3]).Should(ContainSubstring(patch.RawLink))
		})
	})
})
