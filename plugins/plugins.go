package plugins

import (
  "os"
  "io/ioutil"

  "github.com/markelog/eclectica/plugins/nodejs"
  "github.com/markelog/eclectica/plugins/rust"
  "github.com/markelog/eclectica/variables"
)

var (
  List = []string{
    "node",
    "rust",
  }
)

func Versions(name string) (versions []string) {
  versions = []string{}
  path := variables.Home + "/" + name

  if _, err := os.Stat(path); os.IsNotExist(err) {
    return
  }

  folders, _ := ioutil.ReadDir(path)
  for _, folder := range folders {
    versions = append(versions, folder.Name())
  }

  return
}

func Version(langauge, version string) (map[string]string, error) {
  var (
    info map[string]string
    err error
  )

  switch {
    case langauge == "node":
      info, err = nodejs.Version(version)
    case langauge == "rust":
      info, err = rust.Version(version)
  }

  return info, err
}

func Remove(langauge, version string) error {
  switch {
    case langauge == "node":
      return nodejs.Remove(version)
    case langauge == "rust":
      return rust.Remove(version)
  }

  return nil
}

func Activate(data map[string]string) error {
  switch {
    case data["name"] == "node":
      return nodejs.Activate(data)
    case data["name"] == "rust":
      return rust.Activate(data)
  }

  return nil
}

func RemoteList(langauge string) (map[string][]string, error) {
  var (
    versions []string
    err error
  )

  switch {
    case langauge == "node":
      versions, err = nodejs.ListVersions()
    case langauge == "rust":
      versions, err = rust.ListVersions()
  }

  if err != nil {
    return nil, err
  }

  return ComposeVersions(versions), nil
}

func CurrentVersion(langauge string) string {
  switch {
    case langauge == "node":
      return nodejs.CurrentVersion()
    case langauge == "rust":
      return rust.CurrentVersion()
  }

  return ""
}
