package eventregistrations

import (
	"context"

	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"

	. "maragu.dev/gomponents/html"
)

func New(ctx context.Context, event *services.EventDetailsDto) Node {
	lang := ctx.Value("lang").(string)
	title := event.TitleEn
	if lang == "pl" {
		title = event.TitlePl
	}
	l := ctx.Value("localizer").(*i18n.Localizer)

	pageTitle := l.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "event_registrations.new.title",
		TemplateData: map[string]string{
			"Title": title,
		},
	})

	return layout.Layout(ctx, pageTitle,
		Div(
			Class("card"),
			Header(
				Span(
					Class("text-xl"),
					Text(l.MustLocalizeMessage(&i18n.Message{
						ID: "event_registrations.new.heading",
					})),
				),

				H2(
					Class("text-primary font-bold text-4xl"),
					Text(title),
				),
			),
		),
	)
}
