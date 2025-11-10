package layout

import (
	"context"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Layout(ctx context.Context, title string, children ...Node) Node {
	return HTML(
		Head(
			Meta(Charset("UTF-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
			TitleEl(Text(title+" | Homeo sapiens")),
			Link(Rel("stylesheet"), Href("/assets/bundle.css")),
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
				Href("https://fonts.googleapis.com/css2?family=DM+Sans:ital,opsz,wght@0,9..40,100..1000;1,9..40,100..1000&display=swap"),
				Rel("stylesheet"),
			),
		),
		Body(
			Class("flex flex-col min-h-screen"),
			AppHeader(ctx), Main(Class("flex-1 bg-slate-100"), Group(children)), AppFooter(),
		),
	)
}

func NavLink(href string, text string, children ...Node) Node {
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

func AppHeader(ctx context.Context) Node {
	localizer := ctx.Value("localizer").(*i18n.Localizer)
	lang := ctx.Value("lang").(string)
	otherLocale := "pl"
	if lang == "pl" {
		otherLocale = "en"
	}

	return Header(
		Class("h-20 border-b-2"),
		Div(Class("container mx-auto flex justify-between h-full items-center"),
			H1(
				A(
					Class("text-primary no-underline transition-colors hover:text-primary-800 text-4xl font-bold"),
					Href("/"), Text("Homeo sapiens"),
				),
			),
			Nav(
				Class("h-full"),
				Ul(
					Class("flex gap-2 py-4 h-full"),
					LanguageSwitcher(ctx, "/", otherLocale),
					NavLink("/events", localizer.MustLocalize(&i18n.LocalizeConfig{
						DefaultMessage: &i18n.Message{
							ID:    "header.nav.events",
							Other: "Events",
						},
					})),
					NavLink("/videos", localizer.MustLocalize(&i18n.LocalizeConfig{
						DefaultMessage: &i18n.Message{
							ID:    "header.nav.videos",
							Other: "Watch",
						},
					})),
					NavLink("/dashboard", localizer.MustLocalize(&i18n.LocalizeConfig{
						DefaultMessage: &i18n.Message{
							ID:    "header.nav.my_products",
							Other: "My Products",
						},
					})),
				),
			),
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
