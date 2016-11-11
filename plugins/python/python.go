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
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/blang/semver"
	"github.com/markelog/cprf"

	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/variables"
	"github.com/markelog/eclectica/versions"
)

var (
	VersionsLink   = "https://www.python.org/ftp/python"
	versionPattern = "^\\d+\\.\\d+(?:\\.\\d+)?(?:(alpha|beta|rc)(?:\\d*)?)?"

	// Hats off to inconsistent python developers
	noNilVersions, _ = semver.Make("3.3.0")

	bins = []string{"2to3", "idle", "pydoc", "python", "python-config", "easy_install", "pip"}
	dots = []string{".python-version"}
)

type Python struct {
	Version string
}

func (python Python) Install() error {
	prefix := variables.Prefix("python")
	path := variables.Path("python", python.Version)
	tmp := filepath.Join(prefix, "tmp")
	configure := filepath.Join(tmp, "configure")

	// Just in case, tmp might not get removed if this method had an error
	// before we could remove it
	os.RemoveAll(tmp)

	_, err := io.CreateDir(tmp)
	if err != nil {
		return err
	}

	fmt.Println(path+"/", tmp)

	err = cprf.Copy(path+"/", tmp)
	if err != nil {
		return err
	}

	// err = os.RemoveAll(path)
	if err != nil {
		return err
	}

	_, err = io.CreateDir(path)
	if err != nil {
		return err
	}

	fmt.Println("Configuring")
	fmt.Println(path)

	err = command(configure, "--prefix="+path)

	if err != nil {
		os.RemoveAll(tmp)
		return err
	}

	fmt.Println("Preparing")

	err = command("make")

	if err != nil {
		os.RemoveAll(tmp)
		return err
	}

	fmt.Println("Installing")

	err = command("make", "install")

	if err != nil {
		os.RemoveAll(tmp)
		return err
	}

	os.RemoveAll(tmp)

	chosen, err := semver.Make(python.Version)
	if err != nil {
		return err
	}

	if chosen.Major < 3 {
		return nil
	}
	//
	// Since python 3.x versions are naming their binaries with 3 affix
	err = renameLinks(python.Version)
	if err != nil {
		return err
	}

	return nil
}

func (python Python) PostInstall() error {
	return nil
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

	result["version"] = version
	result["extension"] = "tgz"
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
	out, _ := exec.Command(bin, "-V").Output()
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
	version := regexp.MustCompile(versionPattern)

	links := doc.Find("a")

	for i := range links.Nodes {
		content := links.Eq(i).Text()

		content = strings.Replace(content, "/", "", 1)
		if version.MatchString(content) {
			result = append(result, content)
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
				fmt.Println(absPath, newPath)
				return
			}
		}
	}

	return nil
}

func command(args ...interface{}) (err error) {
	var (
		cmd    *exec.Cmd
		errbuf bytes.Buffer
	)

	tmp := filepath.Join(variables.Prefix("python"), "tmp")

	// Required for some reason
	env := append(os.Environ(), "LC_ALL=C")

	if len(args) == 1 {
		fmt.Println(args, "command")
		cmd = exec.Command(args[0].(string))
	} else {
		fmt.Println(args, "command")
		cmd = exec.Command(args[0].(string), args[1].(string))
	}

	// Lots of needless, weird warnings in the Makefile
	cmd.Stderr = &errbuf
	cmd.Env = env
	cmd.Dir = tmp
	_, err = cmd.Output()

	if err != nil {
		return errors.New(errbuf.String())
	}

	return
}
