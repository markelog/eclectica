package list

import (
	"strings"

	"github.com/markelog/list"
)

// List the plugins
func List(header string, options []string, indent int) string {
	strIndent := ""

	if indent > 0 {
		strIndent = strings.Repeat(" ", indent)
	}

	l := list.New(strIndent+header, options)

	if indent > 0 {
		l.SetIndent(indent + len(header)/2)
	}

	// Show the list
	l.Show()

	// Waiting for the user input
	return l.Get()
}
