package base

import (
	"github.com/chuckpreslar/emission"

	"github.com/markelog/eclectica/pkg"
)

var (
	bins = []string{"erb", "gem", "irb", "rake", "rdoc", "ri", "ruby"}
	dots = []string{".ruby-version"}
)

type Ruby struct {
	Version string
	Emitter *emission.Emitter
	pkg.Base
}

func (ruby Ruby) Bins() []string {
	return bins
}

func (ruby Ruby) Dots() []string {
	return dots
}
