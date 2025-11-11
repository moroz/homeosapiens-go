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
		),
		Body(Group(children)),
	)
}

func Layout(ctx context.Context, title string, children ...Node) Node {
	return RootLayout(ctx, title,
		Class("flex flex-col min-h-screen"),
		AppHeader(ctx), Main(Class("flex-1 bg-slate-100 pt-24"), Group(children)), AppFooter(),
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
			Class("inline-flex items-center rounded-sm px-3 font-semibold text-primary transition hover:bg-slate-200 h-full"),
			Href(href), Text(text),
		),
	)
}

func LanguageSwitcher(ctx context.Context, baseUrl string, locale string) Node {
	localizer := ctx.Value("localizer").(*i18n.Localizer)

	langName := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "languages." + locale,
			Other: locale,
		},
	})
	title := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "locale_switcher.switch_to",
			Other: "Przełącz na {{ .Language }}",
		},
		TemplateData: map[string]string{
			"Language": langName,
		},
	})

	return Li(
		A(
			Class("inline-flex items-center rounded-sm px-3 font-semibold text-primary transition hover:bg-slate-200 h-full uppercase"),
			Href(baseUrl+"?lang="+locale), Text(locale), Title(title), Aria("label", title),
		),
	)
}

func AppFooter() Node {
	return Footer(Class("border-t-2 h-30"),
		Div(
			Class("container mx-auto text-center h-full flex items-center justify-center"),
			P(
				Raw("&copy; 2024&ndash;2025 by Wydawnictwo Homeo Sapiens. All rights reserved."),
			),
		),
	)
}
