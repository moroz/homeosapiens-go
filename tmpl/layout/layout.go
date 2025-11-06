package layout

import (
	"fmt"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Layout(title string, children ...Node) Node {
	return HTML(
		Head(
			Meta(Charset("UTF-8")),
			Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
			Title(title),
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
			AppHeader(), Main(Class("flex-1"), Group(children)), AppFooter(),
		),
	)
}

func NavLink(href string, text string) Node {
	return Li(
		A(
			Class("inline-flex items-center rounded-sm px-3 font-semibold text-primary transition hover:bg-slate-200 h-full"),
			Href(href), Text(text),
		),
	)
}

func LanguageSwitcher(baseUrl string, locale string) Node {
	return NavLink(fmt.Sprintf("%s?lang=%s", baseUrl, locale), locale)
}

func AppHeader() Node {
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
					LanguageSwitcher("/", "pl"),
					NavLink("/events", "Events"),
					NavLink("/videos", "Watch"),
					NavLink("/dashboard", "My products"),
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
