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
		Label(For(id), Class("label leading-tight"),
			Span(
				Class("block"),
				Span(
					Class("font-semibold"),
					Text(opts.Label),
				),

				If(opts.Required,
					Span(
						Aria("hidden", "true"),
						Class("ml-1 inline-block text-red-700"),
						Iff(opts.Localizer != nil, func() Node {
							return TitleAttr(opts.Localizer.MustLocalizeMessage(&i18n.Message{
								ID: "components.input_field.required",
							}))
						}),
						Text("*"),
					),
				),

				Iff(opts.Localizer != nil && !opts.Required, func() Node {
					return Span(
						Class("text-sm text-slate-600"),
						Text(" "),
						Text(opts.Localizer.MustLocalizeMessage(&i18n.Message{
							ID: "components.input_field.optional",
						})),
					)
				}),
			),

			If(opts.HelperText != "", Span(
				Class("helper-text block text-sm"),
				ID(opts.ID+"-helper"),
				Text(opts.HelperText),
			)),
		),
		Input(
			Class("input"),
			Type(inputType),
			Name(opts.Name),
			ID(id),
			Value(opts.Value),
			If(opts.Autocomplete != "", AutoComplete(opts.Autocomplete)),
			If(opts.HelperText != "", Aria("describedby", id+"-helper")),
			If(opts.Autofocus, AutoFocus()),
			If(opts.Required, Required()),
		),
		If(opts.Error != "", Span(
			Class("error-explanation"),
			Text(opts.Error),
		)),
	)
}
