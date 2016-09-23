package golang

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/markelog/cprf"

	"github.com/markelog/eclectica/directory"
	"github.com/markelog/eclectica/variables"
)

var (
	VersionsLink = "https://storage.googleapis.com/golang"
	home         = filepath.Join(variables.Home(), "go")
	bin          = filepath.Join(variables.Prefix("go"), "/bin/go")
	files        = [8]string{"api", "bin", "lib", "misc", "pkg", "share", "src"}

	versionPattern = "\\d+\\.\\d+(?:\\.\\d+)?(?:(alpha|beta|rc)(?:\\d*)?)?"
)

type Golang struct{}

func (golang Golang) Install(version string) error {
	var err error

	base := filepath.Join(home, version)
	to := filepath.Join(variables.Prefix("go"), "go")

	files, err := ioutil.ReadDir(base)
	if err != nil {
		return err
	}

	// Remove everything in GOROOT dir in case there was previous versions installed there
	os.RemoveAll(to)

	// Re-create GOROOT
	_, err = directory.Create(to)
	if err != nil {
		return err
	}

	// Copy to GOROOT
	for _, element := range files {
		from := filepath.Join(base, element.Name())

		err = cprf.Copy(from, to)
		if err != nil {
			return err
		}
	}

	to = variables.Prefix("go")
	for _, element := range variables.Files {
		from := filepath.Join(base, element)

		// Some versions might not have certain files
		if _, err := os.Stat(from); os.IsNotExist(err) {
			continue
		}

		err = cprf.Copy(from, to)
		if err != nil {
			return err
		}
	}

	return nil
}

func (golang Golang) PostInstall() (bool, error) {
	return dealWithRc()
}

func (golang Golang) Info(version string) (map[string]string, error) {
	result := make(map[string]string)

	platform, err := getPlatform()
	if err != nil {
		return nil, err
	}

	result["version"] = version
	result["filename"] = fmt.Sprintf("go%s.%s", version, platform)
	result["url"] = fmt.Sprintf("%s/%s.tar.gz", VersionsLink, result["filename"])
	result["unarchive-filename"] = "go"

	return result, nil
}

func (golang Golang) Current() string {
	rVersion := regexp.MustCompile(versionPattern)
	out, _ := exec.Command(bin, "version").Output()
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
