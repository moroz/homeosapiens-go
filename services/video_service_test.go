package services_test

import (
	"errors"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVideoService_ListVideoGroupsForUser(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)

	user, err := mocks.User(db, ctx)
	require.NoError(t, err)

	t.Run("lists public videos with HasAccess = true", func(t *testing.T) {
		_, err = db.Exec(ctx, "truncate video_groups cascade")
		require.NoError(t, err)

		free1, err := mocks.VideoGroup(db, ctx)
		require.NoError(t, err)

		free2, err := mocks.VideoGroup(db, ctx)
		require.NoError(t, err)

		actual, err := services.NewVideoService(db).ListVideoGroupsForUser(ctx, user.ID, nil)
		assert.NoError(t, err)

		var ids []uuid.UUID
		for _, v := range actual {
			ids = append(ids, v.VideoGroup.ID)
			assert.True(t, v.HasAccess)
		}

		assert.Equal(t, []uuid.UUID{free2.ID, free1.ID}, ids)
	})

	t.Run("lists paid videos with HasAccess = false", func(t *testing.T) {
		_, err = db.Exec(ctx, "truncate video_groups cascade")
		require.NoError(t, err)

		free, err := mocks.VideoGroup(db, ctx)
		require.NoError(t, err)

		product, err := mocks.Product(db, ctx)
		require.NoError(t, err)

		paid, err := mocks.VideoGroup(db, ctx, func(params *queries.InsertVideoGroupParams) {
			params.ProductID = &product.ID
		})
		require.NoError(t, err)

		actual, err := services.NewVideoService(db).ListVideoGroupsForUser(ctx, user.ID, nil)
		assert.NoError(t, err)

		var ids []uuid.UUID
		for _, v := range actual {
			ids = append(ids, v.VideoGroup.ID)
			assert.Equal(t, v.VideoGroup.ID != paid.ID, v.HasAccess)
		}

		assert.Equal(t, []uuid.UUID{paid.ID, free.ID}, ids)

		admin, err := mocks.User(db, ctx, func(params *types.SeedUserParams) {
			params.Role = queries.UserRoleAdministrator
		})
		require.NoError(t, err)

		actual, err = services.NewVideoService(db).ListVideoGroupsForUser(ctx, admin.ID, nil)
		assert.NoError(t, err)

		ids = nil
		for _, v := range actual {
			ids = append(ids, v.VideoGroup.ID)
			assert.True(t, v.HasAccess)
		}
		assert.Equal(t, []uuid.UUID{paid.ID, free.ID}, ids)
	})
}

// The paywall on a series page needs the price and the in-cart state to render
// its buy button, so both travel with the group details.
func TestVideoService_GetVideoGroupDetails_Paywall(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)

	_, err = db.Exec(ctx, "truncate video_groups, cart_line_items cascade")
	require.NoError(t, err)

	user, err := mocks.User(db, ctx)
	require.NoError(t, err)

	product, err := mocks.Product(db, ctx, func(params *queries.InsertProductParams) {
		params.BasePriceAmount = decimal.NewFromInt(199)
		params.BasePriceCurrency = "PLN"
	})
	require.NoError(t, err)

	group, err := mocks.VideoGroup(db, ctx, func(params *queries.InsertVideoGroupParams) {
		params.ProductID = &product.ID
	})
	require.NoError(t, err)

	srv := services.NewVideoService(db)

	details, err := srv.GetVideoGroupDetails(ctx, user.ID, &group.Slug, nil)
	require.NoError(t, err)
	assert.False(t, details.HasAccess)
	assert.True(t, details.IsPremium())
	require.NotNil(t, details.Price)
	assert.True(t, details.Price.Equal(decimal.NewFromInt(199)))
	require.NotNil(t, details.Currency)
	assert.Equal(t, "PLN", *details.Currency)
	assert.Zero(t, details.CountInCart)

	cartId, err := mocks.Cart(db, ctx, product.ID)
	require.NoError(t, err)

	details, err = srv.GetVideoGroupDetails(ctx, user.ID, &group.Slug, &cartId)
	require.NoError(t, err)
	assert.Equal(t, 1, details.CountInCart)

	// Granting access to the product is what a paid order does, and it unlocks
	// the group.
	_, err = db.Exec(ctx, "insert into user_product_access (user_id, product_id) values ($1, $2)", user.ID, product.ID)
	require.NoError(t, err)

	details, err = srv.GetVideoGroupDetails(ctx, user.ID, &group.Slug, nil)
	require.NoError(t, err)
	assert.True(t, details.HasAccess)
}

// The YouTube ID column is only meaningful for youtube-provider videos; the DB's
// check constraint enforces that, and the service must turn a violation into a
// field-level validation error instead of bubbling up a raw Postgres error.
func TestVideoService_UpdateVideo_YoutubeId(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)

	srv := services.NewVideoService(db)

	baseParams := func(v *queries.Video) *types.UpdateVideoInput {
		return &types.UpdateVideoInput{
			TitleEn: v.TitleEn,
			TitlePl: v.TitlePl,
			Slug:    v.Slug,
		}
	}

	t.Run("leaves youtube_id nil on a cloudfront video", func(t *testing.T) {
		video, err := mocks.Video(db, ctx)
		require.NoError(t, err)
		require.Equal(t, queries.VideoProviderCloudfront, video.Provider)

		updated, err := srv.UpdateVideo(ctx, video.ID, baseParams(video))
		require.NoError(t, err)
		assert.Nil(t, updated.YoutubeID)
	})

	t.Run("updates youtube_id on a youtube video", func(t *testing.T) {
		// mocks.Video / InsertVideo have no youtube_id column, so a youtube
		// video fixture needs a raw insert. The slug is randomized so reruns
		// against a persistent test DB don't collide with a leftover row.
		var videoID uuid.UUID
		err := db.QueryRow(
			ctx,
			`insert into videos (provider, title_en, title_pl, slug, youtube_id)
			 values ('youtube', 'YT video', 'Wideo YT', $1, 'oldId12345') returning id`,
			"yt-video-"+uuid.NewString(),
		).Scan(&videoID)
		require.NoError(t, err)

		video, err := queries.New(db).GetVideoById(ctx, videoID)
		require.NoError(t, err)

		params := baseParams(video)
		newID := "newId67890"
		params.YoutubeID = &newID

		updated, err := srv.UpdateVideo(ctx, video.ID, params)
		require.NoError(t, err)
		require.NotNil(t, updated.YoutubeID)
		assert.Equal(t, newID, *updated.YoutubeID)

		t.Run("rejects clearing youtube_id", func(t *testing.T) {
			params := baseParams(video)
			params.YoutubeID = nil

			_, err := srv.UpdateVideo(ctx, video.ID, params)
			verrs, ok := errors.AsType[validation.Errors](err)
			require.True(t, ok, "expected validation.Errors, got %v", err)
			assert.Error(t, verrs["youtubeId"])
		})
	})
}
