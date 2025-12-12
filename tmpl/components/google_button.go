package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func GoogleButton(text string) Node {
	return A(
		Class("button w-full outline flex gap-[10px] my-4"),
		Href("/oauth/google/redirect"),
		Img(Src("/assets/google.svg")),
		Text(text),
	)
}
