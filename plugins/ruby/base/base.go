package base

import (
	"github.com/chuckpreslar/emission"

	"github.com/markelog/eclectica/pkg"
)

var (
	bins = []string{"erb", "gem", "irb", "rake", "rdoc", "ri", "ruby"}
	dots = []string{".ruby-version"}
)

// Ruby is base struct for the rest of the Ruby plugin related structs
type Ruby struct {
	Version string
	Emitter *emission.Emitter
	pkg.Base
}

// Bins returns list of the all bins included
// with the distribution of the language
func (ruby Ruby) Bins() []string {
	return bins
}

// Dots returns list of the all available filenames
// which can define versions
func (ruby Ruby) Dots() []string {
	return dots
}
