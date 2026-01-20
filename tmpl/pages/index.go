package pages

import (
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/events"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
)

func Index(ctx *types.CustomContext, eventRows []*services.EventListDto) Node {
	return layout.Layout(ctx, "Events", events.EventList(ctx, eventRows))
}
