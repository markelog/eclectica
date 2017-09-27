package compile

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chuckpreslar/emission"
	"github.com/go-errors/errors"

	"github.com/markelog/eclectica/console"
	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/plugins/ruby/base"
	"github.com/markelog/eclectica/variables"
	"github.com/markelog/eclectica/versions"
)

var (
	VersionLink    = "https://cache.ruby-lang.org/pub/ruby"
	versionPattern = "\\d+\\.\\d+\\.\\d+"
)

type Ruby struct {
	Version string
	Emitter *emission.Emitter
	base.Ruby
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

func (ruby Ruby) PreInstall() (err error) {
	// Just in case
	install := variables.InstallLanguage("ruby", ruby.Version)
	parent := filepath.Dir(install)
	current := variables.Path("ruby", ruby.Version)

	os.RemoveAll(install)

	if _, err = os.Stat(install); err != nil {
		_, err = io.CreateDir(parent)
		if err != nil {
			return
		}

		err = os.Rename(current, install)
		if err != nil {
			return err
		}
	}

	if runtime.GOOS != "linux" {
		return nil
	}

	return dealWithLinuxShell()
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

	return ruby.install()
}

func (ruby Ruby) PostInstall() (err error) {
	return os.RemoveAll(filepath.Join(variables.InstallPath(), "ruby"))
}

func (ruby Ruby) Rollback() (err error) {
	return os.RemoveAll(filepath.Join(variables.InstallPath(), "ruby"))
}

func (ruby Ruby) Info() map[string]string {
	result := make(map[string]string)

	result["filename"] = fmt.Sprintf("ruby-%s", remoteMap(ruby.Version))
	result["url"] = fmt.Sprintf("%s/%s.tar.gz", VersionLink, result["filename"])

	return result
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
	var (
		path      = variables.InstallLanguage("ruby", ruby.Version)
		configure = filepath.Join(path, "configure")
		prefix    = "--prefix=" + variables.Path("ruby", ruby.Version)
		baseruby  = "--with-baseruby="
	)

	ruby.Emitter.Emit("configure")

	bin, err := binRuby()
	if err != nil {
		return
	}
	baseruby = baseruby + bin

	cmd, stdErr, stdOut := ruby.getCmd(configure, prefix, baseruby)

	err = cmd.Run()
	if err != nil {
		return console.GetError(err, stdErr, stdOut)
	}

	return
}

func (ruby Ruby) prepare() (err error) {
	ruby.Emitter.Emit("prepare")

	cmd, stdErr, stdOut := ruby.getCmd("make", "-j", "2")

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
	cmd.Dir = variables.InstallLanguage("ruby", ruby.Version)
	cmd.Stderr = &stdErr
	cmd.Stdout = &stdOut

	if variables.IsDebug() {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	return cmd, &stdOut, &stdErr
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

func remoteMap(version string) string {
	proper := versions.Semverify(version)

	if name, ok := remoteVersions[proper]; ok {
		return name
	}

	return version
}
