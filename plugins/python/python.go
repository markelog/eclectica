package python

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/blang/semver"
	"github.com/cavaliercoder/grab"
	"github.com/chuckpreslar/emission"

	"github.com/markelog/eclectica/console"
	"github.com/markelog/eclectica/plugins/python/patch"
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

func (python Python) PreInstall() error {
	return checkDependencies()
}

func (python Python) Install() (err error) {
	err = python.configure()
	if err != nil {
		return
	}

	err = python.prepare()
	if err != nil {
		return
	}

	err = python.install()
	if err != nil {
		return
	}

	return python.renameLinks()
}

func (python Python) PostInstall() (err error) {
	path := variables.Path("python", python.Version)
	bin := variables.GetBin("python", python.Version)

	if hasTools(python.Version) {
		cmd, stdErr, stdOut := python.getCmd(bin, "-m", "ensurepip")

		err = cmd.Run()
		if err != nil {
			return console.GetError(err, stdErr, stdOut)
		}

		return
	}

	// Setup pip
	pip := filepath.Join(path, pipName)
	out, err := exec.Command(bin, pip).CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}

	// Setup setuptools
	setuptools := filepath.Join(path, setuptoolsName)
	out, err = exec.Command(bin, setuptools).CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}

	return nil
}

func (python Python) Environment() (result []string, err error) {
	return
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

func (python Python) configure() (err error) {
	path := variables.Path("python", python.Version)
	configure := filepath.Join(path, "configure")

	python.Emitter.Emit("configure")

	err = python.externals()
	if err != nil {
		return
	}

	cmd, stdErr, stdOut := python.getCmd(configure, "--prefix="+path, "--with-ensurepip=upgrade")
	cmd.Env = python.getEnvs(cmd.Env)

	err = cmd.Run()
	if err != nil {
		return console.GetError(err, stdErr, stdOut)
	}

	return
}

func (python Python) prepare() (err error) {
	python.Emitter.Emit("prepare")

	cmd, stdErr, stdOut := python.getCmd("make")

	err = cmd.Run()
	if err != nil {
		return console.GetError(err, stdErr, stdOut)
	}

	return
}

func (python Python) install() (err error) {
	python.Emitter.Emit("install")

	cmd, stdErr, stdOut := python.getCmd("make", "install")

	err = cmd.Run()
	if err != nil {
		return console.GetError(err, stdErr, stdOut)
	}

	return
}

func (python Python) getEnvs(original []string) (result []string) {
	if runtime.GOOS == "darwin" {
		result = python.getOSXEnvs(original)
	}

	return
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
	output, _ := exec.Command("xcrun", "--show-sdk-path").CombinedOutput()
	out := strings.TrimSpace(string(output))
	includeFlags += " -I" + filepath.Join(out, "/usr/include")

	original = append(original, "CPPFLAGS="+includeFlags)
	original = append(original, "LDFLAGS="+libFlags)

	return original
}

func (python Python) getCmd(args ...interface{}) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer) {
	var (
		cmd    *exec.Cmd
		stdOut bytes.Buffer
		stdErr bytes.Buffer
	)

	// There is gotta be a better way without reflect module, huh?
	if len(args) == 1 {
		cmd = exec.Command(args[0].(string))
	} else if len(args) == 2 {
		cmd = exec.Command(args[0].(string), args[1].(string))
	} else {
		cmd = exec.Command(args[0].(string), args[1].(string), args[2].(string))
	}

	cmd.Env = append(os.Environ(), "LC_ALL=C") // Required for some reason
	cmd.Dir = variables.Path("python", python.Version)
	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut

	if variables.IsDebug() {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	return cmd, &stdOut, &stdErr
}

func (python Python) externals() (err error) {
	var (
		errStat error
		base    = filepath.Join(variables.Path("python", python.Version), "bin")
		pipBin  = filepath.Join(base, "pip")
		eIBin   = filepath.Join(base, "easy_install")
		path    = variables.Path("python", python.Version)
	)

	// Don't need to do anything if we already have pip and setuptools
	if _, errStat = os.Lstat(pipBin); errStat == nil {
		return
	}
	if _, errStat = os.Lstat(eIBin); errStat == nil {
		return
	}

	// Now try the "hard" way
	err = python.downloadExternals()
	if err != nil {
		return
	}

	err = patch.Apply(path)
	if err != nil {
		return
	}

	return
}

// Since python 3.x versions are naming their binaries with 3 affix
func (python Python) renameLinks() (err error) {
	chosen, _ := semver.Make(python.Version)
	if chosen.Major < 3 {
		return nil
	}

	path := filepath.Join(variables.Path("python", python.Version), "bin")
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

func (python Python) downloadExternals() (err error) {
	path := variables.Path("python", python.Version)

	urls, err := patch.Urls(python.Version)
	if err != nil {
		return
	}

	if hasTools(python.Version) == false {
		urls = append(urls, setuptoolsUrl, pipUrl)
	}

	respch, err := grab.GetBatch(len(urls), path, urls...)
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

// Since 2.7.9 versions we can simplify pip and setuptools install
func hasTools(version string) bool {
	semverVersion, _ := semver.Make(version)

	return semverVersion.Compare(pipAvailable) != -1
}

func checkErrors(out []byte) (err error) {
	output := string(out)

	if strings.Contains(output, "Traceback") {
		err = errors.New(output)
	}

	return err
}

func checkDependencies() (err error) {
	if runtime.GOOS == "linux" {
		return dealWithLinuxShell()
	}

	return dealWithOSXShell()
}
