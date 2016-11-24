package print

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/tj/go-spin"
)

type SpinnerFn func()

type Spinner struct {
	channel chan bool
	Prefix  SpinnerFn
	Postfix SpinnerFn
	Before  SpinnerFn
	After   SpinnerFn
}

func (spinner *Spinner) Start() {
	if os.Getenv("EC_WITHOUT_SPINNER") == "true" {
		return
	}

	s := spin.New()

	spinner.channel = make(chan bool)

	go func() {
		spinner.Before()

		for {
			spinner.Prefix()

			select {
			case <-spinner.channel:
				spinner.After()
				return
			default:
				color.Set(color.FgCyan)
				fmt.Print(s.Next())
				color.Unset()

				spinner.Postfix()
			}
		}
	}()
}

func (spinner *Spinner) Stop() {
	if os.Getenv("EC_WITHOUT_SPINNER") == "true" {
		return
	}

	spinner.channel <- true
}
