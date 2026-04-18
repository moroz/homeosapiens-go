package videos

import (
	"fmt"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Show(ctx *types.CustomContext, video *types.VideoDetailsDTO) Node {
	title := video.TitleEn
	if ctx.Language == "pl" {
		title = video.TitlePl
	}

	return layout.Layout(ctx, title,
		Div(Class("card"),
			Video(
				Controls(),
				Class("video-js vjs-theme-fantasy w-full bg-gray-100 aspect-video"),
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
}
