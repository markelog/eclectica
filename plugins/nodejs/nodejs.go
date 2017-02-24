package nodejs

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
	"github.com/chuckpreslar/emission"

	"github.com/markelog/eclectica/variables"
	"github.com/markelog/eclectica/io"
)

var (
	VersionLink    = "https://nodejs.org/dist"
	versionPattern = "v\\d+\\.\\d+\\.\\d+$"
	removePattern  = "0\\.[0-7]"

	bins = []string{"node", "npm"}
	dots = []string{".nvmrc", ".node-version"}
)

type Node struct {
	Version string
	Emitter *emission.Emitter
}

func (node Node) Events() *emission.Emitter {
	return node.Emitter
}

func (node Node) PreInstall() error {
	return nil
}

func (node Node) Install() error {
	return nil
}

func (node Node) PostInstall() error {
	path := variables.Path("node", node.Version)
	etc := filepath.Join(path, "etc")
	npmrc := filepath.Join(etc, "npmrc")

	// Remove needless warnings from npm output
	io.CreateDir(etc)
	io.WriteFile(npmrc, "scripts-prepend-node-path=false")

	return nil
}

func (node Node) Environment() (result []string, err error) {
	return
}

func (node Node) Info() map[string]string {
	result := make(map[string]string)
	sourcesUrl := fmt.Sprintf("%s/v%s", VersionLink, node.Version)

	result["filename"] = fmt.Sprintf("node-v%s-%s-x64", node.Version, runtime.GOOS)
	result["url"] = fmt.Sprintf("%s/%s.tar.gz", sourcesUrl, result["filename"])

	return result
}

func (node Node) Bins() []string {
	return bins
}

func (node Node) Dots() []string {
	return dots
}

func (node Node) Current() string {
	bin := variables.GetBin("node")
	out, _ := exec.Command(bin, "--version").Output()

	version := strings.TrimSpace(string(out))

	return strings.Replace(version, "v", "", 1)
}

func (node Node) ListRemote() ([]string, error) {
	doc, err := goquery.NewDocument(VersionLink)

	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, errors.New("Can't establish connection")
		}

		return nil, err
	}

	tmp := []string{}
	result := []string{}

	rVersion := regexp.MustCompile(versionPattern)
	rRemove := regexp.MustCompile(removePattern)

	doc.Find("a").Each(func(i int, node *goquery.Selection) {
		href, _ := node.Attr("href")

		href = strings.Replace(href, "/", "", 1)
		if rVersion.MatchString(href) {
			href = strings.Replace(href, "v", "", 1)
			tmp = append(tmp, href)
		}
	})

	// Remove < 0.8 versions
	for _, element := range tmp {
		if rRemove.MatchString(element) == false {
			result = append(result, element)
		}
	}

	return result, nil
}
