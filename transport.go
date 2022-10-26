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
	t     *testing.T
	index int
	calls []call
}

func (transport *transport) assertHeaders(r *http.Request, call call) {
	for name, values := range call.headers {
		requestValues, ok := r.Header[name]
		if !ok {
			transport.t.Errorf("header %q not set in request", name)
		}
		if !reflect.DeepEqual(values, requestValues) {
			transport.t.Errorf("header %q has bad values", name)
		}
	}
}

func (transport *transport) assertQueryParams(r *http.Request, call call) {
	requestQuery := r.URL.Query()
	for name, value := range call.queryParams {
		requestValue := requestQuery.Get(name)
		if value != requestValue {
			transport.t.Errorf("query parameter %q has bad value", name)
		}
	}
}

func (transport *transport) assertJSON(r *http.Request, call call) {
	if len(call.expectedJSON) > 0 {
		if r.Body == nil {
			transport.t.Errorf("expected body but received nothing")
		}
		actual, err := io.ReadAll(r.Body)
		if err != nil {
			transport.t.Errorf("io.ReadAll error: %s", err.Error())
		}

		var expectedJSONAsInterface, actualJSONAsInterface interface{}
		if err := json.Unmarshal(call.expectedJSON, &expectedJSONAsInterface); err != nil {
			transport.t.Errorf("Expected value ('%s') is not valid json.\nJSON parsing error: '%s'", call.expectedJSON, err.Error())
		}
		if err := json.Unmarshal(actual, &actualJSONAsInterface); err != nil {
			transport.t.Errorf("Input ('%s') needs to be valid json.\nJSON parsing error: '%s'", actual, err.Error())
		}

		if !reflect.DeepEqual(expectedJSONAsInterface, actualJSONAsInterface) {
			transport.t.Errorf("objects are not equal.")
		}
	}
}

func (transport *transport) assertBody(r *http.Request, call call) {
	if len(call.expectedBody) > 0 {
		if r.Body == nil {
			transport.t.Errorf("expected body but received nothing")
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			transport.t.Errorf("io.ReadAll error: %s", err.Error())
		}
		if call.expectedBody != string(data) {
			transport.t.Errorf("expected body does not match received body")
		}
	}
}

func (transport *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	transport.t.Helper()
	if transport.index >= len(transport.calls) {
		transport.t.Errorf("unexpected route call #%d", transport.index)
		return nil, fmt.Errorf("unexpected route call #%d", transport.index)
	}
	call := transport.calls[transport.index]
	if r.URL.Path != call.route || r.Method != call.method {
		transport.t.Errorf("unexpected route call #%d on route %s %s, expected %s %s",
			transport.index, r.Method, r.URL.Path, call.method, call.route)
		return nil, fmt.Errorf("unexpected route call #%d on route %s %s, expected %s %s",
			transport.index, r.Method, r.URL.Path, call.method, call.route)
	}

	transport.assertJSON(r, call)
	transport.assertBody(r, call)
	transport.assertHeaders(r, call)
	transport.assertQueryParams(r, call)

	transport.index++
	if call.err != nil {
		return nil, call.err
	}
	return &http.Response{
		Status:     http.StatusText(call.returnStatus),
		StatusCode: call.returnStatus,
		Body:       io.NopCloser(strings.NewReader(call.returnBody)),
	}, nil
}
