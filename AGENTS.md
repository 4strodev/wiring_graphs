# AGENTS.md — Wiring Graphs

## Project Overview

Go dependency injection library using generics and reflection. Provides autowiring
with circular dependency detection via graph-based DFS. Module path:
`github.com/4strodev/wiring_graphs`. Requires Go 1.24.2+.

## Repository Structure

```
pkg/
├── container/        # Core DI container: registration, resolution, struct filling
│   ├── container.go          # Container struct, New(), Singleton/Transient/Token methods
│   ├── container_fill.go     # Fill() for auto-populating struct fields
│   ├── must_container.go     # MustContainer wrapper (panics on error)
│   └── resolve.go            # Generic Resolve[T]() and ResolveToken[T]() functions
├── errors/           # Custom WiringError type with typed error codes
│   └── wiring_error.go
├── internal/         # Internal implementation (not importable externally)
│   ├── collections/
│   │   ├── graph/    # Generic graph with DFS-based circular dependency detection
│   │   ├── queue/    # Generic queue
│   │   ├── set/      # Generic set
│   │   └── stack/    # Generic stack
│   └── testutils/    # Shared test mocks (MyService, MyDeps)
└── resolver/         # DependencyResolver[T] with validation and execution
```

## Build & Task Runner

This project uses [Task](https://taskfile.dev/) (v3) as its task runner. All commands
are defined in `Taskfile.yml`.

### Common Commands

```bash
# Run all tests (verbose, 1 minute timeout)
task test

# Run all tests without cache
task test:nocache

# Run linter (go vet)
task lint

# Update changelog (requires git-cliff)
task changelog:update
```

### Running a Single Test

```bash
# Run a specific test function
go test -v -run TestResolve ./pkg/container/ -timeout 1m

# Run a specific test with subtests
go test -v -run TestGetToken ./pkg/container/ -timeout 1m

# Run tests in a specific package
go test -v ./pkg/internal/collections/graph/ -timeout 1m

# Run a single test without cache
go test -v -run TestDfsVisit ./pkg/internal/collections/graph/ -timeout 1m -count 1
```

### CI Pipeline

GitHub Actions runs on all branch pushes and PRs. The workflow (`.github/workflows/test.yml`)
sets up Go 1.24.2 and Task v3, then runs `task test`.

## Dependencies

- **Only runtime dependency:** None (stdlib + reflection)
- **Test dependency:** `github.com/stretchr/testify` v1.10.0 (`require` and `assert`)

## Code Style Guidelines

### Imports

Use two groups separated by a blank line:
1. Standard library imports
2. Everything else (third-party and local module imports together)

```go
import (
    "reflect"

    "github.com/4strodev/wiring_graphs/pkg/errors"
    "github.com/4strodev/wiring_graphs/pkg/resolver"
)
```

### Naming Conventions

| Element               | Convention               | Example                          |
|-----------------------|--------------------------|----------------------------------|
| Packages              | lowercase, single word   | `container`, `graph`, `testutils`|
| Exported types        | PascalCase               | `Container`, `WiringError`       |
| Unexported types      | camelCase                | `resolverConfig`                 |
| Enum-like constants   | SCREAMING_SNAKE_CASE     | `E_CIRCULAR_DEPENDENCY`, `OUT`   |
| Exported functions    | PascalCase               | `Resolve[T]()`, `New()`          |
| Unexported functions  | camelCase                | `buildConfig()`, `getToken()`    |
| Constructors          | `New()` or `New<Type>()` | `New()`, `NewGraph[T]()`         |
| Receiver variables    | Short (1-2 chars)        | `c` for Container, `n` for Node  |
| Test loop variable    | `ttest`                  | `for _, ttest := range tests`    |

### Generics

The codebase uses generics extensively. Constraints used are `any` and `comparable`
(no custom type constraints). Internal generic types are instantiated with
`resolver.DependencyResolver[any]` at the container level. Use `reflect.TypeFor[T]()`
for type reflection in generic functions.

### Error Handling

- All domain errors use the custom `errors.WiringError` type with typed error codes.
- Construct errors exclusively via `errors.Errorf(code, format, args...)`.
- Error codes are defined as `wiringErrorCode` constants: `E_CIRCULAR_DEPENDENCY`,
  `E_INVALID_RESOLVER`, `E_REDECLARED_DEPENDENCY`, `E_DEPENDENCY_NOT_FOUND`, `E_TYPE_ERROR`.
- Wrap errors with `%w` verb to preserve the error chain.
- `WiringError` implements `Unwrap()` for compatibility with `errors.Is`/`errors.As`.
- The `MustContainer` wrapper converts errors to panics for fluent API chaining.
- Functions return `error` as the last return value (standard Go convention).
- Named return values are used in `resolve()` and `resolveToken()`.

### Types & Receivers

- Use pointer receivers `(c *Container)` for mutating methods.
- Use value receivers `(c Container)` for read-only methods (e.g., `Fill`, `resolve`).
- No interfaces are defined in the codebase; it relies on concrete types and `reflect`.
- The `pkg/errors` package intentionally shadows the stdlib `errors` package.

### File Organization

- Types, constants, and logic are co-located in the same file (no separate types files).
- Split code by responsibility: registration, resolution, fill logic, must-wrapper in
  separate files within the `container` package.
- Test files are co-located with source code using `_test.go` suffix.
- `internal/` is used correctly to hide implementation details from consumers.

## Testing Conventions

### Test Structure

- Use `testing` package with `testify/require` for fatal assertions and `testify/assert`
  for non-fatal assertions.
- Test function naming: `Test<FunctionName>` with `_` separating scenarios
  (e.g., `TestResolve_WithDependencies`, `TestDetectCircularDependencies_SelfReference`).
- Black-box tests use external test package (`package container_test`).
- White-box tests use the source package (`package container`) when testing unexported functions.
- Table-driven tests use anonymous struct slices with `t.Run()` subtests.
- Shared test types live in `pkg/internal/testutils/mocks.go`.

### Test Example Pattern

```go
func TestFeature_Scenario(t *testing.T) {
    // Arrange
    cont := container.New()
    cont.Must().Singleton(func() *MyType { return &MyType{} })

    // Act
    result, err := container.Resolve[*MyType](cont)

    // Assert
    require.NoError(t, err)
    require.NotNil(t, result)
}
```

### Assertions

```go
require.NoError(t, err)       // Fatal if error is non-nil
require.Error(t, err)         // Fatal if error is nil
require.NotNil(t, result)     // Fatal if nil
require.Equal(t, expected, actual)
assert.Equalf(t, expected, actual, "format %s", arg)  // Non-fatal
```
