package layout

import (
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func RootLayout(ctx *types.CustomContext, title string, children ...Node) Node {
	return HTML(
		Lang(ctx.Language),
		Head(
			Meta(Charset("UTF-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
			TitleEl(Text(title+" | Homeo sapiens")),
			AssetEntryPoint(ctx),
			fonts(),
			Meta(Name("user-timezone"), Content(ctx.Timezone.String())),
		),
		Body(Group(children)),
	)
}

func Layout(ctx *types.CustomContext, title string, children ...Node) Node {
	return RootLayout(ctx, title,
		Class("flex min-h-screen max-w-full flex-col overflow-x-hidden"),
		AppHeader(ctx),
		Main(
			Class("flex-1 bg-slate-100 pt-26 pb-6"),
			Div(
				Class("container mx-auto"),
				components.Flash(ctx.Flash),
				Group(children),
			),
		),
		AppFooter(),
	)
}

func fonts() Node {
	return Group{
		Link(
			Rel("preconnect"),
			Href("https://fonts.googleapis.com"),
		),
		Link(
			Rel("preconnect"),
			Href("https://fonts.gstatic.com"),
			CrossOrigin(""),
		),
		Link(
			Href("https://fonts.googleapis.com/css2?family=Noto+Sans:ital,wght@0,100..900;1,100..900&display=swap"),
			Rel("stylesheet"),
		),
	}
}

func LanguageSwitcher(ctx *types.CustomContext) Node {
	l := ctx.Localizer
	otherLocale := "en"
	if ctx.Language == "en" {
		otherLocale = "pl"
	}
	tooltip := l.MustLocalizeMessage(&i18n.Message{ID: "locale_switcher.switch_to"})
	baseUrl := ctx.RequestUrl.Path

	return A(
		Class("button tertiary gap-1 uppercase"),
		Href(baseUrl+"?lang="+otherLocale), Title(tooltip), Aria("label", tooltip),
		SVG(Class("h-5 w-5 fill-current"), Attr("viewBox", "0 0 640 640"), El("use", Href("/assets/language.svg#icon"))),
		Text(ctx.Language),
	)
}

func AppFooter() Node {
	return Footer(Class("border-primary/50 h-30 border-t shadow lg:text-base text-sm text-slate-600"),
		Div(
			Class("container mx-auto flex h-full items-center justify-center text-center"),
			P(
				Raw("&copy; 2024&ndash;2026 by Wydawnictwo Homeo Sapiens.<br/>All rights reserved."),
			),
		),
	)
}
