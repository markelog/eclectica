// Package nodejs provides all needed logic for installation of node.js
package nodejs

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/blang/semver"
	"github.com/chuckpreslar/emission"
	"github.com/go-errors/errors"

	"github.com/markelog/eclectica/pkg"
	"github.com/markelog/eclectica/plugins/nodejs/modules"
	"github.com/markelog/eclectica/variables"
)

var (
	// VersionLink is the URL link from which we can get all possible versions
	VersionLink = "https://nodejs.org/dist"

	versionPattern = "v\\d+\\.\\d+\\.\\d+$"

	minimalVersion, _ = semver.Make("0.10.0")

	bins = []string{"node", "npm", "npx"}
	dots = []string{".nvmrc", ".node-version"}
)

// Node essential struct
type Node struct {
	Version     string
	previous    string
	withModules bool
	Emitter     *emission.Emitter
	pkg.Base
}

// Args is arguments struct for New() method
type Args struct {
	Version     string
	Emitter     *emission.Emitter
	WithModules bool
}

// New returns language struct
func New(args *Args) *Node {
	return &Node{
		Version:     args.Version,
		Emitter:     args.Emitter,
		withModules: args.WithModules,
		previous:    variables.CurrentVersion("node"),
	}
}

// Events returns language related event emitter
func (node Node) Events() *emission.Emitter {
	return node.Emitter
}

// PostInstall hook
func (node Node) PostInstall() (err error) {
	node.Emitter.Emit("postinstall")

	return nil
}

// Switch hook
func (node Node) Switch() (err error) {
	previous := node.previous

	// Should not install modules if not explicitly defined by the user
	if node.withModules == false {
		return
	}

	// If there is no previous version – don't do anything
	if len(previous) == 0 {
		return
	}

	node.Emitter.Emit("reapply modules")

	err = modules.New(node.previous, node.Version).Install()
	if err != nil {
		return
	}

	return
}

// Info provides all the info needed for installation of the plugin
func (node Node) Info() map[string]string {
	result := make(map[string]string)
	sourcesURL := fmt.Sprintf("%s/v%s", VersionLink, node.Version)

	result["filename"] = fmt.Sprintf("node-v%s-%s-%s", node.Version, runtime.GOOS, nodeArchiveArch(node.Version))
	result["url"] = fmt.Sprintf("%s/%s.tar.gz", sourcesURL, result["filename"])

	return result
}

func nodeArchiveArch(version string) string {

	arch := runtime.GOARCH
	if runtime.GOOS == "darwin" && isDarwinArm64Machine() {
		arch = "arm64"
	}
	return NodeArchiveArchFor(version, runtime.GOOS, arch)

}

// NodeArchiveArchFor determines the appropriate architecture for a Node.js archive
func NodeArchiveArchFor(version, goos, arch string) string {
	if arch == "amd64" || arch == "x86_64" {
		return "x64"
	}
	if arch == "arm64" || arch == "aarch64" {
		if goos == "darwin" && !supportsDarwinArm64(version) {
			return "x64"
		}
		return "arm64"
	}
	return arch

}

func isDarwinArm64Machine() bool {
	output, err := exec.Command("sysctl", "-n", "hw.optional.arm64").Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "1"

}

func supportsDarwinArm64(version string) bool {
	parsed, err := semver.Make(version)
	if err != nil {
		return false
	}

	min, _ := semver.Make("16.0.0")
	return parsed.GTE(min)
}

// Bins returns list of the all bins included
// with the distribution of the language
func (node Node) Bins() []string {
	return bins
}

// Dots returns list of the all available filenames
// which can define versions
func (node Node) Dots() []string {
	return dots
}

// ListRemote returns list of the all available remote versions
func (node Node) ListRemote() ([]string, error) {
	doc, err := goquery.NewDocument(VersionLink)

	if err != nil {
		if _, ok := err.(net.Error); ok {
			return nil, errors.New(variables.ConnectionError)
		}

		return nil, errors.New(err)
	}

	tmp := []string{}
	result := []string{}

	rVersion := regexp.MustCompile(versionPattern)

	doc.Find("a").Each(func(i int, node *goquery.Selection) {
		href, _ := node.Attr("href")

		href = strings.Replace(href, "/", "", 1)
		if rVersion.MatchString(href) {
			href = strings.Replace(href, "v", "", 1)
			tmp = append(tmp, href)
		}
	})

	// Remove outdated versions
	for _, element := range tmp {
		version, _ := semver.Make(element)

		if version.GTE(minimalVersion) {
			result = append(result, element)
		}
	}

	return result, nil
}
