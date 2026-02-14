package events

import (
	"fmt"
	"time"

	twmerge "github.com/Oudwins/tailwind-merge-go"
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

func EventAttendanceIcon(going bool, classes ...string) Node {
	icon := "/assets/circle-check-empty.svg#icon"
	if going {
		icon = "/assets/circle-check-solid.svg#icon"
	}

	return SVG(
		Class(twmerge.Merge("h-5 w-5 fill-current", twmerge.Merge(classes...))),
		Attr("viewBox", "0 0 640 640"),
		El("use",
			Href(icon),
		),
	)
}

func EventAttendanceBadge(l *i18n.Localizer) Node {
	messageKey := "common.events.attendance_badge"

	return Span(
		Class("inline-flex items-center justify-center gap-1 rounded-sm border border-primary/20 bg-primary/10 px-2 py-1 text-sm font-semibold text-primary"),
		EventAttendanceIcon(true),
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
				Class("mb-2 flex items-center gap-2 flex-wrap"),
				EventLocationBadge(e.IsVirtual, e.Venue, localizer, ctx.Language),
				If(e.EventRegistration != nil, EventAttendanceBadge(localizer)),

				P(
					Text(helpers.FormatDateRange(e.StartsAt.Time, e.EndsAt.Time, tz, ctx.Language)),
				),
			),

			H3(
				Class("text-primary mobile:text-lg mobile:leading-tight text-2xl font-bold"),
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
				Class("desktop:mb-4 text-gray-600"),
				Text(helpers.TranslateEventType(localizer, e.EventType)),
				Text(", "),
				Text(helpers.FormatHosts(localizer, e.Hosts)),
			),

			Div(
				Class("mobile:grid gap-4 w-full mt-auto flex items-center"),
				If(isFuture && e.EventRegistration == nil, A(
					Href(eventUrl+"/register"),
					Class("button px-6 mobile:w-full"),
					Text(localizer.MustLocalizeMessage(&i18n.Message{
						ID: "common.events.sign_up",
					})),
				)),

				A(
					Href(eventUrl),
					Class("button secondary mobile:w-full"),
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
		Class("desktop:grid leading-tight mobile:row-start-1"),
		Span(Class("desktop:text-xs font-semibold desktop:text-gray-500 desktop:uppercase mobile:after:content-[':_']"), Text(label)),
		Span(Class("text-lg"), Text(priceFormatted)),
	)
}
