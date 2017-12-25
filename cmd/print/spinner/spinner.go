// Package spinner provides essential (empty) spinner
package spinner

import (
	"fmt"
	"os"
	"sync"

	"github.com/mgutz/ansi"
	spin "github.com/tj/go-spin"
)

// Spinner essential struct
type Spinner struct {
	Before, After, Prefix, Postfix Fn

	Spin *spin.Spinner

	isDone bool
	mutex  *sync.Mutex
}

// Fn callback signature
type Fn func()

// New returns new spinner struct
func New(before, after, prefix, postfix Fn) *Spinner {
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

// Start the spinner
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

// Stop the spinner
func (spinner *Spinner) Stop() {
	spinner.mutex.Lock()
	defer spinner.mutex.Unlock()
	if spinner.isDone == true {
		return
	}

	spinner.isDone = true
	spinner.After()
}
