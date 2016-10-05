package nodejs

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
	VersionsLink = "https://nodejs.org/dist"

	bins = []string{"node", "npm"}
)

type Node struct{}

func (node Node) Install(version string) error {
	return nil
}

func (node Node) Environment(version string) (string, error) {
	return "", nil
}

func (node Node) PostInstall(version string) error {
	return nil
}

func (node Node) Info(version string) (map[string]string, error) {
	result := make(map[string]string)
	sourcesUrl := fmt.Sprintf("%s/v%s", VersionsLink, version)

	result["version"] = version
	result["filename"] = fmt.Sprintf("node-v%s-%s-x64", version, runtime.GOOS)
	result["url"] = fmt.Sprintf("%s/%s.tar.gz", sourcesUrl, result["filename"])

	return result, nil
}

func (node Node) Bins() []string {
	return bins
}

func (node Node) Current() string {
	bin := variables.GetBin("node", "")
	out, _ := exec.Command(bin, "--version").Output()
	version := strings.TrimSpace(string(out))

	return strings.Replace(version, "v", "", 1)
}

func (node Node) ListRemote() ([]string, error) {
	doc, err := goquery.NewDocument(VersionsLink)

	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, errors.New("Can't establish connection")
		}

		return nil, err
	}

	tmp := []string{}
	result := []string{}
	version := regexp.MustCompile("v\\d+\\.\\d+\\.\\d+$")
	remove := regexp.MustCompile("0\\.[0-7]")

	links := doc.Find("a")

	for i := range links.Nodes {
		href, _ := links.Eq(i).Attr("href")

		href = strings.Replace(href, "/", "", 1)
		if version.MatchString(href) {
			href = strings.Replace(href, "v", "", 1)
			tmp = append(tmp, href)
		}
	}

	// Remove < 0.8 versions
	for _, element := range tmp {
		if remove.MatchString(element) == false {
			result = append(result, element)
		}
	}

	return result, nil
}
