package events

import (
	"fmt"
	"strconv"

	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Index(ctx *types.CustomContext, events []*services.EventListDto) Node {
	return layout.AdminLayout(ctx, "Events",
		Table(
			Class("index-table"),
			THead(
				Tr(
					Th(Text("Title EN")),
					Th(Text("Title PL")),
					Th(Text("Registrations")),
				),
			),
			TBody(
				Map(events, func(e *services.EventListDto) Node {
					return Tr(
						Data("url", fmt.Sprintf("/admin/events/%s", e.ID)),
						Td(Text(e.TitleEn)),
						Td(Text(e.TitlePl)),
						Td(Text(strconv.Itoa(e.RegistrationCount))),
					)
				}),
			),
		),
	)
}
