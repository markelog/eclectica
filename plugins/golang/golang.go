// Package golang provides all needed logic for installation of Golang
package golang

import (
	"fmt"
	"net"
	"os"
	"os/exec"
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
	// VersionLink is the URL link from which we can get all possible versions
	VersionLink = "https://golang.org/dl"

	// DownloadLink from which we download binaries for golang
	DownloadLink = "https://storage.googleapis.com/golang"

	versionPattern = `\d+\.\d+(?:\.\d+)?(?:(alpha|beta|rc)(?:\d*)?)?`

	bins = []string{"go", "godoc", "gofmt"}
	dots = []string{".go-version"}

	rVersion = regexp.MustCompile(versionPattern)
)

// Golang essential struct
type Golang struct {
	Version string
	Emitter *emission.Emitter
	pkg.Base
}

// New returns language struct
func New(version string, emitter *emission.Emitter) *Golang {
	return &Golang{
		Version: version,
		Emitter: emitter,
	}
}

// Events returns language related event emitter
func (golang Golang) Events() *emission.Emitter {
	return golang.Emitter
}

// PostInstall hook
func (golang Golang) PostInstall() error {
	// In case consumer used go binaries installed without eclectica.
	// So some IDE's use go autocomplete which
	// might use the gocode daemon, which cache the previous path even
	// if you explicitly set the GOROOT path, so we would need to drop that cache
	// Not sure if that would always help
	exec.Command("gocode", "drop-cache")

	return dealWithShell()
}

// Environment returns list of the all needed envionment variables
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

// Info provides all the info needed for installation of the plugin
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

// Bins returns list of the all bins included
// with the distribution of the language
func (golang Golang) Bins() []string {
	return bins
}

// Dots returns list of the all available filenames
// which can define versions
func (golang Golang) Dots() []string {
	return dots
}

// ListRemote returns list of the all available remote versions
func (golang Golang) ListRemote() (result []string, err error) {
	var (
		firstSelector = "#stable + div tr:first-of-type td:first-of-type.filename a"
		selector      = "#archive tr:first-of-type td:first-of-type.filename a"
	)

	doc, err := goquery.NewDocument(VersionLink)

	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, errors.New(variables.ConnectionError)
		}

		return nil, err
	}

	result = augmentSlice(result, doc.Find(firstSelector).Eq(0))

	doc.Find(selector).Each(func(i int, node *goquery.Selection) {
		result = augmentSlice(result, node)
	})

	return result, nil
}

func augmentSlice(result []string, node *goquery.Selection) []string {
	text := node.Text()

	if strings.Contains(text, "bootstrap") {
		return result
	}

	// We not checking for duplicates, since it just might create more errors
	version := rVersion.FindAllStringSubmatch(text, 1)[0][0]
	result = append(result, version)

	return result
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
