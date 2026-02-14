package events

import (
	"fmt"

	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func eventRegistrationIcon(href string) Node {
	return SVG(
		Class("h-5 w-5 fill-current"),
		Attr("viewBox", "0 0 640 640"),
		El("use",
			Href(href),
		),
	)
}

func eventRegistrationButton(event *services.EventDetailsDto) Node {
	if event.EventRegistration == nil {
		return Form(
			Action(fmt.Sprintf("/event_registrations/%s", event.ID)),
			Method("POST"),
			Button(
				Class("inline-flex button bg-slate-100 text-slate-900 gap-1 font-semibold hover:bg-slate-200"),
				Type("submit"),
				eventRegistrationIcon("/assets/circle-check-empty.svg#icon"),
				Text("Wezmę udział"),
			),
		)
	}

	return Form(
		Action(fmt.Sprintf("/event_registrations/%s", event.ID)),
		Method("POST"),
		Input(Type("hidden"), Name("_method"), Value("DELETE")),
		Button(
			Class("inline-flex button bg-blue-100 text-blue-700 gap-1 font-semibold hover:bg-blue-200"),
			Type("submit"),
			eventRegistrationIcon("/assets/circle-check-solid.svg#icon"),
			Text("Wezmę udział"),
		),
	)
}

func Show(ctx *types.CustomContext, event *services.EventDetailsDto) Node {
	lang := ctx.Language
	tz := ctx.Timezone

	title := event.TitleEn
	if lang == "pl" {
		title = event.TitlePl
	}

	description := event.DescriptionEn
	if lang == "pl" && event.DescriptionPl != nil && *event.DescriptionPl != "" {
		description = *event.DescriptionPl
	}

	l := ctx.Localizer

	return layout.Layout(ctx, event.TitleEn, Div(
		Class("card mx-auto"),
		H2(
			Class("text-primary text-2xl leading-normal font-bold"),
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
		EventLocationBadge(event.IsVirtual, event.Venue, l, lang),
		Div(Class("my-4 flex gap-6"),
			eventRegistrationButton(event),
		),
		helpers.MarkdownContent(description, "mt-4 w-full"),
	))
}
