package print

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/tj/go-spin"
)

type fn func()

type Spinner struct {
	channel chan bool
	Prefix  fn
	Postfix fn
	Before  fn
	After   fn
}

func (spinner *Spinner) Start() {
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
	spinner.channel <- true
}
