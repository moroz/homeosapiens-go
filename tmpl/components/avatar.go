package components

import (
	"github.com/moroz/homeosapiens-go/db/queries"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Avatar(user *queries.User) Node {
	initials := user.GivenName[0:1] + user.FamilyName[0:1]

	return Div(
		Class("inline-flex aspect-square h-10 items-center justify-center rounded-full bg-primary text-lg font-semibold text-white uppercase"),
		Text(initials),
	)
}
