package httpmock

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
		data, err := io.ReadAll(r.Body)
		if err != nil {
			transport.t.Errorf("io.ReadAll error: %s", err.Error())
		}
		assert.JSONEq(transport.t, call.expectedJSON, string(data))
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
	if transport.index > len(transport.calls) {
		transport.t.Errorf("unexpected route call #%d", transport.index)
		transport.t.FailNow()
	}
	call := transport.calls[transport.index]
	if r.URL.Path != call.route || r.Method != call.method {
		transport.t.Errorf("unexpected route call #%d on route %s %s, expected %s %s",
			transport.index, r.Method, r.URL.Path, call.method, call.route)
		transport.t.FailNow()
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
