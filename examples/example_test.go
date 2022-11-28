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
	mock.AssertExpectations()
}

func Test_closeRequest(t *testing.T) {
	mock := httpmock.New(t).
		WithRequest(http.MethodPost, "/path",
			httpmock.ExpectQueryParam("query", "param"),
			httpmock.ExpectHeader("hello", []string{"world"}),
			httpmock.ExpectHeader("bonjour", []string{"monde"}),
			httpmock.ExpectHeader("hola", []string{"mundo"}),
			httpmock.ExpectJSON(`{"hello":"world"}`),
			httpmock.ReturnStatus(http.StatusOK),
		)

	simpleGetRequest(mock)
	mock.AssertExpectations()
}

// This test will fail
func Test_simpleRequest_missing_calls(t *testing.T) {
	mock := httpmock.New(t).
		WithRequest(http.MethodGet, "/path").
		WithRequest(http.MethodPost, "/test")

	mock.AssertExpectations()
}

func Test_simplePostRequestWithBody(t *testing.T) {
	mock := httpmock.New(t).
		WithRequest(http.MethodPost, "/path",
			httpmock.ExpectBody(`{"a": "b", "c": "d"}`),
			httpmock.ReturnStatus(http.StatusOK),
		)

	simplePostRequestWithBody(mock)
	mock.AssertExpectations()
}

func Test_simplePostRequestJSON(t *testing.T) {
	mock := httpmock.New(t).
		WithRequest(http.MethodPost, "/path",
			httpmock.ExpectJSON(`{"c": "d", "a": "b"}`),
			httpmock.ReturnStatus(http.StatusOK),
		)

	simplePostRequestWithBody(mock)
	mock.AssertExpectations()
}

func Test_simplePostRequestJSON_bis(t *testing.T) {
	mock := httpmock.New(t)
	mock.On(http.MethodPost, "/path").
		ExpectJSON(`{"c": "d", "a": "b"}`).
		ReturnStatus(http.StatusOK)

	simplePostRequestWithBody(mock)
	mock.AssertExpectations()
}

func Test_simplePostReturnsAnError(t *testing.T) {
	mock := httpmock.New(t).
		WithRequest(http.MethodPost, "/path",
			httpmock.ReturnError(fmt.Errorf("oops")),
		)

	simplePostRequestWithBody(mock)
	mock.AssertExpectations()
}

func Test_simplePostReturnsBodyFromObject(t *testing.T) {
	object := struct {
		Hello string `json:"hello"`
	}{
		Hello: "world",
	}

	mock := httpmock.New(t).
		WithRequest(http.MethodPost, "/path",
			httpmock.ReturnStatus(http.StatusOK),
			httpmock.ReturnBodyFromObject(object),
		)

	simplePostRequestWithBody(mock)
	mock.AssertExpectations()
}

func Test_simpleRequestGoRoutines(t *testing.T) {
	times := 1000

	mock := httpmock.New(t)
	mock.On(http.MethodGet, "/path").ReturnStatus(http.StatusOK).Times(times)

	simpleRequestGoRoutines(mock, times)
	mock.AssertExpectations()
}
