package blog

import (
	"fmt"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func postListItem(ctx *types.CustomContext, post *queries.BlogPost) Node {
	return A(
		Href(fmt.Sprintf("/blog/%s", post.Slug)),
		Class("group grid grid-cols-[9rem_1fr] items-baseline gap-6 border-b border-slate-200 py-6 no-underline transition-colors hover:bg-brand-50 mobile:grid-cols-1 mobile:gap-1 mobile:py-5"),
		Span(
			Class("text-sm text-slate-500 tabular-nums"),
			Iff(post.PublishedAt != nil, func() Node {
				return Text(helpers.FormatDateTime(*post.PublishedAt, ctx.Timezone, ctx.Language))
			}),
		),
		H2(
			Class("text-xl font-bold text-primary transition-colors group-hover:text-primary-hover"),
			Text(post.Title),
		),
	)
}

func Index(ctx *types.CustomContext, posts []*queries.BlogPost) Node {
	l := ctx.Localizer
	title := l.MustLocalizeMessage(&i18n.Message{ID: "blog.index.title"})

	return layout.PageLayout(ctx, title,
		components.PageHeader(
			l.MustLocalizeMessage(&i18n.Message{ID: "blog.index.eyebrow"}),
			title,
			components.Standfirst(l.MustLocalizeMessage(&i18n.Message{ID: "blog.index.standfirst"})),
		),
		components.LastPageSection(
			If(len(posts) == 0,
				components.EmptyState(l.MustLocalizeMessage(&i18n.Message{ID: "blog.index.no_posts"}))),
			If(len(posts) > 0, components.RowList(
				Map(posts, func(post *queries.BlogPost) Node {
					return postListItem(ctx, post)
				}),
			)),
		),
	)
}
