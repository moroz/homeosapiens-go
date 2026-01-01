package eventregistrations

import (
	"context"

	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"

	. "maragu.dev/gomponents/html"
)

func New(ctx context.Context, event *services.EventDetailsDto, params *types.CreateEventRegistrationParams) Node {
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
			Main(
				Form(
					Class("mt-6 space-y-2"),
					components.InputField(&components.InputFieldOptions{
						Label:        "Given name:",
						Name:         "given_name",
						Value:        params.GivenName,
						Autocomplete: "given-name",
						Required:     true,
					}),
					components.InputField(&components.InputFieldOptions{
						Label:        "Family name:",
						Name:         "family_name",
						Value:        params.FamilyName,
						Autocomplete: "family-name",
						Required:     true,
					}),
				),
			),
		),
	)
}
