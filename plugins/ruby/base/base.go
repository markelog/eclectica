package base

import (
	"github.com/chuckpreslar/emission"
)

var (
	bins = []string{"erb", "gem", "irb", "rake", "rdoc", "ri", "ruby"}
	dots = []string{".ruby-version"}
)

type Ruby struct {
	Version string
	Emitter *emission.Emitter
}

func (ruby Ruby) PreDownload() error {
	return nil
}

func (ruby Ruby) PreInstall() error {
	return nil
}

func (ruby Ruby) Install() error {
	return nil
}

func (ruby Ruby) PostInstall() error {
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

func (ruby Ruby) Info() (result map[string]string) {
	return
}

func (ruby Ruby) Bins() []string {
	return bins
}

func (ruby Ruby) Dots() []string {
	return dots
}

func (ruby Ruby) ListRemote() (result []string, err error) {
	return
}
