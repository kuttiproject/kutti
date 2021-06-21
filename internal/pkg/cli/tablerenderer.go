package cli

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"text/template"
)

// TableColumn represents a column of tabular output.
type TableColumn struct {
	// Name should exactly match the name of an input field.
	Name string
	// Title can be left blank. The name will be used.
	Title string
	// DefaultCheck can be set to true to mark a row as default
	// based on the value of this column.
	DefaultCheck bool
	// Width specifies the column width in characters.
	Width int
	// FormatPrefix specifies any go template function that
	// should precede the value.
	FormatPrefix string
}

// TableRenderer generates tabular output given a set of TableColumns.
// It generates and executes a template to print values for
// the columns.
type TableRenderer struct {
	columns  []*TableColumn
	template *template.Template
}

func (f *TableRenderer) preparesimplemap() string {
	for _, column := range f.columns {
		if column.Title == "" {
			column.Title = column.Name
		}
	}
	return "{{ range $key, $value := .}}{{ $key }}\t{{ $value }}\t\n{{ end }}"
}

func (f *TableRenderer) prepare() string {
	var bodybuilder strings.Builder

	bodybuilder.WriteString("{{ range . }}")

	for _, column := range f.columns {
		if column.Title == "" {
			column.Title = column.Name
		}

		fmt.Fprint(
			&bodybuilder,
			"{{ ",
		)

		if column.FormatPrefix != "" {
			fmt.Fprintf(
				&bodybuilder,
				"%v ",
				column.FormatPrefix,
			)
		}

		fmt.Fprintf(
			&bodybuilder,
			".%v ",
			column.Name,
		)

		if column.DefaultCheck {
			fmt.Fprint(
				&bodybuilder,
				"| decoratedefault ",
			)
		}

		fmt.Fprint(
			&bodybuilder,
			" }}\t",
		)
	}

	bodybuilder.WriteString("\n{{ end }}")

	return bodybuilder.String()
}

// Render writes table-formatted output to the specified writer.
// It uses the text/tabwriter package for formatting.
func (f *TableRenderer) Render(out io.Writer, arg interface{}) {
	writer := tabwriter.NewWriter(out, 8, 1, 1, ' ', 0)

	// Write headers
	for _, item := range f.columns {
		formatstring := fmt.Sprintf("%%-%dv\t", item.Width)
		fmt.Fprintf(writer, formatstring, strings.ToUpper(item.Title))
	}
	fmt.Fprintln(writer)

	f.template.Execute(writer, arg)
	writer.Flush()
}

// NewTableRenderer returns a new TableRenderer, and generates a template to
// produce tabular output.
func NewTableRenderer(
	name string,
	columns []*TableColumn,
	defaultvalue string,
) *TableRenderer {

	result := &TableRenderer{
		columns:  columns,
		template: template.New(name),
	}

	result.template.Funcs(template.FuncMap{
		"decoratedefault": func(value string) string {
			if value == defaultvalue {
				return value + "*"
			}

			return value
		},
		"prettytime": prettyTime,
	})

	templatestring := result.prepare()

	template.Must(result.template.Parse(templatestring))

	return result
}

// NewMapTableRenderer returns a TableRenderer, and generates a template, for the
// special case of rendering a simple key/value map.
func NewMapTableRenderer(
	name string,
	columns []*TableColumn,
) *TableRenderer {

	result := &TableRenderer{
		columns:  columns,
		template: template.New(name),
	}

	templatestring := result.preparesimplemap()
	template.Must(result.template.Parse(templatestring))

	return result
}
