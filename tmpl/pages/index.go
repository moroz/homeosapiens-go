package pages

import (
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/events"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Events(ctx *types.CustomContext, eventRows []*services.EventListDto) Node {
	title := "Events"
	if ctx.IsPolish() {
		title = "Wydarzenia"
	}

	return layout.Layout(ctx, title,
		H2(Class("page-title"), Text(title)),
		events.EventList(ctx, eventRows),
	)
}
