package request

import (
  "io/ioutil"
  "net/http"
  "errors"
)

var (
  client = &http.Client{}
)

func Body(url string) (string, error){
  response, _ := client.Get(url)

  if response.StatusCode != 200 {
    return "", errors.New("Can't establish connection")
  }

  defer response.Body.Close()
  contents, err := ioutil.ReadAll(response.Body)

  if err != nil {
    return "", err
  }

  return string(contents), nil
}
