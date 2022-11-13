package httpmock

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

func (r *request) String() string {
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("Request: [%s] %q\n", r.method, r.path))

	if len(r.expectedHeaders) > 0 {
		builder.WriteString("Expected headers:\n")
		for name, values := range r.expectedHeaders {
			builder.WriteString(fmt.Sprintf("\t- %s: %s\n", name, values))
		}
	}

	if len(r.expectedQueryParams) > 0 {
		builder.WriteString("Expected query params:\n")
		for name, values := range r.expectedQueryParams {
			builder.WriteString(fmt.Sprintf("\t- %s: %s\n", name, values))
		}
	}

	if len(r.expectedBody) > 0 {
		builder.WriteString(fmt.Sprintf("Expected body:\n\t%q\n", r.expectedBody))
	}

	if len(r.expectedJSON) > 0 {
		builder.WriteString(fmt.Sprintf("Expected JSON:\n\t%q\n", string(r.expectedJSON)))
	}
	return builder.String()
}
