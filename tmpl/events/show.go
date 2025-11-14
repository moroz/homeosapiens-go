package events

import (
	"context"

	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Show(ctx context.Context, event *services.EventDetailsDto) Node {
	lang := ctx.Value("lang").(string)

	title := event.TitleEn
	if lang == "pl" {
		title = event.TitlePl
	}

	return layout.Layout(ctx, event.TitleEn, Div(
		Class("card"),
		H2(
			Class("text-primary text-4xl font-bold"),
			Text(title),
		),
		Div(
			Class("grid"),
			Time(
				Text(helpers.FormatDateTime(event.StartsAt.Time, lang)),
			),
			Time(
				Text(helpers.FormatDateTime(event.EndsAt.Time, lang)),
			),
		),
	))
}
