package container_test

import (
	"bytes"
	"fmt"

	"github.com/4strodev/wiring_graphs/pkg/container"
)

type ServerConfig struct {
	Host string
	Port int
}

func ExampleContainer_Must() {
	cont := container.New()

	// Must() returns a wrapper that panics on errors,
	// allowing method chaining for multiple registrations.
	cont.Must().
		Singleton(func() *bytes.Buffer {
			return bytes.NewBufferString("buffered")
		}).
		Dependencies(func() string {
			return "a transient string"
		}).
		Token(map[string]any{
			"greeting": func() string { return "hi" },
		})

	buf, _ := container.Resolve[*bytes.Buffer](cont)
	str, _ := container.Resolve[string](cont)
	greeting, _ := container.ResolveToken[string](cont, "greeting")
	fmt.Println(buf.String())
	fmt.Println(str)
	fmt.Println(greeting)
	// Output:
	// buffered
	// a transient string
	// hi
}

func ExampleContainer_Fill() {
	type AppDeps struct {
		Buffer *bytes.Buffer // resolved by type
		Token  *bytes.Buffer `wiring:"appBuffer"` // resolved by token "appBuffer"
		Skip   *bytes.Buffer `wiring:",omit"`     // explicitly skipped
	}

	cont := container.New()
	cont.Must().
		Dependencies(func() *bytes.Buffer {
			return bytes.NewBufferString("by type")
		}).
		Token(map[string]any{
			"appBuffer": func() *bytes.Buffer {
				return bytes.NewBufferString("by token")
			},
		})

	var deps AppDeps
	err := cont.Fill(&deps)
	if err != nil {
		panic(err)
	}

	fmt.Println(deps.Buffer.String())
	fmt.Println(deps.Token.String())
	fmt.Println(deps.Skip == nil)
	// Output:
	// by type
	// by token
	// true
}

func ExampleContainer_Fill_insideResolver() {
	// Fill can be used inside a resolver by injecting *Container itself.
	type Deps struct {
		Greeting string
	}

	cont := container.New()
	cont.Must().
		Dependencies(func() string {
			return "hello from Fill"
		}).
		Dependencies(func(c *container.Container) *ServerConfig {
			var deps Deps
			if err := c.Fill(&deps); err != nil {
				panic(err)
			}
			return &ServerConfig{Host: deps.Greeting, Port: 8080}
		})

	cfg, err := container.Resolve[*ServerConfig](cont)
	if err != nil {
		panic(err)
	}

	fmt.Println(cfg.Host)
	fmt.Println(cfg.Port)
	// Output:
	// hello from Fill
	// 8080
}

func ExampleContainer_Derived() {
	parent := container.New()
	parent.Must().Dependencies(func() string {
		return "from parent"
	})

	// Derived creates a child container that delegates to the parent
	// for dependencies not found locally.
	child := parent.Derived()
	child.Must().Dependencies(func() int {
		return 99
	})

	// Child can resolve its own dependencies.
	num, _ := container.Resolve[int](child)
	fmt.Println(num)

	// Child can also resolve dependencies from the parent.
	str, _ := container.Resolve[string](child)
	fmt.Println(str)
	// Output:
	// 99
	// from parent
}

func ExampleContainer_DetectCircularDependencies() {
	cont := container.New()
	cont.Must().
		Dependencies(func() string { return "hello" }).
		Dependencies(func(s string) int { return len(s) })

	// DetectCircularDependencies runs a DFS-based check on the dependency graph.
	// Returns nil when there are no cycles.
	cycle, err := cont.DetectCircularDependencies()
	if err != nil {
		panic(err)
	}
	fmt.Println(cycle == nil)
	// Output:
	// true
}
