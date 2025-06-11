package extension

import (
	"testing"

	"github.com/duglin/goldmark"
	"github.com/duglin/goldmark/renderer/html"
	"github.com/duglin/goldmark/testutil"
)

func TestTypographer(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			Typographer,
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/typographer.txt", t, testutil.ParseCliCaseArg()...)
}
