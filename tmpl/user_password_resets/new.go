package user_password_resets

import (
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func New(ctx *types.CustomContext, email string, errorMsg string) Node {
	l := ctx.Localizer
	title := l.MustLocalizeMessage(&i18n.Message{
		ID: "user_password_resets.new.title",
	})

	return layout.AuthLayout(ctx, title,
		Form(
			Class("mt-4 grid gap-4"),
			Method("POST"),
			Action("/reset-password"),

			If(errorMsg != "", Div(Class("alert danger"), Text(errorMsg))),

			components.InputField(&components.InputFieldOptions{
				Label: l.MustLocalizeMessage(&i18n.Message{
					ID: "common.users.email",
				}),
				Name:         "email",
				ID:           "email",
				Value:        email,
				Required:     true,
				Autocomplete: "email",
				Localizer:    l,
			}),

			Button(Type("submit"), Class("button font-fallback min-h-10 w-full"),
				Text(l.MustLocalizeMessage(&i18n.Message{ID: "user_password_resets.new.submit"})),
			),
		),

		Footer(
			Class("mt-4 text-center"),
			A(Href("/sign-in"), Text(
				l.MustLocalizeMessage(&i18n.Message{
					ID: "user_password_resets.new.sign_in",
				}),
			)),
			Raw(" &bull; "),
			A(Href("/sign-up"), Text(
				l.MustLocalizeMessage(&i18n.Message{
					ID: "user_password_resets.new.sign_up",
				}),
			)),
		),
	)
}
