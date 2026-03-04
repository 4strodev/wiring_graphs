package container_test

import (
	"fmt"

	"github.com/4strodev/wiring_graphs/pkg/container"
)

func ExampleContainer_Token() {
	cont := container.New()

	// Token registers dependencies by a string key.
	// Useful when you need multiple values of the same type.
	err := cont.Token(map[string]any{
		"greeting": func() string { return "hello" },
		"farewell": func() string { return "goodbye" },
	})
	if err != nil {
		panic(err)
	}

	greeting, _ := container.ResolveToken[string](cont, "greeting")
	farewell, _ := container.ResolveToken[string](cont, "farewell")
	fmt.Println(greeting)
	fmt.Println(farewell)
	// Output:
	// hello
	// goodbye
}

func ExampleContainer_TokenSingleton() {
	cont := container.New()

	callCount := 0
	err := cont.TokenSingleton(map[string]any{
		"counter": func() string {
			callCount++
			return fmt.Sprintf("call #%d", callCount)
		},
	})
	if err != nil {
		panic(err)
	}

	first, _ := container.ResolveToken[string](cont, "counter")
	second, _ := container.ResolveToken[string](cont, "counter")
	fmt.Println(first)
	fmt.Println(second)
	// Output:
	// call #1
	// call #1
}
