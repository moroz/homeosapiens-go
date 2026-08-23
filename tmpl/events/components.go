package events

import (
	"fmt"
	"time"

	twmerge "github.com/Oudwins/tailwind-merge-go"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/internal/countries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func EventLocationBadge(e *services.EventListDto, l *i18n.Localizer, lang string) Node {
	return Span(
		Class("inline-flex items-center gap-1 justify-self-start rounded-sm border border-gray-700/10 bg-gray-100 px-2 py-1 text-sm font-semibold text-gray-700"),
		Iff(e.VenueCityEn != nil, func() Node {
			city := *e.VenueCityEn
			if lang == "pl" && e.VenueCityPl != nil {
				city = *e.VenueCityPl
			}

			return Text(
				fmt.Sprintf("%s, %s", city, countries.CountryDisplayName(*e.VenueCountryCode, lang)),
			)
		}),
		If(e.IsVirtual && e.VenueCityEn != nil, Text(" + ")),
		If(e.IsVirtual, Text("Online")),
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

// HostCard follows the portrait treatment from the event banners: a large
// rounded portrait with the name on a brand-coloured bar beneath it. The host
// is the draw for an event, so it stays prominent — dropping the card chrome
// for the flat pages means losing the shadow and the slate box, not the size.
func HostCard(localizer *i18n.Localizer, host *queries.ListHostsForEventsRow) Node {
	salutation := helpers.TranslateSalutation(localizer, host.Salutation)

	return Div(
		// The frame is the same brand-700 as the name plate, so portrait and
		// name read as one object the way they do on the event banners.
		Class("w-42 overflow-hidden rounded-xl border border-brand-700 mobile:w-32"),
		Div(
			Class("relative aspect-square w-full bg-brand-50"),
			Iff(host.ProfilePictureUrl != nil, func() Node {
				url := fmt.Sprintf("%s/%s", config.AssetCdnBaseUrl, *host.ProfilePictureUrl)

				return Img(
					Src(url),
					Class("absolute inset-0 h-full w-full object-cover"),
					Alt(fmt.Sprintf("Profile picture of %s%s %s", salutation, host.GivenName, host.FamilyName)),
				)
			}),
		),
		Div(
			// brand-700, not the banner's brighter orange: this text is 14px, and
			// brand-600 only reaches 3.9:1 against white.
			Class("bg-brand-700 px-3 py-2 text-center text-sm leading-tight text-white"),
			Text(salutation),
			Strong(Class("font-bold"), Text(host.GivenName+" "+host.FamilyName)),
		),
	)
}

// EventTypeChip labels an event as a seminar or a webinar.
func EventTypeChip(ctx *types.CustomContext, e *services.EventListDto) Node {
	return Span(
		Class("inline-flex rounded-sm border border-brand-300 px-2 py-0.5 text-xs tracking-wider text-primary uppercase"),
		Text(helpers.TranslateEventType(ctx.Localizer, e.EventType)),
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

	isFree := e.BasePriceAmount == nil
	// Registration stays open until the event is over, which is what the sign-up
	// query enforces (GetRegisterableFreeEventById).
	hasEnded := e.EndsAt.Before(time.Now())

	return Article(
		Class("flex justify-between gap-8 border-b border-slate-200 py-8 mobile:gap-6 mobile:py-6 last:border-b-0"),
		Div(Class("flex shrink-0 items-start gap-6 mobile:hidden"),
			Map(e.Hosts, func(host *queries.ListHostsForEventsRow) Node {
				return HostCard(localizer, host)
			}),
		),
		Header(
			Class("flex flex-1 flex-col items-start"),
			Div(
				Class("mb-3 flex flex-wrap items-center gap-x-3 gap-y-2 text-sm text-slate-600"),
				Span(
					Class("tabular-nums"),
					Text(helpers.FormatDateRange(e.StartsAt, e.EndsAt, tz, ctx.Language)),
				),
				EventTypeChip(ctx, e),
				EventLocationBadge(e, localizer, ctx.Language),
				If(e.EventRegistration != nil, EventAttendanceBadge(localizer)),
			),

			H3(
				Class("text-2xl font-bold text-primary mobile:text-lg mobile:leading-tight"),
				A(
					Class("decoration-2 underline-offset-3 transition-colors hover:text-primary-hover hover:underline"),
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
				Class("text-gray-600 desktop:mb-4"),
				Text(helpers.FormatHosts(localizer, ctx.Language, e.Hosts)),
			),

			Div(
				Class("mt-auto flex w-full items-center gap-4 mobile:grid"),
				If(!hasEnded && isFree && e.EventRegistration == nil, A(
					Href(eventUrl+"/register"),
					Class("button px-6 mobile:w-full"),
					Text(localizer.MustLocalizeMessage(&i18n.Message{
						ID: "common.events.sign_up",
					})),
				)),
				If(!hasEnded && !isFree && e.EventRegistration == nil,
					components.AddToCartButton(localizer, e.ListPublishedEventsRow.ID, e.CountInCart),
				),

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
		Class("leading-tight desktop:grid mobile:row-start-1"),
		Span(Class("font-semibold desktop:text-xs desktop:text-gray-500 desktop:uppercase mobile:after:content-[':_']"), Text(label)),
		Span(Class("text-lg"), Text(priceFormatted)),
	)
}
