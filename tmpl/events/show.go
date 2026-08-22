package events

import (
	"fmt"

	twmerge "github.com/Oudwins/tailwind-merge-go"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

const EventRegistrationButtonBaseClasses = "inline-flex button bg-slate-100 text-slate-900 gap-1 font-semibold hover:bg-slate-200"
const EventRegistrationButtonNotGoingClasses = "bg-slate-100 text-slate-900 hover:bg-slate-200"
const EventRegistrationButtonGoingClasses = "bg-primary/10 text-primary hover:bg-primary/20"

// ananymousEventRegistrationButtonLink displays a faux button that redirects the user to the login page if they want to register for an event. If they are already signed in, the GET handler registers the user for the event.
func ananymousEventRegistrationButtonLink(l *i18n.Localizer, event *types.EventDetailsDto) Node {
	return A(
		Href(fmt.Sprintf("/events/%s/register", event.ID)),
		Class(twmerge.Merge(EventRegistrationButtonBaseClasses, EventRegistrationButtonNotGoingClasses)),
		EventAttendanceIcon(false),
		Text(l.MustLocalizeMessage(&i18n.Message{
			ID: "common.events.attendance_badge",
		})),
	)
}

type freeEventCTAProps struct {
	Event     *types.EventDetailsDto
	User      *queries.User
	Localizer *i18n.Localizer
}

// endedEventBadge replaces the registration CTA once the event is over: the
// sign-up endpoint rejects such events, so offering the button would only lead
// to a 404.
func endedEventBadge(l *i18n.Localizer) Node {
	return Span(
		Class("inline-flex items-center justify-center gap-1 rounded-sm border border-slate-300 bg-slate-100 px-2 py-1 text-sm font-semibold text-slate-500"),
		Text(l.MustLocalizeMessage(&i18n.Message{
			ID: "common.events.ended",
		})),
	)
}

// joinMeetingButton links to the virtual meeting. It is only rendered for users
// who are registered for the event, so the link is not handed out publicly.
func joinMeetingButton(l *i18n.Localizer, event *types.EventDetailsDto) Node {
	if event.EventRegistration == nil || event.MeetingUrl == nil || event.HasEnded() {
		return nil
	}

	return A(
		Href(*event.MeetingUrl),
		Target("_blank"),
		Rel("noopener noreferrer"),
		Class("button"),
		Text(l.MustLocalizeMessage(&i18n.Message{
			ID: "events.join_meeting",
		})),
	)
}

func freeEventCTA(props *freeEventCTAProps) Node {
	if props.Event.HasEnded() {
		return endedEventBadge(props.Localizer)
	}

	if props.User == nil {
		return ananymousEventRegistrationButtonLink(props.Localizer, props.Event)
	}

	event := props.Event
	l := props.Localizer

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

func Show(ctx *types.CustomContext, event *types.EventDetailsDto) Node {
	lang := ctx.Language
	tz := ctx.Timezone

	title := event.TitleEn
	if lang == "pl" {
		title = event.TitlePl
	}

	description := event.DescriptionEn
	if lang == "pl" {
		description = event.DescriptionPl
	}

	subtitle := event.SubtitleEn
	if lang == "pl" {
		subtitle = event.SubtitlePl
	}

	l := ctx.Localizer

	return layout.PageLayout(ctx, event.TitleEn,
		components.PageHeader(
			helpers.TranslateEventType(l, event.EventType),
			title,
			Iff(subtitle != nil && *subtitle != "", func() Node {
				return P(Class("mt-4 max-w-prose text-xl text-slate-600"), Text(*subtitle))
			}),
			Div(
				Class("mt-6 grid gap-1 text-slate-700"),
				P(
					Strong(Class("font-fallback"), Text(l.MustLocalizeMessage(&i18n.Message{
						ID: "events.starts_at",
					}))),
					Text(" "),
					Time(
						Text(helpers.FormatDateTime(event.StartsAt, tz, lang)),
					),
				),
				P(
					Strong(Class("font-fallback"), Text(l.MustLocalizeMessage(&i18n.Message{
						ID: "events.ends_at",
					}))),
					Text(" "),
					Time(
						Text(helpers.FormatDateTime(event.EndsAt, tz, lang)),
					),
				),
			),
		),
		components.PageSection(
			Div(Class("flex flex-wrap items-center gap-4"),
				Iff(event.IsFree(), func() Node {
					return freeEventCTA(&freeEventCTAProps{
						Event:     event,
						User:      ctx.User,
						Localizer: ctx.Localizer,
					})
				}),
				// Attendance to a past event is worthless, so it is neither
				// registrable nor purchasable once the event is over.
				If(!event.IsFree() && !event.HasEnded(), components.AddToCartButton(ctx.Localizer, event.Event.ID, event.CountInCart)),
				If(!event.IsFree() && event.HasEnded(), endedEventBadge(l)),
				joinMeetingButton(l, event),
				If(!event.IsFree() && event.CountInCart > 0, A(Href("/cart"), Class("font-semibold underline"), Text(l.MustLocalizeMessage(&i18n.Message{ID: "common.events.view_cart"})))),
				If(
					event.IsFree() && !event.HasEnded() && event.RegistrationCount > 0,
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
				If(event.IsFree() && event.HasEnded() && event.RegistrationCount > 0,
					Text(l.MustLocalize(&i18n.LocalizeConfig{
						DefaultMessage: &i18n.Message{
							ID: "common.events.past_event_attendance_count",
						},
						TemplateData: map[string]any{
							"Count": event.RegistrationCount,
						},
						PluralCount: event.RegistrationCount,
					})),
				),
				If(event.IsFree() && !event.HasEnded() && event.RegistrationCount == 0, Text(l.MustLocalizeMessage(&i18n.Message{ID: "common.events.nobody_attending"}))),
			),
		),
		components.LastPageSection(
			components.Prose(
				Iff(description != nil, func() Node {
					return helpers.RenderMarkdown(*description)
				}),
			),
			Div(
				Class("mt-12 border-t border-slate-200 pt-6"),
				components.TextLink("/events", l.MustLocalizeMessage(&i18n.Message{ID: "events.show.back"})),
			),
		),
	)
}
