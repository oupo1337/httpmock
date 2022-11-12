package httpmock

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_httpMockBody(t *testing.T) {
	mockT := new(testing.T)
	mock := New(mockT).
		WithRequest(http.MethodPost, "/first",
			ExpectHeader("Authorization", []string{"Bearer TOKEN"}),
			ExpectBody("foobar"),
			ExpectQueryParam("param1", "value1"),
			ReturnStatus(http.StatusOK),
			ReturnBodyRaw("hello world"),
		).
		WithRequest(http.MethodPost, "/second",
			ReturnError(fmt.Errorf("oops")),
		)

	assert.Equal(t, []*request{
		{
			path:         "/first",
			method:       http.MethodPost,
			returnStatus: http.StatusOK,
			returnBody:   "hello world",
			expectedHeaders: map[string][]string{
				"Authorization": {"Bearer TOKEN"},
			},
			expectedQueryParams: url.Values{
				"param1": {"value1"},
			},
			expectedBody: "foobar",
		},
		{
			path:        "/second",
			method:      http.MethodPost,
			returnError: fmt.Errorf("oops"),
		},
	}, mock.transport.requests)

	req1, _ := http.NewRequest(http.MethodPost, "/first?param1=value1", strings.NewReader("foobar"))
	req1.Header.Add("Authorization", "Bearer TOKEN")
	response1, err := mock.Do(req1) //nolint:bodyclose
	assert.NoError(t, err)
	assert.True(t, mock.transport.requests[0].called)
	assert.Equal(t, response1.StatusCode, http.StatusOK)
	assert.False(t, mockT.Failed())

	req2, _ := http.NewRequest(http.MethodPost, "/second", nil)
	_, err = mock.Do(req2) //nolint:bodyclose
	assert.Error(t, err)
	assert.True(t, mock.transport.requests[1].called)
	assert.False(t, mockT.Failed())

	req3, _ := http.NewRequest(http.MethodOptions, "/toomuch", nil)
	_, err = mock.Do(req3) //nolint:bodyclose
	assert.Error(t, err)
	assert.True(t, mockT.Failed())
}

func Test_httpMockJSON(t *testing.T) {
	mockT := new(testing.T)
	mock := New(mockT).
		WithRequest(http.MethodPost, "/first",
			ExpectHeader("Authorization", []string{"Bearer TOKEN"}),
			ExpectJSON(`{"foo": "bar"}`),
			ExpectQueryParam("param1", "value1"),
			ReturnStatus(http.StatusNoContent),
			ReturnBodyRaw("hello world"),
		).
		WithRequest(http.MethodPut, "/second",
			ExpectQueryParamValues("param", []string{"value1", "value2"}),
			ReturnStatus(http.StatusOK),
		).
		WithRequest(http.MethodPost, "/third",
			ReturnError(fmt.Errorf("oops")),
		)

	assert.Equal(t, []*request{
		{
			path:         "/first",
			method:       http.MethodPost,
			returnStatus: http.StatusNoContent,
			returnBody:   "hello world",
			expectedHeaders: map[string][]string{
				"Authorization": {"Bearer TOKEN"},
			},
			expectedQueryParams: url.Values{
				"param1": {"value1"},
			},
			expectedJSON: []byte(`{"foo": "bar"}`),
		},
		{
			path:         "/second",
			method:       http.MethodPut,
			returnStatus: http.StatusOK,
			expectedQueryParams: url.Values{
				"param": {"value1", "value2"},
			},
		},
		{
			path:        "/third",
			method:      http.MethodPost,
			returnError: fmt.Errorf("oops"),
		},
	}, mock.transport.requests)

	req1, _ := http.NewRequest(http.MethodPost, "/first?param1=value1", strings.NewReader(`{"foo": "bar"}`))
	req1.Header.Add("Authorization", "Bearer TOKEN")
	response1, err := mock.Do(req1) //nolint:bodyclose
	assert.NoError(t, err)
	assert.Equal(t, response1.StatusCode, http.StatusNoContent)
	assert.True(t, mock.transport.requests[0].called)
	assert.False(t, mockT.Failed())

	req2, _ := http.NewRequest(http.MethodPut, "/second?param=value1&param=value2", nil)
	response2, err := mock.Do(req2) //nolint:bodyclose
	assert.NoError(t, err)
	assert.Equal(t, response2.StatusCode, http.StatusOK)
	assert.True(t, mock.transport.requests[1].called)
	assert.False(t, mockT.Failed())

	req3, _ := http.NewRequest(http.MethodPost, "/third", nil)
	_, err = mock.Do(req3) //nolint:bodyclose
	assert.Error(t, err)
	assert.True(t, mock.transport.requests[2].called)
	assert.False(t, mockT.Failed())

	req4, _ := http.NewRequest(http.MethodOptions, "/toomuch", nil)
	_, err = mock.Do(req4) //nolint:bodyclose
	assert.Error(t, err)
	assert.True(t, mockT.Failed())
}

func Test_httpMock_wrong_call(t *testing.T) {
	mockT := new(testing.T)
	mock := New(mockT)

	assert.Equal(t, []*request{}, mock.transport.requests)

	req0, _ := http.NewRequest(http.MethodGet, "/bad", nil)
	_, err := mock.Do(req0) //nolint:bodyclose
	assert.Error(t, err)
	assert.True(t, mockT.Failed())
}
