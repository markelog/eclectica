package flags

import (
	"github.com/spf13/pflag"
)

// Is action remote?
var IsRemote bool

// Is action local?
var IsLocal bool

// Get config for remote flag
func RemoteFlag() (*bool, string, string, bool, string) {
	return &IsRemote, "remote", "r", false, "Get remote versions"
}

// Get config for local flag
func LocalFlag() (*bool, string, string, bool, string) {
	return &IsLocal, "local", "l", false, "Install as local version"
}

func Parse() {
	pflag.BoolVarP(RemoteFlag())
	pflag.BoolVarP(LocalFlag())
	pflag.Parse()
}
