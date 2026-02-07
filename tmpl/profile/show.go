package profile

import (
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)
import . "maragu.dev/gomponents"
import . "maragu.dev/gomponents/html"

func Show(ctx *types.CustomContext) Node {
	l := ctx.Localizer

	return layout.Layout(ctx, "Profile",
		Div(Class("card"),
			H2(
				Class("card-header"),
				Text(l.MustLocalizeMessage(&i18n.Message{
					ID: "profile.title",
				}))),

			Form(
				Method("POST"),
				Action("/profile"),
				Class("mt-6 space-y-4"),
				Input(Type("hidden"), Name("_method"), Value("PUT")),
				components.InputField(&components.InputFieldOptions{
					Label: l.MustLocalizeMessage(&i18n.Message{
						ID: "profile.form.labels.email",
					}),
					Name:  "email",
					Value: ctx.User.Email.String(),
				}),
			),
		),
	)
}
