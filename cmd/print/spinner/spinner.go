package spinner

import (
	"fmt"
	"os"
	"sync"

	"github.com/fatih/color"
	spin "github.com/tj/go-spin"
)

type SpinnerFn func()

type Spinner struct {
	isDone  bool
	mutex   *sync.Mutex
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
		spinner.mutex = &sync.Mutex{}
		spinner.Spin = spin.New()
	}

	spinner.isDone = false

	spinner.Before()

	go func() {
		for {
			// spinner.mutex.Lock()

			if spinner.isDone {
				break
			}

			spinner.Prefix()

			color.Set(color.FgCyan)
			fmt.Print(spinner.Spin.Next())
			color.Unset()

			spinner.Postfix()
		}

		// spinner.mutex.Unlock()
	}()
}

func (spinner *Spinner) Stop() {
	if spinner.isDone == true {
		return
	}

	spinner.isDone = true
	spinner.After()
}
