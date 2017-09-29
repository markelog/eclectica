package nodejs

import (
	"fmt"
	"net"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/blang/semver"
	"github.com/chuckpreslar/emission"
	"github.com/go-errors/errors"

	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/pkg"
	"github.com/markelog/eclectica/plugins/nodejs/modules"
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
	pkg.Base
}

func New(version string, emitter *emission.Emitter) *Node {
	return &Node{
		Version:  version,
		Emitter:  emitter,
		previous: variables.CurrentVersion("node"),
	}
}

func (node Node) Events() *emission.Emitter {
	return node.Emitter
}

func (node Node) PostInstall() (err error) {
	node.Emitter.Emit("post-install")

	err = node.setNpm()
	if err != nil {
		return errors.New(err)
	}

	ok, err := node.Yarn()
	if err != nil && ok == false {
		return errors.New(err)
	}

	return nil
}

func (node Node) Switch() (err error) {
	previous := node.previous

	if len(previous) == 0 {
		return
	}

	node.Emitter.Emit("reapply modules")

	err = modules.New(node.previous, node.Version).Install()
	if err != nil {
		return
	}

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

func (node Node) ListRemote() ([]string, error) {
	doc, err := goquery.NewDocument(VersionLink)

	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, errors.New("Connection cannot be established")
		}

		return nil, errors.New(err)
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
