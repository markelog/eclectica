package compile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/markelog/archive"
	"gopkg.in/cavaliercoder/grab.v1"

	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/plugins/ruby/rvm"
	"github.com/markelog/eclectica/variables"
)

var (
	versionLink = "https://rvm.io/binaries"
)

func supportBin() string {
	support := variables.Support()

	return filepath.Join(support, "ruby/2.1.5/bin/ruby")
}

func binRuby() (bin string, err error) {
	var (
		tmp        = variables.TempDir()
		archived   = filepath.Join(tmp, "ruby-2.1.5.tar.bz2")
		unarchived = filepath.Join(tmp, "ruby-2.1.5")
		support    = filepath.Join(variables.Support(), "ruby")
		dest       = filepath.Join(support, "2.1.5")
	)

	bin = supportBin()

	_, err = io.CreateDir(support)
	if err != nil {
		return
	}

	if _, statErr := os.Stat(bin); statErr == nil {
		return
	}

	url := fmt.Sprintf("%s/ruby-2.1.5.tar.bz2", rvm.GetUrl(versionLink))
	_, err = grab.Get(tmp, url)
	if err != nil {
		err = errors.New(err)
		return
	}

	err = archive.Extract(archived, tmp)
	if err != nil {
		err = errors.New(err)
		return
	}

	err = os.Rename(unarchived, dest)
	if err != nil {
		err = errors.New(err)
		return
	}

	err = rvm.RemoveArtefacts(dest)
	if err != nil {
		err = errors.New(err)
		return
	}

	return
}
