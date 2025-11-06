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
				Class("flex justify-between rounded-lg border-2 bg-slate-100 p-6"),
				Header(
					Span(
						Class("mb-2 inline-flex items-center gap-1 rounded border-2 border-black bg-white px-2 py-1 text-sm font-semibold text-primary"),
						If(e.IsVirtual,
							Span(Text("+ Online")),
						),
					),

					H3(
						Class("text-4xl font-bold text-primary"),
						Text(e.TitleEn),
					),
				),
			)
		}),
	))
}
