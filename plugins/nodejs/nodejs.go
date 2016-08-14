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
  home = fmt.Sprintf("%s/%s", variables.Home(), "node")
  bin = variables.Prefix() + "/bin/node"
)

// TODO
// Plugin interface –
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

  result["archive-folder"] = os.TempDir()
  result["archive-path"] = fmt.Sprintf("%s%s.tar.gz", result["archive-folder"], result["filename"])

  result["destination-folder"] = fmt.Sprintf("%s/%s/%s", variables.Home(), result["name"], result["version"])

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

  result["archive-folder"] = os.TempDir()
  result["archive-path"] = fmt.Sprintf("%s%s.tar.gz", result["archive-folder"], result["filename"])

  result["destination-folder"] = fmt.Sprintf("%s/%s/%s", variables.Home(), result["name"], result["version"])

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

func Install(data map[string]string) error {
  var err error

  base := fmt.Sprintf("%s/%s", home, data["version"])

  for _, file := range variables.Files {
    from := fmt.Sprintf("%s/%s", base, file)
    to := variables.Prefix()

    // Some versions might not have certain files
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
