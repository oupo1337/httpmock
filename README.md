# httpmock

HTTPMock is a library to easily mock your http clients and describe their behavior.

It's designed to be as simple as possible, first you describe the client behavior :
 - How many calls should it do ? On which path ?
 - What status should it return ?
 - Should it receive headers ? Query parameters ? A body ?

Then you use your mock client, and you finish by checking that all expected calls were done. 

## Examples

### Basic

```go
func Test_basic(t *testing.T) {
    // We declare a mock variable and its behavior.
    // It's an HTTP client waiting for a GET request on /path.
    // It will return a 200 status code.
    mock := httpmock.New(t).
        WithRequest(http.MethodGet, "/path",
            httpmock.ReturnStatus(http.StatusOK), 
        )
	
    // Do something with mock variable (It's a regular http.Client object).
    doSomething(mock)
	
    // We check that all expected requests were done. 
    mock.AssertExpectations(t)
}
```

### Advanced

```go
func Test_advanced(t *testing.T) {
    // We declare a mock variable and its behavior.
    // It's an HTTP client waiting for a POST request on /form.
    // It will return a 201 status code.
    // It expects a body, an authorization header and will return a body
    //
    // It's also waiting for a DELETE request on /route
    // returning a 204 status code.
    // It expects a query param.
    //
    // Calls should be done in the order they are defined.
    mock := httpmock.New(t).
        WithRequest(http.MethodPost, "/form",
            httpmock.ExpectBody(`{"some": "data"}`),
            httpmock.ExpectHeader("Authorization", []string{"Bearer token"}),
            httpmock.ReturnStatus(http.StatusCreated),
            httpmock.ReturnBodyRaw(`{"a": "response"}`),
        ).
        WithRequest(http.MethodDelete, "/route",
            httpmock.ExpectQueryParam("param", "value"),
            httpmock.ReturnStatus(http.StatusNoContent),
        )

    // Do something with mock variable (It's a regular http.Client object).
    doSomething(mock)

    // We check that all expected requests were done. 
    mock.AssertExpectations(t)
}
```

## Documentation

### Options functions
