package prompt

import (
  "github.com/kildevaeld/prompt"
  "github.com/kildevaeld/prompt/terminal"
  "github.com/kildevaeld/prompt/form"
)

type Result struct {
  Language string `form:"list"`
  Version string `form:"list"`
}

func List(name string, choices []string) Result {
  var result Result

  ui := prompt.NewUI()
  ui.Theme = terminal.DefaultTheme

  ui.FormWithFields([]form.Field{
    &form.List{
      Name:    name,
      Choices: choices,
    },
  }, &result)

  return result
}
