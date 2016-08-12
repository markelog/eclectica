package rust

import (
  "runtime"
  "fmt"
  "os"
  "os/exec"
  "strings"
  "errors"
  "regexp"

  "github.com/markelog/eclectica/variables"
)

var (
  versionsLink = "https://static.rust-lang.org/dist"
  listLink = "https://static.rust-lang.org/dist/index.txt"

  home = fmt.Sprintf("%s/%s", variables.Home, "rust")

  dists = [2]string{"cargo", "rustc"}
  files = [4]string{"bin", "lib", "include", "share"}
  prefix = os.Getenv("HOME")
  bin = prefix + "/bin/rustc"

  fullVersionPattern = "[0-9]+\\.[0-9]+(?:\\.[0-9]+)?(?:-(alpha|beta)(?:\\.[0-9]*)?)?"
  nighltyPattern = "nightly(\\.[0-9]+)?"
  betaPattern = "beta(\\.[0-9]+)?"
  defaultPattern = "[0-9]+\\.[0-9]+(\\.[0-9]+)?(-(alpha|beta)(\\.[0-9]*)?)?"
  rcPattern = defaultPattern + "-rc(\\.[0-9]+)?"
  versionPattern = "(" + defaultPattern + "|" + betaPattern + "|" + rcPattern + "|" + betaPattern + ")"
)

// Do not know how to test it :/
func getPlatform() (string, error) {
  if runtime.GOOS == "linux" {
    return "x86_64-unknown-linux-gnu", nil
  }

  if runtime.GOOS == "darwin" {
    return "x86_64-apple-darwin", nil
  }

  return "", errors.New("Not supported envionment")
}

func Keyword(keyword string) (map[string]string, error) {
  return Version(keyword)
}

func Version(version string) (map[string]string, error) {
  result := make(map[string]string)

  platform, err := getPlatform()

  filename := fmt.Sprintf("rust-%s-%s", version, platform)
  sourcesUrl := fmt.Sprintf("%s/%s", versionsLink, filename)

  if err != nil {
    return nil, err
  }

  result["name"] = "rust"
  result["version"] = version
  result["filename"] = filename
  result["url"] = fmt.Sprintf("%s.tar.gz", sourcesUrl)

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
  installer := fmt.Sprintf("%s/%s/%s", home, data["version"], "install.sh")
  _, err := exec.Command(installer, "--prefix=" + prefix).Output()

  return err
}

func CurrentVersion() string {
  vp := regexp.MustCompile(versionPattern)

  out, _ := exec.Command(bin, "--version").Output()
  version := strings.TrimSpace(string(out))
  version = vp.FindAllStringSubmatch(version, 1)[0][0]

  return strings.Replace(version, "v", "", 1)
}
