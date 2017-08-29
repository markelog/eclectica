package compile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/markelog/archive"
	"github.com/markelog/eclectica/io"
	"github.com/markelog/eclectica/plugins/ruby/rvm"
	"github.com/markelog/eclectica/variables"
	"gopkg.in/cavaliercoder/grab.v1"
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

	url := fmt.Sprintf("%s/ruby-2.1.5.tar.bz2", rvm.GetUrl(VersionLink))

	if err != nil {
		return
	}

	_, err = grab.Get(tmp, url)
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
