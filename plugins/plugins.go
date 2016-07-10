package plugins

import (
  "strings"

  "github.com/markelog/eclectica/plugins/nodejs"
)

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

func Activate(data map[string]string) error {
  switch {
    case data["name"] == "node":
      return nodejs.Activate(data)
  }

  return nil
}
