package events

import (
	"context"
	"fmt"
	"time"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Show(ctx context.Context, event *services.EventDetailsDto) Node {
	lang := ctx.Value(config.LangContextName).(string)
	tz := ctx.Value("timezone").(*time.Location)

	title := event.TitleEn
	if lang == "pl" {
		title = event.TitlePl
	}

	l := ctx.Value("localizer").(*i18n.Localizer)

	return layout.Layout(ctx, event.TitleEn, Div(
		Class("card"),
		EventLocationBadge(event.IsVirtual, event.Venue, l, lang),
		H2(
			Class("text-primary my-2 text-4xl font-bold"),
			Text(title),
		),
		Div(
			Class("grid"),
			P(
				Strong(Class("font-fallback"), Text(l.MustLocalizeMessage(&i18n.Message{
					ID: "events.starts_at",
				}))),
				Text(" "),
				Time(
					Text(helpers.FormatDateTime(event.StartsAt.Time, tz, lang)),
				),
			),
			P(
				Strong(Class("font-fallback"), Text(l.MustLocalizeMessage(&i18n.Message{
					ID: "events.ends_at",
				}))),
				Text(" "),
				Time(
					Text(helpers.FormatDateTime(event.EndsAt.Time, tz, lang)),
				),
			),
		),
		If(event.EventRegistration == nil, A(Href(fmt.Sprintf("/events/%s/register", event.Slug)), Text("Register"))),
		Iff(event.EventRegistration != nil, func() Node {
			return Text("You are going")
		}),
	))
}
