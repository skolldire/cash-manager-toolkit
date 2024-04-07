# Módulo web_app_builder

El módulo `web_app_builder` de `op-go-toolkit` ofrece una interfaz y una implementación del patrón de diseño Builder, facilitando la construcción modular y flexible de aplicaciones web en Go. Este módulo permite una configuración detallada de componentes como configuraciones, repositorios, casos de uso y manejadores HTTP, encapsulando la complejidad en un proceso de construcción paso a paso.

## Características

- **Configuración Modular**: Carga y aplica configuraciones específicas para cada parte de la aplicación.
- **Inicialización de Repositorios**: Facilita la creación y configuración de repositorios para la persistencia de datos.
- **Casos de Uso**: Permite definir y configurar los casos de uso, encapsulando la lógica de negocio.
- **Manejadores HTTP**: Configura los manejadores HTTP para exponer los casos de uso a través de una API REST.

## Instalación

Para utilizar el módulo `web_app_builder`, primero asegúrate de tener `op-go-toolkit` disponible en tu proyecto. Si `op-go-toolkit` aún no está instalado, puedes obtenerlo con:

```bash
go get github.com/melisource/fury_op-go-toolkit
```

## Uso

### Paso 1: Implementar el Builder Concreto

Crea tu propia estructura, por ejemplo, `AppBuilder`, que implemente la interfaz `Builder` proporcionada por el módulo.

```go
package main

import (
    "github.com/melisource/fury_op-go-toolkit/kit/web_app_builder"
    "github.com/melisource/fury_op-go-toolkit/kit/web_engine"
    "path/to/your_endpoints"
)

type app struct {
    engine       web_engine.Engine
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

var _ web_app_builder.Builder = (*AppBuilder)(nil)

func NewAppBuilder() *AppBuilder {
    a, err := web_engine.Apply()
    if err != nil {
        panic(err)
    }
    return &AppBuilder{application: &app{engine: a}}
}

func (b *AppBuilder) LoadConfig() web_app_builder.Builder {
    // Implementación de la carga de configuración
    return b
}

func (b *AppBuilder) InitRepositories() web_app_builder.Builder {
    b.application.repositories.TestRepo = test_repository.NewRepository(b.application.engine.HttpClient["with_cb"])
    return b
}

func (b *AppBuilder) InitUseCases() web_app_builder.Builder {
    b.application.useCases.TestCase = test_use_case.NewService(b.application.repositories.TestRepo)
    return b
}

func (b *AppBuilder) InitHandlers() web_app_builder.Builder {
    b.application.handlers.TestHandler = test_handler.NewService(b.application.useCases.TestCase)
    return b
}

func (b *AppBuilder) InitRoutes() web_app_builder.Builder {
    b.application.engine.App.Router.Get("/v1/test/", b.application.handlers.TestHandler.Apply())
    return b
}

func (b *AppBuilder) Build() web_app_builder.App {
    return b.application
}
```

### Paso 2: Construir la Aplicación

Utiliza el director proporcionado por `web_app_builder` para construir y ejecutar tu aplicación.

```go
package main

import (
    "github.com/melisource/fury_op-go-toolkit/kit/web_app_builder"
)

func main() {
    builder := NewAppBuilder()
    application := web_app_builder.Apply(builder)
    err := application.Run()
    if err != nil {
        panic(err)
    }
}
```

> [!TIP]
> El módulo `web_app_builder` está diseñado para trabajar de forma aislada o en combinación con [web_engine](pkg/kit/web_engine), también incluido en `op-go-toolkit`. Esto permite una integración fluida y una configuración extendida cuando se utilizan juntos, aprovechando las capacidades de `web_engine` para la gestión avanzada de los componentes de la aplicación basandose en archivos de configuración.
