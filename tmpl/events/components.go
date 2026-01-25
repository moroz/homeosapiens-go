package events

import (
	"fmt"
	"time"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func EventLocationBadge(isVirtual bool, venue *queries.Venue, l *i18n.Localizer, lang string) Node {
	return Span(
		Class("inline-flex items-center gap-1 justify-self-start rounded-sm border border-gray-700/10 bg-gray-100 px-2 py-1 text-sm font-semibold text-gray-700"),
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
		Class("inline-flex items-center justify-center gap-1 rounded-sm border border-green-900/10 bg-green-100 px-2 py-1 text-sm font-semibold text-green-900"),
		I(Class("h-4 w-4"), Data("lucide", "user-star")),
		Text(l.MustLocalizeMessage(&i18n.Message{
			ID: messageKey,
		})),
	)
}

func HostCard(localizer *i18n.Localizer, host *queries.ListHostsForEventsRow) Node {
	salutation := helpers.TranslateSalutation(localizer, host.Salutation)

	return Div(
		Class("flex h-min w-42 flex-col items-center rounded-sm bg-slate-100 shadow"),
		Div(
			Class("relative aspect-square w-full overflow-hidden rounded-t-sm"),
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
			Class("text-primary border-primary/20 flex h-10 w-full items-center justify-center rounded-b-sm border border-t-0 text-center text-sm"),
			Span(
				Text(salutation),
				Strong(
					Text(host.GivenName+" "+host.FamilyName),
				),
			),
		),
	)
}

func EventCard(ctx *types.CustomContext, e *services.EventListDto) Node {
	localizer := ctx.Localizer

	title := e.TitleEn
	subtitle := e.SubtitleEn
	if ctx.IsPolish() {
		title = e.TitlePl
		subtitle = e.SubtitlePl
	}

	eventUrl := fmt.Sprintf("/events/%s", e.Slug)
	tz := ctx.Timezone

	isFuture := e.StartsAt.Time.After(time.Now())

	return Article(
		Class("card flex justify-between gap-6"),
		Header(
			Class("flex flex-1 flex-col items-start"),
			Div(
				Class("mb-2 flex items-center gap-2"),
				EventLocationBadge(e.IsVirtual, e.Venue, localizer, ctx.Language),
				If(e.EventRegistration != nil, EventAttendanceBadge(isFuture, localizer)),

				P(
					Text(helpers.FormatDateRange(e.StartsAt.Time, e.EndsAt.Time, tz, ctx.Language)),
				),
			),

			H3(
				Class("text-primary mobile:text-xl mobile:leading-tight text-2xl font-bold"),
				A(
					Class("hover:text-primary-hover decoration-2 underline-offset-3 transition-colors hover:underline"),
					Href(eventUrl),
					Text(title),
				),
			),

			Iff(
				subtitle != nil && *subtitle != "",
				func() Node {
					return H4(
						Class("mb-1 font-semibold"),
						Text(*subtitle),
					)
				},
			),

			P(
				Class("mb-4 text-gray-600"),
				Text(helpers.TranslateEventType(localizer, e.EventType)),
				Text(", "),
				Text(helpers.FormatHosts(localizer, e.Hosts)),
			),

			Div(
				If(isFuture, A(
					Href(eventUrl+"/register"),
					Class("button px-6"),
					Text(localizer.MustLocalizeMessage(&i18n.Message{
						ID: "common.events.sign_up",
					})),
				)),

				Class("mt-auto flex items-center gap-4"),
				A(
					Href(eventUrl),
					Class("button secondary font-fallback"),
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

func formatEventPrice(ctx *types.CustomContext, e *services.EventListDto) Node {
	l := ctx.Localizer

	label := l.MustLocalizeMessage(&i18n.Message{
		ID: "common.events.participation_cost",
	})

	var priceFormatted string
	if e.BasePriceAmount == nil {
		priceFormatted = l.MustLocalizeMessage(&i18n.Message{
			ID: "common.events.free",
		})
	} else {
		priceFormatted = helpers.FormatPrice(*e.BasePriceAmount, *e.BasePriceCurrency, ctx.Language)
	}

	return Div(
		Class("grid leading-tight"),
		Span(Class("text-xs font-semibold text-gray-500 uppercase"), Text(label)),
		Span(Class("text-lg"), Text(priceFormatted)),
	)
}
