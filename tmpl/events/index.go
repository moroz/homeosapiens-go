package events

import (
	"context"

	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Index(ctx context.Context, events []*services.EventListDto) Node {
	return layout.Layout(ctx, "Events", EventList(ctx, events))
}

func EventList(ctx context.Context, events []*services.EventListDto) Node {
	return Div(
		Class("mb-8 grid gap-4"),
		Map(events, func(e *services.EventListDto) Node {
			return EventCard(ctx, e)
		}),
	)
}
