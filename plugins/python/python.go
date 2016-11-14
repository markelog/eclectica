package python

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/blang/semver"
	"github.com/cavaliercoder/grab"
	"github.com/chuckpreslar/emission"

	"github.com/markelog/eclectica/variables"
	"github.com/markelog/eclectica/versions"
)

var (
	VersionsLink   = "https://www.python.org/ftp/python"
	versionPattern = "^\\d+\\.\\d+(?:\\.\\d)?"

	setuptoolsName = "ez_setup.py"
	pipName        = "get-pip.py"
	baseUrl        = "https://bootstrap.pypa.io/"
	setuptoolsUrl  = baseUrl + setuptoolsName
	pipUrl         = baseUrl + pipName

	// Hats off to inconsistent python developers
	noNilVersions, _ = semver.Make("3.3.0")
	// When pip began to be available with binaries
	pipAvailable, _ = semver.Make("2.7.9")

	bins = []string{"2to3", "idle", "pydoc", "python", "python-config", "pip", "easy_install"}
	dots = []string{".python-version"}
)

type Python struct {
	Version string
	Emitter *emission.Emitter
}

func (python Python) Events() *emission.Emitter {
	return python.Emitter
}

func (python Python) Install() (err error) {
	path := variables.Path("python", python.Version)
	configure := filepath.Join(path, "configure")

	python.Emitter.Emit("configure")

	cmd := python.getCmd(configure, "--prefix="+path)
	cmd.Env = python.getOSXEnvs(cmd.Env)
	_, err = cmd.CombinedOutput()

	if err != nil {
		os.RemoveAll(path)
		return
	}

	python.Emitter.Emit("prepare")

	_, err = python.getCmd("make").CombinedOutput()
	if err != nil {
		os.RemoveAll(path)
		return
	}

	python.Emitter.Emit("install")

	_, err = python.getCmd("make", "install").CombinedOutput()
	if err != nil {
		os.RemoveAll(path)
		return
	}

	// Since python 3.x versions are naming their binaries with "3" affix
	chosen, _ := semver.Make(python.Version)
	if chosen.Major < 3 {
		return nil
	}

	err = renameLinks(python.Version)
	if err != nil {
		os.RemoveAll(path)
		return
	}

	return nil
}

func (python Python) getOSXEnvs(original []string) []string {
	externals := []string{"readline", "openssl"}

	includeFlags := ""
	libFlags := ""

	for _, name := range externals {
		opt := "/usr/local/opt/"
		libFlags += `-L` + filepath.Join(opt, name, "lib") + " "
		includeFlags += "-I" + filepath.Join(opt, name, "include") + " "
	}

	// For zlib
	// TODO: xcode required
	output, _ := exec.Command("xcrun", "--show-sdk-path").CombinedOutput()
	includeFlags += " -I" + filepath.Join(strings.TrimSpace(string(output)), "/usr/include")

	original = append(original, "CPPFLAGS="+includeFlags)
	original = append(original, "LDFLAGS="+libFlags)

	return original
}

func (python Python) getCmd(args ...interface{}) (cmd *exec.Cmd) {

	// There is gotta be a better way without reflect module, huh?
	if len(args) == 1 {
		cmd = exec.Command(args[0].(string))
	} else if len(args) == 2 {
		cmd = exec.Command(args[0].(string), args[1].(string))
	} else {
		cmd = exec.Command(args[0].(string), args[1].(string), args[2].(string))
	}

	// Lots of needless, Makefile weird warnings in the Makefile
	// cmd.Stderr = os.Stderr
	// cmd.Stdout = os.Stdout

	cmd.Env = append(os.Environ(), "LC_ALL=C") // Required for some reason
	cmd.Dir = variables.Path("python", python.Version)

	return cmd
}

