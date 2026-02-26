package carts

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func deleteItemButton(ctx *types.CustomContext, eventId uuid.UUID) Node {
	return Form(
		Action("/cart_items"),
		Method("POST"),
		Data("confirm", ctx.Localizer.MustLocalizeMessage(&i18n.Message{
			ID: "cart_items.delete.confirmation",
		})),
		Input(Type("hidden"), Name("_method"), Value("DELETE")),
		Input(Type("hidden"), Name("event_id"), Value(eventId.String())),
		Button(Class("button tertiary p-2 font-normal text-xs"),
			SVG(Class("h-4 w-4 fill-current"), Attr("viewBox", "0 0 640 640"), El("use", Href("/assets/xmark.svg#icon"))),
			Span(
				Class("block mb-0.5"),
				Text(ctx.Localizer.MustLocalizeMessage(&i18n.Message{ID: "cart.table.delete"})),
			),
		),
	)
}

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
							Td(
								Class("text-left py-0"),
								Div(Class("flex items-center gap-2"),
									Text(title), deleteItemButton(ctx, item.Event.ID))),

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
				A(Href("/"), Class("button secondary"), Text(l.MustLocalizeMessage(&i18n.Message{
					ID: "cart.continue_shopping",
				}))),
				A(Href("/checkout"), Class("button"), Text(l.MustLocalizeMessage(&i18n.Message{
					ID: "cart.go_to_checkout",
				}))),
			),
		),
	)
}
