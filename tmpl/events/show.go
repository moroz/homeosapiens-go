package events

import (
	"fmt"

	twmerge "github.com/Oudwins/tailwind-merge-go"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/shopspring/decimal"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

const EventRegistrationButtonBaseClasses = "inline-flex button bg-slate-100 text-slate-900 gap-1 font-semibold hover:bg-slate-200"
const EventRegistrationButtonNotGoingClasses = "bg-slate-100 text-slate-900 hover:bg-slate-200"
const EventRegistrationButtonGoingClasses = "bg-primary/10 text-primary hover:bg-primary/20"

func eventRegistrationSignInLink(l *i18n.Localizer, event *services.EventDetailsDto) Node {
	return A(
		Href(fmt.Sprintf("/events/%s/register", event.ID)),
		Class(twmerge.Merge(EventRegistrationButtonBaseClasses, EventRegistrationButtonNotGoingClasses)),
		EventAttendanceIcon(false),
		Text(l.MustLocalizeMessage(&i18n.Message{
			ID: "common.events.attendance_badge",
		})),
	)
}

func eventRegistrationButton(l *i18n.Localizer, event *services.EventDetailsDto) Node {
	classes := EventRegistrationButtonNotGoingClasses
	if event.EventRegistration != nil {
		classes = EventRegistrationButtonGoingClasses
	}

	return Form(
		Action(fmt.Sprintf("/event_registrations/%s", event.ID)),
		Method("POST"),
		If(event.EventRegistration != nil, Input(Type("hidden"), Name("_method"), Value("DELETE"))),
		Button(
			Class(
				twmerge.Merge(EventRegistrationButtonBaseClasses, classes),
			),
			Type("submit"),
			EventAttendanceIcon(event.EventRegistration != nil),
			Text(l.MustLocalizeMessage(&i18n.Message{
				ID: "common.events.attendance_badge",
			})),
		),
	)
}

func addToCartButton(l *i18n.Localizer, event *queries.Event) Node {
	return Form(
		Action(fmt.Sprintf("/cart_line_items/%s", event.ID)),
		Method("POST"),
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
	isFree := event.BasePriceAmount == nil || event.BasePriceAmount.Equal(decimal.Zero)

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
		Div(Class("my-4 flex gap-4 items-center"),
			If(isFree && event.EventRegistration == nil, eventRegistrationSignInLink(l, event)),
			If(isFree && event.EventRegistration != nil, eventRegistrationButton(l, event)),
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
