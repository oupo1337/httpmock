package httpmock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testing"
)

var UnexpectedRequestErr = fmt.Errorf("unexpected request")

type transport struct {
	m        sync.Mutex
	t        *testing.T
	requests []*Request
}

func assertHeaders(r *http.Request, req *Request) bool {
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

func assertQueryParams(r *http.Request, req *Request) bool {
	requestQuery := r.URL.Query()
	for name, values := range req.expectedQueryParams {
		requestValues := requestQuery[name]
		if !reflect.DeepEqual(values, requestValues) {
			return false
		}
	}
	return true
}

func assertJSON(r *http.Request, req *Request) bool {
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

func assertBody(r *http.Request, req *Request) bool {
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

func (t *transport) matchRequest(r *http.Request) (*Request, *Request) {
	var closestReq *Request
	for _, req := range t.requests {
		if req.timesCalled < req.expectedTimesCalled && req.path == r.URL.Path {
			if req.method == r.Method && assertJSON(r, req) && assertBody(r, req) && assertHeaders(r, req) && assertQueryParams(r, req) {
				return req, nil
			}
			closestReq = req
		}
	}
	return nil, closestReq
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	t.t.Helper()
	t.m.Lock()
	defer t.m.Unlock()

	req, closestReq := t.matchRequest(r)
	if closestReq != nil {
		t.t.Errorf("Unexpected request on route [%s] %q the closest request I have is:\n%s", r.Method, r.URL.Path, closestReq.String())
		return nil, UnexpectedRequestErr
	}
	if req == nil {
		t.t.Errorf("Unexpected request on route [%s] %q", r.Method, r.URL.Path)
		return nil, UnexpectedRequestErr
	}
	req.timesCalled += 1

	if req.returnError != nil {
		return nil, req.returnError
	}

	return &http.Response{
		Status:        http.StatusText(req.returnStatus),
		StatusCode:    req.returnStatus,
		Header:        req.returnHeaders,
		ContentLength: req.ContentLength(),
		Body:          io.NopCloser(strings.NewReader(req.returnBody)),
	}, nil
}
