package services_test

import (
	"errors"
	"testing"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventService_CreateEvent(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(ctx, "truncate events, products cascade")
	require.NoError(t, err)

	srv := services.NewEventService(db)

	params := &types.CreateEventInput{
		EventType:     "webinar",
		TitleEn:       "Test Webinar 1",
		TitlePl:       "Webinarium Testowe",
		SubtitleEn:    nil,
		SubtitlePl:    nil,
		Slug:          "test-webinar-1",
		DescriptionEn: new("Test Description"),
		DescriptionPl: new("Opis wydarzenia"),
		Price:         new(decimal.NewFromInt(250)),
		Currency:      new("PLN"),
		HostIds:       nil,
		StartsAt:      time.Now().Add(time.Hour),
		EndsAt:        time.Now().Add(2 * time.Hour),
		IsVirtual:     true,
	}

	t.Run("creates a paid event with valid params", func(t *testing.T) {
		event, err := srv.CreateEvent(ctx, params)
		assert.NoError(t, err)
		assert.NotNil(t, event)

		assert.False(t, event.IsFree())
		assert.Equal(t, "250.00", event.Product.BasePriceAmount.StringFixedBank(2))
		assert.NotNil(t, event.ProductID)
	})

	t.Run("creates a free event without a product", func(t *testing.T) {
		p := *params
		p.Slug = "free-webinar"
		p.TitleEn = "Free Webinar"
		p.Price = nil
		p.Currency = nil

		event, err := srv.CreateEvent(ctx, &p)
		require.NoError(t, err)
		require.NotNil(t, event)

		assert.True(t, event.IsFree())
		assert.Nil(t, event.Product)
		assert.Nil(t, event.ProductID)
	})

	t.Run("treats a zero price as free", func(t *testing.T) {
		p := *params
		p.Slug = "zero-price-webinar"
		p.TitleEn = "Zero Price Webinar"
		p.Price = new(decimal.Zero)

		event, err := srv.CreateEvent(ctx, &p)
		require.NoError(t, err)
		require.NotNil(t, event)

		assert.True(t, event.IsFree())
		assert.Nil(t, event.Product)
	})

	t.Run("persists venue details for a physical event", func(t *testing.T) {
		p := *params
		p.Slug = "physical-seminar"
		p.TitleEn = "Physical Seminar"
		p.EventType = "seminar"
		p.Price = nil
		p.Currency = nil
		p.IsVirtual = false
		p.VenueNameEn = new("Main Hall")
		p.VenueNamePl = new("Sala Główna")
		p.VenueStreet = new("ul. Testowa 1")
		p.VenueCityEn = new("Warsaw")
		p.VenueCityPl = new("Warszawa")
		p.VenuePostalCode = new("00-001")
		p.VenueCountryCode = new("PL")

		event, err := srv.CreateEvent(ctx, &p)
		require.NoError(t, err)
		require.NotNil(t, event)

		assert.False(t, event.IsVirtual)
		assert.Equal(t, "seminar", string(event.EventType))
		require.NotNil(t, event.VenueNameEn)
		assert.Equal(t, "Main Hall", *event.VenueNameEn)
		require.NotNil(t, event.VenueCityPl)
		assert.Equal(t, "Warszawa", *event.VenueCityPl)
		require.NotNil(t, event.VenueCountryCode)
		assert.Equal(t, "PL", *event.VenueCountryCode)
	})

	t.Run("rejects a duplicate slug", func(t *testing.T) {
		p := *params
		p.Slug = "dup-slug-webinar"
		p.TitleEn = "Dup Slug A"

		_, err := srv.CreateEvent(ctx, &p)
		require.NoError(t, err)

		p.TitleEn = "Dup Slug B"
		_, err = srv.CreateEvent(ctx, &p)
		assert.Error(t, err)
	})
}

func TestEventService_UpdateEvent(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(ctx, "truncate events, products cascade")
	require.NoError(t, err)

	srv := services.NewEventService(db)

	event, err := mocks.Event(db, ctx)
	require.NoError(t, err)

	t.Run("Updating titles and descriptions", func(t *testing.T) {
		params := types.PatchEventInput{
			TitleEn:       types.Some("Updated title"),
			TitlePl:       types.Some("Zaktualizowany tytuł"),
			SubtitlePl:    types.Some("Zaktualizowany podtytuł"),
			SubtitleEn:    types.Some("Updated subtitle"),
			DescriptionEn: types.Some("Updated description"),
			DescriptionPl: types.Some("Zaktualizowany opis"),
		}

		updated, err := srv.UpdateEvent(ctx, event.ID, &params)
		assert.NoError(t, err)
		assert.Equal(t, event.ID, updated.ID)

		assert.Equal(t, params.TitleEn.Value, updated.TitleEn)
		assert.Equal(t, params.TitlePl.Value, updated.TitlePl)
		assert.Equal(t, params.SubtitlePl.Value, *updated.SubtitlePl)
		assert.Equal(t, params.SubtitleEn.Value, *updated.SubtitleEn)
		assert.Equal(t, params.DescriptionEn.Value, *updated.DescriptionEn)
		assert.Equal(t, params.DescriptionPl.Value, *updated.DescriptionPl)
		assert.Equal(t, event.Slug, updated.Slug)
	})

	t.Run("Absent fields are left untouched", func(t *testing.T) {
		before, err := srv.GetEventById(ctx, event.ID)
		require.NoError(t, err)

		updated, err := srv.UpdateEvent(ctx, event.ID, &types.PatchEventInput{
			TitleEn: types.Some("Only the English title changes"),
		})
		require.NoError(t, err)

		assert.Equal(t, "Only the English title changes", updated.TitleEn)
		assert.Equal(t, before.TitlePl, updated.TitlePl)
		assert.Equal(t, before.SubtitleEn, updated.SubtitleEn)
		assert.Equal(t, before.DescriptionPl, updated.DescriptionPl)
	})

	t.Run("Explicit null clears a nullable column", func(t *testing.T) {
		updated, err := srv.UpdateEvent(ctx, event.ID, &types.PatchEventInput{
			SubtitleEn: types.Null[string](),
		})
		require.NoError(t, err)

		assert.Nil(t, updated.SubtitleEn)
	})

	t.Run("Empty payload is a no-op", func(t *testing.T) {
		before, err := srv.GetEventById(ctx, event.ID)
		require.NoError(t, err)

		updated, err := srv.UpdateEvent(ctx, event.ID, &types.PatchEventInput{})
		require.NoError(t, err)

		assert.Equal(t, before.UpdatedAt, updated.UpdatedAt)
	})

	t.Run("Updating the slug", func(t *testing.T) {
		updated, err := srv.UpdateEvent(ctx, event.ID, &types.PatchEventInput{
			Slug: types.Some("a-brand-new-slug"),
		})
		require.NoError(t, err)

		assert.Equal(t, "a-brand-new-slug", updated.Slug)
	})

	t.Run("Rejects a malformed slug", func(t *testing.T) {
		_, err := srv.UpdateEvent(ctx, event.ID, &types.PatchEventInput{
			Slug: types.Some("Not A Slug"),
		})

		verrs, ok := errors.AsType[validation.Errors](err)
		require.True(t, ok, "expected validation.Errors, got %v", err)
		assert.Contains(t, verrs, "slug")
	})

	t.Run("Rejects a duplicate slug", func(t *testing.T) {
		other, err := mocks.Event(db, ctx)
		require.NoError(t, err)

		_, err = srv.UpdateEvent(ctx, event.ID, &types.PatchEventInput{
			Slug: types.Some(other.Slug),
		})

		verrs, ok := errors.AsType[validation.Errors](err)
		require.True(t, ok, "expected validation.Errors, got %v", err)
		assert.Contains(t, verrs, "slug")
	})

	t.Run("NOT NULL columns cannot be cleared or blanked", func(t *testing.T) {
		for name, params := range map[string]types.PatchEventInput{
			"explicit null": {TitleEn: types.Null[string]()},
			"blank string":  {TitlePl: types.Some("   ")},
			"null slug":     {Slug: types.Null[string]()},
		} {
			t.Run(name, func(t *testing.T) {
				_, err := srv.UpdateEvent(ctx, event.ID, &params)

				verrs, ok := errors.AsType[validation.Errors](err)
				require.True(t, ok, "expected validation.Errors, got %v", err)
				assert.Len(t, verrs, 1)
			})
		}
	})
}
