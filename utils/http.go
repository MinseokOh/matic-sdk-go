package utils

import (
	"io/ioutil"
	"net/http"
)

func Get(url string, params map[string]string, headers map[string]string) (string, error) {
	if request, err := http.NewRequest("GET", url, nil); err != nil {
		return "", err
	} else {
		query := request.URL.Query()
		for key, value := range params {
			query.Add(key, value)
		}
		request.URL.RawQuery = query.Encode()
		for key, value := range headers {
			request.Header.Add(key, value)
		}

		client := http.DefaultClient
		if response, err := client.Do(request); err != nil {
			return "", err
		} else {
			defer response.Body.Close()

			if body, err := ioutil.ReadAll(response.Body); err != nil {
				return "", err
			} else {
				return string(body), nil
			}
		}
	}
}
