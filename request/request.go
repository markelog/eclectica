package request

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/markelog/eclectica/variables"
)

var (
	client = &http.Client{}
)

func Body(url string) (string, error) {
	response, err := client.Get(url)
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", errors.New(variables.ConnectionError)
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(contents), nil
}
