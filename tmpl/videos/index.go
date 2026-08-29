package videos

import (
	"fmt"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/tmpl/helpers"
	"github.com/moroz/homeosapiens-go/tmpl/layout"
	"github.com/moroz/homeosapiens-go/types"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func DurationBadge(duration *int32) Node {
	if duration == nil {
		return nil
	}

	hours := *duration / 3600
	secs := *duration % 60
	minutes := *duration % 3600 / 60
	text := fmt.Sprintf("%02d:%02d", minutes, secs)
	if hours > 0 {
		text = fmt.Sprintf("%d:%s", hours, text)
	}

	return Span(
		Class("duration-badge"),
		Text(text),
	)
}

func VideoCard(ctx *types.CustomContext, group *types.VideoGroupDetailsDTO, video *queries.Video) Node {
	title := video.TitleEn
	if ctx.Language == "pl" {
		title = video.TitlePl
	}

	card := Article(
		Class("video-card"),
		Header(
			Class("video-card-thumb"),
			Img(Src(fmt.Sprintf("/videos/%s/thumbnail/%s", video.ID, ctx.Language))),
			DurationBadge(video.DurationSeconds),
			If(!group.HasAccess, Span(
				Class("absolute inset-0 flex items-center justify-center bg-slate-900/40 text-white"),
				LockIcon(),
			)),
		),
		Footer(
			Class("video-card-body"),
			H4(
				Class("video-card-title"),
				Text(title)),
		),
	)

	// Without access the video page has nothing to show but the paywall, so the
	// card advertises the series rather than linking into it. The page itself
	// still checks access, so the link is a convenience, not the guard.
	if !group.HasAccess {
		return Div(Class("video-card-link cursor-not-allowed opacity-75"), card)
	}

	return A(
		Class("video-card-link no-underline"),
		Href(fmt.Sprintf("/videos/%s/%s", group.Slug, video.Slug)),
		card,
	)
}

func VideoGroupList(ctx *types.CustomContext, videoGroups []*types.VideoGroupListDTO, active *types.VideoGroupDetailsDTO) Node {
	return Aside(
		Class("mobile:border-r-0 border-r border-slate-300 pr-4 mobile:pr-0"),
		H2(Class("page-title ml-3"), Text("Videos")),
		Nav(
			Ul(
				Class("video-group-menu"),
				Map(videoGroups, func(vg *types.VideoGroupListDTO) Node {
					title := vg.TitleEn
					if ctx.Language == "pl" {
						title = vg.TitlePl
					}

					class := ""
					if active != nil && vg.ID == active.ID {
						class = "is-active"
					}

					dateRange := ""
					if vg.MaxRecordedOn != nil && vg.MinRecordedOn != nil {
						dateRange = helpers.FormatDateRange(*vg.MinRecordedOn, *vg.MaxRecordedOn, ctx.Timezone, ctx.Language)
					}

					return Li(
						A(
							Class(class),
							Href(fmt.Sprintf("/videos/%s", vg.Slug)),
							P(
								Text(title),
								If(vg.IsPremium() && !vg.HasAccess, PaidBadge(ctx.Localizer)),
								If(vg.IsPremium() && vg.HasAccess, CrownIcon("inline")),
							),
							If(dateRange != "", Span(Class("text-sm font-normal"), Text(dateRange))),
						),
					)
				}),
			),
		),
	)
}

func Index(ctx *types.CustomContext, videoGroups []*types.VideoGroupListDTO, activeGroup *types.VideoGroupDetailsDTO) Node {
	return layout.BareLayout(ctx, "Videos",
		Div(
			Class("card mx-auto grid grid-cols-[1fr_3fr] gap-8 lg:w-7xl mt-6 desktop:pl-4"),
			VideoGroupList(ctx, videoGroups, activeGroup),
			Main(
				Iff(activeGroup != nil, func() Node {
					title := activeGroup.TitleEn
					if ctx.Language == "pl" {
						title = activeGroup.TitlePl
					}

					gridClasses := "video-grid"
					if !activeGroup.HasAccess {
						gridClasses += " mt-6"
					}

					return Group{
						H3(Class("mb-6 inline-block border-b-2 border-primary/30 pb-2 text-2xl font-bold text-primary"), Text(title)),
						// Without access the offer comes first, but the thumbnails
						// stay: they are what the series is being sold on. The
						// cards themselves stop linking anywhere.
						If(!activeGroup.HasAccess, LockedPanel(ctx, activeGroup)),
						Div(
							Class(gridClasses),
							Map(activeGroup.Videos, func(video *queries.Video) Node {
								return VideoCard(ctx, activeGroup, video)
							}),
						),
					}
				}),
			),
		),
	)
}
