package plugins

import (
  "github.com/markelog/eclectica/plugins/nodejs"
)

func Detect(plugin string) (map[string]string, error) {
  var (
    version map[string]string
    err error
  )

  switch {
    case plugin == "node":
      version, err = nodejs.Version()
  }

  if err != nil {
    return nil, err
  }

  return version, nil
}

func Activate(name string, data map[string]string) error {
  switch {
    case name == "node":
      return nodejs.Activate(data)
  }

  return nil
}
