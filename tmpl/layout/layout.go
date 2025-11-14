package layout

import (
	"context"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func RootLayout(ctx context.Context, title string, children ...Node) Node {
	return HTML(
		Head(
			Meta(Charset("UTF-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
			TitleEl(Text(title+" | Homeo sapiens")),
			Link(Rel("stylesheet"), Href("/assets/bundle.css")),
			fonts(),
			Script(Src("https://unpkg.com/lucide@latest"), Type("module")),
			Script(Type("module"), Text("lucide.createIcons();")),
		),
		Body(Group(children)),
	)
}

func Layout(ctx context.Context, title string, children ...Node) Node {
	return RootLayout(ctx, title,
		Class("flex min-h-screen flex-col"),
		AppHeader(ctx),
		Main(
			Class("flex-1 bg-slate-100 pt-24"),
			Div(
				Class("container mx-auto"),
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
			Href("https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:ital,wght@0,100..700;1,100..700&display=swap"),
			Rel("stylesheet"),
		),
	}
}

func NavLink(href string, text string) Node {
	return Li(
		A(
			Class("text-primary inline-flex h-full items-center rounded-sm px-3 font-semibold transition hover:bg-slate-200"),
			Href(href), Text(text),
		),
	)
}

func LanguageSwitcher(ctx context.Context, baseUrl string) Node {
	activeLocale := ctx.Value("lang").(string)
	l := ctx.Value("localizer").(*i18n.Localizer)
	otherLocale := "en"
	if activeLocale == "en" {
		otherLocale = "pl"
	}
	tooltip := l.MustLocalizeMessage(&i18n.Message{ID: "locale_switcher.switch_to"})

	return Li(
		A(
			Class("text-primary inline-flex h-full items-center gap-1 rounded-sm px-3 font-semibold uppercase transition hover:bg-slate-200"),
			Href(baseUrl+"?lang="+otherLocale), Title(tooltip), Aria("label", tooltip),
			I(Class("h-5 w-5"), Data("lucide", "languages")),
			Text(activeLocale),
		),
	)
}

func AppFooter() Node {
	return Footer(Class("h-30 border-t-2"),
		Div(
			Class("container mx-auto flex h-full items-center justify-center text-center"),
			P(
				Raw("&copy; 2024&ndash;2025 by Wydawnictwo Homeo Sapiens. All rights reserved."),
			),
		),
	)
}
