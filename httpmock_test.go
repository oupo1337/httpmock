package httpmock

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_httpMockBody(t *testing.T) {
	mockT := new(testing.T)
	mock := New(mockT).
		WithCall(http.MethodPost, "/first",
			ExpectHeader("Authorization", []string{"Bearer TOKEN"}),
			ExpectBody("foobar"),
			ExpectQueryParam("param1", "value1"),
			ReturnStatus(http.StatusOK),
			ReturnBody("hello world"),
		).
		WithCall(http.MethodPost, "/second",
			ReturnError(fmt.Errorf("oops")),
		)

	assert.Equal(t, []call{
		{
			route:        "/first",
			method:       http.MethodPost,
			returnStatus: http.StatusOK,
			returnBody:   "hello world",
			headers: map[string][]string{
				"Authorization": {"Bearer TOKEN"},
			},
			queryParams: map[string]string{
				"param1": "value1",
			},
			expectedBody: "foobar",
		},
		{
			route:  "/second",
			method: http.MethodPost,
			err:    fmt.Errorf("oops"),
		},
	}, mock.transport.calls)
	assert.Equal(t, 0, mock.transport.index)

	req1, _ := http.NewRequest(http.MethodPost, "/first?param1=value1", strings.NewReader("foobar"))
	req1.Header.Add("Authorization", "Bearer TOKEN")
	response1, err := mock.Do(req1) //nolint:bodyclose
	assert.NoError(t, err)
	assert.Equal(t, response1.StatusCode, http.StatusOK)
	assert.Equal(t, 1, mock.transport.index)
	assert.False(t, mockT.Failed())

	req2, _ := http.NewRequest(http.MethodPost, "/second", nil)
	_, err = mock.Do(req2) //nolint:bodyclose
	assert.Error(t, err)
	assert.Equal(t, 2, mock.transport.index)
	assert.False(t, mockT.Failed())

	req3, _ := http.NewRequest(http.MethodOptions, "/toomuch", nil)
	_, err = mock.Do(req3) //nolint:bodyclose
}

func Test_httpMockJSON(t *testing.T) {
	mockT := new(testing.T)
	mock := New(mockT).
		WithCall(http.MethodPost, "/first",
			ExpectHeader("Authorization", []string{"Bearer TOKEN"}),
			ExpectJSON(`{"foo": "bar"}`),
			ExpectQueryParam("param1", "value1"),
			ReturnStatus(http.StatusOK),
			ReturnBody("hello world"),
		).
		WithCall(http.MethodPost, "/second",
			ReturnError(fmt.Errorf("oops")),
		)

	assert.Equal(t, []call{
		{
			route:        "/first",
			method:       http.MethodPost,
			returnStatus: http.StatusOK,
			returnBody:   "hello world",
			headers: map[string][]string{
				"Authorization": {"Bearer TOKEN"},
			},
			queryParams: map[string]string{
				"param1": "value1",
			},
			expectedJSON: []byte(`{"foo": "bar"}`),
		},
		{
			route:  "/second",
			method: http.MethodPost,
			err:    fmt.Errorf("oops"),
		},
	}, mock.transport.calls)
	assert.Equal(t, 0, mock.transport.index)

	req1, _ := http.NewRequest(http.MethodPost, "/first?param1=value1", strings.NewReader(`{"foo": "bar"}`))
	req1.Header.Add("Authorization", "Bearer TOKEN")
	response1, err := mock.Do(req1) //nolint:bodyclose
	assert.NoError(t, err)
	assert.Equal(t, response1.StatusCode, http.StatusOK)
	assert.Equal(t, 1, mock.transport.index)
	assert.False(t, mockT.Failed())

	req2, _ := http.NewRequest(http.MethodPost, "/second", nil)
	_, err = mock.Do(req2) //nolint:bodyclose
	assert.Error(t, err)
	assert.Equal(t, 2, mock.transport.index)
	assert.False(t, mockT.Failed())

	req3, _ := http.NewRequest(http.MethodOptions, "/toomuch", nil)
	_, err = mock.Do(req3) //nolint:bodyclose
	assert.Error(t, err)
	assert.True(t, mockT.Failed())
}

func Test_httpMock_wrong_call(t *testing.T) {
	mockT := new(testing.T)
	mock := New(mockT)

	assert.Equal(t, []call{}, mock.transport.calls)
	assert.Equal(t, 0, mock.transport.index)

	req0, _ := http.NewRequest(http.MethodGet, "/bad", nil)
	_, err := mock.Do(req0) //nolint:bodyclose
	assert.Error(t, err)
	assert.True(t, mockT.Failed())
}
