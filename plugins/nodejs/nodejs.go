package nodejs

import (
  "net/http"
  "io/ioutil"
  "regexp"
  "runtime"
  "fmt"
  "os"

  "github.com/markelog/eclectica/variables"
)

var (
  client = &http.Client{}

  versionsLink = "https://nodejs.org/dist"
  home = fmt.Sprintf("%s/%s", variables.Home, "node")

  bins = [2]string{"node", "npm"}
  prefix = "/usr/local/bin"
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

func Activate(data map[string]string) error {
  base := fmt.Sprintf("%s/%s/bin", home, data["version"])

  if _, err := os.Stat(prefix); os.IsNotExist(err) {
    err := os.MkdirAll(prefix, 0755)

    if err != nil {
      return err
    }
  }

  for _, bin := range bins {
    from := fmt.Sprintf("%s/%s", base, bin)
    to := fmt.Sprintf("%s/%s", prefix, bin)

    if _, err := os.Stat(to); err == nil {
      err := os.Remove(to)
      if err != nil {
        return err
      }
    }

    err := os.Symlink(from, to)

    if err != nil {
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
