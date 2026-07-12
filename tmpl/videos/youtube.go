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
					title := video.TitleEn
					if ctx.Language == "pl" {
						title = video.TitlePl
					}

					desc := video.DescriptionEn
					if ctx.Language == "pl" {
						desc = video.DescriptionPl
					}

					url := fmt.Sprintf("https://youtube.com/watch?v=%s", *video.YoutubeID)

					return Article(
						Class("video-card"),
						A(
							Class("video-card-thumb"),
							Href(url), Target("_blank"), Rel("noopener noreferrer"),
							Img(Src(fmt.Sprintf("https://img.youtube.com/vi/%s/sddefault.jpg", *video.YoutubeID))),
							DurationBadge(video.DurationSeconds),
						),
						Footer(
							Class("video-card-body"),
							H4(Class("video-card-title text-lg"), Text(title)),
							Iff(desc != nil, func() Node {
								return P(Text(*desc))
							}),
							A(Href(url), Target("_blank"), Rel("noopener noreferrer"), Text("Watch on YouTube")),
						),
					)
				}),
			),
		),
	)
}
