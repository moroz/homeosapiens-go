package videos

import (
	"fmt"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Show(ctx *types.CustomContext, group *types.VideoGroupDetailsDTO, video *types.VideoDetailsDTO) Node {
	title := video.TitleEn
	if ctx.Language == "pl" {
		title = video.TitlePl
	}

	groupTitle := group.TitleEn
	if ctx.Language == "pl" {
		groupTitle = group.TitlePl
	}

	var index int
	for i, item := range group.Videos {
		if item.ID == video.ID {
			index = i
			break
		}
	}

	l := ctx.Localizer

	backURL := fmt.Sprintf("/videos/%s", group.Slug)

	return layout.Layout(ctx, title,
		Div(Class("card"),
			A(Href(backURL), Text("<< "), Text(l.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "videos.show.back_to_group",
				TemplateData: map[string]string{
					"GroupName": groupTitle,
				},
			}))),
			H3(Class("font-bold text-primary text-2xl"), Text(title)),
			Raw(l.MustLocalize(&i18n.LocalizeConfig{
				MessageID: "videos.show.video_group_name_html",
				TemplateData: map[string]any{
					"CurrentIndex": index + 1,
					"TotalCount":   len(group.Videos),
					"GroupURL":     backURL,
					"GroupName":    groupTitle,
				},
			})),
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
