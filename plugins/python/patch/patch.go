package patch

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

	"github.com/markelog/eclectica/console"
	eio "github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/variables"
	"github.com/markelog/eclectica/versions"
)

var (
	Link    = "https://github.com/yyuu/pyenv/tree/master/plugins/python-build/share/python-build/patches"
	RawLink = "https://raw.githubusercontent.com/yyuu/pyenv/master/plugins/python-build/share/python-build/patches"
)

func Urls(version string) (result []string, err error) {
	unsemVersion := versions.Unsemverify(version)
	link := fmt.Sprintf("%s/%s/Python-%s", Link, unsemVersion, unsemVersion)
	rawLink := fmt.Sprintf("%s/%s/Python-%s", RawLink, unsemVersion, unsemVersion)

	doc, err := goquery.NewDocument(link)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, errors.New("Can't establish connection")
		}

		return
	}

	// github API has pretty low limits :/
	selector := ".files .js-navigation-item .content a[title]"
	doc.Find(selector).Each(func(i int, node *goquery.Selection) {
		content := node.Text()
		fullPath := fmt.Sprintf("%s/%s", rawLink, content)

		result = append(result, fullPath)
	})
	return
}

func getStrip(path string) string {
	r, _ := regexp.Compile("\\ndiff --git a/")

	text := eio.Read(path)
	isDir := r.MatchString(text)

	if isDir {
		return "1"
	}

	return "0"
}

func Apply(path string) (err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, file := range files {
		name := file.Name()

		if (strings.HasSuffix(name, ".patch") || strings.HasSuffix(name, ".diff")) == false {
			continue
		}

		// There should be port of `patch` command to golang right?
		target := filepath.Join(path, name)
		strip := getStrip(target)

		cmd := exec.Command("patch", "-p", strip, "--force", "-i", target)
		cmd.Dir = path
		stdErr, _ := cmd.StderrPipe()
		stdOut, _ := cmd.StdoutPipe()

		if variables.IsDebug() {
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
		}

		err = cmd.Run()
		if err != nil {
			return console.GetError(err, stdErr, stdOut)
		}
	}

	return
}
