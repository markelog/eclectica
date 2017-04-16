package nodejs

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/blang/semver"
	"github.com/chuckpreslar/emission"

	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/variables"
)

var (
	VersionLink    = "https://nodejs.org/dist"
	versionPattern = "v\\d+\\.\\d+\\.\\d+$"

	minimalVersion, _ = semver.Make("0.10.0")

	bins = []string{"node", "npm"}
	dots = []string{".nvmrc", ".node-version"}
)

type Node struct {
	Version  string
	previous string
	Emitter  *emission.Emitter
}

func (node Node) Events() *emission.Emitter {
	return node.Emitter
}

func (node Node) PreDownload() (err error) {
	return
}

func (node *Node) PreInstall() (err error) {
	node.previous = node.Current()

	return
}

func (node Node) Install() error {
	return nil
}

func (node Node) PostInstall() (err error) {
	err = node.setNpm()
	if err != nil {
		return err
	}

	ok, err := node.Yarn()
	if err != nil && ok == false {
		return err
	}

	return nil
}

func (node Node) Switch() (err error) {
	previous := node.previous

	if len(previous) == 0 {
		return
	}

	node.Emitter.Emit("configure")

	modulesPath := node.modulesPath(previous)

	if node.sameMajors() {
		return node.copyModules(modulesPath)
	}

	return node.installModules(modulesPath)
}

func (node Node) Link() (err error) {
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

	doc.Find("a").Each(func(i int, node *goquery.Selection) {
		href, _ := node.Attr("href")

		href = strings.Replace(href, "/", "", 1)
		if rVersion.MatchString(href) {
			href = strings.Replace(href, "v", "", 1)
			tmp = append(tmp, href)
		}
	})

	// Remove outdated versions
	for _, element := range tmp {
		version, _ := semver.Make(element)

		if version.GTE(minimalVersion) {
			result = append(result, element)
		}
	}

	return result, nil
}

// Removes needless warnings from npm output
func (node Node) setNpm() (err error) {
	path := variables.Path("node", node.Version)
	etc := filepath.Join(path, "etc")
	npmrc := filepath.Join(etc, "npmrc")

	_, err = io.CreateDir(etc)
	if err != nil {
		return
	}

	err = io.WriteFile(npmrc, "scripts-prepend-node-path=false")
	if err != nil {
		return
	}

	return nil
}
