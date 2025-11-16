package videos

import (
	"context"
	"fmt"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Index(ctx context.Context, videos []*services.VideoListDto) Node {
	lang := ctx.Value("lang").(string)

	return layout.Layout(ctx, "Videos",
		H2(Class("text-3xl font-bold text-primary"), Text("Videos")),
		Map(videos, func(video *services.VideoListDto) Node {
			title := video.TitleEn
			if lang == "pl" {
				title = video.TitlePl
			}

			id := "video-" + video.ID.String()

			return Article(
				H3(Text(title)),
				Div(
					Class("video-js"),
					Video(
						ID(id),
						Map(video.Sources, func(source *queries.VideoSource) Node {
							t := source.ContentType
							if source.Codec != nil {
								t = fmt.Sprintf(`%s; codecs="%s"`, t, *source.Codec)
							}

							return Source(Src(config.AssetCdnBaseUrl+source.ObjectKey), Type(t))
						}),
					),
				),
			)
		}),
	)
}
