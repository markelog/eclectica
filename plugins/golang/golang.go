package golang

import (
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chuckpreslar/emission"

	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/variables"
	"github.com/markelog/eclectica/versions"
)

var (
	VersionsLink   = "https://storage.googleapis.com/golang"
	versionPattern = "\\d+\\.\\d+(?:\\.\\d+)?(?:(alpha|beta|rc)(?:\\d*)?)?"

	bins = []string{"go", "godoc", "gofmt"}
	dots = []string{".go-version"}
)

type Golang struct {
	Version string
	Emitter *emission.Emitter
}

func (golang Golang) Events() *emission.Emitter {
	return golang.Emitter
}

func (golang Golang) PreInstall() error {
	return nil
}

func (golang Golang) Install() error {
	return nil
}

func (golang Golang) PostInstall() error {
	return dealWithShell()
}

func (golang Golang) Environment() (string, error) {
	return "GOROOT=" + variables.Path("go", golang.Version), nil
}

func (golang Golang) Info() (map[string]string, error) {
	result := make(map[string]string)

	platform, err := getPlatform()
	if err != nil {
		return nil, err
	}

	version := versions.Unsemverify(golang.Version)
	result["version"] = version
	result["unarchive-filename"] = "go"
	result["filename"] = fmt.Sprintf("go%s.%s", version, platform)
	result["url"] = fmt.Sprintf("%s/%s.tar.gz", VersionsLink, result["filename"])

	return result, nil
}

func (rust Golang) Bins() []string {
	return bins
}

func (rust Golang) Dots() []string {
	return dots
}

func (golang Golang) Current() string {
	path := variables.Path("go")
	version := filepath.Join(path, "VERSION")
	readVersion := strings.Replace(io.Read(version), "go", "", 1)
	semverVersion := versions.Semverify(readVersion)

	return semverVersion
}

func (golang Golang) ListRemote() ([]string, error) {
	doc, err := goquery.NewDocument(VersionsLink)

	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, errors.New("Can't establish connection")
		}

		return nil, err
	}

	result := []string{}
	rVersion := regexp.MustCompile(versionPattern)

	links := doc.Find("Key")
	platform, err := getPlatform()
	if err != nil {
		return nil, err
	}
	platform += "\\.tar\\.gz$"
	rPlatform := regexp.MustCompile(platform)

	for i := range links.Nodes {
		value := links.Eq(i).Text()

		if rPlatform.MatchString(value) {
			result = append(result, rVersion.FindAllStringSubmatch(value, 1)[0][0])
		}
	}

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
