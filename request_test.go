package httpmock

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest_Times(t *testing.T) {
	r := Request{}
	r.Times(1000)

	assert.Equal(t, 1000, r.expectedTimesCalled)
}

func TestTimes(t *testing.T) {
	mock := New(t).WithRequest(http.MethodGet, "/", Times(1000))
	r := mock.transport.requests[0]

	assert.Equal(t, 1000, r.expectedTimesCalled)
}

func TestRequest_ReturnStatus(t *testing.T) {
	r := Request{}
	r.ReturnStatus(http.StatusTeapot)

	assert.Equal(t, http.StatusTeapot, r.returnStatus)
}

func TestReturnStatus(t *testing.T) {
	mock := New(t).WithRequest(http.MethodGet, "/", ReturnStatus(http.StatusTeapot))
	r := mock.transport.requests[0]

	assert.Equal(t, http.StatusTeapot, r.returnStatus)
}

func TestRequest_ReturnBody(t *testing.T) {
	r := Request{}
	r.ReturnBody("this is a test body")

	assert.Equal(t, "this is a test body", r.returnBody)
}

func TestReturnBody(t *testing.T) {
	mock := New(t).WithRequest(http.MethodGet, "/", ReturnBody("this is a test body"))
	r := mock.transport.requests[0]

	assert.Equal(t, "this is a test body", r.returnBody)
}

func TestRequest_ReturnError(t *testing.T) {
	r := Request{}
	r.ReturnError(assert.AnError)

	assert.Equal(t, assert.AnError, r.returnError)
}

func TestReturnError(t *testing.T) {
	mock := New(t).WithRequest(http.MethodGet, "/", ReturnError(assert.AnError))
	r := mock.transport.requests[0]

	assert.Equal(t, assert.AnError, r.returnError)
}

func TestRequest_ReturnHeader(t *testing.T) {
	r := Request{}
	r.ReturnHeader("name", []string{"value"})

	assert.Equal(t, http.Header{"name": {"value"}}, r.returnHeaders)
}

func TestReturnHeader(t *testing.T) {
	mock := New(t).WithRequest(http.MethodGet, "/", ReturnHeader("name", []string{"value"}))
	r := mock.transport.requests[0]

	assert.Equal(t, http.Header{"name": {"value"}}, r.returnHeaders)
}

func TestReturnBodyFromObject(t *testing.T) {
	test := struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}{
		Name:  "name",
		Value: 1000,
	}

	r := Request{}
	r.ReturnBodyFromObject(test)

	assert.Equal(t, `{"name":"name","value":1000}`, r.returnBody)
}

func TestRequest_ReturnBodyFromObject(t *testing.T) {
	test := struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}{
		Name:  "name",
		Value: 1000,
	}

	mock := New(t).WithRequest(http.MethodGet, "/", ReturnBodyFromObject(test))
	r := mock.transport.requests[0]

	assert.Equal(t, `{"name":"name","value":1000}`, r.returnBody)
}

func TestRequest_ExpectBody(t *testing.T) {
	r := Request{}
	r.ExpectBody("this is a test body")

	assert.Equal(t, "this is a test body", r.expectedBody)
}

func TestExpectBody(t *testing.T) {
	mock := New(t).WithRequest(http.MethodGet, "/", ExpectBody("this is a test body"))
	r := mock.transport.requests[0]

	assert.Equal(t, "this is a test body", r.expectedBody)
}

func TestRequest_ExpectJSON(t *testing.T) {
	r := Request{}
	r.ExpectJSON(`{"hello":"world"}`)

	assert.Equal(t, []byte(`{"hello":"world"}`), r.expectedJSON)
}

func TestExpectJSON(t *testing.T) {
	mock := New(t).WithRequest(http.MethodGet, "/", ExpectJSON(`{"hello":"world"}`))
	r := mock.transport.requests[0]

	assert.Equal(t, []byte(`{"hello":"world"}`), r.expectedJSON)
}

func TestRequest_ExpectHeader(t *testing.T) {
	r := Request{}
	r.ExpectHeader("name", []string{"value"})

	assert.Equal(t, http.Header{"name": {"value"}}, r.expectedHeaders)
}

func TestExpectHeader(t *testing.T) {
	mock := New(t).WithRequest(http.MethodGet, "/", ExpectHeader("name", []string{"value"}))
	r := mock.transport.requests[0]

	assert.Equal(t, http.Header{"name": {"value"}}, r.expectedHeaders)
}

func TestRequest_ExpectQueryParam(t *testing.T) {
	r := Request{}
	r.ExpectQueryParam("name", "value")

	assert.Equal(t, url.Values{"name": {"value"}}, r.expectedQueryParams)
}

func TestExpectQueryParam(t *testing.T) {
	mock := New(t).WithRequest(http.MethodGet, "/", ExpectQueryParam("name", "value"))
	r := mock.transport.requests[0]

	assert.Equal(t, url.Values{"name": {"value"}}, r.expectedQueryParams)
}

func TestRequest_ExpectQueryParamValues(t *testing.T) {
	r := Request{}
	r.ExpectQueryParamValues("name", []string{"value1", "value2"})

	assert.Equal(t, url.Values{"name": {"value1", "value2"}}, r.expectedQueryParams)
}

func TestExpectQueryParamValues(t *testing.T) {
	mock := New(t).WithRequest(http.MethodGet, "/", ExpectQueryParamValues("name", []string{"value1", "value2"}))
	r := mock.transport.requests[0]

	assert.Equal(t, url.Values{"name": {"value1", "value2"}}, r.expectedQueryParams)
}

func TestRequest_ContentLength(t *testing.T) {
	r := Request{}
	assert.Equal(t, int64(0), r.ContentLength())

	r.returnBody = "this test is amazing"
	assert.Equal(t, int64(len(r.returnBody)), r.ContentLength())

	r.returnHeaders = make(map[string][]string)
	r.returnHeaders["Content-Length"] = []string{"1000"}
	assert.Equal(t, int64(1000), r.ContentLength())
}
