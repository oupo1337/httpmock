package httpmock

import (
	"net/http"
	"testing"
)

// random comment again and again

type Client struct {
	http.Client
	transport *transport
}

type RequestOption func(*Request)

func (c *Client) WithRequest(method, path string, options ...RequestOption) *Client {
	req := &Request{
		method:              method,
		path:                path,
		expectedTimesCalled: 1,
	}
	for _, option := range options {
		option(req)
	}
	c.transport.requests = append(c.transport.requests, req)
	return c
}

func (c *Client) AssertExpectations() {
	for _, req := range c.transport.requests {
		if req.timesCalled < req.expectedTimesCalled {
			c.transport.t.Errorf("httpmock should have more requests: expected [%s] %q x%d", req.method, req.path, req.expectedTimesCalled-req.timesCalled)
		}
	}
}

func (c *Client) On(method, path string) *Request {
	req := &Request{
		method:              method,
		path:                path,
		expectedTimesCalled: 1,
	}
	c.transport.requests = append(c.transport.requests, req)
	return req
}

func New(t *testing.T) *Client {
	mockTransport := &transport{
		t:        t,
		requests: make([]*Request, 0),
	}

	return &Client{
		transport: mockTransport,
		Client: http.Client{
			Transport: mockTransport,
		},
	}
}
