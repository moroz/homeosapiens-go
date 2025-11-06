package pages

import (
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Index(events []*queries.Event) Node {
	return layout.Layout("Events", Div(
		Class("grid gap-4 container mx-auto"),
		Map(events, func(e *queries.Event) Node {
			return Article(
				Class("event bg-[salmon]"),
				Span(Text(e.TitleEn)),
			)
		}),
	))
}
