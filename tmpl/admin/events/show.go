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

	return layout.AdminLayout(ctx, title,
		H2(Class("font-bold text-2xl"), Text(event.TitleEn)),
		Table(
			Class("data-table"),
			TBody(
				Tr(
					Th(Scope("row"), Text("Title (EN)")),
					Td(Text(event.TitleEn)),
				),
				Tr(
					Th(Scope("row"), Text("Title (PL)")),
					Td(Text(event.TitlePl)),
				),
				Tr(
					Th(Scope("row"), Text("Venue")),
					Td(Iff(event.VenueNamePl != nil, func() Node {
						return Text("PL: " + *event.VenueNamePl)
					}),
					),
				),
			),
		),
	)
}
