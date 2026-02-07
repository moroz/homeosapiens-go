package layout

import (
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func UserHeader(ctx *types.CustomContext) Node {
	title := ctx.User.GivenName.String() + " " + ctx.User.FamilyName.String()
	l := ctx.Localizer

	return Div(
		Class("user-dropdown relative"),
		Button(Class("flex h-full cursor-pointer items-center justify-center rounded-sm px-3 transition-colors hover:bg-slate-200"), Title(title), components.Avatar(ctx.User)),
		Nav(
			Class("dropdown hidden"),
			Section(
				Class("flex items-center justify-between gap-4 px-3 py-2"),
				components.Avatar(ctx.User),
				Div(
					Class("flex flex-col text-right"),
					Strong(Text(title)),
					Span(Class("text-sm"), Text(ctx.User.Email.String())),
				),
			),
			Section(
				A(Class("dropdown-item"), Href("/profile"), Text(l.MustLocalizeMessage(&i18n.Message{
					ID: "header.user_dropdown.profile",
				}))),
				Iff(ctx.User.UserRole == queries.UserRoleAdministrator, func() Node {
					return A(
						Href("/admin"),
						Class("dropdown-item"),
						Text("Admin dashboard"),
					)
				}),
			),
			Section(
				A(
					Href("/sign-out"),
					Class("dropdown-item"),
					Text(l.MustLocalizeMessage(&i18n.Message{
						ID: "header.user_dropdown.sign_out",
					})),
				),
			),
		),
	)
}
