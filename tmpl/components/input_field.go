package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type InputFieldOptions struct {
	Label      string
	Name       string
	ID         string
	Type       string
	Error      string
	HelperText string
	Value      string
}

func InputField(opts *InputFieldOptions) Node {
	id := opts.Name
	if opts.ID != "" {
		id = opts.ID
	}

	inputType := "text"
	if opts.Type != "" {
		inputType = opts.Type
	}

	return Div(
		Class("input-field"),
		Label(For(opts.ID), Class("label"), Text(opts.Label)),
		Input(Class("input"), Type(inputType), Name(opts.Name), ID(id), Value(opts.Value)),
		If(opts.Error != "", Span(
			Class("error-explanation"),
			Text(opts.Error),
		)),
	)
}
