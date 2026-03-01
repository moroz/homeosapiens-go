package orders

import (
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/carts"
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)
import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func New(ctx *types.CustomContext, cart *services.CartViewDto) Node {
	l := ctx.Localizer

	return layout.Layout(ctx, "Checkout",
		Div(
			Class("card"),
			carts.CartTable(ctx, cart),
			Form(
				Action("/orders"),
				Method("POST"),
				Class("grid gap-2"),
				components.InputField(&components.InputFieldOptions{
					Label:        l.MustLocalizeMessage(&i18n.Message{ID: "common.users.given_name"}),
					Name:         "given_name",
					Autocomplete: "given-name",
					Required:     true,
					Localizer:    ctx.Localizer,
				}),
				components.InputField(&components.InputFieldOptions{
					Label:        l.MustLocalizeMessage(&i18n.Message{ID: "common.users.family_name"}),
					Name:         "family_name",
					Autocomplete: "family-name",
					Required:     true,
					Localizer:    ctx.Localizer,
				}),
			),
		),
	)
}
