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
}

func InputField(opts *InputFieldOptions) Node {
	return Div(
		Class("input-field"),
		Label(For(opts.ID), Class("label"), Text(opts.Label)),
		Input(Class("input"), Type(opts.Type)),
		If(opts.Error != "", Span(
			Class("error-explanation"),
			Text(opts.Error),
		)),
	)
}
