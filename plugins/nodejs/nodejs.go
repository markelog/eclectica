package nodejs

import (
  "regexp"
  "runtime"
  "fmt"
  "os"
  "os/exec"
  "strings"

  "github.com/markelog/cprf"

  "github.com/markelog/eclectica/variables"
  "github.com/markelog/eclectica/request"
)

var (
  versionsLink = "https://nodejs.org/dist"
  home = fmt.Sprintf("%s/%s", variables.Home, "node")

  // TODO: move to variables
  files = [4]string{"bin", "lib", "include", "share"}

  // TODO: move to variables
  prefix = os.Getenv("HOME")

  bin = prefix + "/bin/node"
)

// Plugin implementation.
// type Plugin struct{}

func CurrentVersion() string {
  out, _ := exec.Command(bin, "--version").Output()
  version := strings.TrimSpace(string(out))

  return strings.Replace(version, "v", "", 1)
}

// TODO: Info
func Version(version string) (map[string]string, error) {
  if version == "latest" || version == "lts" {
    return Keyword(version)
  }

  result := make(map[string]string)

  sourcesUrl := fmt.Sprintf("%s/v%s", versionsLink, version)

  result["name"] = "node"
  result["version"] = version
  result["filename"] = fmt.Sprintf("node-v%s-%s-x64", version, runtime.GOOS)
  result["url"] = fmt.Sprintf("%s/%s.tar.gz", sourcesUrl, result["filename"])

  return result, nil
}

func Keyword(keyword string) (map[string]string, error) {
  result := make(map[string]string)
  sumUrl := fmt.Sprintf("%s/%s/SHASUMS256.txt", versionsLink, keyword)
  sourcesUrl := fmt.Sprintf("%s/%s", versionsLink, keyword)
  file, err := request.Body(sumUrl)

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

func Remove(version string) error {
  var err error
  base := fmt.Sprintf("%s/%s", home, version)

  err = os.RemoveAll(base)

  if err != nil {
    return err
  }

  return nil
}

func Activate(data map[string]string) error {
  var err error

  base := fmt.Sprintf("%s/%s", home, data["version"])

  for _, file := range files {
    from := fmt.Sprintf("%s/%s", base, file)
    to := prefix

    // Older versions might not have certain files
    if _, err := os.Stat(from); os.IsNotExist(err) {
      continue
    }

    err = cprf.Copy(from, to)

    if err != nil {
      return err
    }
  }

  return nil
}
