// Package install installs the languages
package install

import (
	"os"

	"github.com/schollz/closestmatch"
	"github.com/spf13/cobra"

	"github.com/markelog/eclectica/cmd/info"
	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/plugins"
	"github.com/markelog/eclectica/versions"
)

// Is action remote?
var isRemote bool

// Is action local?
var isLocal bool

// Is action local?
var withModules bool

// Command represents the ls command
var Command = &cobra.Command{
	Use:   "install [<language>@<version>]",
	Short: "install language",
	Run:   run,

	// Rather use "install [<language>@<version>]"
	Hidden: true,
}

// Event type handler
type handleFn func(args ...string)

// getVersion gets version of the language and its correlated version
func getVersion(language, version string) string {
	remoteList, err := info.FullListRemote(language)
	print.Error(err)

	version, err = versions.Complete(version, remoteList)
	print.Error(err)

	return version
}

// Install either globally or locally
func conditionalInstall(plugin *plugins.Plugin) {
	var (
		err error
	)

	SetupEvents(plugin)

	if isLocal {
		err = plugin.LocalInstall()
	} else {
		err = plugin.Install()
	}

	print.Error(err)
	print.LastPrint()
}

// Entry point for installation
func install(language, version string) {
	plugin := plugins.New(&plugins.Args{
		Language:    language,
		Version:     version,
		WithModules: withModules,
	})

	err := plugin.PreDownload()
	print.Error(err)

	response, err := plugin.Download()
	print.Error(err)

	// response == nil means we already downloaded that thing
	if response != nil {
		print.Download(response, plugin.Version)

		err = plugin.Extract()
		print.Error(err)
	}

	conditionalInstall(plugin)
}

// Entry point
func run(cmd *cobra.Command, args []string) {
	var (
		err               error
		language, version = info.GetLanguage(args)
		hasLanguage       = info.HasLanguage(args)
		hasVersion        = info.HasVersion(args)
		cm                = closestmatch.New(plugins.Plugins, []int{2})
	)

	// We don't use cobra here, since we support `ec <language>@<version>` syntax

	// Searching for closest plugin name
	if len(args) > 0 && hasLanguage == false {
		possible := info.PossibleLanguage(args)
		print.ClosestLangWarning(possible, cm.Closest(possible))
		return
	}

	// If nothing was passed - just show list of the local versions
	if isRemote == false && hasLanguage == false {
		lang, ver, errAsk := info.Ask()
		print.Error(errAsk)

		install(lang, ver)
		return
	}

	// In case of `ec <language>@<partial-version like node@5>`
	if hasVersion && versions.IsPartial(version) {
		print.FnInStyleln("langauge:", language)
		version = getVersion(language, version)

		print.InStyleln(" version:", version)

		install(language, version)
		return

		// In case of `ec <language>@<version>`
	} else if hasVersion {
		print.FnInStyleln("langauge:", language)
		print.InStyleln(" version:", version)

		install(language, version)
		return
	}

	if isRemote {

		// In case of `ec -r`
		if hasLanguage {
			print.FnInStyleln("langauge:", language)

			version, err = info.AskRemoteVersion(language)
			print.Error(err)

			install(language, version)
			return

		}

		// In case of `ec -r <language>` or `ec <language> -r`
		language, version, err = info.AskRemote()
		print.Error(err)

		install(language, version)
		return
	}

	// In case of `ec <language>`
	if hasLanguage && hasVersion == false {
		print.FnInStyleln("langauge:", language)

		version, err = info.AskVersion(language)
		print.Error(err)

		install(language, version)
		return
	}

	// We already know it will show an error
	os.Exit(1)
}

// Init
func init() {
	flags := Command.PersistentFlags()
	flags.BoolVarP(&isRemote, "remote", "r", false, "get remote versions")
	flags.BoolVarP(&isLocal, "local", "l", false, "install as local version")
	flags.BoolVarP(&withModules, "with-modules", "w", false, "reinstall global modules from the previous version (currently works only for node.js)")
}
