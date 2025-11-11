package layout

import (
	"context"

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
	localizer := ctx.Value("localizer").(*i18n.Localizer)
	lang := ctx.Value("lang").(string)
	otherLocale := "pl"
	if lang == "pl" {
		otherLocale = "en"
	}

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
					Class("flex gap-2 py-4 h-full"),
					LanguageSwitcher(ctx, "/", otherLocale),
					NavLink("/events", localizer.MustLocalizeMessage(&i18n.Message{
						ID:    "header.nav.events",
						Other: "Events",
					})),
					NavLink("/videos", localizer.MustLocalizeMessage(&i18n.Message{
						ID:    "header.nav.videos",
						Other: "Watch",
					})),
					NavLink("/dashboard", localizer.MustLocalizeMessage(&i18n.Message{
						ID:    "header.nav.my_products",
						Other: "My Products",
					})),
					NavLink("/sign-in", localizer.MustLocalizeMessage(&i18n.Message{
						ID:    "header.nav.sign_in",
						Other: "Sign in",
					})),
				),
			),
		),
	)
}
