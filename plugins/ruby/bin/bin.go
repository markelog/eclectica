package bin

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chuckpreslar/emission"
	"github.com/go-errors/errors"

	"github.com/markelog/eclectica/plugins/ruby/base"
	"github.com/markelog/eclectica/plugins/ruby/rvm"
	"github.com/markelog/eclectica/variables"
)

var (
	VersionLink    = "https://rvm.io/binaries"
	versionPattern = "\\d+\\.\\d+\\.\\d"

	bins = []string{"erb", "gem", "irb", "rake", "rdoc", "ri", "ruby"}
	dots = []string{".ruby-version"}
)

type Ruby struct {
	Version string
	Emitter *emission.Emitter
	base.Ruby
}

func New(version string, emitter *emission.Emitter) *Ruby {
	return &Ruby{
		Version: version,
		Emitter: emitter,
	}
}

func (ruby Ruby) Events() *emission.Emitter {
	return ruby.Emitter
}

func (ruby Ruby) PostInstall() error {
	err := removeRVMArtefacts(variables.Path("ruby", ruby.Version))
	if err != nil {
		return errors.New(err)
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
			return errors.New(err)
		}
	}

	return nil
}
func (ruby Ruby) Info() map[string]string {
	result := make(map[string]string)

	result["filename"] = fmt.Sprintf("ruby-%s", ruby.Version)
	result["extension"] = "tar.bz2"
	result["url"] = fmt.Sprintf("%s/%s.%s", rvm.GetUrl(VersionLink), result["filename"], result["extension"])

	return result
}

func (ruby Ruby) ListRemote() ([]string, error) {
	url := rvm.GetUrl(VersionLink)
	doc, err := goquery.NewDocument(url)

	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, errors.New(variables.ConnectionError)
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
