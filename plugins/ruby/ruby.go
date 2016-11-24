package ruby

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chuckpreslar/emission"
	"github.com/markelog/release"

	"github.com/markelog/eclectica/variables"
)

var (
	VersionLink    = "https://s3.amazonaws.com/travis-rubies"
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
	return dealWithShell()
}

func (ruby Ruby) Environment() (string, error) {
	return "", nil
}

func (ruby Ruby) Info() (map[string]string, error) {
	result := make(map[string]string)

	result["filename"] = fmt.Sprintf("ruby-%s", ruby.Version)
	result["extension"] = "tar.bz2"
	result["url"] = fmt.Sprintf("%s/%s.%s", getUrl(ruby.Version), result["filename"], result["extension"])

	return result, nil
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

func (ruby Ruby) ListRemote() (result []string, err error) {
	doc, err := goquery.NewDocument(VersionLink)

	if err != nil {
		if _, ok := err.(net.Error); ok {
			err = errors.New("Can't establish connection")
		}

		return
	}

	var (
		regPart  = getRegUrl() + "\\/ruby-" + versionPattern + ".tar.bz"
		rPath    = regexp.MustCompile(regPart)
		rVersion = regexp.MustCompile(versionPattern)
	)

	doc.Find("Key").Each(func(i int, node *goquery.Selection) {
		value := node.Text()

		if rPath.MatchString(value) {
			result = append(result, rVersion.FindAllStringSubmatch(value, 1)[0][0])
		}
	})

	return result, nil
}

func getRelativePath() (typa, version, arch string) {
	typa, _, version = release.All()
	arch = "x86_64"

	versions := strings.Split(version, ".")
	version = versions[0] + "." + versions[1]

	return
}

func getUrl(version string) string {
	typa, version, arch := getRelativePath()

	return VersionLink + "/binaries/" + typa + "/" + version + "/" + arch
}

func getRegUrl() string {
	typa, version, arch := getRelativePath()
	return fmt.Sprintf("binaries\\/%s\\/%s\\/%s", typa, version, arch)
}
