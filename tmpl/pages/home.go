package pages

import (
	"fmt"

	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/events"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// tr picks the locale-appropriate string for page-specific marketing copy that
// isn't worth a bundle entry.
func tr(ctx *types.CustomContext, en, pl string) string {
	if ctx.IsPolish() {
		return pl
	}
	return en
}

// eyebrow is a small caps label preceded by a short rule in the brand gradient.
func eyebrow(text string) Node {
	return Div(
		Class("text-primary flex items-center gap-3 text-xs font-semibold tracking-[0.22em] uppercase"),
		Span(Class("brand-rule inline-block h-0.5 w-6 rounded-full")),
		Text(text),
	)
}

// sectionHeading is the shared heading pair used by every section below the hero.
func sectionHeading(label, title string, trailing ...Node) Node {
	return Div(
		Class("mb-8 flex items-end justify-between gap-4"),
		Div(
			eyebrow(label),
			H2(Class("mt-3 text-2xl font-bold text-slate-900 lg:text-3xl"), Text(title)),
		),
		Group(trailing),
	)
}

// Home renders the public landing page. It is deliberately text-only: every
// element on it is either real copy or real data. upcoming is the next few
// published events, soonest first, and may be empty.
func Home(ctx *types.CustomContext, upcoming []*services.EventListDto) Node {
	return layout.BareLayout(ctx, "Homeo sapiens",
		heroSection(ctx),
		If(len(upcoming) > 0, eventsSection(ctx, upcoming)),
		aboutSection(ctx),
	)
}

func heroSection(ctx *types.CustomContext) Node {
	return Section(
		Class("border-b border-slate-200 bg-white"),
		Div(
			Class("mobile:py-14 container mx-auto px-6 py-24"),
			Div(
				Class("max-w-3xl"),
				eyebrow("Similia similibus curentur · Pecunia non olet"),
				H2(
					Class("mobile:text-4xl mt-5 text-5xl leading-[1.05] font-bold tracking-tight text-slate-900"),
					Text("Homeo sapiens"),
				),
				P(
					Class("mobile:text-xl mt-5 max-w-prose text-2xl leading-snug text-slate-700"),
					Text(tr(ctx,
						"Seminars, webinars and recordings on the art of homeopathy.",
						"Seminaria, webinary i nagrania o sztuce homeopatii.")),
				),
				P(
					Class("mt-5 max-w-prose text-lg text-slate-600"),
					Text(tr(ctx,
						"Live sessions with practising clinicians, and a recording library you keep access to.",
						"Sesje na żywo z praktykującymi klinicystami i biblioteka nagrań z dostępem na stałe.")),
				),
				Div(
					Class("mt-8 flex flex-wrap items-center gap-x-6 gap-y-4"),
					A(Href("/events"), Class("button px-6"), Text(tr(ctx, "Browse events", "Zobacz wydarzenia"))),
					A(Href("/watch"), Class("text-primary hover:text-primary-hover font-semibold"),
						Text(tr(ctx, "Explore the library →", "Przeglądaj bibliotekę →"))),
				),
			),
		),
	)
}

func eventsSection(ctx *types.CustomContext, upcoming []*services.EventListDto) Node {
	return Section(
		Class("border-b border-slate-200 bg-white"),
		Div(
			Class("container mx-auto px-6 py-16"),
			sectionHeading(
				tr(ctx, "Agenda", "Agenda"),
				tr(ctx, "Upcoming events", "Najbliższe wydarzenia"),
				A(Href("/events"), Class("button secondary"), Text(tr(ctx, "All events", "Wszystkie wydarzenia"))),
			),
			Div(
				Class("border-t border-slate-200"),
				Map(upcoming, func(e *services.EventListDto) Node {
					return homeEventRow(ctx, e)
				}),
			),
		),
	)
}

func homeEventRow(ctx *types.CustomContext, e *services.EventListDto) Node {
	l := ctx.Localizer
	title := e.TitleEn
	if ctx.IsPolish() {
		title = e.TitlePl
	}

	var price string
	if e.BasePriceAmount == nil {
		price = l.MustLocalizeMessage(&i18n.Message{ID: "common.events.free"})
	} else {
		price = helpers.FormatPrice(*e.BasePriceAmount, *e.BasePriceCurrency, ctx.Language)
	}

	return A(
		Href(fmt.Sprintf("/events/%s", e.Slug)),
		Class("group mobile:grid-cols-1 mobile:gap-3 mobile:py-5 hover:bg-brand-50 grid grid-cols-[9rem_1fr_auto] items-center gap-6 border-b border-slate-200 py-6 no-underline transition-colors"),
		Span(
			Class("text-sm text-slate-500 tabular-nums"),
			Text(helpers.FormatDateRange(e.StartsAt, e.EndsAt, ctx.Timezone, ctx.Language)),
		),
		Div(
			H3(
				Class("text-primary group-hover:text-primary-hover text-xl font-bold transition-colors"),
				Text(title),
			),
			Div(
				Class("mt-2 flex flex-wrap items-center gap-x-3 gap-y-1 text-sm text-slate-600"),
				eventTypeChip(ctx, e),
				events.EventLocationBadge(e, l, ctx.Language),
				Span(Text(helpers.FormatHosts(l, ctx.Language, e.Hosts))),
				Span(Class("font-semibold text-slate-900"), Text(price)),
			),
		),
		Span(
			Class("text-primary mobile:justify-self-start justify-self-end text-xs font-semibold tracking-wider uppercase"),
			Text(tr(ctx, "Details →", "Szczegóły →")),
		),
	)
}

func eventTypeChip(ctx *types.CustomContext, e *services.EventListDto) Node {
	return Span(
		Class("border-brand-300 text-primary inline-flex rounded-sm border px-2 py-0.5 text-xs tracking-wider uppercase"),
		Text(helpers.TranslateEventType(ctx.Localizer, e.EventType)),
	)
}

func aboutSection(ctx *types.CustomContext) Node {
	return Section(
		Class("bg-white"),
		Div(
			Class("container mx-auto px-6 py-16"),
			sectionHeading(tr(ctx, "About us", "O nas"), tr(ctx, "A journal, reestablished", "Pismo reaktywowane")),
			Div(
				Class("max-w-prose space-y-4 text-lg text-slate-700"),
				P(Text(tr(ctx,
					"Homeo sapiens was a Polish homeopathy publisher and journal in the 1990s. Founded in 1993, it published through the nineties, then paused.",
					"Homeo sapiens to polski wydawca i pismo homeopatyczne z lat 90. Założone w 1993 roku, ukazywało się przez lata dziewięćdziesiąte, potem zamilkło."))),
				P(Text(tr(ctx,
					"We've reestablished it as a platform for exchange — open to doctors, naturopaths and professionals practising homeopathy, and to those who don't yet. We share videos, book reviews and a seminar agenda, and keep a Zoom space open for meetings.",
					"Reaktywowaliśmy je jako platformę wymiany — otwartą dla lekarzy, naturopatów i profesjonalistów praktykujących homeopatię, a także dla tych, którzy jeszcze nie praktykują. Udostępniamy filmy, recenzje książek i agendę seminariów, a przestrzeń Zoom pozostaje otwarta na spotkania."))),
			),
			Div(
				Class("mt-8 flex flex-wrap items-center gap-x-6 gap-y-3"),
				A(Href("/blog"), Class("text-primary hover:text-primary-hover font-semibold"),
					Text(tr(ctx, "Read the blog →", "Czytaj blog →"))),
				A(Href("/watch"), Class("text-primary hover:text-primary-hover font-semibold"),
					Text(tr(ctx, "Watch recordings →", "Oglądaj nagrania →"))),
			),
		),
	)
}
