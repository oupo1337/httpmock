package httpmock

import (
	"net/http"
	"testing"
)

type call struct {
	method       string
	route        string
	returnStatus int
	returnBody   string
	expectedBody string
	expectedJSON []byte
	err          error
	headers      map[string][]string
	queryParams  map[string]string
}

type CallOption func(*call)

func ReturnError(err error) CallOption {
	return func(c *call) {
		c.err = err
	}
}

func ReturnStatus(status int) CallOption {
	return func(c *call) {
		c.returnStatus = status
	}
}

func ReturnBody(body string) CallOption {
	return func(c *call) {
		c.returnBody = body
	}
}

func ExpectBody(expectedBody string) CallOption {
	return func(c *call) {
		c.expectedBody = expectedBody
	}
}

func ExpectJSON(expectedJSON string) CallOption {
	return func(c *call) {
		c.expectedJSON = []byte(expectedJSON)
	}
}

func ExpectHeader(name string, values []string) CallOption {
	return func(c *call) {
		if c.headers == nil {
			c.headers = make(map[string][]string)
		}
		c.headers[name] = values
	}
}

func ExpectQueryParam(name, value string) CallOption {
	return func(c *call) {
		if c.queryParams == nil {
			c.queryParams = make(map[string]string)
		}
		c.queryParams[name] = value
	}
}

type Client struct {
	http.Client
	transport *transport
}

func (c *Client) WithCall(method, route string, options ...CallOption) *Client {
	call := call{
		method: method,
		route:  route,
	}
	for _, option := range options {
		option(&call)
	}
	c.transport.calls = append(c.transport.calls, call)
	return c
}

func (c *Client) AssertExpectations(t *testing.T) {
	if c.transport.index != len(c.transport.calls) {
		t.Errorf("missing calls")
	}
}

func New(t *testing.T) *Client {
	transport := &transport{
		t:     t,
		calls: make([]call, 0),
	}
	return &Client{
		transport: transport,
		Client: http.Client{
			Transport: transport,
		},
	}
}
