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
	requests []*request
}

func assertHeaders(r *http.Request, req *request) bool {
	for name, values := range req.expectedHeaders {
		requestValues, ok := r.Header[name]
		if !ok {
			return false
		}
		if !reflect.DeepEqual(values, requestValues) {
			return false
		}
	}
	return true
}

func assertQueryParams(r *http.Request, req *request) bool {
	requestQuery := r.URL.Query()
	for name, values := range req.expectedQueryParams {
		requestValues := requestQuery[name]
		if !reflect.DeepEqual(values, requestValues) {
			return false
		}
	}
	return true
}

func assertJSON(r *http.Request, req *request) bool {
	if len(req.expectedJSON) > 0 {
		if r.Body == nil {
			return false
		}
		actual, err := io.ReadAll(r.Body)
		if err != nil {
			return false
		}

		var expectedJSONAsInterface, actualJSONAsInterface interface{}
		if err := json.Unmarshal(req.expectedJSON, &expectedJSONAsInterface); err != nil {
			return false
		}
		if err := json.Unmarshal(actual, &actualJSONAsInterface); err != nil {
			return false
		}

		if !reflect.DeepEqual(expectedJSONAsInterface, actualJSONAsInterface) {
			return false
		}
	}
	return true
}

func assertBody(r *http.Request, req *request) bool {
	if len(req.expectedBody) > 0 {
		if r.Body == nil {
			return false
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			return false
		}
		if req.expectedBody != string(data) {
			return false
		}
	}
	return true
}

func (t *transport) matchRequest(r *http.Request) *request {
	for _, req := range t.requests {
		if req.method == r.Method && req.path == r.URL.Path &&
			assertJSON(r, req) &&
			assertBody(r, req) &&
			assertHeaders(r, req) &&
			assertQueryParams(r, req) &&
			req.called == false {
			return req
		}
	}
	return nil
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	t.t.Helper()

	req := t.matchRequest(r)
	if req == nil {
		t.t.Errorf("unexpected request on route %s %s", r.Method, r.URL.Path)
		return nil, fmt.Errorf("unexpected request on route %s %s", r.Method, r.URL.Path)
	}

	req.called = true
	if req.returnError != nil {
		return nil, req.returnError
	}
	return &http.Response{
		Status:     http.StatusText(req.returnStatus),
		StatusCode: req.returnStatus,
		Body:       io.NopCloser(strings.NewReader(req.returnBody)),
	}, nil
}
