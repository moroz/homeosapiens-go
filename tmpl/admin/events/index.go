package events

import (
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Index(ctx *types.CustomContext, events []*services.EventListDto) Node {
	return layout.AdminLayout(ctx, "Events",
		Table(
			THead(
				Tr(
					Th(Text("Title")),
				),
			),
			TBody(
				Map(events, func(e *services.EventListDto) Node {
					return Tr(
						Td(Text(e.TitleEn)),
					)
				}),
			),
		),
	)
}
