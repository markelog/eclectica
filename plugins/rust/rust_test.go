package rust_test

import (
	"io/ioutil"
	"os/exec"
	"regexp"
	"runtime"

	"github.com/bouk/monkey"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/markelog/eclectica/plugins/rust"
	"github.com/markelog/eclectica/variables"
)

func Read(path string) string {
	bytes, _ := ioutil.ReadFile(path)

	return string(bytes)
}

var _ = Describe("rust", func() {
	var (
		remotes []string
		err     error
		rust    *Rust
	)

	BeforeEach(func() {
		rust = &Rust{}
	})

	Describe("Install", func() {
		It("should call install script with right arguments", func() {
			program := ""
			firstArg := ""

			monkey.Patch(exec.Command, func(name string, arg ...string) *exec.Cmd {
				program = name
				firstArg = arg[0]

				return &exec.Cmd{}
			})

			rust.Install("1.0.0")

			Expect(program).To(ContainSubstring("versions/rust/1.0.0/install.sh"))
			Expect(firstArg).To(Equal("--prefix=" + variables.Prefix("rustc")))

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
				content := Read("../../testdata/rust/dist.txt")

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
				rp := regexp.MustCompile("[[:digit:]]+\\.[[:digit:]]+")

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
			result, _ := rust.Info("nightly")

			Expect(result["name"]).To(Equal("rust"))
			Expect(result["version"]).To(Equal("nightly"))

			// :/
			if runtime.GOOS == "darwin" {
				Expect(result["filename"]).To(Equal("rust-nightly-x86_64-apple-darwin"))
				Expect(result["url"]).To(Equal("https://static.rust-lang.org/dist/rust-nightly-x86_64-apple-darwin.tar.gz"))
			} else if runtime.GOOS == "linux" {
				Expect(result["filename"]).To(Equal("rust-nightly-x86_64-unknown-linux-gnu"))
				Expect(result["url"]).To(Equal("https://static.rust-lang.org/dist/rust-nightly-x86_64-unknown-linux-gnu.tar.gz"))
			}
		})

		It("should get info about lts version", func() {
			result, _ := rust.Info("beta")

			Expect(result["name"]).To(Equal("rust"))
			Expect(result["version"]).To(Equal("beta"))

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
			result, _ := rust.Info("1.9.0")

			Expect(result["name"]).To(Equal("rust"))
			Expect(result["version"]).To(Equal("1.9.0"))

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

	Describe("Current", func() {
		It("should handle empty output", func() {
			program := ""
			firstArg := ""

			monkey.Patch(exec.Command, func(name string, arg ...string) *exec.Cmd {
				program = name
				firstArg = arg[0]

				return &exec.Cmd{}
			})

			monkey.Patch((*exec.Cmd).Output, func(*exec.Cmd) ([]uint8, error) {
				return []uint8("test"), nil
			})

			rust.Current()

			Expect(program).To(ContainSubstring("bin/rustc"))
			Expect(firstArg).To(Equal("--version"))

			monkey.Unpatch(exec.Command)
		})
	})

	It("outputs version", func() {
		program := ""
		firstArg := ""

		monkey.Patch(exec.Command, func(name string, arg ...string) *exec.Cmd {
			program = name
			firstArg = arg[0]

			return &exec.Cmd{}
		})

		monkey.Patch((*exec.Cmd).Output, func(*exec.Cmd) ([]uint8, error) {
			return []uint8("rustc 1.5.0 (3d7cd77e4 2015-12-04)"), nil
		})

		result := rust.Current()

		Expect(program).To(ContainSubstring("bin/rustc"))
		Expect(firstArg).To(Equal("--version"))
		Expect(result).To(Equal("1.5.0"))

		monkey.Unpatch(exec.Command)
	})
})
