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
	"github.com/markelog/cprf"
	"github.com/markelog/release"

	"github.com/markelog/eclectica/variables"
)

var (
	VersionsLink   = "https://rvm.io/binaries"
	home           = filepath.Join(variables.Home(), "ruby")
	bin            = filepath.Join(variables.Prefix("ruby"), "/bin/ruby")
	Bins           = []string{"erb", "gem", "irb", "rake", "rdoc", "ri", "ruby"}
	versionPattern = "\\d+\\.\\d+\\.\\d"
)

type Ruby struct{}

func (ruby Ruby) Environment(version string) (string, error) {
	return "", nil
}

func (ruby Ruby) Install(version string) error {
	var err error
	rVersion := regexp.MustCompile(versionPattern)
	version = rVersion.FindAllStringSubmatch(version, 1)[0][0]

	base := filepath.Join(home, version)

	removeRVMArtefacts(base)

	for _, file := range variables.Files {
		from := filepath.Join(base, file)
		to := variables.Prefix("ruby")

		// Some versions might not have certain files
		if _, statError := os.Stat(from); os.IsNotExist(statError) {
			continue
		}

		err = cprf.Copy(from, to)

		if err != nil {
			return err
		}
	}

	return nil
}

// Removes RVM artefacts (ignore errors)
func removeRVMArtefacts(base string) {
	gems := filepath.Join(base, "lib/ruby/gems")

	// Remove `cache` folder since it supposed to work with RVM cache
	folders, _ := ioutil.ReadDir(gems)
	for _, folder := range folders {
		os.Remove(filepath.Join(gems, folder.Name(), "cache"))
	}
}

func (ruby Ruby) PostInstall(version string) (bool, error) {
	return dealWithShell()
}

func (ruby Ruby) Info(version string) (map[string]string, error) {
	result := make(map[string]string)

	result["version"] = version
	result["filename"] = fmt.Sprintf("ruby-%s", version)
	result["extension"] = "tar.bz2"
	result["url"] = fmt.Sprintf("%s/%s.%s", getUrl(), result["filename"], result["extension"])

	return result, nil
}

func (ruby Ruby) Bins() []string {
	return Bins
}

func (ruby Ruby) Current() string {
	out, _ := exec.Command(bin, "--version").Output()
	version := strings.TrimSpace(string(out))

	rVersion := regexp.MustCompile(versionPattern)
	version = rVersion.FindAllStringSubmatch(version, 1)[0][0]

	return version
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
