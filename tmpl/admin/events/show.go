package events

import (
	"fmt"

	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Show(ctx *types.CustomContext, event *services.EventDetailsDto) Node {
	title := fmt.Sprintf("Event: %s", event.TitleEn)

	return layout.AdminLayout(ctx, title, Div(
		Text(title)))
}
