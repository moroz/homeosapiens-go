package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Logo(class string) Node {
	return SVG(
		Attr("viewBox", "0 0 1538 361"),
		Class(class),
		El("use",
			Href("/assets/logo.svg#logo"),
		),
	)
}
