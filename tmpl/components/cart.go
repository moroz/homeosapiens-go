package components

import (
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func AddToCartButton(l *i18n.Localizer, event *queries.Event, inCart int) Node {
	if inCart > 0 {
		return nil
	}

	return Form(
		Action("/cart_items"),
		Method("POST"),
		Input(Type("hidden"), Name("event_id"), Value(event.ID.String())),
		Button(
			Class("button"),
			Text(l.MustLocalizeMessage(&i18n.Message{ID: "common.events.add_to_cart"})),
		),
	)
}
