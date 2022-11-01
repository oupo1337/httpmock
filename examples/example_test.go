package example

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/oupo1337/httpmock"
)

func Test_simpleRequest(t *testing.T) {
	mock := httpmock.New(t).
		WithRequest(http.MethodGet, "/path",
			httpmock.ReturnStatus(http.StatusOK),
		)

	simpleGetRequest(mock)
	mock.AssertExpectations(t)
}

// This test will fail
func Test_simpleRequest_missing_calls(t *testing.T) {
	mock := httpmock.New(t).
		WithRequest(http.MethodGet, "/path")

	mock.AssertExpectations(t)
}

func Test_simplePostRequestWithBody(t *testing.T) {
	mock := httpmock.New(t).
		WithRequest(http.MethodPost, "/path",
			httpmock.ExpectBody(`{"a": "b", "c": "d"}`),
			httpmock.ReturnStatus(http.StatusOK),
		)

	simplePostRequestWithBody(mock)
	mock.AssertExpectations(t)
}

func Test_simplePostRequestJSON(t *testing.T) {
	mock := httpmock.New(t).
		WithRequest(http.MethodPost, "/path",
			httpmock.ExpectJSON(`{"c": "d", "a": "b"}`),
			httpmock.ReturnStatus(http.StatusOK),
		)

	simplePostRequestWithBody(mock)
	mock.AssertExpectations(t)
}

func Test_simplePostReturnsAnError(t *testing.T) {
	mock := httpmock.New(t).
		WithRequest(http.MethodPost, "/path",
			httpmock.ReturnError(fmt.Errorf("oops")),
		)

	simplePostRequestWithBody(mock)
	mock.AssertExpectations(t)
}
