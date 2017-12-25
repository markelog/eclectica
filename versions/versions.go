// Package versions provides helpful methods
// for anything related to the versions in eclectica
package versions

import (
	"errors"
	"regexp"
	"sort"
	"strings"

	"github.com/blang/semver"
	hversion "github.com/hashicorp/go-version"
)

// Compose versions to map object of arrays from array
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

// ComposeMajors majors version to map object of arrays from array
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

// ComposeMinors composes minor version to map object of arrays from array
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

// GetKeys returns array of keys
//   map[string][]string{"4.x": []string{}, "0.x": []string{"0.8.2"}}
// gets you:
// 	 string{"0.x", "4.x"}
func GetKeys(versions map[string][]string) []string {
	result := []string{}
	compare := map[string]string{}
	zeroRe := regexp.MustCompile(`\.x$`)

	for version := range versions {
		result = append(result, version)
	}

	// Serialize it to semver list
	semverVersions := make([]*hversion.Version, len(result))
	for i, raw := range result {
		version := zeroRe.ReplaceAllString(raw, ".0-alpha")
		semverVersion, _ := hversion.NewVersion(version)

		compare[semverVersion.String()] = raw
		semverVersions[i] = semverVersion
	}

	// Sort it
	sort.Sort(hversion.Collection(semverVersions))

	// Reverse it
	for i := 0; i < len(semverVersions)/2; i++ {
		j := len(semverVersions) - i - 1
		semverVersions[i], semverVersions[j] = semverVersions[j], semverVersions[i]
	}

	// Bring it back to normal view
	for i, version := range semverVersions {
		result[i] = compare[version.String()]
	}

	return result
}

// GetElements gets all elements for provided range of version in sorted semver format:
//   map[string][]string{
// 		"1.x": string{1.1, 1.1-beta}
// 	}
// Will return:
//   [1.1.0, 1.1.0-beta]
func GetElements(key string, versions map[string][]string) []string {
	for version := range versions {
		if version == key {
			return semverifyList(versions[version])
		}
	}

	return nil
}

// Complete completes the version to semver in case provided value is incomplete
func Complete(version string, vers []string) (string, error) {
	if IsPartial(version) == false {
		return version, nil
	}

	// This shouldn't happen
	if len(vers) == 0 {
		return "", errors.New("No versions available")
	}

	return Latest(version, vers)
}

// IsPartial checks if provided version is full semver version
func IsPartial(version string) bool {
	if version == "latest" {
		return true
	}

	return len(strings.Split(version, ".")) != 3
}

// HasMinor checks if provided version has minor info in it
func HasMinor(version string) bool {
	return len(strings.Split(version, ".")) == 2
}

// getLatest gets last possible version from provided map of array strings
func getLatest(versions map[string][]string) (string, error) {
	latestVersions := GetKeys(versions)[0]
	latestList := semverifyList(versions[latestVersions])

	return latestList[0], nil
}

// Latest returns latest version from provided list
// "1.x" with [1.1.0, 1.1.1-beta, 1.1.1-rc2, 1.0, 1.1.1]
// will return "1.1.1", same for
// "latest" with [1.1.0, 1.1.1-beta, 1.1.1-rc2, 1.0, 1.1.1]
func Latest(version string, versions []string) (string, error) {
	var vers map[string][]string

	if HasMinor(version) {
		vers = ComposeMinors(versions)
	} else {
		vers = ComposeMajors(versions)
	}

	if version == "latest" {
		return getLatest(vers)
	}

	version = version + ".x"

	if _, ok := vers[version]; ok == false {
		return "", errors.New("Incorrect version " + version)
	}

	result := GetElements(version, vers)

	return result[0], nil
}

// semverifyList semverifies the list of incomplete versions
func semverifyList(versions []string) []string {
	semverList := []semver.Version{}
	result := []string{}

	for _, version := range versions {
		parsed, _ := semver.Parse(Semverify(version))
		semverList = append(semverList, parsed)
	}

	semver.Sort(semverList)

	for i := len(semverList) - 1; i != -1; i-- {
		element := semverList[i]
		result = append(result, element.String())
	}

	return result
}

// Semverify will do its best to semverify a string
// "1.8-beta2" -> "1.8.0-beta2"
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

// Unsemverify will do its best to unsemverify a string
// "1.8.0-beta2" -> "1.8-beta2"
func Unsemverify(version string) string {
	rp, _ := regexp.Compile("(\\d+\\.\\d+)\\.0(?:-)?")

	return rp.ReplaceAllString(version, "$1")
}
