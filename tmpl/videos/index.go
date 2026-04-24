package videos

import (
	"fmt"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/components"
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
								A(
									Class(class),
									Href(fmt.Sprintf("/videos/%s", vg.Slug)),
									Text(title),
								),
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
						Div(
							Class("video-grid"),
							Map(group.Videos, func(video *queries.ListVideosForVideoGroupRow) Node {
								title := video.TitleEn
								if ctx.Language == "pl" {
									title = video.TitlePl
								}

								return A(
									Class("no-underline"),
									Href(fmt.Sprintf("/videos/%s/%s", group.Slug, video.Slug)),
									Article(
										Class("card video p-0"),
										Header(
											Class("aspect-video bg-slate-200 p-8 flex items-center justify-center"),
											components.Logo("text-slate-400"),
										),
										Footer(
											Class("p-4"),
											H4(
												Class("font-semibold"),
												Text(title)),
										),
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
