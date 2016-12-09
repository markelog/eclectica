package ruby

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chuckpreslar/emission"
	"github.com/markelog/release"

	"github.com/markelog/eclectica/variables"
)

var (
	VersionsLink   = "https://rvm.io/binaries"
	versionPattern = "\\d+\\.\\d+\\.\\d"

	bins = []string{"erb", "gem", "irb", "rake", "rdoc", "ri", "ruby"}
	dots = []string{".ruby-version"}
)

type Ruby struct {
	Version string
	Emitter *emission.Emitter
}

func (ruby Ruby) Events() *emission.Emitter {
	return ruby.Emitter
}

func (ruby Ruby) PreInstall() error {
	return nil
}

func (ruby Ruby) Install() error {
	return nil
}

func (ruby Ruby) PostInstall() error {
	err := removeRVMArtefacts(variables.Path("ruby", ruby.Version))
	if err != nil {
		return err
	}

	return dealWithShell()
}

// Removes RVM artefacts (ignore errors)
func removeRVMArtefacts(base string) error {
	gems := filepath.Join(base, "lib/ruby/gems")

	// Remove `cache` folder since it supposed to work with RVM cache
	folders, _ := ioutil.ReadDir(gems)
	for _, folder := range folders {
		err := os.RemoveAll(filepath.Join(gems, folder.Name(), "cache"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (ruby Ruby) Environment() (result []string, err error) {
	return
}

func (ruby Ruby) Info() map[string]string {
	result := make(map[string]string)

	result["filename"] = fmt.Sprintf("ruby-%s", ruby.Version)
	result["extension"] = "tar.bz2"
	result["url"] = fmt.Sprintf("%s/%s.%s", getUrl(), result["filename"], result["extension"])

	return result
}

func (ruby Ruby) Bins() []string {
	return bins
}

func (ruby Ruby) Dots() []string {
	return dots
}

func (ruby Ruby) Current() string {
	bin := variables.GetBin("ruby")
	out, _ := exec.Command(bin, "--version").Output()

	if len(out) == 0 {
		return ""
	}

	version := strings.TrimSpace(string(out))
	rVersion := regexp.MustCompile(versionPattern)
	testVersion := rVersion.FindAllStringSubmatch(version, 1)

	if len(testVersion) == 0 {
		return ""
	}

	return testVersion[0][0]
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

	version := regexp.MustCompile("\\d+\\.\\d+\\.\\d+")
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
