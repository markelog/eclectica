package rust

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/markelog/eclectica/request"
	"github.com/markelog/eclectica/variables"
)

var (
	versionsLink = "https://static.rust-lang.org/dist"
	listLink     = "https://static.rust-lang.org/dist/index.txt"

	home = filepath.Join(variables.Home(), "rust")
	bin  = variables.Prefix("rust") + "/bin/rustc"

	versionPattern = "\\d+\\.\\d+(?:\\.\\d+)?(?:-(alpha|beta)(?:\\.\\d*)?)?"
)

type Rust struct{}

func (rust Rust) Install(version string) error {
	installer := fmt.Sprintf("%s/%s/%s", home, version, "install.sh")
	_, err := exec.Command(installer, "--prefix="+variables.Prefix("rustc")).Output()

	return err
}

func (rust Rust) PostInstall() (bool, error) {
	return true, nil
}

func (rust Rust) Info(version string) (map[string]string, error) {
	result := make(map[string]string)

	platform, err := getPlatform()

	if err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("rust-%s-%s", version, platform)
	sourcesUrl := fmt.Sprintf("%s/%s", versionsLink, filename)

	result["version"] = version
	result["filename"] = filename
	result["url"] = fmt.Sprintf("%s.tar.gz", sourcesUrl)

	return result, nil
}

func (rust Rust) Current() string {
	vp := regexp.MustCompile(versionPattern)

	out, _ := exec.Command(bin, "--version").Output()
	version := strings.TrimSpace(string(out))
	versionArr := vp.FindAllStringSubmatch(version, 1)
	if len(versionArr) > 0 {
		version = strings.Replace(versionArr[0][0], "v", "", 1)
	}

	return strings.Replace(version, "v", "", 1)
}

func (rust Rust) ListRemote() ([]string, error) {
	body, err := request.Body(listLink)

	if err != nil {
		return []string{}, err
	}

	return getVersions(body)
}

func getFullPattern() (string, error) {
	platform, err := getPlatform()

	if err != nil {
		return "", err
	}

	result := "/dist/rust-" + versionPattern + "-" + platform + ".tar.gz,"

	return result, nil
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
