package detect

import (
  "github.com/markelog/eclectica/dists/nodejs"
)

func Detect(dist string) (map[string]string, error) {
  var (
    version map[string]string
    err error
  )

  switch {
    case dist == "node":
      version, err = nodejs.Latest()
  }

  if err != nil {
    return nil, err
  }

  return version, nil
}
