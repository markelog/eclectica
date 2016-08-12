package rust

import (
  "regexp"

  "github.com/markelog/eclectica/request"
)

func getFullPattern() (string, error) {
  platform, err := getPlatform()

  if err != nil {
    return "", err
  }

  result := "/dist/rust-" + fullVersionPattern + "-" + platform + ".tar.gz,"

  return result, nil
}

func ListVersions() ([]string, error) {
  body, err := request.Body(listLink)

  if err != nil {
    return []string{}, err
  }

  return getVersions(body)
}

func getVersions(list string) ([]string, error) {
  fullPattern, err := getFullPattern()
  result := []string{}

  if err != nil {
    return result, err
  }

  fullUrlsPattern := regexp.MustCompile(fullPattern)

  fullUrlsTmp := fullUrlsPattern.FindAllStringSubmatch(list, -1)
  var fullUrls []string

  // Flatten them out
  for _, element := range fullUrlsTmp {
    fullUrls = append(fullUrls, element[0])
  }

  vp := regexp.MustCompile(versionPattern)
  for _, element := range fullUrls {
    result = append(result, vp.FindAllStringSubmatch(element, 1)[0][0])
  }

  return result, nil
}
