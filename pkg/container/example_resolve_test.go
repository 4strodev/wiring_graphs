package container_test

import (
	"bytes"
	"fmt"

	"github.com/4strodev/wiring_graphs/pkg/container"
)

func ExampleResolve() {
	cont := container.New()
	cont.Must().Singleton(func() *bytes.Buffer {
		return bytes.NewBufferString("resolved!")
	})

	buf, err := container.Resolve[*bytes.Buffer](cont)
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
	// Output:
	// resolved!
}

func ExampleResolve_withDependencyChain() {
	cont := container.New()
	cont.Must().
		Dependencies(
			func() int { return 42 },
			func(n int) string {
				return fmt.Sprintf("the answer is %d", n)
			},
		)

	result, err := container.Resolve[string](cont)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	// Output:
	// the answer is 42
}

func ExampleResolveToken() {
	cont := container.New()
	cont.Must().Token(map[string]any{
		"message": func() string { return "hello from token" },
	})

	msg, err := container.ResolveToken[string](cont, "message")
	if err != nil {
		panic(err)
	}
	fmt.Println(msg)
	// Output:
	// hello from token
}
