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

func Activate(name string, path string) {
  switch {
    case name == "node":
      nodejs.Activate(path)
  }
}
