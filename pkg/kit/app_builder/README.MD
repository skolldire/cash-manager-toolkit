# Module app_builder

The app_builder module of cash-manager-toolkit provides an interface and an implementation of the Builder design pattern, making it easier to build modular and flexible web applications in Go. This module allows for detailed configuration of components such as configurations, repositories, use cases, and HTTP handlers, encapsulating the complexity in a step-by-step building process.

## Características

- Modular Configuration: Load and apply specific configurations for each part of the application.
- Repository Initialization: Facilitates the creation and configuration of repositories for data persistence.
- Use Cases: Allows defining and configuring use cases, encapsulating business logic.
- HTTP Handlers: Configures HTTP handlers to expose use cases through a REST API.

## Installation

To use the app_builder module, first make sure you have cash-manager-toolkit available in your project. If cash-manager-toolkit is not yet installed, you can get it with:

```bash
go get github.com/skolldire/cash-manager-toolkit
```

## Uso

### Step 1: Implement the Concrete Builder

Create your own structure, for example, AppBuilder, that implements the Builder interface provided by the module.

```go
package main

import (
    "github.com/skolldire/cash-manager-toolkit/kit/app_builder"
    "github.com/skolldire/cash-manager-toolkit/kit/app_engine"
    "path/to/your_endpoints"
)

type app struct {
    engine       app_engine.Engine
    repositories repositories
    useCases     useCases
    handlers     handlers
}

func (a app) Run() error {
    return a.engine.App.Run()
}

type AppBuilder struct {
    application *app
}

var _ app_builder.Builder = (*AppBuilder)(nil)

func NewAppBuilder() *AppBuilder {
    a, err := web_engine.Apply()
    if err != nil {
        panic(err)
    }
    return &AppBuilder{application: &app{engine: a}}
}

func (b *AppBuilder) LoadConfig() app_builder.Builder {
    // Implementación de la carga de configuración
    return b
}

func (b *AppBuilder) InitRepositories() app_builder.Builder {
    b.application.repositories.TestRepo = test_repository.NewRepository(b.application.engine.HttpClient["with_cb"])
    return b
}

func (b *AppBuilder) InitUseCases() app_builder.Builder {
    b.application.useCases.TestCase = test_use_case.NewService(b.application.repositories.TestRepo)
    return b
}

func (b *AppBuilder) InitHandlers() app_builder.Builder {
    b.application.handlers.TestHandler = test_handler.NewService(b.application.useCases.TestCase)
    return b
}

func (b *AppBuilder) InitRoutes() app_builder.Builder {
    b.application.engine.App.Router.Get("/v1/test/", b.application.handlers.TestHandler.Apply())
    return b
}

func (b *AppBuilder) Build() app_builder.App {
    return b.application
}
```

### Step 2: Build the Application

Use the director provided by web_app_builder to build and run your application.

```go
package main

import (
    "github.com/skolldire/cash-manager-toolkit/kit/app_builder"
)

func main() {
    builder := NewAppBuilder()
    application := app_builder.Apply(builder)
    err := application.Run()
    if err != nil {
        panic(err)
    }
}
```

> [!TIP]
> The web_app_builder module is designed to work in isolation or in combination with web_engine, also included in op-go-toolkit. This allows for seamless integration and extended configuration when used together, leveraging the capabilities of app_engine for advanced management of application components based on configuration files.