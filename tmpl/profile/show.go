package profile

import (
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)
import . "maragu.dev/gomponents"
import . "maragu.dev/gomponents/html"

func Show(ctx *types.CustomContext, messages []types.FlashMessage) Node {
	l := ctx.Localizer

	return layout.Layout(ctx, "Profile",
		Div(Class("card"),
			H2(
				Class("card-header"),
				Text(l.MustLocalizeMessage(&i18n.Message{
					ID: "profile.title",
				}))),

			components.Flash(messages),

			Form(
				Method("POST"),
				Action("/profile"),
				Class("mt-6 space-y-4"),
				Input(Type("hidden"), Name("_method"), Value("PUT")),
				components.InputField(&components.InputFieldOptions{
					Label: l.MustLocalizeMessage(&i18n.Message{
						ID: "common.users.email",
					}),
					ID:       "email",
					Value:    ctx.User.Email.String(),
					Readonly: true,
				}),
				components.InputField(&components.InputFieldOptions{
					Label: l.MustLocalizeMessage(&i18n.Message{
						ID: "common.users.given_name",
					}),
					Name:     "given_name",
					Value:    ctx.User.GivenName.String(),
					Required: true,
				}),
				components.InputField(&components.InputFieldOptions{
					Label: l.MustLocalizeMessage(&i18n.Message{
						ID: "common.users.family_name",
					}),
					Name:     "family_name",
					Value:    ctx.User.FamilyName.String(),
					Required: true,
				}),
				components.CountrySelect(&components.CountrySelectOptions{
					Label: l.MustLocalizeMessage(&i18n.Message{
						ID: "common.users.country",
					}),
					Value:    helpers.DerefOrEmpty(ctx.User.Country),
					Language: ctx.Language,
					Required: true,
				}),
				components.InputField(&components.InputFieldOptions{
					Label: l.MustLocalizeMessage(&i18n.Message{
						ID: "common.users.profession",
					}),
					Name:      "profession",
					Value:     helpers.DerefOrEmpty(ctx.User.Profession),
					Localizer: ctx.Localizer,
				}),
				components.InputField(&components.InputFieldOptions{
					Label: l.MustLocalizeMessage(&i18n.Message{
						ID: "common.users.licence_number",
					}),
					Name:      "licence_number",
					Value:     helpers.DerefEncrypted(ctx.User.LicenceNumber),
					Localizer: l,
				}),
				Button(
					Type("submit"),
					Class("button"),
					Text(l.MustLocalizeMessage(&i18n.Message{
						ID: "profile.form.submit",
					})),
				),
			),
		),
	)
}
