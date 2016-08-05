package plugins

import (
  "os"
  "strings"
  "io/ioutil"
  "errors"

  "github.com/markelog/eclectica/plugins/nodejs"
  "github.com/markelog/eclectica/variables"
)

var (
  List = []string{
    "node",
  }
)

func Versions(name string) ([]string, error) {
  path := variables.Home + "/" + name

  if _, err := os.Stat(path); os.IsNotExist(err) {
    return nil, errors.New("There is no installed versions of " + name)
  }

  folders, err := ioutil.ReadDir(variables.Home + "/" + name)
  versions := make([]string, len(folders))

  if err != nil {
    return nil, err
  }

  for _, folder := range folders {
    versions = append(versions, folder.Name())
  }

  return versions, nil
}

func Detect(nameAndVersion string) (map[string]string, error) {
  var (
    info map[string]string
    version = "latest"
    data = strings.Split(nameAndVersion, "@")
    plugin = data[0]
    err error
  )

  if len(data) == 2 {
    version = data[1]
  }

  switch {
    case plugin == "node":
      info, err = nodejs.Version(version)
  }

  if err != nil {
    return nil, err
  }

  return info, nil
}

func Remove(nameAndVersion string) error {
  data := strings.Split(nameAndVersion, "@")

  if len(data) == 1 {
    return errors.New("Can't remove without specific version")
  }

  name := data[0]
  version := data[1]

  switch {
    case name == "node":
      return nodejs.Remove(version)
  }

  return nil
}


func Activate(data map[string]string) error {
  switch {
    case data["name"] == "node":
      return nodejs.Activate(data)
  }

  return nil
}

func CurrentVersion(name string) string {
  switch {
    case name == "node":
      return nodejs.CurrentVersion()
  }

  return ""
}
