package layout

import (
	"context"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func logo(class string) Node {
	return SVG(
		Attr("viewbox", "0 0 1538 361"),
		Class(class),
		El("use",
			Href("/assets/logo.svg#logo"),
		),
	)
}

func UserHeader(ctx context.Context) Node {
	user := ctx.Value(config.CurrentUserContextName).(*queries.User)
	title := user.GivenName + " " + user.FamilyName
	l := ctx.Value("localizer").(*i18n.Localizer)

	return Div(
		Class("user-dropdown relative"),
		Button(Class("flex h-full cursor-pointer items-center justify-center rounded-sm px-3 transition-colors hover:bg-slate-200"), Title(title), components.Avatar(user)),
		Div(
			Class("dropdown absolute right-0 bottom-0 flex hidden translate-y-full flex-col overflow-hidden rounded-sm border border-slate-500 bg-white shadow"),
			Div(
				Class("flex items-center justify-between gap-4 px-3 py-2"),
				components.Avatar(user),
				Div(
					Class("flex flex-col text-right"),
					Strong(Text(title)),
					Span(Class("text-sm"), Text(user.Email)),
				),
			),
			A(
				Href("/sign-out"),
				Class("inline-flex h-10 cursor-pointer items-center justify-center border-t border-slate-400 text-center transition hover:bg-slate-100"),
				Text(l.MustLocalizeMessage(&i18n.Message{
					ID: "header.user_dropdown.sign_out",
				})),
			),
		),
	)
}

func AppHeader(ctx context.Context) Node {
	l := ctx.Value("localizer").(*i18n.Localizer)
	user := ctx.Value(config.CurrentUserContextName).(*queries.User)

	return Header(
		Class("fixed inset-0 z-10 h-20 border-b bg-white shadow"),
		Div(Class("container mx-auto flex h-full items-center justify-between"),
			H1(
				A(
					Class("text-primary hover:text-primary-hover text-4xl font-bold no-underline transition-colors"),
					Href("/"),
					logo("h-15"),
				),
			),
			Nav(
				Class("h-full"),
				Ul(
					Class("flex h-full gap-1 py-4"),
					NavLink("/", l.MustLocalizeMessage(&i18n.Message{
						ID: "header.nav.home",
					})),
					NavLink("/videos", l.MustLocalizeMessage(&i18n.Message{
						ID: "header.nav.videos",
					})),
					NavLink("/dashboard", l.MustLocalizeMessage(&i18n.Message{
						ID: "header.nav.my_products",
					})),
					LanguageSwitcher(ctx),
					If(user == nil, NavLink("/sign-in", l.MustLocalizeMessage(&i18n.Message{
						ID: "header.nav.sign_in",
					}))),
					Iff(user != nil, func() Node {
						return UserHeader(ctx)
					}),
				),
			),
		),
	)
}
