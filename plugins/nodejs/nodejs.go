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

  versionsLink = "https://nodejs.org/dist"

  directories = [1]string{"test"}
  prefix = "/usr/local"
)

func Latest() (map[string]string, error) {
  result := make(map[string]string)
  sumUrl := fmt.Sprintf("%s/latest/SHASUMS256.txt", versionsLink)
  sourcesUrl := fmt.Sprintf("%s/latest", versionsLink)
  file, err := info(sumUrl)

  if err != nil {
    return result, err
  }

  versionReg := regexp.MustCompile(`node-v(\d+\.\d+\.\d)`)

  version := versionReg.FindStringSubmatch(file)[1]
  result["name"] = "node"
  result["version"] = version
  result["filename"] = fmt.Sprintf("node-v%s-%s-x64", version, runtime.GOOS)
  result["url"] = fmt.Sprintf("%s/%s.tar.gz", sourcesUrl, result["filename"])

  return result, nil
}

func Version(params ...string) (map[string]string, error) {
  var version string

  if len(params) == 0 {
    return Latest()
  }

  result := make(map[string]string)

  sourcesUrl := fmt.Sprintf("%s/v%s", versionsLink, version)

  result["name"] = "node"
  result["version"] = version
  result["filename"] = fmt.Sprintf("node-v%s-%s-x64", version, runtime.GOOS)
  result["url"] = fmt.Sprintf("%s/%s.tar.gz", sourcesUrl, result["filename"])

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

func info(url string) (file string, err error){
  response, err := client.Get(url)

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
