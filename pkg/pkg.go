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
