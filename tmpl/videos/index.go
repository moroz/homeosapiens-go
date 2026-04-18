package videos

import (
	"fmt"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Index(ctx *types.CustomContext, videos []*types.VideoGroupListDTO, group *types.VideoGroupDetailsDTO) Node {
	return layout.BareLayout(ctx, "Videos",
		Div(
			Class("mx-auto lg:w-7xl grid grid-cols-[1fr_3fr] gap-8 card"),
			Aside(
				H2(Class("page-title"), Text("Videos")),
				Nav(
					Ul(
						Class("video-group-menu"),
						Map(videos, func(vg *types.VideoGroupListDTO) Node {
							title := vg.TitleEn
							if ctx.Language == "pl" {
								title = vg.TitlePl
							}

							class := ""
							if group != nil && vg.ID == group.ID {
								class = "is-active"
							}

							return Li(
								Class(class),
								Text(title),
							)
						}),
					),
				),
			),
			Main(
				Iff(group != nil, func() Node {
					title := group.TitleEn
					if ctx.Language == "pl" {
						title = group.TitlePl
					}

					return Group{
						H3(Class("text-primary font-bold text-xl mb-4"), Text(title)),
						Ul(
							Map(group.Videos, func(video *queries.Video) Node {
								title := video.TitleEn
								if ctx.Language == "pl" {
									title = video.TitlePl
								}

								return Li(
									A(
										Href(fmt.Sprintf("/videos/%s/%s", group.Slug, video.Slug)),
										Text(title),
									),
								)
							}),
						),
					}
				}),
			),
		),
	)
}
