package blog

import (
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Show(ctx *types.CustomContext, post *queries.BlogPost) Node {
	l := ctx.Localizer

	return layout.Layout(ctx, post.Title, Div(
		Class("card mx-auto"),
		H2(
			Class("text-2xl leading-normal font-bold text-primary"),
			Text(post.Title),
		),
		If(post.PublishedAt != nil, P(
			Strong(Class("font-fallback"), Text(l.MustLocalizeMessage(&i18n.Message{
				ID: "blog.show.published_at",
			}))),
			Text(" "),
			Time(
				Text(helpers.FormatDateTime(*post.PublishedAt, ctx.Timezone, ctx.Language)),
			),
		)),
		Div(
			Class("prose lg:prose-lg mt-4 w-full"),
			Iff(post.Body != nil, func() Node {
				return helpers.RenderMarkdown(*post.Body)
			}),
		),
	))
}
