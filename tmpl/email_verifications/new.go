package email_verifications

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
		ID: "email_verifications.new.title",
	})

	return layout.AuthLayout(ctx, title,
		Form(
			Class("grid gap-4 mt-4"),
			Method("POST"),
			Action("/email-verifications"),

			If(errorMsg != "", Div(Class("alert danger"), Text(errorMsg))),

			components.InputField(&components.InputFieldOptions{
				Label: l.MustLocalizeMessage(&i18n.Message{
					ID: "email_verifications.new.email",
				}),
				Name:         "email",
				ID:           "email",
				Value:        email,
				Required:     true,
				Autocomplete: "email",
				Localizer:    l,
			}),

			Button(Type("submit"), Class("button font-fallback h-10 w-full text-lg"),
				Text(l.MustLocalizeMessage(&i18n.Message{ID: "email_verifications.new.submit"})),
			),
		),
	)
}
