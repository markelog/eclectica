package rust_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/jarcoal/httpmock"
	"github.com/markelog/monkey"

	. "github.com/markelog/eclectica/plugins/rust"

	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/variables"
)

func Read(path string) string {
	bytes, _ := ioutil.ReadFile(path)

	return string(bytes)
}

var _ = Describe("rust", func() {
	var (
		version     = "1.0.0"
		path, _     = filepath.Abs("../../testdata/plugins/versions/")
		versionPath = filepath.Join(path, "rust", version)
	)

	AfterEach(func() {
		os.RemoveAll(versionPath)
	})

	Describe("Install", func() {
		It("should call install script with right arguments", func() {
			program := ""
			firstArg := ""

			io.CreateDir(versionPath)

			monkey.Patch(exec.Command, func(name string, arg ...string) *exec.Cmd {
				program = name
				firstArg = arg[0]

				return &exec.Cmd{}
			})

			monkey.Patch(variables.Home, func() string {
				return path
			})

			(&Rust{Version: version}).Install()

			Expect(program).To(ContainSubstring("versions/rust/" + version + "/install.sh"))

			path := filepath.Join(variables.Path("rust", version), "tmp")
			Expect(firstArg).To(Equal("--prefix=" + path))

			monkey.Unpatch(exec.Command)
		})
	})

	Describe("Info", func() {
		AfterEach(func() {
			defer httpmock.DeactivateAndReset()
		})

		It("should get info about nightly version", func() {
			result := (&Rust{Version: "nightly"}).Info()

			// :/
			if runtime.GOOS == "darwin" {
				Expect(result["filename"]).To(Equal("rust-nightly-x86_64-apple-darwin"))
				Expect(result["url"]).To(Equal("https://static.rust-lang.org/dist/rust-nightly-x86_64-apple-darwin.tar.gz"))
			} else if runtime.GOOS == "linux" {
				Expect(result["filename"]).To(Equal("rust-nightly-x86_64-unknown-linux-gnu"))
				Expect(result["url"]).To(Equal("https://static.rust-lang.org/dist/rust-nightly-x86_64-unknown-linux-gnu.tar.gz"))
			}
		})

		It("should get info about beta version", func() {
			result := (&Rust{Version: "beta"}).Info()

			// :/
			if runtime.GOOS == "darwin" {
				Expect(result["filename"]).To(Equal("rust-beta-x86_64-apple-darwin"))
				Expect(result["url"]).To(Equal("https://static.rust-lang.org/dist/rust-beta-x86_64-apple-darwin.tar.gz"))
			} else if runtime.GOOS == "linux" {
				Expect(result["filename"]).To(Equal("rust-beta-x86_64-unknown-linux-gnu"))
				Expect(result["url"]).To(Equal("https://static.rust-lang.org/dist/rust-beta-x86_64-unknown-linux-gnu.tar.gz"))
			}
		})

		It("should get info about 1.9.0 version", func() {
			result := (&Rust{Version: "1.9.0"}).Info()

			// :/
			if runtime.GOOS == "darwin" {
				Expect(result["filename"]).To(Equal("rust-1.9.0-x86_64-apple-darwin"))
				Expect(result["url"]).To(Equal("https://static.rust-lang.org/dist/rust-1.9.0-x86_64-apple-darwin.tar.gz"))
			} else if runtime.GOOS == "linux" {
				Expect(result["filename"]).To(Equal("rust-1.9.0-x86_64-unknown-linux-gnu"))
				Expect(result["url"]).To(Equal("https://static.rust-lang.org/dist/rust-1.9.0-x86_64-unknown-linux-gnu.tar.gz"))
			}
		})
	})
})
