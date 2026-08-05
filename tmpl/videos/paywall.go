package videos

import (
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// PaidBadge marks a group the visitor has not bought yet in the series menu.
func PaidBadge(l *i18n.Localizer) Node {
	return Span(
		Class("ml-2 inline-flex items-center rounded-sm border border-slate-300 bg-slate-100 px-1.5 py-0.5 text-xs font-semibold text-slate-500"),
		Text(l.MustLocalizeMessage(&i18n.Message{ID: "videos.paywall.badge"})),
	)
}

// LockedPanel stands in for content the visitor has not bought: the video player
// on a video page, and the video grid on a series page. It carries the price and
// an add-to-cart button whenever the group's product is priced.
func LockedPanel(ctx *types.CustomContext, group *types.VideoGroupDetailsDTO) Node {
	l := ctx.Localizer

	return Div(
		Class("flex flex-col items-center gap-3 rounded-sm border border-slate-300 bg-slate-50 px-6 py-12 text-center"),
		H4(
			Class("text-xl font-bold text-primary"),
			Text(l.MustLocalizeMessage(&i18n.Message{ID: "videos.paywall.title"})),
		),
		P(
			Class("max-w-prose text-slate-600"),
			Text(l.MustLocalizeMessage(&i18n.Message{ID: "videos.paywall.body"})),
		),
		Iff(group.Price != nil && group.Currency != nil, func() Node {
			return P(
				Class("font-semibold"),
				Text(l.MustLocalizeMessage(&i18n.Message{ID: "videos.paywall.price"})),
				Text(" "),
				Text(helpers.FormatPrice(*group.Price, *group.Currency, ctx.Language)),
			)
		}),
		Iff(group.ProductID != nil, func() Node {
			return components.AddProductToCartButton(l, *group.ProductID, group.CountInCart)
		}),
	)
}
