package root

import (
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/markelog/eclectica/cmd/flags"
	"github.com/markelog/eclectica/cmd/info"
	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/plugins"
	"github.com/markelog/eclectica/versions"
)

// Command aliases
var aliases = []string{"eclectica"}

// Command config
var Command = &cobra.Command{
	Use:     "ec [<language>@<version>]",
	Aliases: aliases,
	Example: example,
}

// Event type handler
type handleFn func()

// Entry point
func Execute() {
	var (
		err  error
		args = os.Args[1:]
	)

	// If `--remote` or `-r` flag was passed (should go before any other instructions)
	flags.Parse()

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

	language, version := info.GetLanguage(args)
	hasLanguage := info.HasLanguage(args)
	hasVersion := info.HasVersion(args)

	// In case of `ec <language>@<partial-version like node@5>`
	if hasVersion && versions.IsPartial(version) {
		print.InStyleln("Language", language)

		version = getVersion(language, version)

		print.InStyleln("Version", version)

		install(language, version)
		return

		// In case of `ec <language>@<version>`
	} else if hasVersion {
		print.InStyleln("Language", language)
		print.InStyleln("Version", version)

		install(language, version)
		return
	}

	if flags.IsRemote {

		// In case of `ec -r`
		if hasLanguage {
			print.InStyleln("Language", language)

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
		print.InStyleln("Language", language)

		version, err = info.AskVersion(language)
		print.Error(err)

		install(language, version)
		return
	}

	Command.Execute()

	// We already know it will show an error
	os.Exit(1)
}

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
		err     error
		spinner *print.Spinner
	)

	handle := func(note string) handleFn {
		return func() {
			if spinner != nil {
				spinner.Stop()
			}

			time.Sleep(300 * time.Millisecond)
			spinner = print.CustomSpin("Version", plugin.Version, note)
			spinner.Start()
		}
	}

	plugin.Events().On("configure", handle("configuring"))
	plugin.Events().On("prepare", handle("preparing"))
	plugin.Events().On("install", handle("installing"))
	plugin.Events().On("post-install", handle("post-installing"))

	plugin.Events().On("done", func() {
		if spinner == nil {
			return
		}

		spinner.Stop()
		spinner = nil
	})

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
