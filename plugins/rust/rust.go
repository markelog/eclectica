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

	"github.com/chuckpreslar/emission"
	"github.com/markelog/cprf"

	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/request"
	"github.com/markelog/eclectica/variables"
)

var (
	VersionsLink   = "https://static.rust-lang.org/dist"
	versionPattern = "\\d+\\.\\d+(?:\\.\\d+)?(?:-(alpha|beta)(?:\\.\\d*)?)?"
	listLink       = "https://static.rust-lang.org/dist/index.txt"

	bins = []string{"cargo", "rust-gdb", "rustc", "rustdoc"}
	dots = []string{".rust-version"}
)

type Rust struct {
	Version string
	Emitter *emission.Emitter
}

func (rust Rust) Events() *emission.Emitter {
	return rust.Emitter
}

func (rust Rust) PreInstall() error {
	return nil
}

func (rust Rust) Install() error {
	path := variables.Path("rust", rust.Version)
	tmp := filepath.Join(path, "tmp")
	installer := filepath.Join(path, "install.sh")

	// Just in case, tmp might not get removed if this method had an error
	// before we could remove it
	os.RemoveAll(tmp)

	_, err := io.CreateDir(tmp)
	if err != nil {
		return err
	}

	_, err = exec.Command(installer, "--prefix="+tmp).Output()
	if err != nil {
		return err
	}

	err = cprf.Copy(tmp+"/", path)
	os.RemoveAll(tmp)

	return err
}

func (rust Rust) PostInstall() error {
	return nil
}

func (rust Rust) Environment() (string, error) {
	return "", nil
}

func (rust Rust) Info() (map[string]string, error) {
	result := make(map[string]string)

	platform, err := getPlatform()

	if err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("rust-%s-%s", rust.Version, platform)
	sourcesUrl := fmt.Sprintf("%s/%s", VersionsLink, filename)

	result["filename"] = filename
	result["url"] = fmt.Sprintf("%s.tar.gz", sourcesUrl)

	return result, nil
}

func (rust Rust) Bins() []string {
	return bins
}

func (node Rust) Dots() []string {
	return dots
}

func (rust Rust) Current() string {
	bin := variables.GetBin("rust")
	out, _ := exec.Command(bin, "--version").Output()

	version := strings.TrimSpace(string(out))
	rVersion := regexp.MustCompile(versionPattern)
	testVersion := rVersion.FindAllStringSubmatch(version, 1)

	if len(testVersion) == 0 {
		return ""
	}

	return strings.Replace(testVersion[0][0], "v", "", 1)
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
