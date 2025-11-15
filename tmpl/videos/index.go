package videos

import (
	"context"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Index(ctx context.Context, videos []*queries.Video) Node {
	lang := ctx.Value("lang").(string)

	return layout.Layout(ctx, "Videos",
		H2(Class("text-3xl font-bold text-primary"), Text("Videos")),
		Map(videos, func(video *queries.Video) Node {
			title := video.TitleEn
			if lang == "pl" {
				title = video.TitlePl
			}

			id := "video-" + video.ID.String()

			return Article(
				H3(Text(title)),
				Div(
					Class("video-js w-full"),
					Video(
						ID(id),
						Source(Src(config.AssetCdnBaseUrl+video.ObjectKey), Type(`video/mp4; codecs="hev1"`)),
					),
				),
			)
		}),
	)
}
