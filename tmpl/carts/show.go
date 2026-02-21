package carts

import (
	"fmt"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Show(ctx *types.CustomContext, cart *services.CartViewDto) Node {
	l := ctx.Localizer
	title := l.MustLocalizeMessage(&i18n.Message{
		ID: "cart.title",
	})

	return layout.Layout(ctx, title,
		Div(Class("card"),
			H2(Class("page-title"), Text(title)),

			Table(
				Class("cart-table"),
				THead(
					Tr(
						Th(
							Class("text-left"),
							Text(l.MustLocalizeMessage(&i18n.Message{ID: "cart.table.product_title"})),
						),
						Th(
							Class("w-25"),
							Text(l.MustLocalizeMessage(&i18n.Message{ID: "cart.table.unit_price"})),
						),
						Th(
							Class("w-25"),
							Text(l.MustLocalizeMessage(&i18n.Message{ID: "cart.table.quantity"})),
						),
						Th(
							Class("w-25"),
							Text(l.MustLocalizeMessage(&i18n.Message{ID: "cart.table.subtotal"})),
						),
					),
				),
				TBody(
					Map(cart.CartItems, func(item *queries.GetCartItemsByCartIdRow) Node {
						title := item.Event.TitleEn
						if ctx.Language == "pl" {
							title = item.Event.TitlePl
						}

						return Tr(
							Td(Class("text-left"), Text(title)),
							Td(Text(helpers.FormatPrice(*item.Event.BasePriceAmount, "PLN", ctx.Language))),
							Td(Raw(fmt.Sprintf("&times; %d", item.CartLineItem.Quantity))),
							Td(Text(helpers.FormatPrice(item.Subtotal, "PLN", ctx.Language))),
						)
					}),
				),
				TFoot(
					Tr(
						Th(Class("text-right"), Scope("row"), ColSpan("3"), Text(l.MustLocalizeMessage(&i18n.Message{
							ID: "cart.table.grand_total",
						}))),
						Td(Text(helpers.FormatPrice(cart.GrandTotal, "PLN", ctx.Language))),
					),
				),
			),
			Section(
				Class("flex justify-between"),
				A(Href("/"), Class("button secondary"), Text("Continue shopping")),
				A(Href("/checkout"), Class("button"), Text("Go to checkout")),
			),
		),
	)
}
