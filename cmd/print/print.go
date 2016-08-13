package print

import(
  "fmt"

  "github.com/fatih/color"
)

func InStyle(name, entity string) {
  color.Set(color.Bold)
  fmt.Print(name)

  color.Set(color.FgCyan)
  fmt.Print(" ")
  fmt.Print(entity + " ")
  color.Unset()
}

func LaguageOrVersion(language, version string) {
  if language != "" {
    InStyle("Language", language)
    fmt.Println()
  }

  if version != "" {
    InStyle("Version", version)
    fmt.Println()
  }
}
