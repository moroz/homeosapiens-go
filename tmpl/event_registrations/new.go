package eventregistrations

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/moroz/homeosapiens-go/internal/countries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"

	. "maragu.dev/gomponents/html"
)

func mapOptions(options []*countries.CountryOption, lang string) []components.SelectOption {
	var combined []components.SelectOption
	for _, o := range options {
		label := o.LabelEn
		if lang == "pl" {
			label = o.LabelPl
		}

		combined = append(combined, components.SelectOption{
			Label: label,
			Value: o.Value,
		})
	}
	return combined
}

func buildCountryOptions(lang string) []components.SelectOption {
	options := countries.OrderedByEnglish
	popular := countries.PopularRegionsEnglish
	if lang == "pl" {
		options = countries.OrderedByPolish
		popular = countries.PopularRegionsPolish
	}

	var combined []components.SelectOption
	combined = mapOptions(popular, lang)

	combined = append(combined, components.SelectOption{
		Label: "---",
		Value: "",
	})

	all := mapOptions(options, lang)
	return append(combined, all...)
}

func New(ctx *types.CustomContext, event *services.EventDetailsDto, params *types.CreateEventRegistrationParams, validationErrors validation.Errors) Node {
	title := event.TitleEn
	if ctx.Language == "pl" {
		title = event.TitlePl
	}
	l := ctx.Localizer

	pageTitle := l.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "event_registrations.new.title",
		TemplateData: map[string]string{
			"Title": title,
		},
	})

	currentPath := fmt.Sprintf("/events/%s/register", event.Slug)
	countryOptions := buildCountryOptions(ctx.Language)

	return layout.Layout(ctx, pageTitle,
		Div(
			Class("card"),
			Header(
				Span(
					Class("text-xl"),
					Text(l.MustLocalizeMessage(&i18n.Message{
						ID: "event_registrations.new.heading",
					})),
				),

				H2(
					Class("text-primary text-2xl font-bold"),
					A(
						Class("no-underline hover:underline"),
						Href(fmt.Sprintf("/events/%s", event.Slug)),
						Text(title),
					),
				),
			),

			Iff(ctx.User == nil, func() Node {
				return Section(
					components.GoogleButton(l.MustLocalizeMessage(&i18n.Message{
						ID: "event_registrations.new.continue_with_google",
					}), currentPath),
				)
			}),

			Main(
				Form(
					Method("POST"),
					Action("/event_registrations"),
					Class("mt-6 space-y-4"),
					Input(Type("hidden"), Name("event_id"), Value(event.Event.ID.String())),
					components.InputField(&components.InputFieldOptions{
						Label: l.MustLocalizeMessage(&i18n.Message{
							ID: "common.users.email",
						}),
						Name:         "email",
						Value:        params.Email,
						Autocomplete: "email",
						Required:     true,
						Localizer:    l,
						HelperText: func() string {
							if ctx.User != nil {
								return ""
							}
							return l.MustLocalizeMessage(&i18n.Message{
								ID: "event_registrations.new.form.helper_text.email",
							})
						}(),
					}),
					components.InputField(&components.InputFieldOptions{
						Label: l.MustLocalizeMessage(&i18n.Message{
							ID: "common.users.given_name",
						}),
						Name:         "given_name",
						Value:        params.GivenName,
						Autocomplete: "given-name",
						Required:     true,
						Localizer:    l,
					}),
					components.InputField(&components.InputFieldOptions{
						Label: l.MustLocalizeMessage(&i18n.Message{
							ID: "common.users.family_name",
						}),
						Name:         "family_name",
						Value:        params.FamilyName,
						Autocomplete: "family-name",
						Required:     true,
						Localizer:    l,
					}),
					components.SelectComponent(&components.SelectOptions{
						Label: l.MustLocalizeMessage(&i18n.Message{
							ID: "common.users.country",
						}),
						Name:         "country",
						Value:        params.Country,
						Autocomplete: "country",
						Options:      countryOptions,
						Required:     true,
						Localizer:    l,
					}),
					components.InputField(&components.InputFieldOptions{
						Label: l.MustLocalizeMessage(&i18n.Message{
							ID: "common.users.profession",
						}),
						Name:      "profession",
						Value:     params.Profession,
						Localizer: l,
					}),
					components.InputField(&components.InputFieldOptions{
						Label: l.MustLocalizeMessage(&i18n.Message{
							ID: "common.users.licence_number",
						}),
						Name:      "licence_number",
						Value:     params.LicenceNumber,
						Localizer: l,
					}),
					Button(Type("submit"), Class("button mt-2 h-10 w-full text-lg"), Text("Submit")),
				),
			),
		),
	)
}
