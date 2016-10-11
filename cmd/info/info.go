package info

import (
	"fmt"
	"strings"
	"time"

	"github.com/markelog/curse"
	"github.com/markelog/list"

	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/plugins"
	"github.com/markelog/eclectica/variables"
)

// Ask for language and version from the user
func Ask() (language, version string, err error) {
	language = list.GetWith("Language", plugins.Plugins)
	version, err = AskVersion(language)

	return
}

// Ask for version
func AskVersion(language string) (version string, err error) {
	versions, err := plugins.New(language).List()
	if err != nil {
		return
	}

	version = list.GetWith("Version", versions)

	return
}

// Ask for language and remote version
func AskRemote() (language, version string, err error) {
	language = list.GetWith("Language", plugins.Plugins)
	version, err = AskRemoteVersion(language)

	return
}

// Ask for list of remote versions
func AskRemoteVersions(language string) (versions []string, err error) {
	plugin := plugins.New(language)
	c, _ := curse.New()

	prefix := func() {
		c.MoveUp(1)
		c.EraseCurrentLine()
		print.InStyle("Language", language)
	}

	postfix := func() {
		fmt.Println()
		time.Sleep(200 * time.Millisecond)
	}

	s := &print.Spinner{
		Before:  func() { time.Sleep(500 * time.Millisecond) },
		After:   func() { fmt.Println() },
		Prefix:  prefix,
		Postfix: postfix,
	}

	s.Start()
	remoteList, err := plugin.ListRemote()
	s.Stop()

	if err != nil {
		return
	}

	key := list.GetWith("Mask", plugins.GetKeys(remoteList))
	versions = plugins.GetElements(key, remoteList)

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
