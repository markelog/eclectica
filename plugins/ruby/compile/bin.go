package compile

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/markelog/archive"
	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/plugins/ruby/rvm"
	"github.com/markelog/eclectica/variables"
	"gopkg.in/cavaliercoder/grab.v1"
)

var (
	binUbuntuRubyUrl = "https://rvm.io/binaries/ubuntu/16.10/x86_64/ruby-2.2.5.tar.bz2"
	binOSXRubyUrl    = "https://rvm.io/binaries/osx/10.12/x86_64/ruby-2.2.5.tar.bz2"
)

func getUrl() (path string, err error) {
	if runtime.GOOS == "linux" {
		return binUbuntuRubyUrl, nil
	}

	if runtime.GOOS == "darwin" {
		return binOSXRubyUrl, nil
	}

	err = errors.New("Not supported environment")

	return "", err
}

func supportBin() string {
	support := variables.Support()

	return filepath.Join(support, "ruby/2.2.5/bin/ruby")
}

func binRuby() (bin string, err error) {
	var (
		tmp        = variables.TempDir()
		archived   = filepath.Join(tmp, "ruby-2.2.5.tar.bz2")
		unarchived = filepath.Join(tmp, "ruby-2.2.5")
		support    = filepath.Join(variables.Support(), "ruby")
		dest       = filepath.Join(support, "2.2.5")
	)

	bin = supportBin()

	_, err = io.CreateDir(support)
	if err != nil {
		return
	}

	if _, statErr := os.Stat(bin); statErr == nil {
		return
	}

	binRubyUrl, err := getUrl()
	if err != nil {
		return
	}

	_, err = grab.Get(tmp, binRubyUrl)
	if err != nil {
		return
	}

	err = archive.Extract(archived, tmp)
	if err != nil {
		return
	}

	err = os.Rename(unarchived, dest)
	if err != nil {
		return
	}

	err = rvm.RemoveArtefacts(dest)
	if err != nil {
		return
	}

	return
}
