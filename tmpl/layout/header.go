package layout

import (
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/types"
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

func UserHeader(ctx *types.CustomContext) Node {
	title := ctx.User.GivenName.String() + " " + ctx.User.FamilyName.String()
	l := ctx.Localizer

	return Div(
		Class("user-dropdown relative"),
		Button(Class("flex h-full cursor-pointer items-center justify-center rounded-sm px-3 transition-colors hover:bg-slate-200"), Title(title), components.Avatar(ctx.User)),
		Div(
			Class("dropdown absolute right-0 bottom-0 flex hidden translate-y-full flex-col overflow-hidden rounded-sm border border-slate-500 bg-white shadow"),
			Div(
				Class("flex items-center justify-between gap-4 px-3 py-2"),
				components.Avatar(ctx.User),
				Div(
					Class("flex flex-col text-right"),
					Strong(Text(title)),
					Span(Class("text-sm"), Text(ctx.User.Email.String())),
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
		Class("h-full"),
		A(
			Class("button tertiary"),
			Href(href), Text(text),
		),
	)
}

func desktopNav(ctx *types.CustomContext) Node {
	l := ctx.Localizer

	return Group{
		Nav(Class("mobile:hidden absolute inset-0 grid place-items-center"),
			Ul(
				Class("flex items-center gap-1"),
				NavLink("/", l.MustLocalizeMessage(&i18n.Message{
					ID: "header.nav.home",
				})),
				NavLink("/videos", l.MustLocalizeMessage(&i18n.Message{
					ID: "header.nav.videos",
				})),
				NavLink("/dashboard", l.MustLocalizeMessage(&i18n.Message{
					ID: "header.nav.my_products",
				})),
			),
		),
		Div(
			Class("z-20 flex items-center gap-1"),
			LanguageSwitcher(ctx),
			If(ctx.User == nil, A(Href("/sign-in"), Class("button secondary z-20"), Text(l.MustLocalizeMessage(&i18n.Message{
				ID: "header.nav.sign_in",
			})))),
			Iff(ctx.User != nil, func() Node {
				return UserHeader(ctx)
			}),
		),
	}
}

func hamburgerTrigger() Node {
	return El("fade-burger", Attr("size", "lg"), Class("text-primary z-10"), ID("hamburger-toggle"))
}

func HamburgerItem(href string, text string) Node {
	return Li(
		A(
			Href(href),
			Class("text-primary flex h-12 w-full items-center justify-center text-center text-lg font-semibold hover:bg-slate-100"),
			Text(text),
		),
	)
}

func mobileNav(ctx *types.CustomContext) Node {
	l := ctx.Localizer
	//user := ctx.Value(config.CurrentUserContextName).(*queries.User)

	return Div(
		Class("not-mobile:hidden z-10"),
		hamburgerTrigger(),
		Nav(
			Class("hamburger-menu pt-20"),
			// Fake header for shadow
			Div(Class("absolute top-0 right-0 left-0 h-20 border-b bg-white shadow")),
			Ul(
				Class("hamburger-items my-4 space-y-1"),
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

func AppHeader(ctx *types.CustomContext) Node {
	return Header(
		Class("border-primary/50 fixed inset-0 z-10 h-20 border-b bg-white shadow"),
		Div(Class("mobile:px-2 flex h-full items-center justify-between px-6"),
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
