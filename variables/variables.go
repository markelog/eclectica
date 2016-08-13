package variables

import(
  "os"
  "fmt"
)

var (
  Home = fmt.Sprintf("%s/.eclectica/versions", os.Getenv("HOME"))
  Prefix = os.Getenv("HOME")
  Commands = []string{"ls", "rm"}
)


