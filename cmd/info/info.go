package info

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/markelog/curse"
	"github.com/markelog/list"

	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/plugins"
	"github.com/markelog/eclectica/variables"
	"github.com/markelog/eclectica/versions"

	"github.com/markelog/eclectica/cmd/print/spinner"
)

type prefixFn func()

// Ask for language and version from the user
func Ask() (language, version string, err error) {
	language = list.GetWith("Language", plugins.Plugins)
	version, err = AskVersion(language)

	return
}

// Ask for version
func AskVersion(language string) (version string, err error) {
	vers := plugins.New(&plugins.Args{
		Language: language,
	}).List()

	if len(vers) == 0 {
		err = errors.New("There is no installed versions")
		return
	}

	version = list.GetWith("Version", vers)

	return
}

// Ask for language and remote version
func AskRemote() (language, version string, err error) {
	language = list.GetWith("Language", plugins.Plugins)
	version, err = AskRemoteVersion(language)

	return
}

// Ask for list of remote versions
func AskRemoteVersions(language string) (vers []string, err error) {
	remoteList, err := ListRemote(language)
	if err != nil {
		return
	}

	key := list.GetWith("Mask", versions.GetKeys(remoteList))
	vers = versions.GetElements(key, remoteList)

	return
}

// Ask for list of remote version
func AskRemoteVersion(language string) (version string, err error) {
	versions, err := AskRemoteVersions(language)

	if err != nil {
		return
	}

	version = list.GetWith("Version", versions)

	return
}

// Get supported language from args list
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

// Get command from args list
func GetCommand(args []string) string {
	for _, element := range args {
		for _, command := range variables.NonInstallCommands {
			if command == element {
				return command
			}
		}
	}

	return ""
}

func GetSpinner(language string, prefix spinner.SpinnerFn) *spinner.Spinner {
	c, _ := curse.New()

	before := func() {}

	postfix := func() {
		fmt.Println()
		time.Sleep(200 * time.Millisecond)
	}

	after := func() {
		c.MoveUp(1)
		c.EraseCurrentLine()
		print.InStyleln("Language", language)
	}

	return spinner.New(before, after, prefix, postfix)
}

func ListRemote(language string) (versions map[string][]string, err error) {
	plugin := plugins.New(&plugins.Args{
		Language: language,
	})
	c, _ := curse.New()
	s := GetSpinner(language, func() {
		c.MoveUp(1)
		c.EraseCurrentLine()
		print.InStyle("Language", language)
	})

	s.Start()
	versions, err = plugin.ListRemote()
	s.Stop()

	return
}

func FullListRemote(language string) (versions []string, err error) {
	plugin := plugins.New(&plugins.Args{
		Language: language,
	})
	c, _ := curse.New()
	s := GetSpinner(language, func() {
		c.MoveUp(1)
		c.EraseCurrentLine()
		print.InStyle("Language", language)
	})

	s.Start()
	versions, err = plugin.Pkg.ListRemote()
	s.Stop()

	return
}

// Is there an langauge in args list?
func HasLanguage(args []string) bool {
	language, _ := GetLanguage(args)

	return language != ""
}

// Is there an version in args list?
func HasVersion(args []string) bool {
	_, version := GetLanguage(args)

	return version != ""
}

// Is this is non-install command in args list?
func NonInstallCommand(args []string) bool {
	return GetCommand(args) != ""
}
