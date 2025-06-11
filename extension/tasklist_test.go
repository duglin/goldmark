package extension

import (
	"testing"

	"github.com/duglin/goldmark"
	"github.com/duglin/goldmark/renderer/html"
	"github.com/duglin/goldmark/testutil"
)

func TestTaskList(t *testing.T) {
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			TaskList,
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/tasklist.txt", t, testutil.ParseCliCaseArg()...)
}
