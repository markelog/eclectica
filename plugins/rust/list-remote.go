package rust

import (
  "errors"
  "runtime"
  "regexp"

  "github.com/markelog/eclectica/request"
)

var (
  fullVersionPattern = "[0-9]+\\.[0-9]+(?:\\.[0-9]+)?(?:-(alpha|beta)(?:\\.[0-9]*)?)?"
  nighltyPattern = "nightly(\\.[0-9]+)?"
  betaPattern = "beta(\\.[0-9]+)?"
  defaultPattern = "[0-9]+\\.[0-9]+(\\.[0-9]+)?(-(alpha|beta)(\\.[0-9]*)?)?"
  rcPattern = defaultPattern + "-rc(\\.[0-9]+)?"
  versionPattern = "(" + defaultPattern + "|" + betaPattern + "|" + rcPattern + "|" + betaPattern + ")"
)

// Do not know how to test it :/
func getPlatfrom() (string, error) {
  if runtime.GOOS == "linux" {
    return "x86_64-unknown-linux-gnu", nil
  }

  if runtime.GOOS == "darwin" {
    return "x86_64-apple-darwin", nil
  }

  return "", errors.New("Not supported envionment")
}

func getFullPattern() (string, error) {
  platform, err := getPlatfrom()
  result := ""

  if err != nil {
    return result, err
  }

  result = "/dist/rust-" + fullVersionPattern + "-" + platform + ".tar.gz,"

  return result, nil
}

func ListVersions() ([]string, error) {
  body, err := request.Body(versionsLink)
  result := []string{}

  if err != nil {
    return result, err
  }

  fullPattern, err := getFullPattern()

  if err != nil {
    return result, err
  }

  fullUrlsPattern := regexp.MustCompile(fullPattern)
  vp := regexp.MustCompile(versionPattern)

  fullUrlsTmp := fullUrlsPattern.FindAllStringSubmatch(body, -1)
  var fullUrls []string

  // Flatten them out
  for _, element := range fullUrlsTmp {
    fullUrls = append(fullUrls, element[0])
  }

  for _, element := range fullUrls {
    result = append(result, vp.FindAllStringSubmatch(element, 1)[0][0])
  }

  return result, nil
}
