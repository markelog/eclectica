package golang_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"runtime"

	"github.com/bouk/monkey"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/markelog/eclectica/plugins/golang"
)

func Read(path string) string {
	bytes, _ := ioutil.ReadFile(path)

	return string(bytes)
}

var _ = Describe("golang", func() {
	var (
		remotes []string
		err     error
	)

	golang := &Golang{}

	Describe("ListRemote", func() {
		old := VersionsLink

		AfterEach(func() {
			VersionsLink = old
		})

		Describe("success", func() {
			BeforeEach(func() {
				content := Read("../../testdata/golang/download.xml")

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

				remotes, err = golang.ListRemote()
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("should have correct version values", func() {
				// :/
				if runtime.GOOS == "darwin" {
					Expect(remotes[0]).To(Equal("1.4.3"))
				}

				if runtime.GOOS == "linux" {
					Expect(remotes[0]).To(Equal("1.2.2"))
				}
			})

			It("shouldn't have duplicates", func() {
				var i int

				for _, element := range remotes {
					i = 0
					for _, secondRound := range remotes {
						if element == secondRound {
							i++
						}
					}

					Expect(i).To(Equal(1))
				}
			})
		})

		Describe("fail", func() {
			BeforeEach(func() {
				VersionsLink = ""
				remotes, err = golang.ListRemote()
			})

			It("should return an error", func() {
				Expect(err).Should(MatchError("Can't establish connection"))
			})
		})
	})

	Describe("Info", func() {
		BeforeEach(func() {
			content := Read("../../testdata/golang/latest.txt")

			httpmock.Activate()

			httpmock.RegisterResponder(
				"GET",
				"https://golang.org/dist/latest/SHASUMS256.txt",
				httpmock.NewStringResponder(200, content),
			)

			httpmock.RegisterResponder(
				"GET",
				"https://golang.org/dist/lts/SHASUMS256.txt",
				httpmock.NewStringResponder(200, content),
			)
		})

		AfterEach(func() {
			defer httpmock.DeactivateAndReset()
		})

		It("should get info about 1.7 version", func() {
			result, _ := golang.Info("1.7")

			Expect(result["version"]).To(Equal("1.7"))

			// :/
			if runtime.GOOS == "darwin" {
				Expect(result["filename"]).To(Equal("go1.7.darwin-amd64"))
				Expect(result["url"]).To(Equal("https://storage.googleapis.com/golang/go1.7.darwin-amd64.tar.gz"))
			} else if runtime.GOOS == "linux" {
				Expect(result["filename"]).To(Equal("go1.7.linux-amd64"))
				Expect(result["url"]).To(Equal("https://storage.googleapis.com/golang/go1.7.linux-amd64.tar.gz"))
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

			golang.Current()

			Expect(program).To(ContainSubstring("bin/go"))
			Expect(firstArg).To(Equal("version"))

			monkey.Unpatch(exec.Command)
		})

		It("should report correct version", func() {
			program := ""
			firstArg := ""

			monkey.Patch(exec.Command, func(name string, arg ...string) *exec.Cmd {
				program = name
				firstArg = arg[0]

				return &exec.Cmd{}
			})

			monkey.Patch((*exec.Cmd).Output, func(*exec.Cmd) ([]uint8, error) {
				return []uint8("go version go1.6.2 darwin/amd64"), nil
			})

			result := golang.Current()

			Expect(result).To(Equal("1.6.2"))
			monkey.Unpatch(exec.Command)
		})
	})
})
