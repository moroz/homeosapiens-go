package blog

import (
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Show(ctx *types.CustomContext, post *queries.BlogPost) Node {
	l := ctx.Localizer

	return layout.PageLayout(ctx, post.Title,
		components.PageHeader(
			l.MustLocalizeMessage(&i18n.Message{ID: "blog.index.eyebrow"}),
			post.Title,
			Iff(post.PublishedAt != nil, func() Node {
				return P(
					Class("mt-4 text-sm text-slate-500"),
					Text(l.MustLocalizeMessage(&i18n.Message{ID: "blog.show.published_at"})+" "),
					Time(Text(helpers.FormatDateTime(*post.PublishedAt, ctx.Timezone, ctx.Language))),
				)
			}),
		),
		components.LastPageSection(
			components.Prose(
				Iff(post.Body != nil, func() Node {
					return helpers.RenderMarkdown(*post.Body)
				}),
			),
			Div(
				Class("mt-12 border-t border-slate-200 pt-6"),
				components.TextLink("/blog", l.MustLocalizeMessage(&i18n.Message{ID: "blog.show.back"})),
			),
		),
	)
}
