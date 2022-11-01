package httpmock

import (
	"net/http"
	"net/url"
)

type request struct {
	method              string
	route               string
	returnStatus        int
	returnBody          string
	returnError         error
	expectedBody        string
	expectedJSON        []byte
	expectedHeaders     http.Header
	expectedQueryParams url.Values
}

func newMockRequest(method, route string) request {
	return request{
		method:       method,
		route:        route,
		returnStatus: http.StatusOK,
	}
}
