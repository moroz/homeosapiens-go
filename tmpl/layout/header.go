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
	title := user.GivenName.String() + " " + user.FamilyName.String()
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
					Span(Class("text-sm"), Text(user.Email.String())),
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

func NavLink(href string, text string) Node {
	return Li(
		A(
			Class("text-primary inline-flex h-full items-center rounded-sm px-3 font-semibold transition hover:bg-slate-200"),
			Href(href), Text(text),
		),
	)
}

func desktopNav(ctx context.Context) Node {
	l := ctx.Value("localizer").(*i18n.Localizer)
	user := ctx.Value(config.CurrentUserContextName).(*queries.User)

	return Nav(
		Class("h-full mobile:hidden"),
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
	)
}

func hamburgerTrigger() Node {
	return El("fade-burger", Attr("size", "lg"), Class("z-10 text-primary"), ID("hamburger-toggle"))
}

func HamburgerItem(href string, text string) Node {
	return Li(
		A(
			Href(href),
			Class("text-primary w-full h-12 text-center font-semibold flex items-center justify-center text-lg hover:bg-slate-100"),
			Text(text),
		),
	)
}

func mobileNav(ctx context.Context) Node {
	l := ctx.Value("localizer").(*i18n.Localizer)
	//user := ctx.Value(config.CurrentUserContextName).(*queries.User)

	return Div(
		Class("not-mobile:hidden z-10"),
		hamburgerTrigger(),
		Nav(
			Class("hamburger-menu pt-20"),
			// Fake header for shadow
			Div(Class("absolute top-0 left-0 right-0 h-20 border-b bg-white shadow")),
			Ul(
				Class("hamburger-items space-y-1 my-4"),
				HamburgerItem("/", l.MustLocalizeMessage(&i18n.Message{
					ID: "header.nav.home",
				})),
				HamburgerItem("/videos", l.MustLocalizeMessage(&i18n.Message{
					ID: "header.nav.videos",
				})),
				HamburgerItem("/dashboard", l.MustLocalizeMessage(&i18n.Message{
					ID: "header.nav.my_products",
				})),
			),
		),
	)
}

func AppHeader(ctx context.Context) Node {
	return Header(
		Class("fixed inset-0 z-10 h-20 border-b bg-white shadow"),
		Div(Class("container mx-auto flex h-full items-center justify-between mobile:px-2"),
			H1(
				Class("z-20"),
				A(
					Class("text-primary hover:text-primary-hover text-4xl font-bold no-underline transition-colors"),
					Href("/"),
					logo("h-15 mobile:h-12"),
				),
			),
			mobileNav(ctx),
			desktopNav(ctx),
		),
	)
}
