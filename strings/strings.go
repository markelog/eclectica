package strings

import "golang.org/x/crypto/ssh/terminal"

func Elipsis(str string, num int) string {
	var (
		result = str
		length = len(str)
	)

	if num < 1 {
		return result
	}

	if length > num {
		if num > 3 {
			num -= 3
		}
		result = str[0:num] + "..."
	}
	return result
}

func ElipsisForTerminal(str string) (result string) {
	var (
		width, _, _ = terminal.GetSize(0)

		// Needed to compensate for the reset of the data already
		// printed to stdout/stderr
		subjectiveValue = 50
	)

	return Elipsis(str, width-subjectiveValue)
}
