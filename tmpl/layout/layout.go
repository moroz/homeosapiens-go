package layout

import (
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/components/icons"
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
			Meta(Name("robots"), Content("noindex, nofollow")),
			TitleEl(Text(title+" | Homeo sapiens")),
			AssetEntryPoint(ctx),
			// fonts(),
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
			Class("flex-1 bg-slate-100 pt-20 pb-6"),
			Div(
				Class("container mx-auto"),
				components.Flash(ctx.Flash),
				Group(children),
			),
		),
		AppFooter(),
	)
}

// PageLayout is the flat, full-width white shell used by the landing, blog and
// events pages. Unlike Layout it adds no container of its own: each child is a
// full-bleed band that supplies its own container, so pages can run edge to
// edge and separate themselves with hairline rules instead of cards.
func PageLayout(ctx *types.CustomContext, title string, children ...Node) Node {
	return RootLayout(ctx, title,
		Class("flex min-h-screen max-w-full flex-col overflow-x-hidden"),
		AppHeader(ctx),
		Main(
			Class("flex-1 bg-white pt-20"),
			Iff(len(ctx.Flash) > 0, func() Node {
				return Div(Class("container mx-auto px-6 pt-6"), components.Flash(ctx.Flash))
			}),
			Group(children),
		),
		AppFooter(),
	)
}

func BareLayout(ctx *types.CustomContext, title string, children ...Node) Node {
	return RootLayout(ctx, title,
		Class("flex min-h-screen max-w-full flex-col overflow-x-hidden"),
		AppHeader(ctx),
		Main(
			Class("flex-1 bg-slate-100 pt-20 pb-6"),
			Group(children),
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
			Href("https://fonts.googleapis.com/css2?family=DM+Sans:ital,opsz,wght@0,9..40,100..1000;1,9..40,100..1000&display=swap"),
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

const (
	YoutubeUrl   = "https://www.youtube.com/@Homeosapiens-p7z"
	InstagramUrl = "https://www.instagram.com/homeosapiens.eu"
)

// socialLink is an icon-only link to one of our social media profiles. The
// icons are monochrome so that they take the colour of the surrounding text.
func socialLink(href, label string, icon Node) Node {
	return A(
		Class("text-slate-500 transition-colors hover:text-primary"),
		Href(href), Target("_blank"), Rel("noopener noreferrer"),
		Title(label), Aria("label", label),
		icon,
	)
}

func AppFooter() Node {
	return Footer(Class("relative bg-white py-8 text-sm text-slate-600 shadow lg:text-base"),
		Div(Class("absolute inset-x-0 top-0 h-0.5 brand-rule")),
		Div(
			Class("container mx-auto flex flex-col items-center justify-center gap-4 text-center"),
			Ul(
				Class("flex items-center gap-5"),
				Li(socialLink(YoutubeUrl, "YouTube", icons.Icon(&icons.IconProps{
					Name: "youtube-mono", ViewBox: "0 0 576 512", Classes: "h-6 w-6",
				}))),
				Li(socialLink(InstagramUrl, "Instagram", icons.Icon(&icons.IconProps{
					Name: "instagram", ViewBox: "0 0 448 512", Classes: "h-6 w-6",
				}))),
			),
			P(
				Raw("&copy; 2024&ndash;2026 by Wydawnictwo Homeo Sapiens.<br/>All rights reserved."),
			),
		),
	)
}
