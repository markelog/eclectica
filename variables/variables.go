package variables

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/markelog/eclectica/io"
)

var (

	// DefaultInstall is a default path to the general bin folder
	DefaultInstall = filepath.Join(Base(), "bin")

	// ConnectionError is a general connection error text
	ConnectionError = "Connection cannot be established"
)

// TempDir gets OS consistent folder path
// I am crying over here :/
func TempDir() (tmp string) {
	tmp = os.TempDir()
	if runtime.GOOS == "linux" {
		tmp += "/"
	}

	return
}

// IsDebug checks if eclectica in the debug state
// i.e. will print more info when executing commands
func IsDebug() bool {
	return os.Getenv("EC_DEBUG") == "true"
}

// GetBin returns path to the bin folder of the provided language
func GetBin(args ...interface{}) string {
	name, version := nameAndVersion(args)

	base := Path(name, version)

	// FIXME: should look better somehow
	if name == "rust" {
		name = "rustc"
	}

	return filepath.Join(base, "bin", name)
}

// GetShellName gets name of the used shell
func GetShellName() string {
	path := GetShellPath()
	parts := strings.Split(path, "/")

	return parts[len(parts)-1]
}

// GetShellPath gets path to current shell binary
func GetShellPath() string {
	path := os.Getenv("SHELL")

	if len(path) == 0 {
		return "/bin/bash"
	}

	return path
}

// Base provides path where eclectica stores everything filesystem related
func Base() string {
	usr, _ := user.Current()
	username := usr.Username

	if username == "root" {
		usr, _ = user.Lookup(os.Getenv("SUDO_USER"))
	}

	return filepath.Join(usr.HomeDir, ".eclectica")
}

// Prefix gets path to the language install folder
func Prefix(name string) string {
	return filepath.Join(Home(), name)
}

// Path gets path of the install folder for the
// specific version or to the current one
func Path(args ...interface{}) string {
	name, version := nameAndVersion(args)

	return filepath.Join(Home(), name, version)
}

// Home gets path to the place where eclectica installs their languages
func Home() string {
	return filepath.Join(Base(), "versions")
}

// Support get path to support folder
func Support() string {
	return filepath.Join(Base(), "support")
}

// InstallPath get path to install folder
func InstallPath() string {
	return filepath.Join(Support(), "install")
}

// InstallLanguage get path to dist folder in install folder
func InstallLanguage(language, version string) string {
	return filepath.Join(InstallPath(), language, version)
}

// nameAndVersion helper method to get name and version
func nameAndVersion(args []interface{}) (string, string) {
	var (
		name    = args[0].(string)
		version string
	)

	if len(args) == 2 {
		version = args[1].(string)
	} else {
		version = "current"
	}

	return name, version
}

// CurrentVersion get current version for the specific language
func CurrentVersion(name string) string {
	base := Path(name)
	path := filepath.Join(base, ".eclectica")

	return io.Read(path)
}

// WriteVersion writes version to the language install folder path
// under the name ".eclectica"
func WriteVersion(name, version string) error {
	base := Path(name, version)
	path := filepath.Join(base, ".eclectica")

	return io.WriteFile(path, version)
}

// IsInstalled checks if this version was already installed
func IsInstalled(name, version string) bool {
	base := Path(name, version)
	path := filepath.Join(base, ".eclectica")

	// If binary for this plugin already exist then we can assume it was installed before;
	// which means we can bail out this point
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}
