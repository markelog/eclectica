package nodejs

import (
  "net/http"
  "io/ioutil"
  "regexp"
  "runtime"
  "fmt"
  "os"
  "path/filepath"

  "github.com/termie/go-shutil"
)

var (
  client = &http.Client{}

  partLatest = "https://nodejs.org/dist/latest"
  latest = fmt.Sprintf("%s/SHASUMS256.txt", partLatest)

  directories = [1]string{"test"}
  prefix = "/usr/local"
)

func Latest() (map[string]string, error) {
  file, err := info(latest)
  result := make(map[string]string)

  if err != nil {
    return result, err
  }

  versionReg := regexp.MustCompile(`node-v(\d+\.\d+\.\d)`)

  version := versionReg.FindStringSubmatch(file)[1]
  result["name"] = "node"
  result["version"] = version
  result["filename"] = fmt.Sprintf("node-v%s-%s-x64", version, runtime.GOOS)
  result["url"] = fmt.Sprintf("%s/%s.tar.gz", partLatest, result["filename"])

  return result, nil
}

func Activate(path string) error {
  var err error

  for _, directory := range directories {
    to := fmt.Sprintf("%s/%s", prefix, directory)
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

func exists(name string) bool {
  _, err := os.Stat(name)
  return !os.IsNotExist(err)
}

func info(url string) (file string, err error){
  response, err := client.Get(latest)

  if err != nil {
    return "", err
  }

  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)

  if err != nil {
    return "", err
  }

  return string(contents), nil
}
