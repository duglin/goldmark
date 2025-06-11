package goldmark_test

import (
	"testing"

	. "github.com/duglin/goldmark"
	"github.com/duglin/goldmark/parser"
	"github.com/duglin/goldmark/testutil"
)

func TestAttributeAndAutoHeadingID(t *testing.T) {
	markdown := New(
		WithParserOptions(
			parser.WithAttribute(),
			parser.WithAutoHeadingID(),
		),
	)
	testutil.DoTestCaseFile(markdown, "_test/options.txt", t, testutil.ParseCliCaseArg()...)
}
