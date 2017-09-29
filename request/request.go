package request

import (
	"errors"
	"io/ioutil"
	"net/http"
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
		return "", errors.New("Connection cannot be established")
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(contents), nil
}