func (python Python) PostInstall() (err error) {
	var (
		errStat error
		base    = filepath.Join(variables.Path("python", python.Version), "bin")
		pipBin  = filepath.Join(base, "pip")
		eIBin   = filepath.Join(base, "easy_install")
	)

	// Don't need to do anything if we already have pip and setuptools
	if _, errStat = os.Lstat(pipBin); errStat == nil {
		return
	}
	if _, errStat = os.Lstat(eIBin); errStat == nil {
		return
	}

	python.Emitter.Emit("post-install")

	// Since 2.7.9 versions we can simplify pip and setuptools install
	semverVersion, _ := semver.Make(python.Version)
	if semverVersion.Compare(pipAvailable) != -1 {
		cmd := python.getCmd("python", "-m", "ensurepip")
		_, err = cmd.CombinedOutput()

		return
	}

	// Now try the "hard" way
	path, err := downloadExternals()
	if err != nil {
		return
	}

	// Setup pip
	pip := filepath.Join(path, pipName)
	_, err = exec.Command("python", pip).CombinedOutput()
	if err != nil {
		return err
	}

	// Setup setuptools
	setuptools := filepath.Join(path, setuptoolsName)
	_, err = exec.Command("python", setuptools).CombinedOutput()
	if err != nil {
		return err
	}

	return
}

func (python Python) Environment() (string, error) {
	return "", nil
}

func (python Python) Info() (map[string]string, error) {
	result := make(map[string]string)
	version := python.Version
	chosen, err := semver.Make(python.Version)

	if err != nil {
		return nil, err
	}

	// Hats off to inconsistent python developers
	if chosen.LT(noNilVersions) {
		version = versions.Unsemverify(python.Version)
	}

	// Python 2.0 has different format and its not supported
	result["extension"] = "tgz"
	result["version"] = version
	result["filename"] = "Python-" + version
	result["url"] = fmt.Sprintf(
		"%s/%s/%s.%s",
		VersionsLink,
		version,
		result["filename"],
		result["extension"],
	)

	return result, nil
}

func (rust Python) Bins() []string {
	return bins
}

func (rust Python) Dots() []string {
	return dots
}

func (python Python) Current() string {
	bin := variables.GetBin("python")

	out, _ := exec.Command(bin, "--version").CombinedOutput()

	readVersion := strings.Replace(string(out), "Python ", "", 1)
	version := strings.TrimSpace(readVersion)

	return versions.Semverify(version)
}

func (python Python) ListRemote() ([]string, error) {
	doc, err := goquery.NewDocument(VersionsLink)

	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, errors.New("Can't establish connection")
		}

		return nil, err
	}

	result := []string{}
	tmp := []string{}
	version := regexp.MustCompile(versionPattern)

	links := doc.Find("a")

	for i := range links.Nodes {
		content := links.Eq(i).Text()

		content = strings.Replace(content, "/", "", 1)
		if version.MatchString(content) {
			tmp = append(tmp, content)
		}
	}

	// Remove < 2.7 versions
	for _, element := range tmp {
		smr, _ := semver.Make(versions.Semverify(element))

		if smr.Major > 2 || smr.Minor > 5 {
			result = append(result, element)
		}
	}

	// Latest version is a development one
	result = result[:len(result)-1]

	return result, nil
}

// Since python 3.x versions are naming their binaries with 3 affix
func renameLinks(version string) (err error) {
	path := filepath.Join(variables.Path("python", version), "bin")
	rp := regexp.MustCompile("(-?)3\\.\\w")

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		name := file.Name()
		absPath := filepath.Join(path, name)

		if rp.MatchString(name) {
			pathPart := rp.ReplaceAllString(name, "")
			newPath := filepath.Join(path, pathPart)

			// Since python install creates some links with version and some without
			if _, errStat := os.Lstat(newPath); errStat == nil {
				continue
			}

			err = os.Symlink(absPath, newPath)
			if err != nil {
				return
			}
		}
	}

	return nil
}

func downloadExternals() (path string, err error) {
	path = os.TempDir()
	urls := []string{setuptoolsUrl, pipUrl}

	respch, err := grab.GetBatch(2, path, urls...)
	if err != nil {
		return
	}

	// Start a ticker to update progress every 200ms
	ticker := time.NewTicker(200 * time.Millisecond)

	// Monitor downloads
	completed := 0
	responses := make([]*grab.Response, 0)
	for completed < len(urls) {
		select {
		case resp := <-respch:

			// When done
			if resp != nil {
				responses = append(responses, resp)
			}

		case <-ticker.C:

			// Update completed downloads
			for i, resp := range responses {
				if resp != nil && resp.IsComplete() {

					if resp.Error != nil {
						err = resp.Error
						return
					}

					// Mark completed
					responses[i] = nil
					completed++
				}
			}
		}
	}

	ticker.Stop()

	return
}
