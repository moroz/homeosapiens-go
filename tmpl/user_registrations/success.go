package userregistrations

import (
	"strings"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Success(ctx *types.CustomContext, user *queries.User, tokenParam string) Node {
	l := ctx.Localizer
	title := l.MustLocalizeMessage(&i18n.Message{ID: "user_registrations.success.title"})

	isGoogle := user != nil && strings.HasSuffix(user.Email.String(), "@gmail.com")

	email := ""
	if user != nil {
		email = user.Email.String()
	}

	return layout.Layout(ctx, title,
		Div(
			Class("card"),
			H2(Class("page-title"), Text(title)),
			Div(Class("prose"),
				Iff(user != nil, func() Node {
					return Raw(l.MustLocalize(&i18n.LocalizeConfig{
						MessageID: "user_registrations.success.body_html",
						TemplateData: map[string]string{
							"Email": email,
						},
					}))
				}),
				If(isGoogle, Raw(l.MustLocalizeMessage(&i18n.Message{
					ID: "user_registrations.success.gmail_notice_html",
				}))),
			),
		),
	)
}
