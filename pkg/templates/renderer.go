package templates

import (
	"os"
	"text/tabwriter"
	"text/template"
)

// Render takes a layout string and data model, then prints it to stdout.
// It automatically handles tab alignment for tables.
func Render(layout string, data interface{}) error {
	// 1. Setup TabWriter (minwidth, tabwidth, padding, padchar, flags)
	// This ensures columns line up perfectly like a spreadsheet.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// 2. Parse the Template
	t := template.Must(template.New("output").Parse(layout))

	// 3. Execute
	if err := t.Execute(w, data); err != nil {
		return err
	}

	// 4. Flush to terminal
	return w.Flush()
}
