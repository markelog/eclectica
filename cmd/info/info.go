package info

import (
  "strings"

  "github.com/markelog/list"

  "github.com/markelog/eclectica/plugins"
  "github.com/markelog/eclectica/variables"
)

func Ask() (language, version string) {
  language = list.GetWith("Language", plugins.Plugins)
  version = AskVersion(language)

  return
}

func AskVersion(language string) (version string) {
  version = list.GetWith("Version", plugins.New(language).List())

  return
}

func AskRemote() (language, version string) {
  language = list.GetWith("Language", plugins.Plugins)
  version = AskRemoteVersion(language)

  return
}

func AskRemoteVersion(language string) (version string) {
  remoteList, _ := plugins.New(language).ListRemote()
  key := list.GetWith("Mask", plugins.GetKeys(remoteList))
  versions := plugins.GetElements(key, remoteList)
  version = list.GetWith("Version", versions)

  return
}

func GetLanguage(args []string) (language, version string) {
  for _, element := range args {
    data := strings.Split(element , "@")
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

func GetCommand(args []string) string {
  for _, element := range args {
    for _, command := range variables.Commands {
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

func HasCommand(args []string) bool {
  return GetCommand(args) != ""
}
