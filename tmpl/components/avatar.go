package components

import (
	"github.com/moroz/homeosapiens-go/db/queries"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Avatar(user *queries.User) Node {
	if user.ProfilePicture != nil {
		return Div(
			Class("inline-flex aspect-square h-10 overflow-hidden rounded-full"),
			Img(Src(*user.ProfilePicture)),
		)
	}

	initials := user.GivenName.Plaintext()[0:1] + user.FamilyName.Plaintext()[0:1]

	return Div(
		Class("inline-flex aspect-square h-10 items-center justify-center rounded-full bg-primary text-lg font-semibold text-white uppercase"),
		Text(initials),
	)
}
