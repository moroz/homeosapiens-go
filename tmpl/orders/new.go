package orders

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/internal/countries"
	"github.com/moroz/homeosapiens-go/internal/phone"
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

func New(ctx *types.CustomContext, cart *services.CartViewDto, params *types.OrderParams, errors validation.Errors) Node {
	l := ctx.Localizer
	title := l.MustLocalizeMessage(&i18n.Message{
		ID: "cart.title",
	})

	examplePhone := phone.ExamplePhoneNumber(params.BillingCountry)
	if params.BillingCountry == "PL" {
		examplePhone = phone.ExamplePhoneNumber("GB")
	}

	if cart.IsEmpty() {
		return layout.Layout(ctx, title,
			Div(
				Class("card"),
				H2(Class("page-title"), Text(title)),

				P(
					Data("testid", "empty-message"),
					Text(l.MustLocalizeMessage(&i18n.Message{ID: "cart.cart_empty"})),
				),
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

				H3(Class("text-2xl font-bold text-primary mt-8 mb-4"), Text(l.MustLocalizeMessage(&i18n.Message{ID: "orders.contact_information"}))),

				If(ctx.User == nil,
					Group{
						components.GoogleButton(l.MustLocalizeMessage(&i18n.Message{ID: "orders.checkout_with_google"}), "/cart", "my-0"),
						P(Class("my-4"), Text(l.MustLocalizeMessage(&i18n.Message{ID: "orders.or_continue_with_email"}))),
					},
				),

				Section(Class("gap-4 grid mb-4"),
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
							HelperText:   l.MustLocalizeMessage(&i18n.Message{ID: "orders.form.billing_email_helper_text"}),
							Error:        errors,
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
							Error:        errors,
						}),

						components.InputField(&components.InputFieldOptions{
							Label:        l.MustLocalizeMessage(&i18n.Message{ID: "common.users.family_name"}),
							Name:         "billing_family_name",
							Value:        params.BillingFamilyName,
							Autocomplete: "family-name",
							Required:     true,
							Localizer:    l,
							Error:        errors,
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
						HelperText: l.MustLocalize(&i18n.LocalizeConfig{
							DefaultMessage: &i18n.Message{ID: "orders.form.billing_phone_helper_text"},
							TemplateData: map[string]string{
								"Example": examplePhone,
							},
						}),
						Localizer: l,
						Value:     params.BillingPhone,
						Error:     errors,
					}),

					components.InputField(&components.InputFieldOptions{
						Label:        l.MustLocalizeMessage(&i18n.Message{ID: "orders.form.billing_address_line1"}),
						Name:         "billing_address_line1",
						Autocomplete: "address-line1",
						Required:     true,
						Localizer:    l,
						Value:        params.BillingAddressLine1,
						Error:        errors,
					}),

					components.InputField(&components.InputFieldOptions{
						Label:        l.MustLocalizeMessage(&i18n.Message{ID: "orders.form.billing_address_line2"}),
						Name:         "billing_address_line2",
						Autocomplete: "address-line2",
						Localizer:    l,
						Value:        params.BillingAddressLine2,
						Error:        errors,
					}),

					components.InputGroup(
						components.InputField(&components.InputFieldOptions{
							Label:        l.MustLocalizeMessage(&i18n.Message{ID: "orders.form.billing_city"}),
							Name:         "billing_city",
							Autocomplete: "billing address-level2",
							Required:     true,
							Localizer:    l,
							Value:        params.BillingCity,
							Error:        errors,
						}),

						components.InputField(&components.InputFieldOptions{
							Label:        l.MustLocalizeMessage(&i18n.Message{ID: "orders.form.billing_postal_code"}),
							Name:         "billing_postal_code",
							Autocomplete: "billing postal-code",
							Required:     true,
							Localizer:    l,
							Value:        params.BillingPostalCode,
							Error:        errors,
						}),
					),
				),
			),
		),
	)
}
