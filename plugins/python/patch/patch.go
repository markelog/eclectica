// Package patch provides logic for applying patches for the python installation
package patch

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-errors/errors"

	"github.com/markelog/eclectica/console"
	eIO "github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/variables"
	"github.com/markelog/eclectica/versions"
)

var (
	// Link holds link for list of all patches
	Link = "https://github.com/pyenv/pyenv/tree/master/plugins/python-build/share/python-build/patches"

	// RawLink is a part of the url for the patches so we can download them
	RawLink = "https://raw.githubusercontent.com/pyenv/pyenv/master/plugins/python-build/share/python-build/patches"
)

// URLs returns list of all needed patch urls
func URLs(version string) (urls []string, err error) {
	var (
		unsemVersion = versions.Unsemverify(version)
		link         = fmt.Sprintf("%s/%s/Python-%s", Link, unsemVersion, unsemVersion)
		rawLink      = fmt.Sprintf("%s/%s/Python-%s", RawLink, unsemVersion, unsemVersion)
	)

	doc, err := goquery.NewDocument(link)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, errors.New(variables.ConnectionError)
		}

		return nil, errors.New(err)
	}

	// github API has pretty low limits :/
	selector := ".files .js-navigation-item .content a[title]"
	doc.Find(selector).Each(func(i int, node *goquery.Selection) {
		content := node.Text()
		fullPath := fmt.Sprintf("%s/%s", rawLink, content)

		urls = append(urls, fullPath)
	})
	return
}

func getStrip(path string) string {
	r, _ := regexp.Compile("\\ndiff --git a/")

	text := eIO.Read(path)
	isDir := r.MatchString(text)

	if isDir {
		return "1"
	}

	return "0"
}

// Apply the patch
func Apply(path string) (err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return errors.New(err)
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

		if variables.IsDebug() {
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			return
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return errors.New(err)
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return errors.New(err)
		}

		err = cmd.Run()
		if err != nil {
			return console.Error(err, stdout, stderr)
		}
	}

	return
}
