package users

import (
	"fmt"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Index(ctx *types.CustomContext, users []*queries.User) Node {
	return layout.AdminLayout(ctx, "Users",
		Table(
			Class("index-table"),
			THead(
				Tr(
					Th(Text("Given name")),
					Th(Text("Family name")),
					Th(Text("Email")),
					Th(Text("Country")),
				),
			),
			TBody(
				Map(users, func(u *queries.User) Node {
					return Tr(
						Data("url", fmt.Sprintf("/admin/users/%s", u.ID)),
						Td(Text(u.GivenName.String())),
						Td(Text(u.FamilyName.String())),
						Td(Text(u.Email.String())),
						Td(Text(helpers.DerefOrEmpty(u.Country))),
					)
				}),
			),
		),
	)
}
