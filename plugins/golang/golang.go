package golang

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/markelog/cprf"

	"github.com/markelog/eclectica/variables"
)

var (
	VersionsLink = "https://storage.googleapis.com/golang"
	home         = fmt.Sprintf("%s/%s", variables.Home(), "go")
	bin          = variables.Prefix("go") + "/bin/go"
	files        = [8]string{"api", "bin", "lib", "misc", "pkg", "share", "src"}

	versionPattern = "[[:digit:]]+\\.[[:digit:]]+(?:\\.[[:digit:]]+)?(?:(alpha|beta|rc)(?:[[:digit:]]*)?)?"
)

type Golang struct{}

func (golang Golang) Install(version string) error {
	var err error

	base := fmt.Sprintf("%s/%s", home, version)

	for _, element := range files {
		from := fmt.Sprintf("%s/%s", base, element)
		to := variables.Prefix("go")

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
	out, _ := exec.Command(bin, "--version").Output()
	version := strings.TrimSpace(string(out))

	return strings.Replace(version, "v", "", 1)
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
