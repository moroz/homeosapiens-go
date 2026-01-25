package sessions

import (
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func New(ctx *types.CustomContext, email string, msg string) Node {
	l := ctx.Localizer
	title := l.MustLocalizeMessage(&i18n.Message{
		ID: "sessions.new.header",
	})

	return layout.AuthLayout(ctx, title,
		components.GoogleButton(l.MustLocalizeMessage(&i18n.Message{
			ID: "sessions.new.sign_in_with_google",
		}), ""),
		Hr(Class("my-6")),
		Form(
			Class("grid gap-4"),
			Method("POST"),
			Action("/sessions"),

			If(msg != "", Div(Class("alert danger"), Text(msg))),

			components.InputField(&components.InputFieldOptions{
				Label: l.MustLocalizeMessage(&i18n.Message{
					ID: "sessions.new.email",
				}),
				Name:         "email",
				ID:           "email",
				Value:        email,
				Required:     true,
				Autocomplete: "email",
				Localizer:    l,
			}),

			components.InputField(&components.InputFieldOptions{
				Label: l.MustLocalizeMessage(&i18n.Message{
					ID: "sessions.new.password",
				}),
				Name:         "password",
				Type:         "password",
				ID:           "password",
				Autocomplete: "current-password",
				Required:     true,
				Localizer:    l,
			}),

			Button(Type("submit"), Class("button font-fallback h-10 w-full text-lg"), Text(l.MustLocalizeMessage(&i18n.Message{ID: "sessions.new.submit"}))),
		),

		Footer(
			Class("mt-4 text-center"),
			Raw(
				l.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "sessions.new.no_account_yet_html",
					TemplateData: map[string]string{
						"Path": "/sign-up",
					},
				}),
			),
		),
	)
}
