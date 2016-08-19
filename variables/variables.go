package variables

import(
  "os"
  "fmt"
)

var (
  Commands = []string{"ls", "rm", "--help", "-h"}
  Files = [4]string{"bin", "lib", "include", "share"}
)

func Prefix() string {
  return os.Getenv("HOME")
}

func Home() string {
  return fmt.Sprintf("%s/.eclectica/versions", Prefix())
}


