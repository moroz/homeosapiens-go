package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type DisplayFieldOptions struct {
	Label string
	Value string
	// Name, if set, renders a hidden input so the value is included in form submissions.
	Name string
}

func DisplayField(opts *DisplayFieldOptions) Node {
	return Div(Class("input-field"),
		Label(Class("label leading-tight font-semibold"), Text(opts.Label)),
		P(Text(opts.Value)),
		If(opts.Name != "", Input(Type("hidden"), Name(opts.Name), Value(opts.Value))),
	)
}
