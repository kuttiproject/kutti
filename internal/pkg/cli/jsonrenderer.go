package cli

import (
	"encoding/json"
	"io"
	"strings"
)

// JSONRenderer generates prettified JSON.
type JSONRenderer struct {
	indent int
}

// Render writes prettified JSON to the specified writer.
// It uses json.MarshalIndent for formatting.
func (j *JSONRenderer) Render(out io.Writer, arg interface{}) {
	result, _ := json.MarshalIndent(
		arg,
		"",
		strings.Repeat(" ", j.indent),
	)
	out.Write(result)
	out.Write([]byte("\n"))
}

// NewJSONRenderer returns a new JSON Renderer, which will render
// prettified JSON indented by the specified number of spaces.
func NewJSONRenderer(indentspaces int) *JSONRenderer {
	return &JSONRenderer{
		indent: indentspaces,
	}
}
