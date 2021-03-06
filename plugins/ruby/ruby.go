// Package ruby provides all needed logic for installation of ruby
package ruby

import (
	"runtime"

	"github.com/chuckpreslar/emission"

	"github.com/markelog/eclectica/pkg"

	"github.com/markelog/eclectica/plugins/ruby/bin"
	"github.com/markelog/eclectica/plugins/ruby/compile"
)

// New returns either compile or bin Ruby struct
func New(version string, emitter *emission.Emitter) pkg.Pkg {
	if hasBin(version, emitter) {
		return bin.New(version, emitter)
	}

	return compile.New(version, emitter)
}

func hasBin(version string, emitter *emission.Emitter) bool {
	bin := bin.New(version, emitter)

	if runtime.GOOS == "darwin" {
		return false
	}

	remotes, err := bin.ListRemote()
	if err != nil {
		return false
	}

	for _, remote := range remotes {
		if version == remote {
			return true
		}
	}

	return false
}
