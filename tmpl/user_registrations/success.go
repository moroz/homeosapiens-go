package userregistrations

import (
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Success(ctx *types.CustomContext, user *types.UserDecorator) Node {
	l := ctx.Localizer
	title := l.MustLocalizeMessage(&i18n.Message{ID: "user_registrations.success.title"})

	return layout.Layout(ctx, title,
		Div(
			Class("card"),
			H2(Class("page-title"), Text(title)),
			Div(Class("prose"),
				Raw(l.MustLocalize(&i18n.LocalizeConfig{
					MessageID: "user_registrations.success.body_html",
					TemplateData: map[string]string{
						"Email": user.Email.String(),
					},
				})),
				If(user.IsGoogleAccount(), Raw(l.MustLocalizeMessage(&i18n.Message{
					ID: "user_registrations.success.gmail_notice_html",
				}))),
			),
		),
	)
}
