package pages

import (
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/events"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
)

func Events(ctx *types.CustomContext, eventRows []*services.EventListDto) Node {
	l := ctx.Localizer
	title := l.MustLocalizeMessage(&i18n.Message{ID: "events.index.title"})

	return layout.PageLayout(ctx, title,
		components.PageHeader(
			l.MustLocalizeMessage(&i18n.Message{ID: "events.index.eyebrow"}),
			title,
			components.Standfirst(l.MustLocalizeMessage(&i18n.Message{ID: "events.index.standfirst"})),
		),
		components.PageSection(
			If(len(eventRows) == 0,
				components.EmptyState(l.MustLocalizeMessage(&i18n.Message{ID: "events.index.no_events"}))),
			If(len(eventRows) > 0, events.EventList(ctx, eventRows)),
		),
	)
}
