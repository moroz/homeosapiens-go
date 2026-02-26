package components

import (
	"github.com/moroz/homeosapiens-go/tmpl/components/icons"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type SelectOptions struct {
	Label        string
	Name         string
	ID           string
	Value        string
	Required     bool
	Autocomplete string
	Options      []SelectOption
	Localizer    *i18n.Localizer
}

type SelectOption struct {
	Label string
	Value string
}

func SelectComponent(opts *SelectOptions) Node {
	id := opts.Name
	if opts.ID != "" {
		id = opts.ID
	}

	class := "input-field"
	if opts.Required {
		class += " required"
	}

	// Only render Selected() for the first matching value
	var selectedFound bool

	return Div(
		Class(class),
		Label(
			For(id), Class("label grid gap-1"),
			Span(
				Span(
					Class("font-semibold"),
					Text(opts.Label),
				),

				If(opts.Required,
					Span(
						Aria("hidden", "true"),
						Class("ml-1 inline-block text-red-700"),
						Iff(opts.Localizer != nil, func() Node {
							return Title(opts.Localizer.MustLocalizeMessage(&i18n.Message{
								ID: "components.input_field.required",
							}))
						}),
						Text("*"),
					),
				),
			),

			Div(
				Class("wrapper relative"),

				Select(
					ID(id),
					Name(opts.Name),
					Class("h-10 w-full appearance-none rounded-sm border px-2 pr-8 font-normal shadow"),
					If(opts.Autocomplete != "", AutoComplete(opts.Autocomplete)),

					Map(opts.Options, func(option SelectOption) Node {
						return Option(Value(option.Value), Iff(option.Value == opts.Value && !selectedFound, func() Node {
							selectedFound = true
							return Selected()
						}), Text(option.Label))
					}),
				),

				icons.Icon(&icons.IconProps{
					Name:    "chevron-down",
					Classes: "absolute top-1/2 right-1 h-5 w-5 -translate-y-1/2 fill-slate-600",
				}),
			),
		),
	)
}
