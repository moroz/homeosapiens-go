package layout

import (
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func navLink(children ...Node) Node {
	return Li(
		Class("-mx-3"),
		A(
			Class("flex h-10 w-full hover:bg-white/10 items-center px-3"),
			Group(children),
		),
	)
}

func AdminLayout(ctx *types.CustomContext, title string, children ...Node) Node {
	return RootLayout(ctx, title,
		Div(
			Class("grid h-screen grid-cols-[300px_1fr]"),
			Aside(
				Class("bg-slate-700 text-white p-6 gap-6 flex flex-col"),
				Header(
					H1(Class("font-bold text-2xl"), Text("Homeo sapiens")),
					H2(Text("Admin panel")),
				),
				Nav(
					Ul(
						Class("grid gap-2"),
						navLink(Href("/admin"), Text("Events")),
						navLink(Href("/admin/users"), Text("Users")),
						navLink(Href("/admin/blog"), Text("Blog")),
					),
				),
				Footer(
					Class("mt-auto"),
					A(Class("button w-full"), Href("/sign-out"), Text("Sign out")),
				),
			),
			Main(
				Class("p-6"),
				Group(children),
			),
		),
	)
}
