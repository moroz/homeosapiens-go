package events

import (
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Index(ctx *types.CustomContext, events []*services.EventListDto) Node {
	return layout.Layout(ctx, "Events", EventList(ctx, events))
}

func EventList(ctx *types.CustomContext, events []*services.EventListDto) Node {
	return Div(
		Class("space-y-4"),
		Map(events, func(e *services.EventListDto) Node {
			return EventCard(ctx, e)
		}),
	)
}
