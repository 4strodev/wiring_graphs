package container_test

import (
	"bytes"
	"fmt"

	"github.com/4strodev/wiring_graphs/pkg/container"
)

func ExampleNew() {
	cont := container.New()
	fmt.Printf("%T\n", cont)
	// Output:
	// *container.Container
}

func ExampleContainer_Singleton() {
	cont := container.New()

	// Singleton resolvers are executed once and the result is cached.
	err := cont.Singleton(func() *bytes.Buffer {
		return bytes.NewBufferString("singleton value")
	})
	if err != nil {
		panic(err)
	}

	buf1, _ := container.Resolve[*bytes.Buffer](cont)
	buf2, _ := container.Resolve[*bytes.Buffer](cont)

	fmt.Println(buf1.String())
	fmt.Println(buf1 == buf2) // same pointer, resolved once
	// Output:
	// singleton value
	// true
}

func ExampleContainer_Dependencies() {
	cont := container.New()

	// Dependencies registers transient resolvers.
	// A new instance is created on every resolution.
	err := cont.Dependencies(func() *bytes.Buffer {
		return bytes.NewBufferString("transient")
	})
	if err != nil {
		panic(err)
	}

	buf, _ := container.Resolve[*bytes.Buffer](cont)
	fmt.Println(buf.String())
	// Output:
	// transient
}

func ExampleContainer_Dependencies_withDependencies() {
	cont := container.New()

	// Resolver input parameters are resolved automatically from the container.
	cont.Must().
		Dependencies(
			func() *bytes.Buffer {
				return bytes.NewBufferString("hello")
			}, func(buf *bytes.Buffer) string {
				return fmt.Sprintf("got: %s", buf.String())
			},
		)

	result, _ := container.Resolve[string](cont)
	fmt.Println(result)
	// Output:
	// got: hello
}
