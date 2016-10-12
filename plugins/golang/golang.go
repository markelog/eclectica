package golang

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/markelog/eclectica/variables"
)

var (
	VersionsLink   = "https://storage.googleapis.com/golang"
	versionPattern = "\\d+\\.\\d+(?:\\.\\d+)?(?:(alpha|beta|rc)(?:\\d*)?)?"

	bins = []string{"go", "godoc", "gofmt"}
)

type Golang struct {
	Version string
}

func (golang Golang) Install() error {
	return nil
}

func (golang Golang) PostInstall() error {
	return nil
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

	result["unarchive-filename"] = "go"
	result["filename"] = fmt.Sprintf("go%s.%s", golang.Version, platform)
	result["url"] = fmt.Sprintf("%s/%s.tar.gz", VersionsLink, result["filename"])

	return result, nil
}

func (rust Golang) Bins() []string {
	return bins
}

func (golang Golang) Current() string {
	bin := variables.GetBin("go")
	out, _ := exec.Command(bin, "version").Output()

	rVersion := regexp.MustCompile(versionPattern)
	version := strings.TrimSpace(string(out))
	testVersion := rVersion.FindAllStringSubmatch(version, 1)

	if len(testVersion) == 0 {
		return ""
	}

	return testVersion[0][0]
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
