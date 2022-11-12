package example

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oupo1337/httpmock"
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
	if errors.Is(err, httpmock.UnexpectedRequestErr) {
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", response)
}

func simplePostRequestWithBody(client doer) {
	body := strings.NewReader(`{"a": "b", "c": "d"}`)
	req, err := http.NewRequest(http.MethodPost, "https://fake.url/path", body)
	if err != nil {
		return
	}

	response, err := client.Do(req)
	if errors.Is(err, httpmock.UnexpectedRequestErr) {
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", response)
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	fmt.Printf("responseData: %s\n", string(responseData))
}
