package events

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func MarkdownContent(content string, classes ...string) Node {
	extensions := parser.CommonExtensions | parser.Autolink | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(content))

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	innerHTML := markdown.Render(doc, renderer)

	class := strings.Join(append([]string{"prose lg:prose-lg max-w-[unset] "}, classes...), " ")

	return Div(Class(class), Raw(string(innerHTML)))
}

func Show(ctx context.Context, event *services.EventDetailsDto) Node {
	lang := ctx.Value(config.LangContextName).(string)
	tz := ctx.Value("timezone").(*time.Location)
	isFuture := event.StartsAt.Time.After(time.Now())

	title := event.TitleEn
	if lang == "pl" {
		title = event.TitlePl
	}

	description := event.DescriptionEn
	if lang == "pl" && event.DescriptionPl != nil && *event.DescriptionPl != "" {
		description = *event.DescriptionPl
	}

	l := ctx.Value("localizer").(*i18n.Localizer)

	return layout.Layout(ctx, event.TitleEn, Div(
		Class("card"),
		Div(
			Class("mb-2 flex items-center gap-2"),
			EventLocationBadge(event.IsVirtual, event.Venue, l, lang),
			If(event.EventRegistration != nil, EventAttendanceBadge(isFuture, l)),
		),
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
		MarkdownContent(description, "mt-4"),
	))
}
