// Package info provides ways to acquiring additional info from the user
package info

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/markelog/curse"

	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/list"
	"github.com/markelog/eclectica/plugins"
	"github.com/markelog/eclectica/versions"

	"github.com/markelog/eclectica/cmd/print/spinner"
)

// Ask for language and version from the user
func Ask() (language, version string, err error) {
	fmt.Println()

	language = list.List("langauge:", plugins.Plugins, 0)
	version, err = AskVersion(language)

	return
}

// AskVersion asks version from the user
func AskVersion(language string) (version string, err error) {
	vers := plugins.New(&plugins.Args{
		Language: language,
	}).List()

	if len(vers) == 0 {
		err = errors.New("There is no installed versions")
		return
	}

	version = list.List("version:", vers, 1)

	return
}

// AskRemote asks for remote version from the user
func AskRemote() (language, version string, err error) {
	fmt.Println()

	language = list.List("langauge:", plugins.Plugins, 0)
	version, err = AskRemoteVersion(language)

	return
}

// AskRemoteVersions asks for list of remote versions
func AskRemoteVersions(language string) (vers []string, err error) {
	remoteList, err := ListRemote(language)
	if err != nil {
		return
	}

	key := list.List("mask:", versions.GetKeys(remoteList), 4)
	vers = versions.GetElements(key, remoteList)

	return
}

// AskRemoteVersion asks for list of remote version
func AskRemoteVersion(language string) (version string, err error) {
	versions, err := AskRemoteVersions(language)

	if err != nil {
		return
	}

	version = list.List("version:", versions, 1)

	return
}

// PossibleLanguage gets the most possible language that user probably meant
func PossibleLanguage(args []string) (language string) {
	for _, element := range args {
		data := strings.Split(element, "@")
		language = data[0]

		if strings.Contains(language, "-") == false {
			return
		}
	}

	return ""
}

// GetLanguage gets supported language from args list
func GetLanguage(args []string) (language, version string) {
	for _, element := range args {
		data := strings.Split(element, "@")
		language = data[0]

		if len(data) == 2 {
			version = data[1]
		}

		for _, plugin := range plugins.Plugins {
			if language == plugin {
				return
			}
		}
	}

	return "", ""
}

// GetSpinner gets the spinner, just easier that way
func GetSpinner(language string, prefix spinner.Fn) *spinner.Spinner {
	c, _ := curse.New()

	before := func() {}

	postfix := func() {
		fmt.Println()
		time.Sleep(200 * time.Millisecond)
	}

	after := func() {
		c.MoveUp(1)
		c.EraseCurrentLine()
		print.InStyleln("langauge:", language)
	}

	return spinner.New(before, after, prefix, postfix)
}

// ListRemote lists remote version from the plugin indirectly
func ListRemote(language string) (versions map[string][]string, err error) {
	plugin := plugins.New(&plugins.Args{
		Language: language,
	})
	c, _ := curse.New()
	s := GetSpinner(language, func() {
		c.MoveUp(1)
		c.EraseCurrentLine()
		print.InStyle("langauge:", language)
	})

	s.Start()
	versions, err = plugin.ListRemote()
	s.Stop()

	return
}

// FullListRemote lists remote version from the plugin directly
func FullListRemote(language string) (versions []string, err error) {
	plugin := plugins.New(&plugins.Args{
		Language: language,
	})
	c, _ := curse.New()
	s := GetSpinner(language, func() {
		c.MoveUp(1)
		c.EraseCurrentLine()
		print.InStyle("langauge:", language)
	})

	s.Start()
	versions, err = plugin.Pkg.ListRemote()
	s.Stop()

	return
}

// HasLanguage do we have language in args list?
func HasLanguage(args []string) bool {
	language, _ := GetLanguage(args)

	return language != ""
}

// HasVersion do we have version in args list?
func HasVersion(args []string) bool {
	_, version := GetLanguage(args)

	return version != ""
}
