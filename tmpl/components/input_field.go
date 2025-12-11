package components

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type InputFieldOptions struct {
	Label        string
	Name         string
	ID           string
	Type         string
	Error        string
	HelperText   string
	Value        string
	Autocomplete string
	Autofocus    bool
	Required     bool
	Localizer    *i18n.Localizer
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

	class := "input-field"
	if opts.Required {
		class += " required"
	}

	return Div(
		Class(class),
		Label(For(opts.ID), Class("label"),
			Text(opts.Label),

			If(opts.Localizer == nil && opts.Required,
				Span(
					Class("text-red-700"),
					Text(" *"),
				),
			),

			Iff(opts.Localizer != nil && !opts.Required, func() Node {
				return Span(
					Class("text-slate-700"),
					Text(opts.Localizer.MustLocalizeMessage(&i18n.Message{
						ID: "components.input_field.optional",
					})),
				)
			}),
		),
		Input(
			Class("input"),
			Type(inputType),
			Name(opts.Name),
			ID(id),
			Value(opts.Value),
			If(opts.Autocomplete != "", AutoComplete(opts.Autocomplete)),
			If(opts.Autofocus, AutoFocus()),
			If(opts.Required, Required()),
		),
		If(opts.Error != "", Span(
			Class("error-explanation"),
			Text(opts.Error),
		)),
	)
}
