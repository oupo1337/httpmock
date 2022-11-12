package httpmock

import (
	"net/http"
	"net/url"
)

type request struct {
	method              string
	path                string
	returnStatus        int
	returnBody          string
	returnError         error
	expectedBody        string
	expectedJSON        []byte
	expectedHeaders     http.Header
	expectedQueryParams url.Values
	called              bool
}

func newMockRequest(method, path string) request {
	return request{
		method: method,
		path:   path,
	}
}
