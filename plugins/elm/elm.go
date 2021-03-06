// Package elm provides all needed logic for installation of Elm
package elm

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/blang/semver"
	"github.com/chuckpreslar/emission"
	"github.com/go-errors/errors"
	"github.com/markelog/cprf"

	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/pkg"
	"github.com/markelog/eclectica/variables"
)

var (
	// VersionLink is the URL link from which we can get all possible versions
	VersionLink = "https://dl.bintray.com/elmlang/elm-platform"

	versionPattern = "\\d+\\.\\d+\\.\\d+"

	diffFolderBinaryName, _ = semver.Make("0.17.1")

	bins = []string{"elm", "elm-make", "elm-package", "elm-reactor", "elm-repl"}
	dots = []string{".elm-version"}
)

// Elm essential struct
type Elm struct {
	Version string
	Emitter *emission.Emitter
	pkg.Base
}

// New returns language struct
func New(version string, emitter *emission.Emitter) *Elm {
	return &Elm{
		Version: version,
		Emitter: emitter,
	}
}

// Events returns language related event emitter
func (elm Elm) Events() *emission.Emitter {
	return elm.Emitter
}

// PreDownload hook
func (elm Elm) PreDownload() (err error) {
	path := elm.getTmpPath()

	if _, errStat := os.Stat(path); os.IsNotExist(errStat) {
		_, err = io.CreateDir(path)
	}

	return
}

// PostInstall hook
func (elm Elm) PostInstall() error {
	path := variables.Path("elm", elm.Version)
	binPath := filepath.Join(path, "bin")

	_, err := io.CreateDir(binPath)
	if err != nil {
		return err
	}

	for _, name := range elm.Bins() {
		currentBinPath := filepath.Join(path, name)

		cprf.Copy(currentBinPath, binPath)

		err = os.RemoveAll(currentBinPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// Info provides all the info needed for installation of the plugin
func (elm Elm) Info() map[string]string {
	var (
		result     = make(map[string]string)
		sourcesURL = fmt.Sprintf("%s/%s", VersionLink, elm.Version)
		chosen, _  = semver.Make(elm.Version)
	)

	// Man, why?!
	if chosen.LT(diffFolderBinaryName) {
		result["unarchive-filename"] = "dist_binaries"
	}
	if runtime.GOOS == "linux" {
		result["unarchive-filename"] = "dist_binaries"
	}
	if elm.Version == "0.15.1" && runtime.GOOS == "darwin" {
		result["unarchive-filename"] = "osx"
	}
	if elm.Version == "0.15.1" && runtime.GOOS == "linux" {
		result["unarchive-filename"] = "linux64"
	}

	result["filename"] = fmt.Sprintf("%s-x64", runtime.GOOS)
	result["url"] = fmt.Sprintf("%s/%s.tar.gz", sourcesURL, result["filename"])
	result["archive-folder"] = filepath.Join(variables.TempDir(), "elm-archive-"+elm.Version) + "/"

	return result
}

// Bins returns list of the all bins included
// with the distribution of the language
func (elm Elm) Bins() []string {
	return bins
}

// Dots returns list of the all available filenames
// which can define versions
func (elm Elm) Dots() []string {
	return dots
}

// ListRemote returns list of the all available remote versions
func (elm Elm) ListRemote() ([]string, error) {
	doc, err := goquery.NewDocument(VersionLink)

	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, errors.New(variables.ConnectionError)
		}

		return nil, err
	}

	result := []string{}
	rVersion := regexp.MustCompile(versionPattern)

	doc.Find("a").Each(func(i int, elm *goquery.Selection) {
		value := elm.Text()

		if rVersion.MatchString(value) {
			result = append(result, strings.Replace(value, "/", "", 1))
		}
	})

	return result, nil
}

func (elm Elm) getTmpPath() string {
	return filepath.Join(variables.TempDir(), "elm-archive-"+elm.Version) + "/"
}
