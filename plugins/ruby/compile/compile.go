// Package compile provides ruby compilation plugin
package compile

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/chuckpreslar/emission"
	"github.com/go-errors/errors"
	"github.com/kr/pty"

	"github.com/markelog/eclectica/console"
	eIO "github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/plugins/ruby/base"
	eStrings "github.com/markelog/eclectica/strings"
	"github.com/markelog/eclectica/variables"
	"github.com/markelog/eclectica/versions"
)

var (
	// VersionLink is the URL link from which we can get all possible versions
	VersionLink = "https://cache.ruby-lang.org/pub/ruby"

	versionHref    = `\d+\.\d+\.\d+\.tar\.gz`
	versionPattern = `\d+\.\d+\.\d+`
)

// Ruby compile essential struct
type Ruby struct {
	base.Ruby
	Version   string
	Emitter   *emission.Emitter
	waitGroup *sync.WaitGroup
}

// New returns language struct
func New(version string, emitter *emission.Emitter) *Ruby {
	return &Ruby{
		Version:   version,
		Emitter:   emitter,
		waitGroup: &sync.WaitGroup{},
	}
}

// Events returns language related event emitter
func (ruby Ruby) Events() *emission.Emitter {
	return ruby.Emitter
}

// PreInstall hook
func (ruby Ruby) PreInstall() (err error) {
	install := variables.InstallLanguage("ruby", ruby.Version)
	parent := filepath.Dir(install)
	current := variables.Path("ruby", ruby.Version)

	// Just in case
	os.RemoveAll(install)

	if _, err = os.Stat(install); err != nil {
		_, err = eIO.CreateDir(parent)
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

// Install hook
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

// PostInstall hook
func (ruby Ruby) PostInstall() (err error) {
	return os.RemoveAll(filepath.Join(variables.InstallPath(), "ruby"))
}

// Rollback hook
func (ruby Ruby) Rollback() (err error) {
	return os.RemoveAll(filepath.Join(variables.InstallPath(), "ruby"))
}

// Info provides all the info needed for installation of the plugin
func (ruby Ruby) Info() map[string]string {
	result := make(map[string]string)

	result["filename"] = fmt.Sprintf("ruby-%s", remoteMap(ruby.Version))
	result["url"] = fmt.Sprintf("%s/%s.tar.gz", VersionLink, result["filename"])

	return result
}

// ListRemote returns list of the all available remote versions
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
	rHref := regexp.MustCompile(versionHref)
	rVersion := regexp.MustCompile(versionPattern)

	doc.Find("a").Each(func(i int, node *goquery.Selection) {
		href, _ := node.Attr("href")

		if rHref.MatchString(href) == false {
			return
		}

		version := rVersion.FindAllStringSubmatch(href, -1)[0][0]

		if _, ok := tmp[version]; ok == false {
			tmp[version] = true
		}
	})

	for key := range tmp {
		result = append(result, key)
	}

	return result, nil
}

func (ruby Ruby) configure() (err error) {
	ruby.Emitter.Emit("configure")

	cmd, stderr, stdout, err := ruby.configureArgs()
	if err != nil {
		return
	}

	ruby.listen("configure", stderr, false)
	ruby.listen("configure", stdout, true)

	err = cmd.Start()
	if err != nil {
		return console.Error(err, stderr, stdout)
	}

	cmd.Wait()

	return
}

func (ruby Ruby) prepare() (err error) {
	ruby.Emitter.Emit("prepare")

	cmd, stderr, stdout, err := ruby.getCmd("make", "-j")
	if err != nil {
		return
	}

	ruby.listen("prepare", stderr, false)
	ruby.listen("prepare", stdout, true)

	err = cmd.Start()
	if err != nil {
		return console.Error(err, stderr, stdout)
	}

	ruby.waitGroup.Wait()
	cmd.Wait()

	return
}

func (ruby Ruby) install() (err error) {
	ruby.Emitter.Emit("install")

	cmd, stderr, stdout, err := ruby.getCmd("make", "install", "-j")
	if err != nil {
		return
	}

	ruby.listen("install", stderr, false)
	ruby.listen("install", stdout, true)

	err = cmd.Start()
	if err != nil {
		return console.Error(err, stderr, stdout)
	}

	ruby.waitGroup.Wait()
	cmd.Wait()

	return
}

func (ruby Ruby) listen(event string, pipe io.ReadCloser, emit bool) {
	if pipe == nil {
		return
	}

	scanner := bufio.NewScanner(pipe)

	ruby.waitGroup.Add(1)
	go func() {
		defer ruby.waitGroup.Done()

		if emit == false {
			for scanner.Scan() {
			}

			return
		}

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if len(line) == 0 {
				continue
			}

			line = eStrings.ElipsisForTerminal(line)

			ruby.Emitter.Emit(event, line)
		}
	}()
}

func (ruby Ruby) getCmd(args ...string) (
	cmd *exec.Cmd,
	stderr, stdout io.ReadCloser,
	err error,
) {
	cmd = exec.Command(args[0], args[1:]...)

	cmd.Env = append(os.Environ(), "LC_ALL=C") // Required for some reason
	cmd.Dir = variables.InstallLanguage("ruby", ruby.Version)

	stderr, err = cmd.StderrPipe()
	if err != nil {
		err = errors.New(err)
		return
	}

	// In order to preserve colors output -
	// trick the command into thinking this is real tty.
	// Works properly only with "configure" command
	if path.Base(args[0]) == "configure" {
		stdout, tty, ptyErr := pty.Open()
		if ptyErr != nil {
			return nil, nil, nil, errors.New(ptyErr)
		}
		cmd.Stdout = tty
		cmd.Stdin = tty

		return cmd, stderr, stdout, nil
	}

	stdout, err = cmd.StdoutPipe()
	if err != nil {
		err = errors.New(err)
		return
	}
	return
}

func remoteMap(version string) string {
	proper := versions.Semverify(version)

	if name, ok := remoteVersions[proper]; ok {
		return name
	}

	return version
}

func (ruby Ruby) configureArgs() (
	cmd *exec.Cmd,
	stdout, stderr io.ReadCloser,
	err error,
) {
	var (
		path      = variables.InstallLanguage("ruby", ruby.Version)
		configure = filepath.Join(path, "configure")
		prefix    = "--prefix=" + variables.Path("ruby", ruby.Version)
		baseruby  = "--with-baseruby="
		shared    = "--enable-shared"
	)

	bin, err := binRuby()
	if err != nil {
		return
	}
	baseruby = baseruby + bin

	if runtime.GOOS != "darwin" {
		cmd, stdout, stderr, err = ruby.getCmd(configure, prefix, baseruby)
		return
	}

	libyaml, openssl, err := brewDependencies()
	if err != nil {
		err = errors.New(err)
		return
	}

	opensslDir := "--with-openssl-dir=" + openssl
	libyamlDir := "--with-libyaml-dir=" + libyaml

	cmd, stdout, stderr, err = ruby.getCmd(
		configure,
		prefix,
		baseruby,
		libyamlDir,
		opensslDir,
		shared,
	)
	return
}

func brewDependencies() (libyaml, openssl string, err error) {
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
