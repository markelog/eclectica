package variables

import(
  "os"
  "fmt"
)

var (
  Home = fmt.Sprintf("%s/.eclectica/versions", os.Getenv("HOME"))
  Directories = [1]string{"bin", "include", "share", "lib"}
  Prefix = "/usr/local"
)


