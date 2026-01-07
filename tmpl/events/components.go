package events

import (
	"context"
	"fmt"
	"time"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func EventLocationBadge(isVirtual bool, venue *queries.Venue, l *i18n.Localizer, lang string) Node {
	return Span(
		Class("font-fallback inline-flex items-center gap-1 justify-self-start rounded border border-blue-900/50 bg-blue-100 px-2 py-1 text-sm font-semibold text-blue-900"),
		Iff(venue != nil, func() Node {
			city := venue.CityEn
			if lang == "pl" && venue.CityPl != nil {
				city = *venue.CityPl
			}

			return Text(
				fmt.Sprintf("%s, %s", city, helpers.TranslateCountry(l, venue.CountryCode)),
			)
		}),
		If(isVirtual && venue != nil, Text(" + ")),
		If(isVirtual, Text("Online")),
	)
}

func EventAttendanceBadge(isFuture bool, l *i18n.Localizer) Node {
	messageKey := "common.events.attendance_badge.past"
	if isFuture {
		messageKey = "common.events.attendance_badge.upcoming"
	}

	return Span(
		Class("font-fallback inline-flex items-center justify-center gap-1 rounded border border-green-900/50 bg-green-100 px-2 py-1 text-sm font-semibold text-green-900"),
		I(Class("h-4 w-4"), Data("lucide", "user-star")),
		Text(l.MustLocalizeMessage(&i18n.Message{
			ID: messageKey,
		})),
	)
}

func HostCard(localizer *i18n.Localizer, host *queries.ListHostsForEventsRow) Node {
	salutation := helpers.TranslateSalutation(localizer, host.Salutation)

	return Div(
		Class("bg-primary flex h-min w-42 flex-col items-center overflow-hidden rounded-lg border-2"),
		Div(
			Class("relative aspect-square w-full overflow-hidden"),
			Iff(host.ProfilePictureUrl != nil, func() Node {
				url := fmt.Sprintf("%s/%s", config.AssetCdnBaseUrl, *host.ProfilePictureUrl)

				return Img(
					Src(url),
					Class("absolute inset-0 h-full w-full object-cover"),
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

func EventCard(ctx context.Context, e *services.EventListDto) Node {
	localizer := ctx.Value("localizer").(*i18n.Localizer)
	lang := ctx.Value(config.LangContextName).(string)

	title := e.TitleEn
	if lang == "pl" {
		title = e.TitlePl
	}

	eventUrl := fmt.Sprintf("/events/%s", e.Slug)
	tz := ctx.Value("timezone").(*time.Location)

	isFuture := e.StartsAt.Time.After(time.Now())

	return Article(
		Class("card flex justify-between gap-6"),
		Header(
			Class("flex flex-1 flex-col items-start"),
			Div(
				Class("mb-2 flex items-center gap-2"),
				EventLocationBadge(e.IsVirtual, e.Venue, localizer, lang),
				If(e.EventRegistration != nil, EventAttendanceBadge(isFuture, localizer)),
			),

			H3(
				Class("text-primary mobile:text-xl mobile:leading-tight text-3xl font-bold"),
				A(
					Class("hover:text-primary-hover decoration-2 underline-offset-3 transition-colors hover:underline"),
					Href(eventUrl),
					Text(title),
				),
			),

			P(
				Text(helpers.TranslateEventType(localizer, e.EventType)),
				Text(", "),
				Text(helpers.FormatDateRange(e.StartsAt.Time, e.EndsAt.Time, tz, lang)),
			),

			P(
				Class("mb-4"),
				Raw(
					localizer.MustLocalize(&i18n.LocalizeConfig{
						MessageID: "common.events.hosts",
						TemplateData: map[string]string{
							"Hosts": helpers.FormatHosts(localizer, e.Hosts),
						},
					}),
				),
			),

			Div(
				If(isFuture, A(
					Href(eventUrl+"/register"),
					Class("button secondary px-6"),
					Text(localizer.MustLocalizeMessage(&i18n.Message{
						ID: "common.events.sign_up",
					})),
				)),

				Class("mt-auto flex items-center gap-4"),
				A(
					Href(eventUrl),
					Class("button font-fallback"),
					Text(localizer.MustLocalizeMessage(&i18n.Message{
						ID:    "common.events.learn_more",
						Other: "Learn more…",
					})),
				),

				formatEventPrice(ctx, e),
			),
		),

		Div(Class("mobile:hidden flex items-center gap-6"),
			Map(e.Hosts, func(host *queries.ListHostsForEventsRow) Node {
				return HostCard(localizer, host)
			}),
		),
	)
}

func formatEventPrice(ctx context.Context, e *services.EventListDto) Node {
	localizer := ctx.Value("localizer").(*i18n.Localizer)
	lang := ctx.Value(config.LangContextName).(string)

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
