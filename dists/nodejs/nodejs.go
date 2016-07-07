package nodejs

import (
  "net/http"
  "io/ioutil"
  "regexp"
  "runtime"
  "fmt"
)

var (
  client = &http.Client{}

  partLatest = "https://nodejs.org/dist/latest"
  latest = fmt.Sprintf("%s/SHASUMS256.txt", partLatest)
)

func system() string {
  name := runtime.GOOS

  switch {
    case name == "darwin":
      return "darwin-x64"
    default:
      return "darwin-x64"
  }
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
  result["filename"] = fmt.Sprintf("node-v%s-%s", version, system())
  result["url"] = fmt.Sprintf("%s/%s.tar.gz", partLatest, result["filename"])

  return result, nil
}
