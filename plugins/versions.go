package plugins

import (
  "regexp"
  "sort"
)

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

  sort.Strings(result)

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

  sort.Strings(result)

  return result
}
