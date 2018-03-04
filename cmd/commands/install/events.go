package install

import (
	"github.com/markelog/eclectica/cmd/print/custom-spinner"
	"github.com/markelog/eclectica/plugins"
)

// SetupEvents sets events
func SetupEvents(plugin *plugins.Plugin) {
	var spinner *CustomSpinner.Spin

	handle := func(note string) handleFn {
		return func(args ...string) {
			var (
				message string
			)

			if len(args) > 0 {
				message = args[0]
			}

			if spinner != nil {
				spinner.Set(&CustomSpinner.SpinArgs{
					Item:    plugin.Version,
					Note:    note,
					Message: message,
				})

				return
			}

			spinner = CustomSpinner.New(&CustomSpinner.SpinArgs{
				Header:  "version",
				Item:    plugin.Version,
				Note:    note,
				Message: message,
			})
			spinner.Start()
		}
	}

	plugin.Events().On("configure", handle("configure"))
	plugin.Events().On("prepare", handle("prepare"))
	plugin.Events().On("install", handle("install"))
	plugin.Events().On("post-install", handle("post-install"))
	plugin.Events().On("reapply modules", handle("reapply global modules"))

	plugin.Events().On("done", func() {
		if spinner == nil {
			return
		}

		spinner.Stop()
		spinner = nil
	})
}
