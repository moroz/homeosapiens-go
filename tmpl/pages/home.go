package pages

import (
	"fmt"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/events"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

const gold = "#c9921f"

// tr picks the locale-appropriate string for page-specific marketing copy that
// isn't worth a bundle entry.
func tr(ctx *types.CustomContext, en, pl string) string {
	if ctx.IsPolish() {
		return pl
	}
	return en
}

func iff(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

func eyebrow(text string) Node {
	return Div(
		Class("flex items-center gap-3 font-mono text-xs font-medium tracking-[0.22em] uppercase text-primary"),
		Span(Class("inline-block h-px w-6"), Style("background:"+gold)),
		Text(text),
	)
}

func eyebrowCentered(text string) Node {
	return Div(
		Class("flex items-center justify-center gap-3 font-mono text-xs font-medium tracking-[0.22em] uppercase text-primary"),
		Span(Class("inline-block h-px w-6"), Style("background:"+gold)),
		Text(text),
	)
}

// Home renders the public landing page. featured is the next upcoming event (or
// nil), upcoming is the list of further upcoming events to preview.
func Home(ctx *types.CustomContext, featured *services.EventListDto, upcoming []*services.EventListDto) Node {
	return layout.BareLayout(ctx, "Homeo sapiens",
		heroSection(ctx, featured),
		tenetsSection(ctx),
		If(len(upcoming) > 0, eventsSection(ctx, upcoming)),
		librarySection(ctx),
		aboutSection(ctx),
	)
}

func heroCopy(ctx *types.CustomContext, centered bool) Node {
	headlineSize := "text-5xl"
	if centered {
		headlineSize = "text-6xl"
	}
	return Div(
		If(centered, Class("mx-auto max-w-3xl text-center")),
		If(centered, eyebrowCentered("Similia similibus curentur · "+tr(ctx, "since 1993", "od 1993"))),
		If(!centered, eyebrow("Similia similibus curentur · "+tr(ctx, "since 1993", "od 1993"))),
		H2(
			Class("mt-5 font-bold leading-[1.05] tracking-tight text-slate-900 mobile:text-4xl "+headlineSize),
			Text(tr(ctx, "The art of homeopathy", "Sztuka homeopatii")),
			Text(" — "),
			Span(Style("color:"+gold), Text(tr(ctx, "seminars, webinars, recordings.", "seminaria, webinary, nagrania."))),
		),
		P(
			Class("mt-6 max-w-prose text-lg text-slate-600 "+iff(centered, "mx-auto", "")),
			Text(tr(ctx,
				"Live sessions with practising clinicians, and a recording library you keep access to.",
				"Sesje na żywo z praktykującymi klinicystami i biblioteka nagrań z dostępem na stałe.")),
		),
		Div(
			Class("mt-8 flex flex-wrap items-center gap-4 "+iff(centered, "justify-center", "")),
			A(Href("/events"), Class("button px-6"), Text(tr(ctx, "Browse events", "Zobacz wydarzenia"))),
			A(Href("/videos"), Class("font-semibold text-primary hover:text-primary-hover"),
				Text(tr(ctx, "Explore the library →", "Przeglądaj bibliotekę →"))),
		),
	)
}

func sunDecoration() Node {
	glow := "background:radial-gradient(circle at 50% 50%, rgba(201,146,31,0.30), rgba(201,146,31,0) 65%);"
	rays := "background:repeating-conic-gradient(from 0deg at 50% 50%, #c9921f 0deg 1.5deg, transparent 1.5deg 9deg);" +
		"-webkit-mask:radial-gradient(circle at 50% 50%, #000 5%, transparent 60%);" +
		"mask:radial-gradient(circle at 50% 50%, #000 5%, transparent 60%);"
	return Div(
		Class("pointer-events-none absolute -top-40 -right-40 h-[680px] w-[680px] opacity-70 mobile:hidden"),
		Div(Class("absolute inset-0"), Style(glow)),
		Div(Class("absolute inset-0 opacity-50"), Style(rays)),
	)
}

func heroSection(ctx *types.CustomContext, featured *services.EventListDto) Node {
	return Section(
		Class("relative overflow-hidden border-b border-slate-200 bg-white"),
		sunDecoration(),
		Div(
			Class("relative z-10 container mx-auto px-6 py-24 mobile:py-14"),
			Iff(featured != nil, func() Node {
				return Div(
					Class("grid items-center gap-14 lg:grid-cols-[1.1fr_0.9fr]"),
					heroCopy(ctx, false),
					featuredTicket(ctx, featured),
				)
			}),
			If(featured == nil, heroCopy(ctx, true)),
		),
	)
}

func featuredTicket(ctx *types.CustomContext, e *services.EventListDto) Node {
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
		Class("card block no-underline transition-shadow hover:shadow-md"),
		Div(
			Class("flex items-center justify-between"),
			Span(
				Class("inline-flex rounded-sm border px-2 py-1 font-mono text-xs tracking-widest uppercase"),
				Style("color:"+gold+";border-color:"+gold),
				Text(tr(ctx, "Next event", "Najbliższe wydarzenie")),
			),
			Span(Class("font-mono text-sm text-slate-500"),
				Text(helpers.TranslateEventType(l, e.EventType))),
		),
		H3(Class("mt-4 text-2xl font-bold text-primary"), Text(title)),
		Ul(
			Class("mt-5 grid gap-2 text-slate-700"),
			ticketRow(tr(ctx, "Host", "Prowadzący"), helpers.FormatHosts(l, e.Hosts)),
			ticketRow(tr(ctx, "When", "Termin"), helpers.FormatDateRange(e.StartsAt, e.EndsAt, ctx.Timezone, ctx.Language)),
		),
		Div(
			Class("mt-6 flex items-center justify-between gap-4 border-t border-slate-200 pt-5"),
			Div(
				Class("mb-2"),
				EventLocationBadgeForHome(ctx, e),
			),
			Span(Class("text-lg font-semibold text-slate-900"), Text(price)),
		),
	)
}

func ticketRow(k, v string) Node {
	return Li(
		Class("grid grid-cols-[6rem_1fr] items-baseline gap-3"),
		Span(Class("font-mono text-xs tracking-wider uppercase text-slate-500"), Text(k)),
		Span(Text(v)),
	)
}

// EventLocationBadgeForHome reuses the events package badge.
func EventLocationBadgeForHome(ctx *types.CustomContext, e *services.EventListDto) Node {
	return events.EventLocationBadge(e, ctx.Localizer, ctx.Language)
}

type tenet struct{ no, titleEn, titlePl, bodyEn, bodyPl string }

func tenetsSection(ctx *types.CustomContext) Node {
	items := []tenet{
		{"i.", "Live seminars & webinars", "Seminaria i webinary na żywo",
			"Case-taking, remedy differentiation and provings with practising clinicians — on site in Poland and on Zoom.",
			"Wywiady, różnicowanie leków i provingi z praktykującymi klinicystami — na żywo w Polsce i przez Zoom."},
		{"ii.", "Recordings you keep", "Dostęp do nagrań",
			"Talks and full seminars stay in your library, in Polish and English, to watch again.",
			"Wykłady i całe seminaria zostają w bibliotece — po polsku i angielsku — do ponownego obejrzenia."},
		{"iii.", "For practitioners and beginners", "Dla praktyków i początkujących",
			"Doctors, naturopaths and people new to homeopathy. A continuation of the Polish Homeopathic Journal.",
			"Lekarze, naturopaci i osoby zaczynające z homeopatią. Kontynuacja Polskiego Pisma Homeopatycznego."},
	}

	return Section(
		Class("border-b border-slate-200 bg-white"),
		Div(
			Class("container mx-auto px-6 py-16"),
			eyebrow(tr(ctx, "Overview", "Przegląd")),
			H2(Class("mt-4 mb-10 text-3xl font-bold text-slate-900"),
				Text(tr(ctx, "What's on the platform", "Co znajdziesz na platformie"))),
			Div(
				Class("grid gap-px overflow-hidden rounded-lg border border-slate-200 bg-slate-200 desktop:grid-cols-3"),
				Map(items, func(t tenet) Node {
					return Div(
						Class("bg-white p-7"),
						Span(Class("font-mono text-sm"), Style("color:"+gold), Text(t.no)),
						H3(Class("mt-4 mb-2 text-xl font-bold text-primary"), Text(tr(ctx, t.titleEn, t.titlePl))),
						P(Class("text-slate-600"), Text(tr(ctx, t.bodyEn, t.bodyPl))),
					)
				}),
			),
		),
	)
}

func eventsSection(ctx *types.CustomContext, upcoming []*services.EventListDto) Node {
	return Section(
		Class("border-b border-slate-200 bg-white"),
		Div(
			Class("container mx-auto px-6 py-16"),
			Div(
				Class("mb-6 flex items-end justify-between gap-4"),
				Div(
					eyebrow(tr(ctx, "Agenda", "Agenda")),
					H2(Class("mt-4 text-3xl font-bold text-slate-900"),
						Text(tr(ctx, "Upcoming events", "Najbliższe wydarzenia"))),
				),
				A(Href("/events"), Class("button secondary"), Text(tr(ctx, "All events", "Wszystkie wydarzenia"))),
			),
			homeEventList(ctx, upcoming),
		),
	)
}

func eventTypeChip(ctx *types.CustomContext, e *services.EventListDto) Node {
	color := gold
	if e.EventType == queries.EventTypeWebinar {
		color = "#2f6e5e"
	}
	return Span(
		Class("inline-flex rounded-sm border px-2 py-0.5 font-mono text-xs tracking-wider uppercase"),
		Style("color:"+color+";border-color:"+color),
		Text(helpers.TranslateEventType(ctx.Localizer, e.EventType)),
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
		Class("group grid grid-cols-[9rem_1fr_auto] items-center gap-6 border-b border-slate-200 py-6 no-underline transition-colors hover:bg-slate-50 mobile:grid-cols-1 mobile:gap-3 mobile:py-5"),
		Span(
			Class("font-mono text-sm text-slate-500"),
			Text(helpers.FormatDateRange(e.StartsAt, e.EndsAt, ctx.Timezone, ctx.Language)),
		),
		Div(
			H3(
				Class("text-xl font-bold text-primary transition-colors group-hover:text-primary-hover"),
				Text(title),
			),
			Div(
				Class("mt-2 flex flex-wrap items-center gap-x-3 gap-y-1 text-sm text-slate-600"),
				eventTypeChip(ctx, e),
				events.EventLocationBadge(e, l, ctx.Language),
				Span(Text(helpers.FormatHosts(l, e.Hosts))),
				Span(Class("font-semibold text-slate-900"), Text(price)),
			),
		),
		Span(
			Class("justify-self-end font-mono text-xs font-semibold tracking-wider uppercase mobile:justify-self-start"),
			Style("color:"+gold),
			Text(tr(ctx, "Details →", "Szczegóły →")),
		),
	)
}

func homeEventList(ctx *types.CustomContext, rows []*services.EventListDto) Node {
	return Div(
		Class("border-t border-slate-200"),
		Map(rows, func(e *services.EventListDto) Node {
			return homeEventRow(ctx, e)
		}),
	)
}

type libItem struct{ titleEn, titlePl, metaEn, metaPl, accent string }

func librarySection(ctx *types.CustomContext) Node {
	items := []libItem{
		{"Why the Organon is still important", "Dlaczego Organon wciąż jest ważny",
			"Conversation · 41 min", "Rozmowa · 41 min", "#2f6e5e"},
		{"Peripheral circulation & Secale cornutum", "Krążenie obwodowe i Secale cornutum",
			"Lecture · 58 min", "Wykład · 58 min", gold},
		{"Dutiful remedies: differential diagnosis", "Leki o wygórowanym poczuciu obowiązku",
			"Webinar · 22.03.2025", "Webinar · 22.03.2025", "#a33a22"},
	}

	return Section(
		Class("border-b border-slate-200 bg-white"),
		Div(
			Class("container mx-auto px-6 py-16"),
			Div(
				Class("mb-8 flex items-end justify-between gap-4"),
				Div(
					eyebrow(tr(ctx, "From the library", "Z biblioteki")),
					H2(Class("mt-4 text-3xl font-bold text-slate-900"),
						Text(tr(ctx, "Watch now", "Obejrzyj teraz"))),
				),
				A(Href("/videos"), Class("button secondary"), Text(tr(ctx, "Browse all", "Przeglądaj wszystko"))),
			),
			Div(
				Class("grid gap-6 desktop:grid-cols-3"),
				Map(items, func(it libItem) Node {
					return A(
						Href("/videos"),
						Class("video-card block no-underline transition-shadow hover:shadow-md"),
						Div(
							Class("video-card-thumb"),
							Style("background:linear-gradient(150deg,"+it.accent+",#1b1c2e)"),
							Span(Class("px-4 text-lg font-bold text-white"), Text(tr(ctx, it.titleEn, it.titlePl))),
						),
						Div(
							Class("video-card-body"),
							Span(Class("font-mono text-xs tracking-wider uppercase text-primary"),
								Text(tr(ctx, it.metaEn, it.metaPl))),
						),
					)
				}),
			),
		),
	)
}

type remedy struct{ name, srcEn, srcPl string }

func aboutSection(ctx *types.CustomContext) Node {
	remedies := []remedy{
		{"Aurum met.", "Dutiful remedies", "Leki obowiązku"},
		{"Nux vomica", "Dutiful remedies", "Leki obowiązku"},
		{"Silicea", "Dutiful remedies", "Leki obowiązku"},
		{"Secale cornutum", "Peripheral circulation", "Krążenie obwodowe"},
		{"Spongia tosta", "German New Medicine", "Germańska Nowa Medycyna"},
	}

	return Section(
		Class("bg-primary text-white"),
		Div(
			Class("container mx-auto grid items-center gap-14 px-6 py-20 lg:grid-cols-2"),
			Div(
				Div(
					Class("flex items-center gap-3 font-mono text-xs font-medium tracking-[0.22em] uppercase"),
					Style("color:"+gold),
					Span(Class("inline-block h-px w-6"), Style("background:"+gold)),
					Text(tr(ctx, "About us", "O nas")),
				),
				H2(Class("mt-4 text-4xl font-bold leading-tight mobile:text-3xl"),
					Text(tr(ctx,
						"Homeo sapiens was a Polish homeopathy publisher and journal in the 1990s.",
						"Homeo sapiens to polski wydawca i pismo homeopatyczne z lat 90."))),
				P(Class("mt-5 text-slate-200"),
					Text(tr(ctx,
						"Founded in 1993, it published through the nineties, then paused. We've reestablished it as a platform for exchange — open to doctors, naturopaths and professionals practising homeopathy, and to those who don't yet. We share videos, book reviews and a seminar agenda, and keep a Zoom space open for meetings.",
						"Założone w 1993 roku, ukazywało się przez lata dziewięćdziesiąte, potem zamilkło. Reaktywowaliśmy je jako platformę wymiany — otwartą dla lekarzy, naturopatów i profesjonalistów praktykujących homeopatię, a także dla tych, którzy jeszcze nie praktykują. Udostępniamy filmy, recenzje książek i agendę seminariów, a przestrzeń Zoom pozostaje otwarta na spotkania."))),
			),
			Div(
				Class("overflow-hidden rounded-lg border border-white/15"),
				Div(
					Class("border-b border-white/15 px-4 py-3 font-mono text-xs tracking-widest uppercase"),
					Style("color:"+gold),
					Text(tr(ctx, "Remedies covered in the library", "Leki omawiane w bibliotece")),
				),
				Div(
					Class("grid grid-cols-3"),
					Map(remedies, func(r remedy) Node {
						return Div(
							Class("border-r border-b border-white/15 px-3 py-5 text-center"),
							Div(Class("text-lg italic"), Text(r.name)),
							Div(Class("mt-1 font-mono text-xs text-slate-300"), Text(tr(ctx, r.srcEn, r.srcPl))),
						)
					}),
					A(
						Href("/videos"),
						Class("flex items-center justify-center border-b border-white/15 px-3 py-5 text-center text-sm font-semibold no-underline"),
						Style("color:"+gold),
						Text(tr(ctx, "More in the library →", "Więcej w bibliotece →")),
					),
				),
			),
		),
	)
}
