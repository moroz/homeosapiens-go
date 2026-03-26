package components

import (
	"net/url"

	twmerge "github.com/Oudwins/tailwind-merge-go"
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

func GoogleButton(text, redirectTo string, classes ...string) Node {
	class := twmerge.Merge("button font-fallback my-4 inline-flex gap-2.5 outline", twmerge.Merge(classes...))

	return A(
		Class(class),
		Href(googleRedirectPath(redirectTo)),
		Img(Src("/assets/google.svg")),
		Text(text),
	)
}
