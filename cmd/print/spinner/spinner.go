package spinner

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	spin "github.com/tj/go-spin"
)

type SpinnerFn func()

type Spinner struct {
	isDone  bool
	channel chan bool
	Spin    *spin.Spinner
	Prefix  SpinnerFn
	Postfix SpinnerFn
	Before  SpinnerFn
	After   SpinnerFn
}

func (spinner *Spinner) Start() {
	if os.Getenv("EC_WITHOUT_SPINNER") == "true" {
		spinner.isDone = true
		return
	}

	if spinner.Spin == nil {
		spinner.Spin = spin.New()
	}

	spinner.isDone = false

	spinner.Before()

	go func() {
		for spinner.isDone == false {
			spinner.Prefix()

			color.Set(color.FgCyan)
			fmt.Print(spinner.Spin.Next())
			color.Unset()

			spinner.Postfix()
		}
	}()
}

func (spinner *Spinner) Stop() {
	if spinner.isDone == true {
		return
	}

	spinner.isDone = true
	spinner.After()
}
