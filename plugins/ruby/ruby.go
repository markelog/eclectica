package ruby

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
	// "sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/blang/semver"
	"github.com/chuckpreslar/emission"

	"github.com/markelog/eclectica/console"
	"github.com/markelog/eclectica/variables"
)

var (
	VersionLink    = "https://cache.ruby-lang.org/pub/ruby"
	versionPattern = "\\d+\\.\\d+\\.\\d+"

	bins = []string{"erb", "gem", "irb", "rake", "rdoc", "ri", "ruby"}
	dots = []string{".ruby-version"}
)

type Ruby struct {
	Version string
	Emitter *emission.Emitter
}

func New(version string, emitter *emission.Emitter) *Ruby {
	return &Ruby{
		Version: version,
		Emitter: emitter,
	}
}

func (ruby Ruby) Events() *emission.Emitter {
	return ruby.Emitter
}

func (ruby Ruby) PreDownload() (err error) {
	return
}

func (ruby Ruby) PreInstall() error {
	return dealWithShell()
}

func (ruby Ruby) Install() (err error) {
	err = ruby.configure()
	if err != nil {
		return
	}

	err = ruby.prepare()
	if err != nil {
		return
	}

	err = ruby.install()
	if err != nil {
		return
	}

	return ruby.renameLinks()
}

func (ruby Ruby) PostInstall() (err error) {
	return nil
}

func (ruby Ruby) Switch() error {
	return nil
}

func (ruby Ruby) Link() error {
	return nil
}

func (ruby Ruby) Environment() (result []string, err error) {
	return
}

func test(version string) string {
	tmp := []string{}
	doc, _ := goquery.NewDocument(VersionLink)
	rVersion := regexp.MustCompile(versionPattern)

	doc.Find("a").Each(func(i int, node *goquery.Selection) {
		href, _ := node.Attr("href")

		if strings.Contains(href, ".tar.gz") == false {
			return
		}

		if rVersion.MatchString(href) == false {
			return
		}

		if strings.Contains(href, version) == false {
			return
		}

		if strings.Contains(href, "-rc") == true {
			return
		}

		tmp = append(tmp, href)
	})

	return tmp[len(tmp)-1]
}

func (ruby Ruby) Info() map[string]string {
	result := make(map[string]string)
	a := test(ruby.Version)

	result["filename"] = strings.Replace(a, ".tar.gz", "", -1)
	result["url"] = fmt.Sprintf("%s/%s", VersionLink, a)

	return result
}

func (rust Ruby) Bins() []string {
	return bins
}

func (rust Ruby) Dots() []string {
	return dots
}

func (ruby Ruby) ListRemote() ([]string, error) {
	doc, err := goquery.NewDocument(VersionLink)

	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, errors.New("Can't establish connection")
		}

		return nil, err
	}

	result := []string{}
	tmp := make(map[string]bool)
	rVersion := regexp.MustCompile(versionPattern)

	doc.Find("a").Each(func(i int, node *goquery.Selection) {
		href, _ := node.Attr("href")

		if strings.Contains(href, ".tar.gz") == false {
			return
		}

		if rVersion.MatchString(href) == false {
			return
		}

		version := rVersion.FindAllStringSubmatch(href, -1)[0][0]

		if _, ok := tmp[version]; ok == false {
			tmp[version] = true
		}
	})

	for key, _ := range tmp {
		result = append(result, key)
	}

	return result, nil
}

func (ruby Ruby) configure() (err error) {
	path := variables.Path("ruby", ruby.Version)
	configure := filepath.Join(path, "configure")

	ruby.Emitter.Emit("configure")

	err = ruby.externals()
	if err != nil {
		return
	}

	cmd, stdErr, stdOut := ruby.getCmd(configure, "--prefix="+path)
	cmd.Env = ruby.getEnvs(cmd.Env)

	err = cmd.Run()
	if err != nil {
		return console.GetError(err, stdErr, stdOut)
	}

	return
}

func (ruby Ruby) touch() (err error) {
	cmd, stdErr, stdOut := ruby.getCmd("make", "touch")

	err = cmd.Run()
	if err != nil {
		return console.GetError(err, stdErr, stdOut)
	}

	return
}

func (ruby Ruby) prepare() (err error) {
	ruby.Emitter.Emit("prepare")

	// Ignore touch errors since newest ruby makefile doesn't have this task
	ruby.touch()

	cmd, stdErr, stdOut := ruby.getCmd("make")

	err = cmd.Run()
	if err != nil {
		return console.GetError(err, stdErr, stdOut)
	}

	return
}

func (ruby Ruby) install() (err error) {
	ruby.Emitter.Emit("install")

	cmd, stdErr, stdOut := ruby.getCmd("make", "install")

	err = cmd.Run()
	if err != nil {
		return console.GetError(err, stdErr, stdOut)
	}

	return
}

func (ruby Ruby) getEnvs(original []string) (result []string) {
	if runtime.GOOS == "darwin" {
		result = ruby.getOSXEnvs(original)
	}

	return
}

func (ruby Ruby) getOSXEnvs(original []string) []string {
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

func (ruby Ruby) getCmd(args ...interface{}) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer) {
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
	cmd.Dir = variables.Path("ruby", ruby.Version)
	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut

	if variables.IsDebug() {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	return cmd, &stdOut, &stdErr
}

func (ruby Ruby) externals() (err error) {
	return
}

// Since ruby 3.x versions are naming their binaries with 3 affix
func (ruby Ruby) renameLinks() (err error) {
	chosen, _ := semver.Make(ruby.Version)
	if chosen.Major < 3 {
		return nil
	}

	path := filepath.Join(variables.Path("ruby", ruby.Version), "bin")
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

			// Since ruby install creates some links with version and some without
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

func getRemoteVersions() ([]string, error) {
	doc, err := goquery.NewDocument(VersionLink)

	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, errors.New("Can't establish connection")
		}

		return nil, err
	}

	result := []string{}
	tmp := make(map[string]bool)
	rVersion := regexp.MustCompile(versionPattern)

	doc.Find("a").Each(func(i int, node *goquery.Selection) {
		href, _ := node.Attr("href")

		if strings.Contains(href, ".tar.gz") == false {
			return
		}

		if rVersion.MatchString(href) == false {
			return
		}

		version := rVersion.FindAllStringSubmatch(href, -1)[0][0]

		if _, ok := tmp[version]; ok == false {
			tmp[version] = true
		}
	})

	for key, _ := range tmp {
		result = append(result, key)
	}

	return result, nil
}

func checkErrors(out []byte) (err error) {
	output := string(out)

	if strings.Contains(output, "Traceback") {
		err = errors.New(output)
	}

	return err
}
