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

func AppHeader(ctx context.Context) Node {
	l := ctx.Value("localizer").(*i18n.Localizer)
	user := ctx.Value(config.CurrentUserContextName).(*queries.User)

	return Header(
		Class("fixed inset-0 z-10 h-20 border-b-2 bg-white"),
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
