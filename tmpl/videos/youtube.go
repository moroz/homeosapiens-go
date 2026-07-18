package videos

import (
	"fmt"
	"strings"

	"github.com/moroz/homeosapiens-go/tmpl/components/icons"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Youtube(ctx *types.CustomContext, videos []*types.VideoListDTO) Node {
	l := ctx.Localizer
	title := l.MustLocalizeMessage(&i18n.Message{
		ID: "videos.youtube.title",
	})

	return layout.BareLayout(ctx, title,
		Div(
			Class("container mx-auto py-6"),
			H2(Class("page-title"), Text(title)),
			Div(
				Class("video-grid"),
				Map(videos, func(video *types.VideoListDTO) Node {
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
							H4(Class("video-card-title text-lg"), helpers.RenderMarkdown(title)),
							Iff(len(video.Hosts) != 0, func() Node {
								names := make([]string, len(video.Hosts))
								for i, h := range video.Hosts {
									names[i] = fmt.Sprintf("%s&nbsp;%s", h.GivenName, h.FamilyName)
								}

								return P(
									Strong(Text(l.MustLocalizeMessage(&i18n.Message{
										ID: "videos.youtube.featuring",
									}))),
									Text(" "),
									Raw(strings.Join(names, ", ")),
								)
							}),
							Iff(desc != nil, func() Node {
								return Div(Class("text-sm"), helpers.RenderMarkdown(*desc))
							}),
							A(
								Class("youtube-button mt-auto w-full"),
								Href(url), Target("_blank"), Rel("noopener noreferrer"),
								icons.Icon(&icons.IconProps{Name: "youtube", ViewBox: "0 0 28.57 20", Classes: "w-auto"}),
								Span(Text(l.MustLocalizeMessage(&i18n.Message{
									ID: "videos.youtube.watch_on_youtube",
								}))),
							),
						),
					)
				}),
			),
		),
	)
}
