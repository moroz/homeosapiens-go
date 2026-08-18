package blog

import (
	"fmt"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func postListItem(ctx *types.CustomContext, post *queries.BlogPost) Node {
	return Li(
		A(
			Href(fmt.Sprintf("/blog/%s", post.Slug)),
			Class("block card"),
			H3(
				Class("text-xl font-bold text-primary"),
				Text(post.Title),
			),
			If(post.PublishedAt != nil, P(
				Class("text-sm text-slate-500"),
				Text(helpers.FormatDateTime(*post.PublishedAt, ctx.Timezone, ctx.Language)),
			)),
		),
	)
}

func Index(ctx *types.CustomContext, posts []*queries.BlogPost) Node {
	l := ctx.Localizer
	title := l.MustLocalizeMessage(&i18n.Message{
		ID: "blog.index.title",
	})

	return layout.Layout(ctx, title, Div(
		H2(
			Class("page-title text-center mt-4"),
			Text(title),
		),
		If(len(posts) == 0, Div(
			Class("text-center p-4"),
			Text(l.MustLocalizeMessage(&i18n.Message{
				ID: "blog.index.no_posts",
			}))),
		),
		If(len(posts) > 0, Ul(
			Class("grid gap-4 mt-4"),
			Map(posts, func(post *queries.BlogPost) Node {
				return postListItem(ctx, post)
			}),
		)),
	))
}
