package orders

import (
	"strconv"

	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Success(ctx *types.CustomContext, order *types.OrderDTO) Node {
	l := ctx.Localizer
	title := l.MustLocalizeMessage(&i18n.Message{
		ID: "orders.success.title",
	})
	body := l.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "orders.success.body_html",
		TemplateData: map[string]string{
			"Amount":      helpers.FormatPrice(order.GrandTotal, order.Currency, ctx.Language),
			"OrderNumber": strconv.Itoa(int(order.OrderNumber)),
			"Email":       order.Email.String(),
		},
	})

	return layout.Layout(ctx, title,
		Div(
			Class("card"),
			H2(Class("page-title mb-4"), Text(title)),
			Div(Class("prose"),
				Raw(body),
			),
		),
	)
}
