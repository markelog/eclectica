package versions

import (
	"errors"
	"regexp"
	"sort"
	"strings"

	"github.com/blang/semver"
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

	// In revese order
	// Should we use semver sort?
	sort.Sort(sort.Reverse(sort.StringSlice(result)))

	return result
}

func GetElements(key string, versions map[string][]string) []string {
	result := []string{}
	semverList := []semver.Version{}

	for version, _ := range versions {
		if version == key {
			for _, element := range versions[version] {
				parsed, _ := semver.Parse(Semverify(element))

				semverList = append(semverList, parsed)
			}
		}
	}

	semver.Sort(semverList)

	for i := len(semverList) - 1; i != -1; i-- {
		element := semverList[i]
		result = append(result, element.String())
	}

	return result
}

func IsPartialVersion(version string) bool {
	return len(strings.Split(version, ".")) != 3
}

func HasMinor(version string) bool {
	return len(strings.Split(version, ".")) == 2
}

func GetLatest(version string, versions []string) (string, error) {
	var vers map[string][]string

	if HasMinor(version) {
		vers = ComposeMinors(versions)
	} else {
		vers = ComposeMajors(versions)
	}

	version = version + ".x"

	if _, ok := vers[version]; ok == false {
		return "", errors.New("Incorrect version " + version)
	}

	result := GetElements(version, vers)

	return result[0], nil
}

func Semverify(version string) string {
	if HasMinor(version) == false {
		return version
	}

	rp, _ := regexp.Compile("[a-z](.+)?")

	if rp.MatchString(version) {
		start := rp.ReplaceAllString(version, "")
		end := rp.FindAllString(version, 1)[0]

		return start + ".0-" + end
	}

	version += ".0"

	return version
}

func Unsemverify(version string) string {
	rp, _ := regexp.Compile("(\\d+\\.\\d+)\\.0(?:-)?")

	return rp.ReplaceAllString(version, "$1")
}
