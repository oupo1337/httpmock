package example

import (
	"net/http"
	"strings"
)

type doer interface {
	Do(req *http.Request) (*http.Response, error)
}

func simpleGetRequest(client doer) {
	req, err := http.NewRequest(http.MethodGet, "https://fake.url/path", nil)
	if err != nil {
		return
	}

	response, err := client.Do(req)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		return
	}
}

func simplePostRequestWithBody(client doer) {
	body := strings.NewReader(`{"a": "b", "c": "d"}`)
	req, err := http.NewRequest(http.MethodPost, "https://fake.url/path", body)
	if err != nil {
		return
	}

	response, err := client.Do(req)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		return
	}
}