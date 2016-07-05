package detect

import (
  "github.com/markelog/eclectica/download"
  "github.com/markelog/eclectica/dists/nodejs"
)

func Detect(dist string) error {
  var (
    version map[string]string
    err error
  )

  switch {
    case dist == "node":
      version, err = nodejs.Latest()
  }

  if err != nil {
    return err
  }

  download.Download(version["url"])

  return nil
}
