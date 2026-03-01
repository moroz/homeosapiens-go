package carts

import (
	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/internal/countries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/components/icons"
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
			icons.Icon(&icons.IconProps{Name: "xmark", Classes: "h-4 w-4"}),
			Text(ctx.Localizer.MustLocalizeMessage(&i18n.Message{ID: "cart.table.delete"})),
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

			CartTable(ctx, cart),

			H3(Class("text-2xl font-bold text-primary my-3"), Text("Billing address")),

			Form(
				Action("/orders"),
				Method("POST"),
				Class("grid gap-4"),
				components.InputGroup(
					components.InputField(&components.InputFieldOptions{
						Label:        l.MustLocalizeMessage(&i18n.Message{ID: "common.users.given_name"}),
						Name:         "billing_given_name",
						Autocomplete: "given-name",
						Required:     true,
						Localizer:    l,
					}),

					components.InputField(&components.InputFieldOptions{
						Label:        l.MustLocalizeMessage(&i18n.Message{ID: "common.users.family_name"}),
						Name:         "billing_family_name",
						Autocomplete: "family-name",
						Required:     true,
						Localizer:    l,
					}),
				),
				components.CountrySelect(&components.CountrySelectOptions{
					Name:       "billing_country",
					Language:   ctx.Language,
					Value:      "",
					Required:   true,
					Label:      l.MustLocalizeMessage(&i18n.Message{ID: "orders.form.billing_country"}),
					HelperText: l.MustLocalizeMessage(&i18n.Message{ID: "orders.form.billing_country_helper_text"}),
					Countries:  countries.EuMemberStates(),
				}),

				components.InputField(&components.InputFieldOptions{
					Label:        l.MustLocalizeMessage(&i18n.Message{ID: "orders.form.billing_phone"}),
					Name:         "billing_phone",
					Autocomplete: "tel",
					Required:     false,
					HelperText:   l.MustLocalizeMessage(&i18n.Message{ID: "orders.form.billing_phone_helper_text"}),
					Localizer:    l,
				}),

				components.InputField(&components.InputFieldOptions{
					Label:        l.MustLocalizeMessage(&i18n.Message{ID: "orders.form.billing_street"}),
					Name:         "billing_street",
					Autocomplete: "billing street-address",
					Required:     true,
					Localizer:    l,
				}),
			),
		),
	)
}
