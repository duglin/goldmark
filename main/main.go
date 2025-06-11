// GO111MODULE=off go run ~/go/src/github.com/duglin/goldmark/main/main.go

package main

import (
	"bytes"
	"fmt"
	"os"

	// anchor "github.com/abhinav/goldmark-anchor"
	"github.com/duglin/goldmark"
	"github.com/duglin/goldmark/extension"
	"github.com/duglin/goldmark/parser"
	ghtml "github.com/duglin/goldmark/renderer/html"
)

func main() {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.AnchorExtender{
				// Texter: extension.Text("🔗"),
				Texter:   extension.Text("☍"),
				Position: extension.Before, // or extension.After
			},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			// html.WithHardWraps(),
			ghtml.WithUnsafe(),
		),
	)

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s FILE\n", os.Args[0])
		os.Exit(1)
	}
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	html := bytes.Buffer{}
	md.Convert(data, &html)
	fmt.Printf("%s", html.String())
}
