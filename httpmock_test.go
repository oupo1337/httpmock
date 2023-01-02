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
			ReturnBody("hello world"),
		).
		WithRequest(http.MethodPost, "/second",
			ReturnError(fmt.Errorf("oops")),
		)

	assert.Equal(t, []*Request{
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
			expectedBody:        "foobar",
			expectedTimesCalled: 1,
		},
		{
			path:                "/second",
			method:              http.MethodPost,
			returnError:         fmt.Errorf("oops"),
			expectedTimesCalled: 1,
		},
	}, mock.transport.requests)

	req1, _ := http.NewRequest(http.MethodPost, "/first?param1=value1", strings.NewReader("foobar"))
	req1.Header.Add("Authorization", "Bearer TOKEN")
	response1, err := mock.Do(req1)
	if response1 != nil && response1.Body != nil {
		_ = response1.Body.Close()
	}

	assert.NoError(t, err)
	assert.Equal(t, 1, mock.transport.requests[0].timesCalled)
	assert.Equal(t, response1.StatusCode, http.StatusOK)
	assert.False(t, mockT.Failed())

	req2, _ := http.NewRequest(http.MethodPost, "/second", nil)
	response2, err := mock.Do(req2)
	if response2 != nil && response2.Body != nil {
		_ = response2.Body.Close()
	}

	assert.Error(t, err)
	assert.Equal(t, 1, mock.transport.requests[1].timesCalled)
	assert.False(t, mockT.Failed())

	req3, _ := http.NewRequest(http.MethodOptions, "/toomuch", nil)
	response3, err := mock.Do(req3)
	if response3 != nil && response3.Body != nil {
		_ = response3.Body.Close()
	}

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
			ReturnBody("hello world"),
		).
		WithRequest(http.MethodPut, "/second",
			ExpectQueryParamValues("param", []string{"value1", "value2"}),
			ReturnStatus(http.StatusOK),
		).
		WithRequest(http.MethodPost, "/third",
			ReturnError(fmt.Errorf("oops")),
		)

	assert.Equal(t, []*Request{
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
			expectedJSON:        []byte(`{"foo": "bar"}`),
			expectedTimesCalled: 1,
		},
		{
			path:         "/second",
			method:       http.MethodPut,
			returnStatus: http.StatusOK,
			expectedQueryParams: url.Values{
				"param": {"value1", "value2"},
			},
			expectedTimesCalled: 1,
		},
		{
			path:                "/third",
			method:              http.MethodPost,
			returnError:         fmt.Errorf("oops"),
			expectedTimesCalled: 1,
		},
	}, mock.transport.requests)

	req1, _ := http.NewRequest(http.MethodPost, "/first?param1=value1", strings.NewReader(`{"foo": "bar"}`))
	req1.Header.Add("Authorization", "Bearer TOKEN")
	response1, err := mock.Do(req1)
	if response1.Body != nil {
		_ = response1.Body.Close()
	}

	assert.NoError(t, err)
	assert.Equal(t, response1.StatusCode, http.StatusNoContent)
	assert.Equal(t, 1, mock.transport.requests[0].timesCalled)
	assert.False(t, mockT.Failed())

	req2, _ := http.NewRequest(http.MethodPut, "/second?param=value1&param=value2", nil)
	response2, err := mock.Do(req2)
	if response2.Body != nil {
		_ = response2.Body.Close()
	}

	assert.NoError(t, err)
	assert.Equal(t, response2.StatusCode, http.StatusOK)
	assert.Equal(t, 1, mock.transport.requests[1].timesCalled)
	assert.False(t, mockT.Failed())

	req3, _ := http.NewRequest(http.MethodPost, "/third", nil)
	response3, err := mock.Do(req3)
	if response3 != nil && response3.Body != nil {
		_ = response3.Body.Close()
	}

	assert.Error(t, err)
	assert.Equal(t, 1, mock.transport.requests[2].timesCalled)
	assert.False(t, mockT.Failed())

	req4, _ := http.NewRequest(http.MethodOptions, "/toomuch", nil)
	response4, err := mock.Do(req4)
	if response4 != nil && response4.Body != nil {
		_ = response4.Body.Close()
	}

	assert.Error(t, err)
	assert.True(t, mockT.Failed())
}

func Test_httpMock_wrong_call(t *testing.T) {
	mockT := new(testing.T)
	mock := New(mockT)

	assert.Equal(t, []*Request{}, mock.transport.requests)

	req0, _ := http.NewRequest(http.MethodGet, "/bad", nil)
	response, err := mock.Do(req0)
	if response != nil && response.Body != nil {
		_ = response.Body.Close()
	}

	assert.Error(t, err)
	assert.True(t, mockT.Failed())
}

func Test_httpMock_return_headers(t *testing.T) {
	mockT := new(testing.T)
	mock := New(mockT).
		WithRequest(http.MethodPost, "/first",
			ExpectHeader("Authorization", []string{"Bearer TOKEN"}),
			ExpectJSON(`{"foo": "bar"}`),
			ExpectQueryParam("param1", "value1"),
			ReturnStatus(http.StatusNoContent),
			ReturnBody("hello world"),
			ReturnHeader("Content-Type", []string{"application/json"}),
			ReturnHeader("Content-Length", []string{"1024"}),
		)

	assert.Equal(t, []*Request{
		{
			path:         "/first",
			method:       http.MethodPost,
			returnStatus: http.StatusNoContent,
			returnBody:   "hello world",
			returnHeaders: map[string][]string{
				"Content-Type":   {"application/json"},
				"Content-Length": {"1024"},
			},
			expectedHeaders: map[string][]string{
				"Authorization": {"Bearer TOKEN"},
			},
			expectedQueryParams: url.Values{
				"param1": {"value1"},
			},
			expectedJSON:        []byte(`{"foo": "bar"}`),
			expectedTimesCalled: 1,
		},
	}, mock.transport.requests)

	req, _ := http.NewRequest(http.MethodPost, "/first?param1=value1", strings.NewReader(`{"foo": "bar"}`))
	req.Header.Add("Authorization", "Bearer TOKEN")
	response, err := mock.Do(req)
	if response.Body != nil {
		_ = response.Body.Close()
	}

	assert.NoError(t, err)
	assert.Equal(t, response.StatusCode, http.StatusNoContent)
	assert.Equal(t, 1, mock.transport.requests[0].timesCalled)
	assert.False(t, mockT.Failed())
}
