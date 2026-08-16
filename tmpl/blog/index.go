package blog

import (
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

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
	))
}
