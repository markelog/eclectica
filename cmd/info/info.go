package info

import (
  "os"
  "fmt"
  "io/ioutil"

  "github.com/markelog/list"

  "github.com/markelog/eclectica/plugins"
  "github.com/markelog/eclectica/variables"
)

func Ask() string {
  language := list.GetWith("Language", plugins.List)
  version := AskVersion(language)

  return language + "@" + version
}

func AskVersion(language string) string {
  version := list.GetWith("Version", Versions(language))

  return version
}

func AskRemote() string {
  language := list.GetWith("Language", plugins.List)
  version := AskRemoteVersion(language)

  return language + "@" + version
}

func AskRemoteVersion(language string) string {
  remoteList, _ := plugins.RemoteList(language)
  key := list.GetWith("Mask", plugins.GetKeys(remoteList))
  versions := plugins.GetElements(key, remoteList)
  version := list.GetWith("Version", versions)

  return language + "@" + version
}

func Versions(name string) []string {
  path := variables.Home + "/" + name

  if _, err := os.Stat(path); os.IsNotExist(err) {
    fmt.Println("There is no installed versions of " + name)
    os.Exit(1)
  }

  folders, err := ioutil.ReadDir(variables.Home + "/" + name)
  versions := []string{}

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  for _, folder := range folders {
    if folder.IsDir() {
      versions = append(versions, folder.Name())
    }
  }

  return versions
}
