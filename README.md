# Wiring graphs
This package is an autowiring and DI injection library. Designed to make use of generics and type safe resolvers.
It is a container based DI framework. That means that you will need access to the container to being able to resolve
dependencies.

## Features
- Circular dependencies detection: Thanks to the usage of graphs and DFS this library is able to detect circular
dependencies. That detection is meant to be produced early avoiding to crash on production. Circular dependencies
should be detected while developing the app.

- Type safe resolving: The container only holds dependencies all the logic of resolving is encapsulated in a
`container.Resolve` function that manages the type detection successfully.

## Considerations
- Reflection: If you need an extreme performance this library is not for you. All the execution of adding and resolving
dependencies is done with reflection. That consumes time and CPU, which is not ideal for ultra performant services.

## Usage

### Define dependencies
First you need to create a container and declare dependencies. There are two lifetimes:
- **Singleton**: The resolver runs once and the result is cached for subsequent resolutions.
- **Transient**: The resolver runs every time the dependency is resolved.

```go
import (
    "bytes"
    "log/slog"
    "os"

    "github.com/4strodev/wiring_graphs/pkg/container"
)

func main() {
    cont := container.New()

    // Must() returns a wrapper that panics on errors, allowing method chaining.
    cont.Must().
        // Singleton: resolved once and cached.
        Singleton(func() *slog.Logger {
            return slog.New(slog.NewJSONHandler(os.Stdout, nil))
        }).
        // Transient: a new instance is created on every resolution.
        Dependencies(func() *bytes.Buffer {
            return bytes.NewBufferString("hello")
        })

    // You can also register without Must(), handling errors explicitly.
    err := cont.Singleton(func() *bytes.Buffer {
        return bytes.NewBufferString("hello")
    })
    if err != nil {
        panic(err)
    }
}
```

Resolvers can depend on other registered types. Input parameters are resolved
automatically by the container:

```go
cont.Must().
    Singleton(func() *slog.Logger {
        return slog.New(slog.NewJSONHandler(os.Stdout, nil))
    }).
    // The *slog.Logger parameter is resolved from the container.
    Dependencies(func(logger *slog.Logger) *bytes.Buffer {
        logger.Info("creating buffer")
        return bytes.NewBufferString("hello")
    })
```

Dependencies can also be registered by **token** (a string key) instead of by type.
This is useful when you need multiple values of the same type:

```go
cont.Must().
    // Register by token (transient).
    Token(map[string]any{
        "stdout": func() *bytes.Buffer { return bytes.NewBufferString("stdout") },
        "stderr": func() *bytes.Buffer { return bytes.NewBufferString("stderr") },
    }).
    // Register by token (singleton).
    TokenSingleton(map[string]any{
        "config": func() *bytes.Buffer { return bytes.NewBufferString("config data") },
    })
```

### Resolve dependencies
Use the generic `Resolve` and `ResolveToken` package-level functions to retrieve
dependencies from the container. Circular dependencies are automatically detected
before resolution.

```go
import (
    "bytes"
    "fmt"
    "log/slog"

    "github.com/4strodev/wiring_graphs/pkg/container"
)

func main() {
    cont := container.New()
    cont.Must().
        Singleton(func() *slog.Logger {
            return slog.New(slog.NewJSONHandler(os.Stdout, nil))
        }).
        Token(map[string]any{
            "buffer": func() *bytes.Buffer { return bytes.NewBufferString("hello") },
        })

    // Resolve by type.
    logger, err := container.Resolve[*slog.Logger](cont)
    if err != nil {
        panic(err)
    }
    logger.Info("resolved logger")

    // Resolve by token.
    buf, err := container.ResolveToken[*bytes.Buffer](cont, "buffer")
    if err != nil {
        panic(err)
    }
    fmt.Println(buf.String()) // hello
}
```

### Fill structs
You can auto-populate a struct's exported fields from the container using `Fill()`.
Fields are resolved by type by default. Use the `wiring` struct tag to resolve by
token or to skip a field.

**Tag format:** `wiring:"tokenName,omit"`

| Tag | Behavior |
|-----|----------|
| *(no tag)* | Resolve by field type |
| `wiring:"myToken"` | Resolve by token `"myToken"` |
| `wiring:",omit"` | Skip this field |

Unexported fields are always skipped automatically.

```go
import (
    "bytes"
    "fmt"
    "log/slog"
    "os"

    "github.com/4strodev/wiring_graphs/pkg/container"
)

// Define a struct with optional wiring tags.
type AppDeps struct {
    Logger  *slog.Logger                        // resolved by type
    Buffer  *bytes.Buffer `wiring:"appBuffer"`  // resolved by token "appBuffer"
    Skipped *bytes.Buffer `wiring:",omit"`       // explicitly skipped
    private *bytes.Buffer                        // unexported, automatically skipped
}

func main() {
    cont := container.New()
    cont.Must().
        Singleton(func() *slog.Logger {
            return slog.New(slog.NewJSONHandler(os.Stdout, nil))
        }).
        Token(map[string]any{
            "appBuffer": func() *bytes.Buffer { return bytes.NewBufferString("filled!") },
        })

    // Fill takes a pointer to a struct.
    var deps AppDeps
    err := cont.Fill(&deps)
    if err != nil {
        panic(err)
    }

    deps.Logger.Info("logger resolved via Fill")
    fmt.Println(deps.Buffer.String()) // filled!
    fmt.Println(deps.Skipped == nil)  // true
    fmt.Println(deps.private == nil)  // true
}
```

`Fill` can also be used inside a resolver to inject multiple dependencies at once:

```go
type ServiceDeps struct {
    Logger *slog.Logger
    Buffer *bytes.Buffer `wiring:"buffer"`
}

cont.Must().
    Singleton(
        func() *slog.Logger {
            return slog.New(slog.NewJSONHandler(os.Stdout, nil))
        },
        // The *Container is automatically available for injection.
        func(c *container.Container) *MyService {
            var deps ServiceDeps
            if err := c.Fill(&deps); err != nil {
                panic(err)
            }
            return NewMyService(deps.Logger, deps.Buffer)
        },
    ).
    Token(map[string]any{
        "buffer": func() *bytes.Buffer { return bytes.NewBufferString("data") },
    })
```
