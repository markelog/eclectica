package info

import (
  "os"
  "fmt"
  "io/ioutil"

  "github.com/markelog/eclectica/plugins"
  "github.com/markelog/eclectica/cmd/prompt"
  "github.com/markelog/eclectica/variables"
)

func Ask() string {
  language := prompt.List("Language", plugins.List).Language
  version := AskVersion(language)

  return language + "@" + version
}

func AskVersion(language string) string {
  version := prompt.List("Version", Versions(language)).Version

  return version
}

func AskRemote() string {
  language := prompt.List("Language", plugins.List).Language
  version := AskRemoteVersion(language)

  return language + "@" + version
}

func AskRemoteVersion(language string) string {
  list, _ := plugins.RemoteList(language)
  key := prompt.List("Mask", plugins.GetKeys(list)).Mask
  versions := plugins.GetElements(key, list)
  version := prompt.List("Version", versions).Version

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
