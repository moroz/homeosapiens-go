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
		H2(Class("text-primary text-4xl font-bold"), Text("Videos")),
		Div(Class("mt-6 mb-12 grid gap-6"),
			Map(videos, func(video *services.VideoListDto) Node {
				title := video.TitleEn
				if lang == "pl" {
					title = video.TitlePl
				}

				id := "video-" + video.ID.String()

				return Article(
					H3(Class("text-2xl font-bold"), Text(title)),
					Div(
						Video(
							Class("video-js vjs-theme-fantasy"),
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
		),
	)
}
