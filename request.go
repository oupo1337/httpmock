package httpmock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	method              string
	path                string
	returnStatus        int
	returnBody          string
	returnError         error
	expectedBody        string
	expectedJSON        []byte
	expectedHeaders     http.Header
	expectedQueryParams url.Values
	expectedTimesCalled int
	timesCalled         int
}

func Times(times int) RequestOption {
	return func(r *Request) {
		r.Times(times)
	}
}

func (r *Request) Times(times int) *Request {
	r.expectedTimesCalled = times
	return r
}

func ReturnStatus(status int) RequestOption {
	return func(r *Request) {
		r.ReturnStatus(status)
	}
}

func (r *Request) ReturnStatus(status int) *Request {
	r.returnStatus = status
	return r
}

func ReturnBody(body string) RequestOption {
	return func(r *Request) {
		r.ReturnBody(body)
	}
}

func (r *Request) ReturnBody(body string) *Request {
	r.returnBody = body
	return r
}

func ReturnBodyFromObject(object interface{}) RequestOption {
	return func(r *Request) {
		r.ReturnBodyFromObject(object)
	}
}

func (r *Request) ReturnBodyFromObject(object interface{}) *Request {
	body, _ := json.Marshal(&object)
	r.returnBody = string(body)
	return r
}

func ReturnError(err error) RequestOption {
	return func(r *Request) {
		r.ReturnError(err)
	}
}

func (r *Request) ReturnError(err error) *Request {
	r.returnError = err
	return r
}

func ExpectBody(expectedBody string) RequestOption {
	return func(r *Request) {
		r.ExpectBody(expectedBody)
	}
}

func (r *Request) ExpectBody(body string) *Request {
	r.expectedBody = body
	return r
}

func ExpectJSON(expectedJSON string) RequestOption {
	return func(r *Request) {
		r.ExpectJSON(expectedJSON)
	}
}

func (r *Request) ExpectJSON(data string) *Request {
	r.expectedJSON = []byte(data)
	return r
}

func ExpectHeader(name string, values []string) RequestOption {
	return func(r *Request) {
		r.ExpectHeader(name, values)
	}
}

func (r *Request) ExpectHeader(name string, values []string) *Request {
	if r.expectedHeaders == nil {
		r.expectedHeaders = make(map[string][]string)
	}
	r.expectedHeaders[name] = values
	return r
}

func ExpectQueryParamValues(name string, values []string) RequestOption {
	return func(r *Request) {
		r.ExpectQueryParamValues(name, values)
	}
}

func (r *Request) ExpectQueryParamValues(name string, values []string) *Request {
	if r.expectedQueryParams == nil {
		r.expectedQueryParams = make(url.Values)
	}
	r.expectedQueryParams[name] = values
	return r
}

func ExpectQueryParam(name, value string) RequestOption {
	return func(r *Request) {
		r.ExpectQueryParam(name, value)
	}
}

func (r *Request) ExpectQueryParam(name, value string) *Request {
	if r.expectedQueryParams == nil {
		r.expectedQueryParams = make(url.Values)
	}
	r.expectedQueryParams[name] = []string{value}
	return r
}

func (r *Request) String() string {
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
