package golang_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"runtime"

	"github.com/bouk/monkey"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/markelog/eclectica/cmd/print"
	. "github.com/markelog/eclectica/plugins/golang"

	eio "github.com/markelog/eclectica/io"
)

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
				content := eio.Read("../../testdata/plugins/golang/download.xml")

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
			content := eio.Read("../../testdata/plugins/golang/latest.txt")

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
			result := (&Golang{Version: "1.7"}).Info()

			// :/
			if runtime.GOOS == "darwin" {
				Expect(result["filename"]).To(Equal("go1.7.darwin-amd64"))
				Expect(result["url"]).To(Equal("https://storage.googleapis.com/golang/go1.7.darwin-amd64.tar.gz"))
			} else if runtime.GOOS == "linux" {
				Expect(result["filename"]).To(Equal("go1.7.linux-amd64"))
				Expect(result["url"]).To(Equal("https://storage.googleapis.com/golang/go1.7.linux-amd64.tar.gz"))
			}
		})

		It("should get info about 1.7.0 version", func() {
			result := (&Golang{Version: "1.7.0"}).Info()

			Expect(result["version"]).To(Equal("1.7"))
		})

		It("should get info about 1.7beta1 version", func() {
			result := (&Golang{Version: "1.7.0-beta1"}).Info()

			Expect(result["version"]).To(Equal("1.7beta1"))
		})
	})

	Describe("PostInstall", func() {
		BeforeEach(func() {
			monkey.Patch(exec.Command, func(name string, arg ...string) *exec.Cmd {
				return &exec.Cmd{}
			})
		})

		It("should print warning for absent git", func() {
			var msg, cmd string

			monkey.Patch(print.Warning, func(message, command string) {
				msg = message
				cmd = command
			})

			golang.PostInstall()

			Expect(msg).Should(ContainSubstring("Golang has been installed"))
			Expect(msg).Should(ContainSubstring("you need to do it only"))

			if runtime.GOOS == "linux" {
				Expect(cmd).To(Equal("sudo apt-get update && sudo apt-get install -y git"))
			}

			if runtime.GOOS == "darwin" {
				Expect(cmd).To(Equal("brew update && brew install git"))
			}
		})
	})

	Describe("Environment", func() {
		It("should set GOROOT and GOPATH environment variables", func() {
			monkey.Patch(os.Getenv, func(name string) string {
				return ""
			})

			result, _ := golang.Environment()

			Expect(result[0]).To(Equal("GOROOT=.eclectica/versions/go"))
			Expect(result[1]).To(Equal("GOPATH=go"))

			monkey.Unpatch(os.Getenv)
		})

		It("should set GOROOT and GOPATH environment variables", func() {
			monkey.Patch(os.Getenv, func(name string) string {
				if name == "GOPATH" {
					return "test"
				}

				return ""
			})

			result, _ := golang.Environment()

			Expect(len(result)).To(Equal(1))
			Expect(result[0]).To(Equal("GOROOT=.eclectica/versions/go"))

			monkey.Unpatch(os.Getenv)
		})
	})
})
