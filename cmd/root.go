package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/markelog/eclectica/cmd/info"
	"github.com/markelog/eclectica/cmd/print"
	"github.com/markelog/eclectica/plugins"
)

var isRemote bool
var isLocal bool

var RootCmd = &cobra.Command{
	Use:     "ec [<language>@<version>]",
	Aliases: []string{"eclectica"},
	Example: example,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd
func Execute() {
	var err error
	args := os.Args[1:]

	if info.NonInstallCommand(args) {
		// Initialize cobra for other commands
		if err := RootCmd.Execute(); err != nil {
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

	pflag.BoolVarP(remoteInfo())
	pflag.BoolVarP(localInfo())
	pflag.Parse()

	language, version := info.GetLanguage(args)
	hasLanguage := info.HasLanguage(args)
	hasVersion := info.HasVersion(args)

	print.LaguageOrVersion(language, version)

	// In case of `ec <language>@<version>`
	if hasLanguage && hasVersion {
		install(language, version)
		return
	}

	// If `--remote` or `-r` flag was passed
	if isRemote {

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

	RootCmd.Execute()

	// We already know it will show an error
	os.Exit(1)
}

func conditionalInstall(plugin *plugins.Plugin) {
	var err error

	if isLocal {
		err = plugin.LocalInstall()
	} else {
		err = plugin.Install()
	}

	print.Error(err)
}

func install(language, version string) {
	plugin := plugins.New(language, version)

	response, err := plugin.Download()
	print.Error(err)

	// response == nil means we already downloaded that thing
	if response != nil {
		print.Download(response, version)

		err = plugin.Extract()
		print.Error(err)
	}

	conditionalInstall(plugin)
}

func remoteInfo() (*bool, string, string, bool, string) {
	return &isRemote, "remote", "r", false, "Get remote versions"
}

func localInfo() (*bool, string, string, bool, string) {
	return &isLocal, "local", "l", false, "Install as local version"
}

func init() {
	RootCmd.SetHelpTemplate(help)
	RootCmd.SetUsageTemplate(usage)

	cobra.OnInitialize()
}
