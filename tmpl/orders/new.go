package orders

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

func New(ctx *types.CustomContext, cart *services.CartViewDto, params *types.OrderParams) Node {
	l := ctx.Localizer
	title := l.MustLocalizeMessage(&i18n.Message{
		ID: "cart.title",
	})

	if cart.IsEmpty() {
		return layout.Layout(ctx, title,
			Div(
				Class("card"),
				H2(Class("page-title"), Text(title)),

				Text(l.MustLocalizeMessage(&i18n.Message{ID: "cart.cart_empty"})),
			),
		)
	}

	return layout.Layout(ctx, title,
		Div(Class("card"),
			H2(Class("page-title"), Text(title)),

			CartTable(ctx, cart),

			Form(
				Action("/orders"),
				Method("POST"),

				H3(Class("text-2xl font-bold text-primary mt-8 mb-0"), Text(l.MustLocalizeMessage(&i18n.Message{ID: "orders.contact_information"}))),

				If(ctx.User == nil,
					components.GoogleButton(l.MustLocalizeMessage(&i18n.Message{ID: "orders.checkout_with_google"}), "/sign-in?ref=%2Fcart"),
				),

				Section(Class("gap-4 grid my-4"),
					If(ctx.User == nil,
						P(Raw(l.MustLocalize(&i18n.LocalizeConfig{
							DefaultMessage: &i18n.Message{
								ID: "orders.already_have_an_account_html",
							},
							TemplateData: map[string]string{
								"Url": "/sign-in?ref=%2Fcart",
							},
						}))),
					),

					If(ctx.User != nil,
						components.DisplayField(&components.DisplayFieldOptions{
							Label: l.MustLocalizeMessage(&i18n.Message{ID: "orders.form.billing_email"}),
							Value: params.Email,
							Name:  "email",
						}),
					),
					If(ctx.User == nil,
						components.InputField(&components.InputFieldOptions{
							Label:        l.MustLocalizeMessage(&i18n.Message{ID: "orders.form.billing_email"}),
							Name:         "email",
							Autocomplete: "email",
							Required:     true,
							Localizer:    l,
							Value:        params.Email,
						}),
					),

					components.InputGroup(
						components.InputField(&components.InputFieldOptions{
							Label:        l.MustLocalizeMessage(&i18n.Message{ID: "common.users.given_name"}),
							Name:         "billing_given_name",
							Value:        params.BillingGivenName,
							Autocomplete: "given-name",
							Required:     true,
							Localizer:    l,
						}),

						components.InputField(&components.InputFieldOptions{
							Label:        l.MustLocalizeMessage(&i18n.Message{ID: "common.users.family_name"}),
							Name:         "billing_family_name",
							Value:        params.BillingFamilyName,
							Autocomplete: "family-name",
							Required:     true,
							Localizer:    l,
						}),
					),
				),

				H3(Class("text-2xl font-bold text-primary mt-8 mb-0"), Text(l.MustLocalizeMessage(&i18n.Message{ID: "orders.billing_address"}))),

				Section(Class("grid gap-4 my-4"),
					components.CountrySelect(&components.CountrySelectOptions{
						Name:       "billing_country",
						Language:   ctx.Language,
						Value:      params.BillingCountry,
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
						Value:        params.BillingPhone,
					}),

					components.InputField(&components.InputFieldOptions{
						Label:        l.MustLocalizeMessage(&i18n.Message{ID: "orders.form.billing_address_line1"}),
						Name:         "billing_address_line1",
						Autocomplete: "address-line1",
						Required:     true,
						Localizer:    l,
						Value:        params.BillingAddressLine1,
					}),

					components.InputField(&components.InputFieldOptions{
						Label:        l.MustLocalizeMessage(&i18n.Message{ID: "orders.form.billing_address_line2"}),
						Name:         "billing_address_line2",
						Autocomplete: "address-line2",
						Localizer:    l,
						Value:        params.BillingAddressLine2,
					}),

					components.InputGroup(
						components.InputField(&components.InputFieldOptions{
							Label:        l.MustLocalizeMessage(&i18n.Message{ID: "orders.form.billing_city"}),
							Name:         "billing_city",
							Autocomplete: "billing address-level2",
							Required:     true,
							Localizer:    l,
							Value:        params.BillingCity,
						}),

						components.InputField(&components.InputFieldOptions{
							Label:        l.MustLocalizeMessage(&i18n.Message{ID: "orders.form.billing_postal_code"}),
							Name:         "billing_postal_code",
							Autocomplete: "billing postal-code",
							Required:     true,
							Localizer:    l,
							Value:        params.BillingPostalCode,
						}),
					),
				),
			),
		),
	)
}
