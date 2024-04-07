# op-go-toolkit

## Description

This repository contains the **Error Wrapper** feature for the Fury Op Go Toolkit, developed by Open-Platform team. The Fury Op Go Toolkit is a set of tools designed to enhance the development process and improve the productivity of Go-based projects.

## Error Wrapper Package

Contain an abstraction to handle errors in a more friendly way, similar to the one used in Java, the main idea is to set
the properties and handle them by a wrapper of Go Error implementation

Today we have only the CommonApiError for microservices based on HTTP and REST traffic, but in the future we can add
another implementations for CLI, GRPC etc.

> **_NOTE:_**  code and msg should be coherent through the whole microservice, also you can choose if place them in constants or in the same line of the error.

### How to gen errors

```go
package one_repository

import (
    "github.com/melisource/fury_op-go-toolkit/pkg/kit/error_wrapper"
)
// Implementation of repository logic
// ...
// ...
// ...
if somethingHappendAndIWantFail {
    return nil, error_wrapper.NewCommonApiError("R1-001", "[one_repository]::DoSomething Error trying to do something in the repository", nil, http.StatusInternalServerError)
}
return result, nil
```
#### Attention:
The `NewCommonApiError` function accepts an implementation of the error interface as its third parameter.
Although it is possible to send nil values in this parameter, it is advisable to avoid this, as the library will create a generic error to replace the nil error.

The creation of this generic error to replace the nil error is done to avoid the occurrence of nil pointer in the application.

Recommended creation:
```go
error_wrapper.NewCommonApiError("Error code", "message", fmt.Errorf("New error created"), httpcode)
```
Creation not recommended:
```go
error_wrapper.NewCommonApiError("Error code", "message", nil, httpcode)
```

### How to wrap errors in repositories

```go
package one_repository

import (
    "github.com/melisource/fury_op-go-toolkit/pkg/kit/error_wrapper"
)
// Implementation of repository logic
// ...
// ...
// ...

result, err := s.client.DoSomething(ctx, param1, param2)
if err != nil {
    return nil, error_wrapper.NewCommonApiError("R1-001", "[one_repository]::DoSomething Error trying to do something in the repository", err, http.StatusInternalServerError)
}
return result, nil
```

### How to gen errors with internal and external messages

> **_NOTE:_**  code and msg are transformed to JSON format and used at the end for the handler to the response body, if you want to distinguish between internal and external messages, you can use the following format:
```go
package one_repository

import (
    "github.com/melisource/fury_op-go-toolkit/pkg/kit/error_wrapper"
)
// Implementation of repository logic
// ...
// ...
// ...

result, err := s.client.DoSomething(ctx, param1, param2)
if err != nil {
    return nil, error_wrapper.NewCommonApiError("R1-001", "I'm and external and very formal message for the client", fmt.Errorf("[one_repository]::DoSomething Error trying to do something in the repository: %w", err), http.StatusInternalServerError)
}
return result, nil
```

### How to get wrapped errors and use it in use case (basic case)
```go
package one_use_case

import (
    "github.com/melisource/fury_op-go-toolkit/pkg/kit/error_wrapper"
    "github.com/melisource/fury_my-app/internal/platform/one_repository"
)
// Implementation of use case logic
// ...
// ...
// ...

result, err := one_repository.DoSomething(ctx, param1, param2)
if err != nil {
    return nil, err
}
return result, nil
```

### How to use get wrapped errors and create another one in the use case (to generate a StackTrace)
```go
package one_use_case

import (
    "github.com/melisource/fury_op-go-toolkit/pkg/kit/error_wrapper"
    "github.com/melisource/fury_my-app/internal/platform/one_repository"
)
// Implementation of use case logic
// ...
// ...
// ...

result, err := one_repository.DoSomething(ctx, param1, param2)
if err != nil {
    return nil, error_wrapper.NewCommonApiError("UC-001", "[one_use_case]::DoSomething Error trying to do something in the usecase", err, http.StatusInternalServerError)
}
return result, nil
```

### How to use wrapped errors in endpoint handler (CHI and HTTP) 
```go
package one_endpoint

import (
	"encoding/json"
	"github.com/melisource/fury_go-core/pkg/web"
	"github.com/go-chi/chi/v5"
	"github.com/melisource/fury_op-go-toolkit/pkg/kit/error_wrapper"
	"github.com/melisource/fury_my-app/internal/one_use_case"
)

// Implementation of handler logic
// ...
// ...
// ...

result, err := one_use_case.DoSomething(r.Context(), param1, param2)
if err != nil {
    _ = error_wrapper.HandleApiErrorResponse(err, w)
}
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOk)
w.Write(result)
```

> **_NOTE:_**  This function abstract the implementation of how to respond to the client, based on the error type, it will return the correct HTTP status code, code and msg in the response body.

### How to generate Custom Errors for border logics

1. First create a custom error in your entity.go file, like this:

```go
// entity.go file
package two_repository

import (
    "fmt"
)

type MyCustomHappyError struct {
    field string
}

func (m *MyCustomHappyError) Error() string {
    return fmt.Sprintf("%s Something happend and CustomHappyError appears :)", m.field)
}

func NewMyCustomHappyError(field string) error {
    return &MyCustomHappyError{
        field: field,
    }
}

```

2. Then use it in your service.go file, like this:

```go

package two_repository

import (
	"github.com/melisource/fury_op-go-toolkit/pkg/kit/error_wrapper"
)
// Implementation of repository logic
// ...
// ...
// ...


if somethingHappend {
    return nil, error_wrapper.NewCommonApiError("R2-001", "[two_repository:: Error trying to do something in the repository]", two_repository.NewMyCustomHappyError("something"), http.StatusInternalServerError)
}
return result, nil
```

### How to use Custom Errors for border logics in use case
```go
package two_use_case

import (
    "github.com/melisource/fury_op-go-toolkit/pkg/kit/error_wrapper"
    "github.com/melisource/fury_my-app/internal/platform/two_repository"
)
// Implementation of use case logic
// ...
// ...
// ...
result, err := two_repository.DoSomething(ctx, param1, param2)
if err != nil {
    if errors.Is(err, two_repository.CustomHappyError) {
        // Do something with the error
    }
}
```





