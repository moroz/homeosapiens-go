package events

import (
	"fmt"

	twmerge "github.com/Oudwins/tailwind-merge-go"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func eventRegistrationButton(l *i18n.Localizer, event *services.EventDetailsDto) Node {
	classes := "bg-slate-100 text-slate-900 hover:bg-slate-200"
	if event.EventRegistration != nil {
		classes = "bg-primary/10 text-primary hover:bg-primary/20"
	}

	return Form(
		Action(fmt.Sprintf("/event_registrations/%s", event.ID)),
		Method("POST"),
		If(event.EventRegistration != nil, Input(Type("hidden"), Name("_method"), Value("DELETE"))),
		Button(
			Class(
				twmerge.Merge("inline-flex button bg-slate-100 text-slate-900 gap-1 font-semibold hover:bg-slate-200", classes),
			),
			Type("submit"),
			EventAttendanceIcon(event.EventRegistration != nil),
			Text(l.MustLocalizeMessage(&i18n.Message{
				ID: "common.events.attendance_badge",
			})),
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
		Div(Class("my-4 flex gap-4 items-center"),
			eventRegistrationButton(l, event),
			If(
				event.RegistrationCount > 0,
				Text(l.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID: "common.events.attendance_count",
					},
					TemplateData: map[string]any{
						"Count": event.RegistrationCount,
					},
					PluralCount: event.RegistrationCount,
				})),
			),
			If(event.RegistrationCount == 0, Text(l.MustLocalizeMessage(&i18n.Message{ID: "common.events.nobody_attending"}))),
		),
		helpers.MarkdownContent(description, "mt-4 w-full"),
	))
}
