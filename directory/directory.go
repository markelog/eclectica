package directory

import (
  "os"
  "fmt"

  "github.com/markelog/eclectica/variables"
)

func Create(name string) (string, error) {
  path := fmt.Sprintf("%s/%s", variables.Home, name)

  err := os.MkdirAll(path, 0700)

  if err != nil {
    return "", err
  }

  return path, nil
}
