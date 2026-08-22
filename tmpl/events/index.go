package events

import (
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
)

func EventList(ctx *types.CustomContext, events []*services.EventListDto) Node {
	return components.RowList(
		Map(events, func(e *services.EventListDto) Node {
			return EventCard(ctx, e)
		}),
	)
}
