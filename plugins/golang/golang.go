package golang

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chuckpreslar/emission"
	"github.com/go-errors/errors"

	"github.com/markelog/eclectica/pkg"
	"github.com/markelog/eclectica/variables"
	"github.com/markelog/eclectica/versions"
)

var (
	VersionsLink   = "https://golang.org/dl"
	DownloadLink   = "https://storage.googleapis.com/golang"
	versionPattern = "\\d+\\.\\d+(?:\\.\\d+)?(?:(alpha|beta|rc)(?:\\d*)?)?"

	bins = []string{"go", "godoc", "gofmt"}
	dots = []string{".go-version"}
)

type Golang struct {
	Version string
	Emitter *emission.Emitter
	pkg.Base
}

func New(version string, emitter *emission.Emitter) *Golang {
	return &Golang{
		Version: version,
		Emitter: emitter,
	}
}

func (golang Golang) Events() *emission.Emitter {
	return golang.Emitter
}

func (golang Golang) PostInstall() error {
	return dealWithShell()
}

func (golang Golang) Environment() (result []string, err error) {
	result = append(result, "GOROOT="+variables.Path("go", golang.Version))

	// Go versions lower then 1.7 do not have default `GOPATH` environment variable.
	// Starting from 1.7 `GOPATH` is now set to `~/go` path (see `go help gopath` for more)
	// We do the same if for other versions as a default, but only if user didn't set themselves
	if os.Getenv("GOPATH") == "" {
		result = append(result, "GOPATH="+filepath.Join(os.Getenv("HOME"), "go"))
	}

	return
}

func (golang Golang) Info() map[string]string {
	result := make(map[string]string)

	platform, _ := getPlatform()

	version := versions.Unsemverify(golang.Version)
	result["version"] = version
	result["unarchive-filename"] = "go"
	result["filename"] = fmt.Sprintf("go%s.%s", version, platform)
	result["url"] = fmt.Sprintf("%s/%s.tar.gz", DownloadLink, result["filename"])

	return result
}

func (rust Golang) Bins() []string {
	return bins
}

func (rust Golang) Dots() []string {
	return dots
}

func (golang Golang) ListRemote() ([]string, error) {
	doc, err := goquery.NewDocument(VersionsLink)

	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, errors.New(variables.ConnectionError)
		}

		return nil, err
	}

	result := []string{}
	rVersion := regexp.MustCompile(versionPattern)

	doc.Find("tr:first-of-type td:first-of-type.filename a").Each(func(i int, node *goquery.Selection) {
		text := node.Text()

		if strings.Contains(text, "bootstrap") {
			return
		}

		// We not checking for duplicates, since it just might create more errors
		version := rVersion.FindAllStringSubmatch(text, 1)[0][0]
		result = append(result, versions.Semverify(version))
	})

	return result, nil
}

func getPlatform() (string, error) {
	if runtime.GOOS == "linux" {
		return "linux-amd64", nil
	}

	if runtime.GOOS == "darwin" {
		return "darwin-amd64", nil
	}

	return "", errors.New("Not supported envionment")
}
