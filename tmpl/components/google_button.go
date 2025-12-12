package components

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func GoogleButton(text string) Node {
	return A(
		Class("button my-4 flex w-full gap-[10px] outline"),
		Href("/oauth/google/redirect"),
		Img(Src("/assets/google.svg")),
		Text(text),
	)
}
