package pkg

import "github.com/chuckpreslar/emission"

type Pkg interface {
	PreDownload() error
	PreInstall() error
	Install() error
	PostInstall() error
	Switch() error
	Link() error
	Events() *emission.Emitter
	Environment() ([]string, error)
	ListRemote() ([]string, error)
	Info() map[string]string
	Bins() []string
	Dots() []string
}

type Base struct {
	Version string
	Emitter *emission.Emitter
}

func (base Base) PreDownload() error {
	return nil
}

func (base Base) PreInstall() error {
	return nil
}

func (base Base) Install() error {
	return nil
}

func (base Base) PostInstall() error {
	return nil
}

func (base Base) Switch() error {
	return nil
}

func (base Base) Link() error {
	return nil
}

func (base Base) Environment() (result []string, err error) {
	return
}

func (base Base) Info() (result map[string]string) {
	return
}

func (base Base) ListRemote() (result []string, err error) {
	return
}
