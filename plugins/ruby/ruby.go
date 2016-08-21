package ruby

import (
	"regexp"
	"runtime"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"net"
	"errors"

	"github.com/markelog/cprf"
	"github.com/markelog/release"
	"github.com/PuerkitoBio/goquery"

	"github.com/markelog/eclectica/variables"
)

var (
	VersionsLink = "https://rvm.io/binaries"
	home = fmt.Sprintf("%s/%s", variables.Home(), "ruby")
	bin = variables.Prefix() + "/bin/ruby"
)

type Ruby struct {}

func (ruby Ruby) Install(version string) error {
	var err error

	base := fmt.Sprintf("%s/%s", home, version)

	for _, file := range variables.Files {
		from := fmt.Sprintf("%s/%s", base, file)
		to := variables.Prefix()

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

func (ruby Ruby) Info(version string) (map[string]string, error) {
	result := make(map[string]string)

	sourcesUrl := fmt.Sprintf("%s/v%s", VersionsLink, version)

	result["name"] = "node"
	result["version"] = version
	result["filename"] = fmt.Sprintf("node-v%s-%s-x64", version, runtime.GOOS)
	result["url"] = fmt.Sprintf("%s/%s.tar.gz", sourcesUrl, result["filename"])

	return result, nil
}

func (ruby Ruby) Current() string {
	out, _ := exec.Command(bin, "--version").Output()
	version := strings.TrimSpace(string(out))

	return strings.Replace(version, "v", "", 1)
}

func (ruby Ruby) ListRemote() ([]string, error) {
	url := getUrl()
	doc, err := goquery.NewDocument(url)

	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, errors.New("Can't establish connection")
		}

		return nil, err
	}

	version := regexp.MustCompile("[[:digit:]]+\\.[[:digit:]]+\\.[[:digit:]]+")
	result := []string{}
	links := doc.Find("a")

	for i := range links.Nodes {
		href, _ := links.Eq(i).Attr("href")

		href = strings.Replace(href, "ruby-", "", 1)
		href = strings.Replace(href, ".tar.bz2", "", 1)

		if version.MatchString(href) {
			result = append(result, href)
		}
	}

	return result, nil
}

func getUrl() string {
	typa, _, version := release.All()
	arch := "x86_64"
	versions := strings.Split(version, ".")
	version = versions[0] + "." + versions[1]

	return fmt.Sprintf("%s/%s/%s/%s", VersionsLink, typa, version, arch)
}
