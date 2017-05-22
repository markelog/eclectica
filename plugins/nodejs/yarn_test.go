package nodejs_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bouk/monkey"
	"github.com/markelog/archive"
	"github.com/markelog/cprf"
	"gopkg.in/cavaliercoder/grab.v1"

	"github.com/markelog/eclectica/io"
	. "github.com/markelog/eclectica/plugins/nodejs"
)

var _ = Describe("yarn", func() {
	var (
		grabGet        bool
		archiveExtract bool
		cprfCopy       bool
		osRemovaAll    bool
		ioSymlink      bool
	)

	BeforeEach(func() {
		grabGet = false
		archiveExtract = false
		cprfCopy = false
		osRemovaAll = false
		ioSymlink = false

		monkey.Patch(grab.Get, func(path, url string) (*grab.Response, error) {
			grabGet = true
			return nil, nil
		})

		monkey.Patch(archive.Extract, func(from, to string) error {
			archiveExtract = true
			return nil
		})

		monkey.Patch(cprf.Copy, func(from, to string) error {
			cprfCopy = true
			return nil
		})

		monkey.Patch(os.RemoveAll, func(path string) error {
			osRemovaAll = true
			return nil
		})

		monkey.Patch(io.Symlink, func(from, to string) error {
			ioSymlink = true
			return nil
		})
	})

	AfterEach(func() {
		monkey.Unpatch(grab.Get)
		monkey.Unpatch(archive.Extract)
		monkey.Unpatch(cprf.Copy)
		monkey.Unpatch(os.RemoveAll)
		monkey.Unpatch(io.Symlink)
	})

	It("should not try to install yarn for unsupported node version", func() {
		working, err := (&Node{Version: "0.10.0"}).Yarn()

		Expect(grabGet).To(Equal(false))
		Expect(archiveExtract).To(Equal(false))
		Expect(cprfCopy).To(Equal(false))
		Expect(osRemovaAll).To(Equal(false))
		Expect(ioSymlink).To(Equal(false))

		Expect(working).To(Equal(true))
		Expect(err.Error()).To(Equal("\"0.10.0\" version is not supported by yarn"))
	})

	It("installs yarn", func() {
		working, err := (&Node{Version: "6.10.0"}).Yarn()

		Expect(grabGet).To(Equal(true))
		Expect(archiveExtract).To(Equal(true))
		Expect(cprfCopy).To(Equal(true))
		Expect(osRemovaAll).To(Equal(true))
		Expect(ioSymlink).To(Equal(true))

		Expect(working).To(Equal(true))
		Expect(err).To(BeNil())
	})
})
