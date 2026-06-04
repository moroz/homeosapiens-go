package videos

import (
	"fmt"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/components"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func VideoThumbnail(video *queries.Video, locale string) Node {
	thumbnailID := video.ThumbnailEnID
	if locale == "pl" && video.ThumbnailPlID != nil {
		thumbnailID = video.ThumbnailPlID
	}

	if thumbnailID == nil {
		return components.Logo("text-slate-400")
	}

	baseURL := fmt.Sprintf("%s/images/%s", config.AssetCdnBaseUrl, thumbnailID)

	return Picture(
		Source(SrcSet(fmt.Sprintf("%s/1x.webp 1x, %s/2x.webp 2x", baseURL, baseURL)), Type("image/webp")),
		Img(SrcSet(fmt.Sprintf("%s/1x.png 1x, %s/2x.png 2x", baseURL, baseURL)), Src(fmt.Sprintf("%s/1x.png", baseURL)), Width("320"), Height("180")),
	)
}

func VideoCard(ctx *types.CustomContext, group *types.VideoGroupDetailsDTO, video *queries.Video) Node {
	title := video.TitleEn
	if ctx.Language == "pl" {
		title = video.TitlePl
	}

	return A(
		Class("video-card-link no-underline"),
		Href(fmt.Sprintf("/videos/%s/%s", group.Slug, video.Slug)),
		Article(
			Class("video-card"),
			Header(
				Class("video-card-thumb"),
				VideoThumbnail(video, ctx.Language),
				Iff(video.DurationSeconds != nil, func() Node {
					hours := *video.DurationSeconds / 3600
					secs := *video.DurationSeconds % 60
					minutes := *video.DurationSeconds % 3600 / 60
					text := fmt.Sprintf("%02d:%02d", minutes, secs)
					if hours > 0 {
						text = fmt.Sprintf("%d:%s", hours, text)
					}

					return Span(
						Class("duration-badge"),
						Text(text),
					)
				}),
			),
			Footer(
				Class("video-card-body"),
				H4(
					Class("video-card-title"),
					Text(title)),
			),
		),
	)
}

func Index(ctx *types.CustomContext, videos []*types.VideoGroupListDTO, group *types.VideoGroupDetailsDTO) Node {
	return layout.BareLayout(ctx, "Videos",
		Div(
			Class("card mx-auto grid grid-cols-[1fr_3fr] gap-8 lg:w-7xl"),
			Aside(
				Class("mobile:border-r-0 border-r border-slate-300 pr-6 mobile:pr-0"),
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
						H3(Class("mb-6 inline-block border-b-2 border-primary/30 pb-2 text-2xl font-bold text-primary"), Text(title)),
						Div(
							Class("video-grid"),
							Map(group.Videos, func(video *queries.Video) Node {
								return VideoCard(ctx, group, video)
							}),
						),
					}
				}),
			),
		),
	)
}
