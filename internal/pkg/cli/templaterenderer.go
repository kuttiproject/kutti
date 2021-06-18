package cli

import (
	"io"
	"text/template"
)

// TemplateRenderer renders a template.
type TemplateRenderer struct {
	template *template.Template
}

// Render executes this renderer's template on the specified writer.
func (t *TemplateRenderer) Render(out io.Writer, arg interface{}) {
	t.template.Execute(out, arg)
}

// NewTemplateRenderer returns a new TemplateRender, or an error.
func NewTemplateRenderer(name string, templatesource string) (*TemplateRenderer, error) {
	result := &TemplateRenderer{
		template: template.New(name),
	}

	_, err := result.template.Parse(templatesource)

	return result, err
}
