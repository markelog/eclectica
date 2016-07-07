package activation

import (
  "fmt"
  "os"
  "path/filepath"

  "github.com/termie/go-shutil"

  "github.com/markelog/eclectica/variables"
)

func exists(name string) bool {
  _, err := os.Stat(name)
  return !os.IsNotExist(err)
}

func Activate(path string) error {
  var err error

  for _, directory := range variables.Directories {
    to := fmt.Sprintf("%s/%s", variables.Prefix, directory)
    from := fmt.Sprintf("%s/%s", path, directory)

    filepath.Walk(from, func(path string, info os.FileInfo, err error) error {
      if from == path {
        return nil
      }

      newPath := filepath.Join(to, info.Name())

      if info.IsDir() {
        fmt.Println()

        err = os.MkdirAll(newPath, info.Mode())

        if err != nil {
          return err
        }

        return nil
      }

      err = shutil.CopyFile(path, newPath, true)
      fmt.Println(err)

      return nil
    })

    if err != nil {
      fmt.Println(err)
      return err
    }
  }

  return nil
}
