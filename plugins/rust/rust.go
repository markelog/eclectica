package rust

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/markelog/cprf"

	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/request"
	"github.com/markelog/eclectica/variables"
)

var (
	listLink       = "https://static.rust-lang.org/dist/index.txt"
	versionsLink   = "https://static.rust-lang.org/dist"
	versionPattern = "\\d+\\.\\d+(?:\\.\\d+)?(?:-(alpha|beta)(?:\\.\\d*)?)?"

	bins = []string{"cargo", "rust-gdb", "rustc", "rustdoc"}
)

type Rust struct{}

func (rust Rust) Install(version string) error {
	path := variables.Path("rust", version)
	tmp := filepath.Join(variables.Home(), "rust", "tmp")
	source := filepath.Join(path, "source")
	installer := filepath.Join(source, "install.sh")

	// Just in case, tmp might not get removed if this method had an error
	// before we could remove it
	os.RemoveAll(tmp)

	err := os.Rename(path, tmp)
	if err != nil {
		return err
	}

	_, err = io.CreateDir(filepath.Join(path, "source"))
	if err != nil {
		return err
	}

	cprf.Copy(tmp+"/", source)

	_, err = exec.Command(installer, "--prefix="+path).Output()
	os.RemoveAll(tmp)

	return err
}

func (rust Rust) PostInstall(version string) error {
	return nil
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

func (rust Rust) Bins() []string {
	return bins
}

func (rust Rust) Current() string {
	vp := regexp.MustCompile(versionPattern)
	out, _ := exec.Command(variables.GetBin("rust", ""), "--version").Output()

	version := strings.TrimSpace(string(out))
	versionArr := vp.FindAllStringSubmatch(version, 1)
	if len(versionArr) > 0 {
		version = strings.Replace(versionArr[0][0], "v", "", 1)
	}

	return strings.Replace(version, "v", "", 1)
}

func (rust Rust) Environment(version string) (string, error) {
	return "", nil
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
