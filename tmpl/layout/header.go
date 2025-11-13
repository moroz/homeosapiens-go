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
			Href("/assets/logo.svg"),
		),
	)
}

func AppHeader(ctx context.Context) Node {
	l := ctx.Value("localizer").(*i18n.Localizer)
	user := ctx.Value(config.CurrentUserContextName).(*queries.User)

	return Header(
		Class("h-20 border-b-2 fixed inset-0 bg-white z-10"),
		Div(Class("container mx-auto flex justify-between h-full items-center"),
			H1(
				A(
					Class("text-primary no-underline transition-colors hover:text-primary-hover text-4xl font-bold outline"),
					Href("/"),
					logo("h-15"),
				),
			),
			Nav(
				Class("h-full"),
				Ul(
					Class("flex gap-1 py-4 h-full"),
					NavLink("/", l.MustLocalizeMessage(&i18n.Message{
						ID: "header.nav.home",
					})),
					NavLink("/events", l.MustLocalizeMessage(&i18n.Message{
						ID: "header.nav.events",
					})),
					NavLink("/videos", l.MustLocalizeMessage(&i18n.Message{
						ID: "header.nav.videos",
					})),
					NavLink("/dashboard", l.MustLocalizeMessage(&i18n.Message{
						ID: "header.nav.my_products",
					})),
					LanguageSwitcher(ctx, "/"),
					If(user == nil, NavLink("/sign-in", l.MustLocalizeMessage(&i18n.Message{
						ID: "header.nav.sign_in",
					}))),
					Iff(user != nil, func() Node {
						return Button(Class("flex h-full cursor-pointer items-center justify-center rounded-sm px-3 transition-colors hover:bg-slate-200"), components.Avatar(user))
					}),
				),
			),
		),
	)
}
