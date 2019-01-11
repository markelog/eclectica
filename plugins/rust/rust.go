// Package rust provides all needed logic for installation of rust
package rust

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/chuckpreslar/emission"
	"github.com/go-errors/errors"
	"github.com/markelog/cprf"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/storage/memory"

	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/pkg"
	"github.com/markelog/eclectica/variables"
)

var (
	// VersionLink is the URL link from which we can get all possible versions
	VersionLink = "https://static.rust-lang.org/dist"

	versionPattern = "\\d+\\.\\d+(?:\\.\\d+)?(?:-(alpha|beta)(?:\\.\\d*)?)?"
	listLink       = "https://github.com/rust-lang/rust.git"

	bins = []string{"cargo", "rust-gdb", "rustc", "rustdoc"}
	dots = []string{".rust-version"}
)

// Rust essential struct
type Rust struct {
	Version string
	Emitter *emission.Emitter
	pkg.Base
}

// New returns language struct
func New(version string, emitter *emission.Emitter) *Rust {
	return &Rust{
		Version: version,
		Emitter: emitter,
	}
}

// Events returns language related event emitter
func (rust Rust) Events() *emission.Emitter {
	return rust.Emitter
}

// Install hook
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
		return errors.New(err)
	}

	err = cprf.Copy(tmp+"/", path)
	os.RemoveAll(tmp)

	return err
}

// Info provides all the info needed for installation of the plugin
func (rust Rust) Info() map[string]string {
	var (
		result      = make(map[string]string)
		platform, _ = getPlatform()
		filename    = fmt.Sprintf("rust-%s-%s", rust.Version, platform)
		sourcesURL  = fmt.Sprintf("%s/%s", VersionLink, filename)
	)

	result["filename"] = filename
	result["url"] = fmt.Sprintf("%s.tar.gz", sourcesURL)

	return result
}

// Bins returns list of the all bins included
// with the distribution of the language
func (rust Rust) Bins() []string {
	return bins
}

// Dots returns list of the all available filenames
// which can define versions
func (rust Rust) Dots() []string {
	return dots
}

// ListRemote returns list of the all available remote versions
func (rust Rust) ListRemote() ([]string, error) {
	// Get stuff from git, since it's the only way to get it for rust.
	// See https://github.com/rust-lang/www.rust-lang.org/issues/662
	return gitRemote(listLink)
}

func gitRemote(link string) ([]string, error) {
	remote, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		return nil, err
	}

	repo, err := remote.CreateRemote(&config.RemoteConfig{
		Name: "rust",
		URLs: []string{link},
	})
	if err != nil {
		return nil, err
	}

	list, err := repo.List(&git.ListOptions{})
	if err != nil {
		return nil, err
	}

	tag := "refs/tags/"
	versions := []string{}
	for _, item := range list {
		name := item.Name().String()

		if strings.Contains(name, tag) == false {
			continue
		}

		versions = append(versions, strings.Replace(name, tag, "", 1))
	}

	return versions, nil
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
