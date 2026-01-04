package eventregistrations

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
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

func New(ctx context.Context, event *services.EventDetailsDto, params *types.CreateEventRegistrationParams, validationErrors validation.Errors) Node {
	lang := ctx.Value("lang").(string)
	title := event.TitleEn
	if lang == "pl" {
		title = event.TitlePl
	}
	l := ctx.Value("localizer").(*i18n.Localizer)

	pageTitle := l.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "event_registrations.new.title",
		TemplateData: map[string]string{
			"Title": title,
		},
	})

	user := ctx.Value(config.CurrentUserContextName).(*queries.User)
	currentPath := fmt.Sprintf("/events/%s/register", event.Slug)
	countryOptions := buildCountryOptions(lang)

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
					Class("text-primary text-4xl font-bold"),
					Text(title),
				),
			),

			Iff(user == nil, func() Node {
				return Section(
					components.GoogleButton(l.MustLocalizeMessage(&i18n.Message{
						ID: "event_registrations.new.continue_with_google",
					}), currentPath),
				)
			}),

			Main(
				Form(
					Class("mt-6 space-y-2"),
					Input(Type("hidden"), Name("event_id"), Value(event.Event.ID.String())),
					components.InputField(&components.InputFieldOptions{
						Label: l.MustLocalizeMessage(&i18n.Message{
							ID: "event_registrations.new.form.labels.email",
						}),
						Name:         "email",
						Value:        params.Email,
						Autocomplete: "email",
						Required:     true,
						Localizer:    l,
					}),
					components.InputField(&components.InputFieldOptions{
						Label: l.MustLocalizeMessage(&i18n.Message{
							ID: "event_registrations.new.form.labels.given_name",
						}),
						Name:         "given_name",
						Value:        params.GivenName,
						Autocomplete: "given-name",
						Required:     true,
						Localizer:    l,
					}),
					components.InputField(&components.InputFieldOptions{
						Label: l.MustLocalizeMessage(&i18n.Message{
							ID: "event_registrations.new.form.labels.family_name",
						}),
						Name:         "family_name",
						Value:        params.FamilyName,
						Autocomplete: "family-name",
						Required:     true,
						Localizer:    l,
					}),
					components.SelectComponent(&components.SelectOptions{
						Label: l.MustLocalizeMessage(&i18n.Message{
							ID: "event_registrations.new.form.labels.country",
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
							ID: "event_registrations.new.form.labels.profession",
						}),
						Name:      "profession",
						Value:     params.Profession,
						Localizer: l,
					}),
				),
			),
		),
	)
}
