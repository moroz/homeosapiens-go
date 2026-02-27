package components

import (
	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/tmpl/components/icons"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AddToCartButton(l *i18n.Localizer, eventId uuid.UUID, inCart int) Node {
	if inCart > 0 {
		return InCartButton(l)
	}

	return Form(
		Action("/cart_items"),
		Method("POST"),
		Input(Type("hidden"), Name("event_id"), Value(eventId.String())),
		Button(
			Class("button"),
			icons.Icon(&icons.IconProps{Name: "cart-arrow-down"}),
			Text(l.MustLocalizeMessage(&i18n.Message{ID: "common.events.add_to_cart"})),
		),
	)
}

func InCartButton(l *i18n.Localizer) Node {
	return A(
		Href("/cart"),
		Class("button success gap-1"),
		icons.Icon(&icons.IconProps{Name: "check", ViewBox: "0 0 448 512", Classes: "h-4 w-4"}),
		Text(l.MustLocalizeMessage(&i18n.Message{
			ID: "common.events.added_to_cart",
		})),
	)
}
