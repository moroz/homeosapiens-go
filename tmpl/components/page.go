package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// The flat page vocabulary shared by the landing, blog and events pages: a
// full-width white canvas divided into bands by hairline rules, with no card
// chrome. Every band is a Section that supplies its own container, so pages
// built from these run edge to edge.

// Eyebrow is a small caps label preceded by a short rule in the brand gradient.
func Eyebrow(text string) Node {
	return Div(
		Class("flex items-center gap-3 text-sm font-semibold tracking-[0.1em] text-primary uppercase font-heading"),
		Span(Class("inline-block h-0.5 w-6 rounded-full brand-rule mobile:hidden")),
		Text(text),
	)
}

// PageSection is one full-width band. The bottom rule drops off automatically
// when it's the final band on a page.
func PageSection(children ...Node) Node {
	return El("section",
		Class("border-b border-slate-200 bg-white last:border-b-0"),
		Div(
			Class("container mx-auto py-16 px-2 mobile:px-5"),
			Group(children),
		),
	)
}

// PageHeader is the opening band of a page: eyebrow, H1, optional standfirst,
// and an optional trailing action aligned to the right on wide screens.
func PageHeader(eyebrow, title string, children ...Node) Node {
	return El("section",
		Class("border-b border-slate-200 bg-white"),
		Div(
			Class("container mx-auto px-2 pt-16 pb-10 mobile:pt-10 mobile:pb-8 mobile:px-5"),
			Eyebrow(eyebrow),
			H1(
				Class("mt-4 max-w-3xl text-4xl leading-[1.1] font-bold tracking-tight text-slate-900 mobile:text-3xl"),
				Text(title),
			),
			Group(children),
		),
	)
}

// Standfirst is the short intro paragraph that may follow a PageHeader title.
func Standfirst(text string) Node {
	return P(
		Class("mt-4 max-w-prose text-lg text-slate-600"),
		Text(text),
	)
}

// SectionHeading is the label/title pair that opens a band below the header,
// with optional actions pushed to the right.
func SectionHeading(label, title string, trailing ...Node) Node {
	return Div(
		Class("mb-8 flex items-end justify-between gap-4"),
		Div(
			Eyebrow(label),
			H2(Class("mt-3 text-2xl font-bold text-slate-900 lg:text-3xl"), Text(title)),
		),
		Group(trailing),
	)
}

// TextLink is the understated inline call to action used beside the primary
// buttons. The arrow is decoration, so it stays out of the message bundle.
func TextLink(href, label string) Node {
	return A(
		Href(href),
		Class("font-semibold text-primary hover:text-primary-hover font-heading"),
		Text(label+" →"),
	)
}

// RowList wraps hairline-separated rows. The rows carry their own bottom rule;
// this supplies the top one so the first row is closed off too.
func RowList(children ...Node) Node {
	return Div(
		Class("border-t border-slate-200 first:border-t-0"),
		Group(children),
	)
}

// Prose is the reading column for markdown bodies.
func Prose(children ...Node) Node {
	return Div(
		Class("prose lg:prose-lg max-w-prose"),
		Group(children),
	)
}

// EmptyState is the muted placeholder shown when a listing has no rows.
func EmptyState(text string) Node {
	return P(
		Class("py-16 text-center text-slate-500"),
		Text(text),
	)
}
