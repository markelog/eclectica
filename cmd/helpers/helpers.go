package helpers

import(
  "fmt"

  "github.com/fatih/color"
)


func PrintInStyle(name, entity string) {
  color.Set(color.FgBlack)
  fmt.Print(name)

  color.Set(color.FgCyan)
  fmt.Print(" ")
  fmt.Print(entity + " ")
  color.Unset()
}
