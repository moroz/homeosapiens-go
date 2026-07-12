package layout

import (
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/shopspring/decimal"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func NavLink(href string, text string) Node {
	return Li(
		Class("h-full"),
		A(
			Class("no-underline inline-block p-3 hover:bg-slate-100 transition-colors text-base rounded-sm"),
			Href(href), Text(text),
		),
	)
}

func desktopNav(ctx *types.CustomContext) Node {
	l := ctx.Localizer

	return Group{
		Nav(Class("grid mobile:hidden"),
			Ul(
				Class("flex items-center gap-1 ml-2"),
				NavLink("/events", l.MustLocalizeMessage(&i18n.Message{
					ID: "header.nav.events",
				})),
				NavLink("/watch", l.MustLocalizeMessage(&i18n.Message{
					ID: "header.nav.videos",
				})),
				NavLink("/videos", l.MustLocalizeMessage(&i18n.Message{
					ID: "header.nav.my_products",
				})),
			),
		),
		Div(
			Class("z-20 flex items-center gap-1 mobile:hidden ml-auto"),
			Iff(ctx.Cart != nil && !ctx.Cart.ProductTotal.Equal(decimal.Zero), func() Node {
				return A(Href("/cart"), Class("button tertiary z-20 gap-1"),
					Title(l.MustLocalizeMessage(&i18n.Message{
						ID: "header.cart",
					})),
					SVG(Class("h-5 w-5 fill-current"), Attr("viewBox", "0 0 640 640"), El("use", Href("/assets/cart-shopping.svg#icon"))),
					Text(helpers.FormatPrice(ctx.Cart.ProductTotal, "PLN", ctx.Language)),
				)
			}),
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
	return El("fade-burger", Attr("size", "lg"), Class("z-10 text-primary"), ID("hamburger-toggle"))
}

func HamburgerItem(href string, text string) Node {
	return Li(
		A(
			Href(href),
			Class("flex h-12 w-full items-center justify-center text-center text-lg font-semibold text-primary hover:bg-slate-100"),
			Text(text),
		),
	)
}

func mobileNav(ctx *types.CustomContext) Node {
	l := ctx.Localizer
	//user := ctx.Value(config.CurrentUserContextName).(*queries.User)

	return Div(
		Class("z-10 not-mobile:hidden"),
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
				HamburgerItem("/events", l.MustLocalizeMessage(&i18n.Message{
					ID: "header.nav.events",
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
		Class("fixed inset-0 z-20 h-20 border-b border-primary/15 bg-white shadow-sm backdrop-blur-md"),
		Div(Class("container mx-auto flex h-full items-center px-6 mobile:px-2"),
			H1(
				Class("z-20"),
				A(
					Class("text-4xl font-bold text-primary no-underline transition-colors hover:text-primary-hover"),
					Href("/"),
					components.Logo("h-15 mobile:h-12"),
				),
			),
			mobileNav(ctx),
			desktopNav(ctx),
		),
	)
}
