package pages

import (
	"context"
	"fmt"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

const AssetCdnBaseUrl = "https://d3n1g0yg3ja4p3.cloudfront.net"

func hostCard(localizer *i18n.Localizer, host *queries.ListHostsForEventsRow) Node {
	salutation := helpers.TranslateSalutation(localizer, host.Salutation)

	return Div(
		Class("flex w-35 flex-col items-center overflow-hidden rounded-lg border-2 bg-primary"),
		Div(
			Class("aspect-square overflow-hidden w-full"),
			Iff(host.ProfilePictureUrl != nil, func() Node {
				url := fmt.Sprintf("%s/%s", AssetCdnBaseUrl, *host.ProfilePictureUrl)

				return Img(
					Src(url),
					Class("scale-120"),
					Alt(fmt.Sprintf("Profile picture of %s%s %s", salutation, host.GivenName, host.FamilyName)),
				)
			}),
		),
		Footer(
			Class("flex h-10 items-center justify-center text-center text-white"),
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

	return Article(
		Class("flex justify-between rounded-lg border-2 bg-slate-100 p-6"),
		Header(
			Span(
				Class("mb-2 inline-flex items-center gap-1 rounded border-2 border-black bg-white px-2 py-1 text-sm font-semibold text-primary"),
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
				Class("text-4xl font-bold text-primary"),
				Text(title),
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

			Iff(e.BasePriceAmount != nil, func() Node {
				return Text(helpers.FormatPrice(*e.BasePriceAmount, *e.BasePriceCurrency, lang))
			}),

			If(e.BasePriceAmount == nil, Text(
				localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "common.events.free",
						Other: "Free",
					},
				}),
			)),

			Ul(
				Map(e.Prices, func(price *queries.EventPrice) Node {
					return Li(Text(
						helpers.FormatPrice(price.PriceAmount, price.PriceCurrency, lang),
					))
				}),
			),
		),

		Map(e.Hosts, func(host *queries.ListHostsForEventsRow) Node {
			return hostCard(localizer, host)
		}),
	)
}

func Index(ctx context.Context, events []*services.EventListDto) Node {
	return layout.Layout(ctx, "Events", Div(
		Class("grid gap-4 container mx-auto mt-4"),
		Map(events, func(e *services.EventListDto) Node {
			return eventCard(ctx, e)
		}),
	))
}
