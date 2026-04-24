package components

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type InputFieldOptions struct {
	Label        string
	Name         string
	ID           string
	Type         string
	HelperText   string
	Value        string
	Autocomplete string
	Autofocus    bool
	Required     bool
	Readonly     bool
	Error        any
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

	errMessage := resolveErrorMessage(opts.Error, opts.Name)

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
						Class("ml-1 inline-block text-danger"),
						Iff(opts.Localizer != nil, func() Node {
							return Title(opts.Localizer.MustLocalizeMessage(&i18n.Message{
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
			If(opts.Readonly, ReadOnly()),
		),
		If(errMessage != "", Span(
			Class("error-explanation"),
			Text(errMessage),
		)),
	)
}

func resolveErrorMessage(err any, name string) string {
	switch err := err.(type) {
	case string:
		return err
	case validation.Errors:
		value := err[name]
		if value != nil {
			return value.Error()
		}
	case error:
		return err.Error()
	}
	return ""
}

func InputGroup(children ...Node) Node {
	return Div(
		Class("grid gap-4 desktop:flex desktop:gap-6"),
		Group(children),
	)
}
