package api_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/api"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVideoGroupServer(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(ctx, "truncate video_groups cascade")
	require.NoError(t, err)

	srv := api.NewVideoGroupServer(db)

	t.Run("creates a free series without a product", func(t *testing.T) {
		resp, err := srv.CreateVideoGroup(ctx, api.CreateVideoGroupRequestObject{
			Body: &api.VideoGroupInput{
				TitleEn: "Free series",
				TitlePl: "Seria bezpłatna",
				Slug:    "free-series",
			},
		})
		require.NoError(t, err)

		out, ok := resp.(api.CreateVideoGroup201JSONResponse)
		require.True(t, ok, "expected 201, got %T", resp)
		assert.False(t, out.Body.IsPremium)
		assert.Nil(t, out.Body.Price)
	})

	t.Run("creates a paid series and reports its price", func(t *testing.T) {
		resp, err := srv.CreateVideoGroup(ctx, api.CreateVideoGroupRequestObject{
			Body: &api.VideoGroupInput{
				TitleEn:  "Paid series",
				TitlePl:  "Seria płatna",
				Slug:     "paid-series",
				Price:    new("199.00"),
				Currency: new("PLN"),
			},
		})
		require.NoError(t, err)

		out, ok := resp.(api.CreateVideoGroup201JSONResponse)
		require.True(t, ok, "expected 201, got %T", resp)
		assert.True(t, out.Body.IsPremium)
		require.NotNil(t, out.Body.Price)
		assert.Equal(t, "199.00", *out.Body.Price)
		require.NotNil(t, out.Body.Currency)
		assert.Equal(t, "PLN", *out.Body.Currency)
	})

	t.Run("rejects a duplicate slug", func(t *testing.T) {
		resp, err := srv.CreateVideoGroup(ctx, api.CreateVideoGroupRequestObject{
			Body: &api.VideoGroupInput{
				TitleEn: "Another series",
				TitlePl: "Kolejna seria",
				Slug:    "free-series",
			},
		})
		require.NoError(t, err)

		out, ok := resp.(api.CreateVideoGroup422JSONResponse)
		require.True(t, ok, "expected 422, got %T", resp)
		assert.Contains(t, out.Errors, "slug")
	})

	t.Run("pricing a free series puts it behind the paywall", func(t *testing.T) {
		group, err := mocks.VideoGroup(db, ctx)
		require.NoError(t, err)

		resp, err := srv.UpdateVideoGroup(ctx, api.UpdateVideoGroupRequestObject{
			Id: group.ID,
			Body: &types.PatchVideoGroupInput{
				Price:    types.Some(decimal.NewFromInt(249)),
				Currency: types.Some("PLN"),
			},
		})
		require.NoError(t, err)

		out, ok := resp.(api.UpdateVideoGroup200JSONResponse)
		require.True(t, ok, "expected 200, got %T", resp)
		assert.True(t, out.IsPremium)
		require.NotNil(t, out.Price)
		assert.Equal(t, "249.00", *out.Price)

		// The product is what the paywall reads, so it has to exist on the row.
		stored, err := queries.New(db).GetVideoGroupById(ctx, group.ID)
		require.NoError(t, err)
		assert.NotNil(t, stored.VideoGroup.ProductID)
	})

	t.Run("a price without a currency is rejected on a free series", func(t *testing.T) {
		group, err := mocks.VideoGroup(db, ctx)
		require.NoError(t, err)

		resp, err := srv.UpdateVideoGroup(ctx, api.UpdateVideoGroupRequestObject{
			Id:   group.ID,
			Body: &types.PatchVideoGroupInput{Price: types.Some(decimal.NewFromInt(249))},
		})
		require.NoError(t, err)

		out, ok := resp.(api.UpdateVideoGroup422JSONResponse)
		require.True(t, ok, "expected 422, got %T", resp)
		assert.Contains(t, out.Errors, "currency")
	})

	t.Run("zeroing the price keeps the product", func(t *testing.T) {
		product, err := mocks.Product(db, ctx, func(p *queries.InsertProductParams) {
			p.BasePriceAmount = decimal.NewFromInt(99)
		})
		require.NoError(t, err)

		group, err := mocks.VideoGroup(db, ctx, func(p *queries.InsertVideoGroupParams) {
			p.ProductID = &product.ID
		})
		require.NoError(t, err)

		resp, err := srv.UpdateVideoGroup(ctx, api.UpdateVideoGroupRequestObject{
			Id:   group.ID,
			Body: &types.PatchVideoGroupInput{Price: types.Some(decimal.Zero)},
		})
		require.NoError(t, err)

		out, ok := resp.(api.UpdateVideoGroup200JSONResponse)
		require.True(t, ok, "expected 200, got %T", resp)
		assert.True(t, out.IsPremium)
		require.NotNil(t, out.Price)
		assert.Equal(t, "0.00", *out.Price)
	})

	t.Run("404s for an unknown id", func(t *testing.T) {
		resp, err := srv.GetVideoGroup(ctx, api.GetVideoGroupRequestObject{Id: uuid.New()})
		require.NoError(t, err)
		assert.IsType(t, api.GetVideoGroup404Response{}, resp)
	})
}
