<h1><img src="https://www.dicebear.com/logo-readme.svg" width="28" /> DiceBear Core (Go)</h1>

Go implementation of the DiceBear avatar library. Generates deterministic SVG
avatars from style definitions and a seed string.

DiceBear is available for multiple languages. All implementations share the same
PRNG and rendering pipeline, producing identical SVG output for the same seed,
style, and options — regardless of the language used.

[Playground](https://www.dicebear.com/playground) |
[Documentation](https://www.dicebear.com/introduction)

## Installation

```sh
go get github.com/dicebear/dicebear-go/v10
```

Requires Go 1.23 or newer.

## Usage

```go
package main

import (
	"fmt"

	dicebear "github.com/dicebear/dicebear-go/v10"
	"github.com/dicebear/styles/v10"
)

func main() {
	// From a style definition (raw JSON, e.g. the pure-data styles module)
	style, err := dicebear.NewStyle([]byte(styles.Lorelei))
	if err != nil {
		panic(err)
	}

	avatar, err := dicebear.NewAvatar(style, map[string]any{
		"seed": "John Doe",
		"size": 128,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(avatar.SVG())     // SVG string
	fmt.Println(avatar.DataURI()) // data:image/svg+xml;charset=utf-8,...
}
```

### With options

```go
avatar, err := dicebear.NewAvatar(style, map[string]any{
	"seed":                "Jane",
	"size":                64,
	"flip":                "horizontal",
	"backgroundColor":     []string{"#0077b6", "#00b4d8"},
	"backgroundColorFill": "linear",
	"scale":               []float64{0.8, 1.0},
	"borderRadius":        10,
})
```

### Inspecting a style

```go
style, err := dicebear.NewStyle(definitionJSON)
if err != nil {
	panic(err)
}

// Inspect the available options for a style
descriptor := dicebear.NewOptionsDescriptor(style).ToJSON()

// Create multiple avatars from the same style
alice, _ := dicebear.NewAvatar(style, map[string]any{"seed": "Alice"})
bob, _ := dicebear.NewAvatar(style, map[string]any{"seed": "Bob"})
```

The style definitions ship as the pure-data
[`github.com/dicebear/styles/v10`](https://pkg.go.dev/github.com/dicebear/styles/v10)
module; the two draft-07 JSON Schemas come from
[`github.com/dicebear/schema`](https://pkg.go.dev/github.com/dicebear/schema).
Color helpers live in the `github.com/dicebear/dicebear-go/v10/color`
sub-package.

## License

The source code of DiceBear is released under the **MIT License**.

## Sponsors

Advertisement: Many thanks to our sponsors who provide us with free or
discounted products.

<a href="https://bunny.net/" target="_blank" rel="noopener noreferrer">
    <picture>
        <source media="(prefers-color-scheme: dark)" srcset="https://www.dicebear.com/sponsors/bunny-light.svg">
        <source media="(prefers-color-scheme: light)" srcset="https://www.dicebear.com/sponsors/bunny-dark.svg">
        <img alt="bunny.net" src="https://www.dicebear.com/sponsors/bunny-dark.svg" height="64">
    </picture>
</a>
