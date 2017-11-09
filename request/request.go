package request

import (
	"io/ioutil"
	"net/http"

	"github.com/go-errors/errors"

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
		err = errors.New(err)
		return "", err
	}

	return string(contents), nil
}
