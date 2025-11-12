package sessions

import (
	"context"

	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func New(ctx context.Context, email string, msg string) Node {
	l := ctx.Value("localizer").(*i18n.Localizer)
	title := l.MustLocalizeMessage(&i18n.Message{
		ID: "sessions.new.header",
	})

	return layout.AuthLayout(ctx, title,
		Form(
			Class("grid gap-4"),
			Method("POST"),
			Action("/sessions"),

			If(msg != "", Div(Class("px-4 py-3 border-2 text-red-900 rounded-sm bg-red-100 mt-4"), Text(msg))),

			components.InputField(&components.InputFieldOptions{
				Label: l.MustLocalizeMessage(&i18n.Message{
					ID: "sessions.new.email",
				}),
				Name: "email",
				ID:   "email",
			}),

			components.InputField(&components.InputFieldOptions{
				Label: l.MustLocalizeMessage(&i18n.Message{
					ID: "sessions.new.password",
				}),
				Name: "password",
				Type: "password",
				ID:   "password",
			}),

			Button(Type("submit"), Class("button w-full h-10 text-lg"), Text(l.MustLocalizeMessage(&i18n.Message{ID: "sessions.new.submit"}))),
		),

		Footer(
			Class("text-center mt-4"),
			Raw(l.MustLocalizeMessage(&i18n.Message{ID: "sessions.new.no_account_yet_html"})),
		),
	)
}
