package directory

import (
  "os"
  "fmt"
)

func Create(name string) (string, error) {
  path := fmt.Sprintf("%s/.electica/versions/%s", os.Getenv("HOME"), name)

  err := os.MkdirAll(path, 0700)

  if err != nil {
    return "", err
  }

  return path, nil
}
