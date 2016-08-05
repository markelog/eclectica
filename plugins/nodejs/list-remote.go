package nodejs

import (
  "regexp"
  "strings"

  "github.com/PuerkitoBio/goquery"
)

func ListVersions() ([]string, error) {
  doc, err := goquery.NewDocument(versionsLink)
  tmp := []string{}
  result := []string{}
  version := regexp.MustCompile("v[[:digit:]]+\\.[[:digit:]]+\\.[[:digit:]]+$")
  remove := regexp.MustCompile("0\\.[0-7]")

  if err != nil {
    return nil, err
  }

  links := doc.Find("a")

  for i := range links.Nodes {
    href, _ := links.Eq(i).Attr("href")

    href = strings.Replace(href, "/", "", 1)
    if version.MatchString(href) {
      href = strings.Replace(href, "v", "", 1)
      tmp = append(tmp, href)
    }
  }

  // Remove < 0.8 versions
  for _, element := range tmp {
    if remove.MatchString(element) == false {
      result = append(result, element)
    }
  }

  return result, nil
}

func ComposeVersions(versions []string) map[string][]string {
  result := map[string][]string{}
  firstPart := regexp.MustCompile("([[:digit:]]+)\\.")

  for _, version := range versions {
    major := firstPart.FindAllStringSubmatch(version, 1)[0][1]
    major += ".x"

    if _, ok := result[major]; ok == false {
      result[major] = []string{}
    }

    result[major] = append(result[major], version)
  }

  return result
}

func GetKeys(versions map[string][]string) []string {
  result := []string{}

  for version, _ := range versions {
    result = append(result, version)
  }

  return result
}

func GetElements(key string, versions map[string][]string) []string {
  result := []string{}

  for version, _ := range versions {
    if version == key {
      for _, element := range versions[version] {
        result = append(result, element)
      }
    }
  }

  return result
}
