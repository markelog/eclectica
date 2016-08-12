package rust

import (
  "regexp"

  "github.com/markelog/eclectica/request"
)

func getFullPattern() (string, error) {
  platform, err := getPlatform()
  result := ""

  if err != nil {
    return result, err
  }

  result = "/dist/rust-" + fullVersionPattern + "-" + platform + ".tar.gz,"

  return result, nil
}

func ListVersions() ([]string, error) {
  body, err := request.Body(listLink)
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
