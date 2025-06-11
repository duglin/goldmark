package extension

import (
	"testing"

	"github.com/duglin/goldmark"
	"github.com/duglin/goldmark/renderer/html"
	"github.com/duglin/goldmark/testutil"
)

func TestStrikethrough(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			Strikethrough,
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/strikethrough.txt", t, testutil.ParseCliCaseArg()...)
}
