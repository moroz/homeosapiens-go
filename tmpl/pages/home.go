package pages

import (
	"fmt"

	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/events"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// Home renders the public landing page. It is deliberately text-only: every
// element on it is either real copy or real data. upcoming is the next few
// published events, soonest first, and may be empty.
func Home(ctx *types.CustomContext, upcoming []*services.EventListDto) Node {
	return layout.PageLayout(ctx, "Homeo sapiens",
		heroSection(ctx),
		If(len(upcoming) > 0, eventsSection(ctx, upcoming)),
		aboutSection(ctx),
	)
}

func heroSection(ctx *types.CustomContext) Node {
	l := ctx.Localizer

	return El("section",
		Class("border-b border-slate-200 bg-white"),
		Div(
			Class("container mx-auto px-2 py-24 mobile:py-14"),
			Div(
				Class("max-w-3xl"),
				components.Eyebrow(l.MustLocalizeMessage(&i18n.Message{ID: "home.eyebrow"})),
				H1(
					Class("mt-5 text-5xl leading-[1.05] font-bold tracking-tight text-slate-900 mobile:text-4xl"),
					Text("Homeo sapiens"),
				),
				P(
					Class("mt-5 max-w-prose text-2xl leading-snug text-slate-700 mobile:text-xl font-heading"),
					Text(l.MustLocalizeMessage(&i18n.Message{ID: "home.tagline"})),
				),
				P(
					Class("mt-5 max-w-prose text-lg text-slate-600"),
					Text(l.MustLocalizeMessage(&i18n.Message{ID: "home.intro"})),
				),
				Div(
					Class("mt-8 flex flex-wrap items-center gap-x-6 gap-y-4"),
					A(Href("/events"), Class("button px-6"),
						Text(l.MustLocalizeMessage(&i18n.Message{ID: "home.browse_events"}))),
					components.TextLink("/watch", l.MustLocalizeMessage(&i18n.Message{ID: "home.explore_library"})),
				),
			),
		),
	)
}

func eventsSection(ctx *types.CustomContext, upcoming []*services.EventListDto) Node {
	l := ctx.Localizer

	return components.PageSection(
		components.SectionHeading(
			l.MustLocalizeMessage(&i18n.Message{ID: "home.agenda.label"}),
			l.MustLocalizeMessage(&i18n.Message{ID: "home.agenda.title"}),
			A(Href("/events"), Class("button secondary"),
				Text(l.MustLocalizeMessage(&i18n.Message{ID: "home.agenda.all_events"}))),
		),
		components.RowList(
			Map(upcoming, func(e *services.EventListDto) Node {
				return homeEventRow(ctx, e)
			}),
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
		Class("group grid grid-cols-[9rem_1fr_auto] items-center gap-6 border-b border-slate-200 py-6 no-underline transition-colors hover:bg-brand-50 mobile:grid-cols-1 mobile:gap-3 mobile:py-5"),
		Span(
			Class("text-sm text-slate-500 tabular-nums"),
			Text(helpers.FormatDateRange(e.StartsAt, e.EndsAt, ctx.Timezone, ctx.Language)),
		),
		Div(
			H3(
				Class("text-xl font-bold text-primary transition-colors group-hover:text-primary-hover"),
				Text(title),
			),
			Div(
				Class("mt-2 flex flex-wrap items-center gap-x-3 gap-y-1 text-sm text-slate-600"),
				events.EventTypeChip(ctx, e),
				events.EventLocationBadge(e, l, ctx.Language),
				Span(Text(helpers.FormatHosts(l, ctx.Language, e.Hosts))),
				Span(Class("font-semibold text-slate-900"), Text(price)),
			),
		),
		Span(
			Class("justify-self-end text-xs font-semibold tracking-wider text-primary uppercase mobile:justify-self-start"),
			Text(l.MustLocalizeMessage(&i18n.Message{ID: "home.agenda.details"})+" →"),
		),
	)
}

func aboutSection(ctx *types.CustomContext) Node {
	l := ctx.Localizer

	return components.PageSection(
		components.SectionHeading(
			l.MustLocalizeMessage(&i18n.Message{ID: "home.about.label"}),
			l.MustLocalizeMessage(&i18n.Message{ID: "home.about.title"}),
		),
		Div(
			Class("max-w-prose space-y-4 text-lg text-slate-700"),
			P(Text(l.MustLocalizeMessage(&i18n.Message{ID: "home.about.history"}))),
			P(Text(l.MustLocalizeMessage(&i18n.Message{ID: "home.about.mission"}))),
		),
		Div(
			Class("mt-8 flex flex-wrap items-center gap-x-6 gap-y-3"),
			components.TextLink("/blog", l.MustLocalizeMessage(&i18n.Message{ID: "home.about.read_blog"})),
			components.TextLink("/watch", l.MustLocalizeMessage(&i18n.Message{ID: "home.about.watch_recordings"})),
		),
	)
}
