package plugins

import (
	"regexp"
	"sort"
)

func Compose(versions []string) map[string][]string {
	majors := ComposeMajors(versions)
	if len(majors) > 1 {
		return majors
	}

	minors := ComposeMinors(versions)
	if len(minors) > 1 {
		return minors
	}

	return majors
}

func ComposeMajors(versions []string) map[string][]string {
	result := map[string][]string{}
	firstPart := regexp.MustCompile("(\\d+)\\.")

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

func ComposeMinors(versions []string) map[string][]string {
	result := map[string][]string{}
	firstPart := regexp.MustCompile("(\\d)+\\.(\\d+)")

	for _, version := range versions {
		checkVersions := firstPart.FindAllStringSubmatch(version, 1)

		// Just in case
		if len(checkVersions) == 0 {
			continue
		}

		versions := checkVersions[0]

		part := versions[1] + "." + versions[2] + ".x"

		if _, ok := result[part]; ok == false {
			result[part] = []string{}
		}

		result[part] = append(result[part], version)
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
