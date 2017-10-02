package compile

import (
	"bufio"
	"fmt"
	"io"
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
	"github.com/kr/pty"

	"github.com/markelog/eclectica/console"
	eio "github.com/markelog/eclectica/io"
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
	install := variables.InstallLanguage("ruby", ruby.Version)
	parent := filepath.Dir(install)
	current := variables.Path("ruby", ruby.Version)

	// Just in case
	os.RemoveAll(install)

	if _, err = os.Stat(install); err != nil {
		_, err = eio.CreateDir(parent)
		if err != nil {
			return
		}

		err = os.Rename(current, install)
		if err != nil {
			return err
		}
	}

	if runtime.GOOS == "linux" {
		return dealWithLinuxShell()
	}

	if runtime.GOOS == "darwin" {
		return dealWithOSXShell()
	}

	return
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
			return nil, errors.New(variables.ConnectionError)
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
	ruby.Emitter.Emit("configure")

	err, cmd, stderr, stdout := ruby.configureArgs()
	if err != nil {
		return
	}

	ruby.listen("configure", stdout)

	err = cmd.Run()
	if err != nil {
		return console.GetError(err, stderr, stdout)
	}

	return
}

func (ruby Ruby) prepare() (err error) {
	ruby.Emitter.Emit("prepare")

	err, cmd, stderr, stdout := ruby.getCmd("make", "-j", "2")
	if err != nil {
		return
	}

	ruby.listen("prepare", stdout)

	err = cmd.Run()
	if err != nil {
		return console.GetError(err, stderr, stdout)
	}

	return
}

func (ruby Ruby) install() (err error) {
	ruby.Emitter.Emit("install")

	err, cmd, stderr, stdout := ruby.getCmd("make", "install")
	if err != nil {
		return
	}

	ruby.listen("install", stdout)

	err = cmd.Run()
	if err != nil {
		return console.GetError(err, stderr, stdout)
	}

	return
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

func (ruby Ruby) listen(event string, pipe io.ReadCloser) {
	if pipe == nil {
		return
	}

	scanner := bufio.NewScanner(pipe)
	go func() {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if len(line) == 0 {
				continue
			}

			line = truncateString(line, 80)

			ruby.Emitter.Emit(event, line)
		}
	}()
}

func (ruby Ruby) getCmd(args ...string) (err error, cmd *exec.Cmd, stderr, stdout io.ReadCloser) {
	cmd = exec.Command(args[0], args[1:]...)

	cmd.Env = append(os.Environ(), "LC_ALL=C") // Required for some reason
	cmd.Dir = variables.InstallLanguage("ruby", ruby.Version)

	if variables.IsDebug() {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return
	}

	stderr, err = cmd.StderrPipe()
	if err != nil {
		err = errors.New(err)
		return
	}

	// In order to preserve colors output -
	// trick the command into thinking this is real tty
	stdout, tty, err := pty.Open()
	if err != nil {
		err = errors.New(err)
		return
	}
	cmd.Stdout = tty

	return
}

func getRemoteVersions() ([]string, error) {
	doc, err := goquery.NewDocument(VersionLink)

	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, errors.New(variables.ConnectionError)
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

func (ruby Ruby) configureArgs() (err error, cmd *exec.Cmd, out, outErr io.ReadCloser) {
	var (
		path      = variables.InstallLanguage("ruby", ruby.Version)
		configure = filepath.Join(path, "configure")
		prefix    = "--prefix=" + variables.Path("ruby", ruby.Version)
		baseruby  = "--with-baseruby="
	)

	bin, err := binRuby()
	if err != nil {
		return
	}
	baseruby = baseruby + bin

	if runtime.GOOS != "darwin" {
		err, cmd, out, outErr = ruby.getCmd(configure, prefix, baseruby)
		return
	}

	err, libyaml, openssl := brewDependencies()
	if err != nil {
		return
	}

	opensslDir := "--with-openssl-dir=" + openssl
	libyamlDir := "--with-libyaml-dir=" + libyaml

	err, cmd, out, outErr = ruby.getCmd(configure, prefix, baseruby, libyamlDir, opensslDir)
	return
}

func brewDependencies() (err error, libyaml string, openssl string) {
	out, err := exec.Command("brew", "--prefix", "libyaml").Output()
	libyaml = strings.TrimSpace(string(out))
	if err != nil {
		return
	}

	out, err = exec.Command("brew", "--prefix", "openssl").Output()
	openssl = strings.TrimSpace(string(out))
	if err != nil {
		return
	}

	return
}
