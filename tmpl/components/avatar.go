package components

import (
	"github.com/moroz/homeosapiens-go/db/queries"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Avatar(user *queries.User) Node {
	if user.ProfilePicture != nil {
		return Div(
			Class("inline-flex aspect-square h-10 rounded-full overflow-hidden"),
			Img(Src(*user.ProfilePicture)),
		)
	}

	initials := user.GivenName[0:1] + user.FamilyName[0:1]

	return Div(
		Class("bg-primary inline-flex aspect-square h-10 items-center justify-center rounded-full text-lg font-semibold text-white uppercase"),
		Text(initials),
	)
}
