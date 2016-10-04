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

func Ask() (language, version string, err error) {
	language = list.GetWith("Language", plugins.Plugins)
	version, err = AskVersion(language)

	return
}

func AskVersion(language string) (version string, err error) {
	versions, err := plugins.New(language).List()
	if err != nil {
		return
	}

	version = list.GetWith("Version", versions)

	return
}

func AskRemote() (language, version string, err error) {
	language = list.GetWith("Language", plugins.Plugins)
	version, err = AskRemoteVersion(language)

	return
}

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

func AskRemoteVersion(language string) (version string, err error) {
	versions, err := AskRemoteVersions(language)

	if err != nil {
		return
	}

	version = list.GetWith("Version", versions)

	return
}

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

	return
}

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

func HasLanguage(args []string) bool {
	language, _ := GetLanguage(args)

	return language != ""
}

func HasVersion(args []string) bool {
	_, version := GetLanguage(args)

	return version != ""
}

func NonInstallCommand(args []string) bool {
	return GetCommand(args) != ""
}
