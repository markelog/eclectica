package spinner

import (
	"fmt"
	"os"
	"sync"

	"github.com/mgutz/ansi"
	spin "github.com/tj/go-spin"
)

type SpinnerFn func()

type Spinner struct {
	isDone bool

	mutex *sync.Mutex
	Spin  *spin.Spinner

	Before, After, Prefix, Postfix SpinnerFn
}

func New(before, after, prefix, postfix SpinnerFn) *Spinner {
	return &Spinner{
		isDone: false,

		mutex: &sync.Mutex{},
		Spin:  spin.New(),

		Before:  before,
		After:   after,
		Prefix:  prefix,
		Postfix: postfix,
	}
}

func (spinner *Spinner) Start() {
	spinner.mutex.Lock()
	defer spinner.mutex.Unlock()

	if os.Getenv("EC_WITHOUT_SPINNER") == "true" {
		spinner.isDone = true
		return
	}

	spinner.isDone = false

	spinner.Before()

	go func() {
		for {
			spinner.mutex.Lock()

			if spinner.isDone {
				break
			}

			spinner.Prefix()

			fmt.Print(ansi.Color(spinner.Spin.Next(), "cyan"))

			spinner.Postfix()

			spinner.mutex.Unlock()
		}

		spinner.mutex.Unlock()
	}()
}

func (spinner *Spinner) Stop() {
	spinner.mutex.Lock()
	defer spinner.mutex.Unlock()
	if spinner.isDone == true {
		return
	}

	spinner.isDone = true
	spinner.After()
}
