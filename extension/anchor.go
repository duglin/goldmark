package extension

import (
	"strconv"

	"github.com/duglin/goldmark"
	"github.com/duglin/goldmark/ast"
	"github.com/duglin/goldmark/parser"
	"github.com/duglin/goldmark/renderer"
	"github.com/duglin/goldmark/renderer/html"
	"github.com/duglin/goldmark/text"
	"github.com/duglin/goldmark/util"
)

// Extender adds support for anchors to a Goldmark Markdown parser.
//
// Use it by installing it into the [goldmark.Markdown] object upon creation.
// For example:
//
//	goldmark.New(
//		// ...
//		goldmark.WithExtensions(
//			// ...
//			&anchor.Extender{},
//		),
//	)
type AnchorExtender struct {
	// Texter determines the anchor text.
	//
	// Defaults to '¶' if unspecified.
	Texter Texter

	// Position specifies where the anchor will be placed in a header.
	//
	// Defaults to After.
	Position Position

	// Attributer determines the attributes
	// that will be associated with the anchor link.
	//
	// Defaults to adding a 'class="anchor"' attribute.
	Attributer Attributer

	// Unsafe specifies whether the Texter values will be escaped or not.
	// Setting this to true can lead to HTML injection if you don't handle
	// Texter values with care.
	//
	// Defaults to false.
	Unsafe bool
}

var _ goldmark.Extender = (*AnchorExtender)(nil)

// Extend extends the provided Goldmark Markdown.
func (e AnchorExtender) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&Transformer{
				Texter:     e.Texter,
				Position:   e.Position,
				Attributer: e.Attributer,
			}, 100),
		),
	)
	md.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&Renderer{
				Position: e.Position,
				Unsafe:   e.Unsafe,
			}, 100),
		),
	)
}

// Kind is the NodeKind used by anchor nodes.
var Kind = ast.NewNodeKind("Anchor")

// Node is an anchor node in the Markdown AST.
type Node struct {
	ast.BaseInline

	// ID of the header this anchor is for.
	ID []byte

	// Level of the header that this anchor is for.
	Level int

	// Value is the text inside the anchor.
	// Typically this is a fixed string
	// like '¶' or '#'.
	Value []byte
}

// Kind reports that this is a Anchor node.
func (*Node) Kind() ast.NodeKind { return Kind }

// Dump dumps this node to stdout for debugging.
func (n *Node) Dump(src []byte, level int) {
	ast.DumpHelper(n, src, level, map[string]string{
		"ID":    string(n.ID),
		"Value": string(n.Value),
		"Level": strconv.Itoa(n.Level),
	}, nil)
}

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[After-0]
	_ = x[Before-1]
}

const _Position_name = "AfterBefore"

var _Position_index = [...]uint8{0, 5, 11}

