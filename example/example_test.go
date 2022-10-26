package example

import (
	"net/http"
	"testing"

	"github.com/oupo1337/httpmock"
)

func Test_simpleRequest(t *testing.T) {
	mock := httpmock.New(t).
		WithCall(http.MethodGet, "/path")

	simpleGetRequest(mock)
	mock.AssertExpectations(t)
}

func Test_simpleRequest_missing_calls(t *testing.T) {
	mock := httpmock.New(t).
		WithCall(http.MethodGet, "/path")

	mock.AssertExpectations(t)
}

func Test_simplePostRequestWithBody(t *testing.T) {
	mock := httpmock.New(t).
		WithCall(http.MethodPost, "/path", httpmock.ExpectBody(`{"hello": "world"}`))

	simplePostRequestWithBody(mock)
	mock.AssertExpectations(t)
}
