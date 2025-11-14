package pages

import (
	"context"

	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/events"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	. "maragu.dev/gomponents"
)

func Index(ctx context.Context, eventRows []*services.EventListDto) Node {
	return layout.Layout(ctx, "Events", events.EventList(ctx, eventRows))
}
