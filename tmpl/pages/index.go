package pages

import (
	"context"
	"fmt"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Index(ctx context.Context, events []*services.EventListDto) Node {
	return layout.Layout(ctx, "Events", Div(
		Class("container mx-auto grid gap-4"),
		Map(events, func(e *services.EventListDto) Node {
			return eventCard(ctx, e)
		}),
	))
}

func hostCard(localizer *i18n.Localizer, host *queries.ListHostsForEventsRow) Node {
	salutation := helpers.TranslateSalutation(localizer, host.Salutation)

	return Div(
		Class("bg-primary flex h-min w-35 flex-col items-center overflow-hidden rounded-lg border-2"),
		Div(
			Class("aspect-square w-full overflow-hidden"),
			Iff(host.ProfilePictureUrl != nil, func() Node {
				url := fmt.Sprintf("%s/%s", config.AssetCdnBaseUrl, *host.ProfilePictureUrl)

				return Img(
					Src(url),
					Class("scale-120"),
					Alt(fmt.Sprintf("Profile picture of %s%s %s", salutation, host.GivenName, host.FamilyName)),
				)
			}),
		),
		Footer(
			Class("flex h-10 w-full items-center justify-center text-center text-white"),
			Span(
				Text(salutation),
				Strong(
					Text(host.GivenName+" "+host.FamilyName),
				),
			),
		),
	)
}

func eventCard(ctx context.Context, e *services.EventListDto) Node {
	localizer := ctx.Value("localizer").(*i18n.Localizer)
	lang := ctx.Value("lang").(string)

	title := e.TitleEn
	if lang == "pl" {
		title = e.TitlePl
	}

	eventUrl := fmt.Sprintf("/events/%s", e.Slug)

	return Article(
		Class("card flex justify-between gap-6"),
		Header(
			Class("flex flex-1 flex-col items-start"),
			Span(
				Class("text-primary mb-2 inline-flex items-center gap-1 justify-self-start rounded border-2 border-black bg-white px-2 py-1 text-sm font-semibold"),
				Iff(e.VenueID.Valid, func() Node {
					city := *e.VenueCityEn
					if lang == "pl" && e.VenueCityPl != nil {
						city = *e.VenueCityPl
					}

					return Text(
						fmt.Sprintf("%s, %s", city, helpers.TranslateCountry(localizer, *e.VenueCountryCode)),
					)
				}),
				If(e.IsVirtual && e.VenueID.Valid, Text(" + ")),
				If(e.IsVirtual, Text("Online")),
			),

			H3(
				Class("text-primary text-4xl font-bold"),
				A(
					Class("hover:text-primary-hover decoration-2 underline-offset-3 transition-colors hover:underline"),
					Href(eventUrl),
					Text(title),
				),
			),

			P(
				Text(helpers.TranslateEventType(localizer, e.EventType)),
				Text(", "),
				Text(helpers.FormatDateRange(e.StartsAt.Time, e.EndsAt.Time, "Europe/Warsaw", lang)),
			),

			P(
				Raw(
					localizer.MustLocalize(&i18n.LocalizeConfig{
						DefaultMessage: &i18n.Message{
							ID:    "common.events.hosts",
							Other: "Hosts: {{ .Hosts }}.",
						},
						TemplateData: map[string]string{
							"Hosts": helpers.FormatHosts(localizer, e.Hosts),
						},
					}),
				),
			),

			Div(
				Class("mt-6 flex items-center gap-4"),
				A(
					Href(eventUrl),
					Class("button"),
					Text(localizer.MustLocalizeMessage(&i18n.Message{
						ID:    "common.events.learn_more",
						Other: "Learn more…",
					})),
				),

				formatEventPrice(ctx, e),
			),
		),

		Div(Class("flex items-center gap-6"),
			Map(e.Hosts, func(host *queries.ListHostsForEventsRow) Node {
				return hostCard(localizer, host)
			}),
		),
	)
}

func formatEventPrice(ctx context.Context, e *services.EventListDto) Node {
	localizer := ctx.Value("localizer").(*i18n.Localizer)
	lang := ctx.Value("lang").(string)

	label := localizer.MustLocalizeMessage(&i18n.Message{
		ID:    "common.events.participation_cost",
		Other: "Participation cost:",
	})

	var priceFormatted string
	if e.BasePriceAmount == nil {
		priceFormatted = localizer.MustLocalizeMessage(&i18n.Message{
			ID:    "common.events.free",
			Other: "Free",
		})
	} else {
		priceFormatted = helpers.FormatPrice(*e.BasePriceAmount, *e.BasePriceCurrency, lang)
	}

	return Div(
		Class("grid leading-tight"),
		Span(Class("font-semibold"), Text(label)),
		Span(Class("text-lg"), Text(priceFormatted)),
	)
}
