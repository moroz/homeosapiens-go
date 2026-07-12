package videos

import (
	"fmt"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Youtube(ctx *types.CustomContext, videos []*queries.Video) Node {
	return layout.BareLayout(ctx, "Watch",
		Div(
			Class("container mx-auto py-6"),
			H2(Class("page-title"), Text("Watch")),
			Div(
				Class("video-grid"),
				Map(videos, func(video *queries.Video) Node {
					return Article(
						Class("video-card"),
						Header(
							Class("video-card-thumb"),
							Img(Src(fmt.Sprintf("https://img.youtube.com/vi/%s/sddefault.jpg", *video.YoutubeID))),
							DurationBadge(video.DurationSeconds),
						),
						Footer(
							Class("video-card-body"),
							H4(Class("video-card-title text-lg"), Text(video.TitleEn)),
							Iff(video.DescriptionEn != nil, func() Node {
								return P(Text(*video.DescriptionEn))
							}),
						),
					)
				}),
			),
		),
	)
}