func (i Position) String() string {
	if i < 0 || i >= Position(len(_Position_index)-1) {
		return "Position(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Position_name[_Position_index[i]:_Position_index[i+1]]
}

var (
	_defaultTexter     = Text("¶")
	_defaultAttributer = Attributes{"class": "anchor"}
)

// HeaderInfo holds information about a header
// for which an anchor is being considered.
type HeaderInfo struct {
	// Level of the header.
	Level int

	// Identifier for the header on the page.
	// This will typically become part of the URL fragment.
	ID []byte
}

// Texter determines the anchor text.
//
// This is the clickable text displayed next to the header
// which tells readers that they can use it as an anchor to the header.
//
// By default, we will use the string '¶'.
type Texter interface {
	// AnchorText returns the anchor text
	// that should be used for the provided header info.
	//
	// If AnchorText returns an empty slice or nil,
	// an anchor will not be generated for this header.
	AnchorText(*HeaderInfo) []byte
}

// Text builds a Texter that uses a constant string
// as the anchor text.
//
// Pass this into [Extender] or [Transformer]
// to specify a custom anchor text.
//
//	anchor.Extender{
//		Texter: Text("#"),
//	}
func Text(s string) Texter {
	return textTexter(s)
}

type textTexter []byte

func (t textTexter) AnchorText(*HeaderInfo) []byte {
	return []byte(t)
}

// Position specifies where inside a heading we should place an anchor [Node].
type Position int

//go:generate stringer -type Position

const (
	// After places the anchor node after the heading text.
	//
	// This is the default.
	After Position = iota

	// Before places the anchor node before the heading text.
	Before
)

// Attributer determines attributes that will be attached to an anchor node.
//
// By default, we will add 'class="anchor"' to all nodes.
type Attributer interface {
	// AnchorAttributes returns the attributes
	// that should be attached to the anchor node
	// for the given header.
	//
	// If AnchorAttributes returns an empty map or nil,
	// no attributes will be added.
	AnchorAttributes(*HeaderInfo) map[string]string
}

// Attributes is an Attributer that uses a constant set of attributes
// for all anchor nodes.
//
// Pass this into [Extender] or [Transformer] to specify custom attributes.
//
//	anchor.Extender{
//		Attributer: Attributes{"class": "permalink"},
//	}
type Attributes map[string]string

var _ Attributer = Attributes{}

// AnchorAttributes reports the attributes associated with this object
// for all headers.
func (as Attributes) AnchorAttributes(*HeaderInfo) map[string]string {
	return as
}

// Transformer transforms a Goldmark Markdown AST,
// adding anchor [Node] objects for headers across the document.
type Transformer struct {
	// Texter determines the anchor text.
	//
	// Defaults to '¶' for all headers if unset.
	Texter Texter

	// Position specifies where the anchor will be placed in a header.
	//
	// Defaults to After.
	Position Position

	// Attributer determines the attributes
	// that will be associated with the anchor link.
	//
	// Defaults to adding a 'class="anchor"' attribute
	// for all headers if unset.
	Attributer Attributer
}

var _ parser.ASTTransformer = (*Transformer)(nil)

// Transform traverses and transforms the provided Markdown document.
//
// This method is typically called by Goldmark
// and should not need to be invoked directly.
func (t *Transformer) Transform(doc *ast.Document, _ text.Reader, _ parser.Context) {
	tr := transform{
		Attributer: t.Attributer,
		Position:   t.Position,
		Texter:     t.Texter,
	}
	if tr.Attributer == nil {
		tr.Attributer = _defaultAttributer
	}
	if tr.Texter == nil {
		tr.Texter = _defaultTexter
	}

	_ = ast.Walk(doc, tr.Visit)
	// Visit always returns a nil error.
}

// transform holds state for a single transformation traversal.
type transform struct {
	Texter     Texter
	Position   Position
	Attributer Attributer
}

func (t *transform) Visit(n ast.Node, enter bool) (ast.WalkStatus, error) {
	if !enter {
		return ast.WalkContinue, nil
	}
	h, ok := n.(*ast.Heading)
	if !ok {
		return ast.WalkContinue, nil
	}

	t.transform(h)
	return ast.WalkSkipChildren, nil
}

func (t *transform) transform(h *ast.Heading) {
	idattr, ok := h.AttributeString("id")
	if !ok {
		return
	}

	id, ok := idattr.([]byte)
	if !ok {
		return
	}

	info := HeaderInfo{
		Level: h.Level,
		ID:    id,
	}

	text := t.Texter.AnchorText(&info)
	if len(text) == 0 {
		return
	}

	n := &Node{
		ID:    id,
		Level: h.Level,
		Value: text,
	}

	for name, value := range t.Attributer.AnchorAttributes(&info) {
		n.SetAttributeString(name, []byte(value))
	}

	// If the header has no children yet, just append the anchor.
	if h.ChildCount() == 0 {
		h.AppendChild(h, n)
		return
	}

	if t.Position == Before {
		h.InsertBefore(h, h.FirstChild(), n)
	} else {
		h.InsertAfter(h, h.LastChild(), n)
	}
}

// Renderer renders anchor [Node]s.
type Renderer struct {
	// Position specifies where in the header text
	// the anchor is being added.
	Position Position
	// Unsafe specifies whether the Texter values will be HTML escaped or
	// not.
	Unsafe bool
}

var _ renderer.NodeRenderer = (*Renderer)(nil)

// RegisterFuncs registers functions against the provided goldmark Registerer.
func (r *Renderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(Kind, r.RenderNode)
}

// RenderNode renders an anchor node.
// Goldmark will invoke this method when it encounters a Node.
func (r *Renderer) RenderNode(w util.BufWriter, _ []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// If position is Before, we need to add the anchor when entering;
	// otherwise when exiting.
	if (r.Position == Before) != entering {
		return ast.WalkContinue, nil
	}

	n := node.(*Node)
	if len(n.ID) == 0 {
		return ast.WalkContinue, nil
	}

	// Add leading/trailing ' ' depending on position.
	if r.Position == Before {
		defer func() {
			_ = w.WriteByte(' ')
		}()
	} else {
		_ = w.WriteByte(' ')
	}

	_, _ = w.WriteString("<a")
	html.RenderAttributes(w, node, nil)
	_, _ = w.WriteString(` href="#`)
	_, _ = w.Write(util.EscapeHTML(n.ID))
	_, _ = w.WriteString(`">`)
	if r.Unsafe {
		_, _ = w.Write(n.Value)
	} else {
		_, _ = w.Write(util.EscapeHTML(n.Value))
	}
	_, _ = w.WriteString("</a>")

	return ast.WalkContinue, nil
}
