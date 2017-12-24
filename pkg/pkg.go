// Package provides helpful base interfaces and struct definitions
package pkg

import "github.com/chuckpreslar/emission"

// Pkg plugin interface
type Pkg interface {
	PreDownload() error
	PreInstall() error
	Install() error
	PostInstall() error
	Switch() error
	Link() error
	Rollback() error
	Events() *emission.Emitter
	Environment() ([]string, error)
	ListRemote() ([]string, error)
	Info() map[string]string
	Bins() []string
	Dots() []string
}

// Base struct from which every plugin should inherit
type Base struct {
	Version string
	Emitter *emission.Emitter
}

// PreDownload hook
func (base Base) PreDownload() error {
	return nil
}

// PreInstall hook
func (base Base) PreInstall() error {
	return nil
}

// Install hook
func (base Base) Install() error {
	return nil
}

// PostInstall hook
func (base Base) PostInstall() error {
	return nil
}

// Switch hook
func (base Base) Switch() error {
	return nil
}

// Link hook
func (base Base) Link() error {
	return nil
}

// Rollback hook
func (base Base) Rollback() error {
	return nil
}

// Environment get needed environment variables and values
func (base Base) Environment() (result []string, err error) {
	return
}

// Info gets plugin related info
func (base Base) Info() (result map[string]string) {
	return
}

// ListRemote gets availiable remote version for this plugin
func (base Base) ListRemote() (result []string, err error) {
	return
}
