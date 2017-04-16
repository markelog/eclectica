package rust_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bouk/monkey"
	"github.com/jarcoal/httpmock"

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
		remotes     []string
		err         error
		rust        *Rust
		version     = "1.0.0"
		path, _     = filepath.Abs("../../testdata/plugins/versions/")
		versionPath = filepath.Join(path, "rust", version)
	)

	BeforeEach(func() {
		rust = &Rust{}
	})

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

	Describe("ListRemote", func() {
		Describe("fail", func() {
			BeforeEach(func() {
				httpmock.Activate()

				httpmock.RegisterResponder(
					"GET",
					"https://static.rust-lang.org/dist/index.txt",
					httpmock.NewStringResponder(500, ""),
				)
			})

			AfterEach(func() {
				defer httpmock.DeactivateAndReset()
			})

			It("should return an error", func() {
				remotes, err = rust.ListRemote()

				Expect(err).Should(MatchError("Can't establish connection"))
			})
		})

		Describe("success", func() {
			BeforeEach(func() {
				content := Read("../../testdata/plugins/rust/dist.txt")

				httpmock.Activate()

				httpmock.RegisterResponder(
					"GET",
					"https://static.rust-lang.org/dist/index.txt",
					httpmock.NewStringResponder(200, content),
				)
			})

			AfterEach(func() {
				defer httpmock.DeactivateAndReset()
			})

			BeforeEach(func() {
				remotes, err = rust.ListRemote()
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("gets list of versions", func() {
				Expect(remotes[0]).To(Equal("0.10"))
				Expect(remotes[2]).To(Equal("0.12.0"))
				Expect(remotes[8]).To(Equal("1.0.0-beta.4"))
			})

			It("should have correct version values", func() {
				rp := regexp.MustCompile("\\d+\\.\\d+")

				for _, element := range remotes {
					Expect(rp.MatchString(element)).To(Equal(true))
				}
			})
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
