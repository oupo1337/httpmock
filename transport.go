package httpmock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type transport struct {
	t        *testing.T
	index    int
	requests []request
}

func (transport *transport) assertHeaders(r *http.Request, req request) {
	for name, values := range req.expectedHeaders {
		requestValues, ok := r.Header[name]
		if !ok {
			transport.t.Errorf("header %q not set in request", name)
		}
		if !reflect.DeepEqual(values, requestValues) {
			transport.t.Errorf("header %q has bad values", name)
		}
	}
}

func (transport *transport) assertQueryParams(r *http.Request, req request) {
	requestQuery := r.URL.Query()
	for name, values := range req.expectedQueryParams {
		requestValues := requestQuery[name]
		if !reflect.DeepEqual(values, requestValues) {
			transport.t.Errorf("query parameter %q has bad value", name)
		}
	}
}

func (transport *transport) assertJSON(r *http.Request, req request) {
	if len(req.expectedJSON) > 0 {
		if r.Body == nil {
			transport.t.Errorf("expected body but received nothing")
		}
		actual, err := io.ReadAll(r.Body)
		if err != nil {
			transport.t.Errorf("io.ReadAll error: %s", err.Error())
		}

		var expectedJSONAsInterface, actualJSONAsInterface interface{}
		if err := json.Unmarshal(req.expectedJSON, &expectedJSONAsInterface); err != nil {
			transport.t.Errorf("Expected value ('%s') is not valid json.\nJSON parsing error: '%s'", req.expectedJSON, err.Error())
		}
		if err := json.Unmarshal(actual, &actualJSONAsInterface); err != nil {
			transport.t.Errorf("Input ('%s') needs to be valid json.\nJSON parsing error: '%s'", actual, err.Error())
		}

		if !reflect.DeepEqual(expectedJSONAsInterface, actualJSONAsInterface) {
			transport.t.Errorf("objects are not equal.")
		}
	}
}

func (transport *transport) assertBody(r *http.Request, req request) {
	if len(req.expectedBody) > 0 {
		if r.Body == nil {
			transport.t.Errorf("expected body but received nothing")
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			transport.t.Errorf("io.ReadAll error: %s", err.Error())
		}
		if req.expectedBody != string(data) {
			transport.t.Errorf("expected body does not match received body")
		}
	}
}

func (transport *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	transport.t.Helper()
	if transport.index >= len(transport.requests) {
		transport.t.Errorf("unexpected route request #%d", transport.index)
		return nil, fmt.Errorf("unexpected route request #%d", transport.index)
	}

	req := transport.requests[transport.index]
	if r.URL.Path != req.route || r.Method != req.method {
		transport.t.Errorf("unexpected route request #%d on route %s %s, expected %s %s",
			transport.index, r.Method, r.URL.Path, req.method, req.route)
		return nil, fmt.Errorf("unexpected route request #%d on route %s %s, expected %s %s",
			transport.index, r.Method, r.URL.Path, req.method, req.route)
	}

	transport.assertJSON(r, req)
	transport.assertBody(r, req)
	transport.assertHeaders(r, req)
	transport.assertQueryParams(r, req)

	transport.index++
	if req.returnError != nil {
		return nil, req.returnError
	}
	return &http.Response{
		Status:     http.StatusText(req.returnStatus),
		StatusCode: req.returnStatus,
		Body:       io.NopCloser(strings.NewReader(req.returnBody)),
	}, nil
}
