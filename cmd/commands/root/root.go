package root

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/markelog/eclectica/cmd/flags"
	"github.com/markelog/eclectica/cmd/info"
	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/plugins"
)

// Command aliases
var aliases = []string{"eclectica"}

// Command config
var Command = &cobra.Command{
	Use:     "ec [<language>@<version>]",
	Aliases: aliases,
	Example: example,
}

// Entry point
func Execute() {
	var err error
	args := os.Args[1:]

	if info.NonInstallCommand(args) {

		// Initialize cobra for other commands
		if err := Command.Execute(); err != nil {
			os.Exit(1)
		}

		return
	}

	// We don't use cobra here, since we support `ec <language>@<version>` syntax

	// If nothing was passed - just show list of the local versions
	if len(args) == 0 {
		language, version, err := info.Ask()
		print.Error(err)

		install(language, version)
		return
	}

	flags.Parse()

	language, version := info.GetLanguage(args)
	hasLanguage := info.HasLanguage(args)
	hasVersion := info.HasVersion(args)

	print.InStyleln("Language", language)

	// In case of `ec <language>@<version>`
	if hasLanguage && hasVersion {
		install(language, version)
		return
	}

	// If `--remote` or `-r` flag was passed
	if flags.IsRemote {

		// In case of `ec -r`
		if hasLanguage {
			version, err = info.AskRemoteVersion(language)
			print.Error(err)

			install(language, version)
			return

			// In case of `ec -r <language>` or `ec <language> -r`
		} else {
			language, version, err = info.AskRemote()
			print.Error(err)

			install(language, version)
			return
		}
	}

	// In case of `ec <language>`
	if hasLanguage && hasVersion == false {
		version, err = info.AskVersion(language)
		print.Error(err)

		install(language, version)
		return
	}

	Command.Execute()

	// We already know it will show an error
	os.Exit(1)
}

// Install either globally or locally
func conditionalInstall(plugin *plugins.Plugin) {
	var err error

	if flags.IsLocal {
		err = plugin.LocalInstall()
	} else {
		err = plugin.Install()
	}

	print.Error(err)
}

// Entry point for installation
func install(language, version string) {
	plugin := plugins.New(language, version)

	remoteList, err := info.FullListRemote(language)
	print.Error(err)

	err = plugin.SetFullVersion(remoteList)
	print.Error(err)

	print.InStyleln("Version", plugin.Version)

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

// Add command to root command
func Register(cmd *cobra.Command) {
	Command.AddCommand(cmd)
}

// Init
func init() {
	Command.SetHelpTemplate(help)
	Command.SetUsageTemplate(usage)

	cobra.OnInitialize()
}
