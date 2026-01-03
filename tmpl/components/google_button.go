package components

import (
	"net/url"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func googleRedirectPath(redirectTo string) string {
	base := "/oauth/google/redirect"
	if redirectTo == "" {
		return base
	}
	qs := url.Values{
		"ref": {redirectTo},
	}
	return base + "?" + qs.Encode()
}

func GoogleButton(text, redirectTo string) Node {
	return A(
		Class("button my-4 flex w-full gap-[10px] outline"),
		Href(googleRedirectPath(redirectTo)),
		Img(Src("/assets/google.svg")),
		Text(text),
	)
}
