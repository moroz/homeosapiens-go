package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/types"
	"github.com/moroz/homeosapiens-go/web/api"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initDB(ctx context.Context) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, config.MustGetenv("TEST_DATABASE_URL"))
}

// validEventInput is a virtual, paid event that passes CreateEventInput.Validate.
func validEventInput(slug string) *api.EventInput {
	return &api.EventInput{
		EventType:     "webinar",
		TitleEn:       "API Webinar",
		TitlePl:       "Webinarium API",
		Slug:          slug,
		DescriptionEn: new("Description"),
		DescriptionPl: new("Opis"),
		Price:         new("250.00"),
		Currency:      new("PLN"),
		StartsAt:      time.Now().Add(time.Hour).UTC(),
		EndsAt:        time.Now().Add(2 * time.Hour).UTC(),
		IsVirtual:     true,
	}
}

func TestEventServer_ListEvents(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(ctx, "truncate events, products cascade")
	require.NoError(t, err)

	srv := api.NewEventServer(db)

	// Rows come back ordered by starts_at desc, so the offsets are deterministic.
	base := time.Now().Truncate(time.Second)
	created := make([]*queries.Event, 3)
	for i := range created {
		event, err := mocks.Event(db, ctx, func(p *queries.UpsertEventParams) {
			p.StartsAt = base.Add(time.Duration(i) * time.Hour)
			p.EndsAt = p.StartsAt.Add(time.Hour)
			p.TitleEn = fmt.Sprintf("Event %d", i)
		})
		require.NoError(t, err)
		created[i] = event
	}
	newest, middle, oldest := created[2], created[1], created[0]

	list := func(params api.ListEventsParams) api.ListEvents200JSONResponse {
		resp, err := srv.ListEvents(ctx, api.ListEventsRequestObject{Params: params})
		require.NoError(t, err)
		out, ok := resp.(api.ListEvents200JSONResponse)
		require.True(t, ok, "expected 200, got %T", resp)
		return out
	}

	t.Run("defaults to the first page and the configured page size", func(t *testing.T) {
		out := list(api.ListEventsParams{})

		assert.Equal(t, int32(1), out.Pagination.Page)
		assert.Equal(t, int32(config.DefaultPageSize), out.Pagination.PerPage)
		assert.Equal(t, int64(3), out.Pagination.Total)
		assert.Equal(t, int32(1), out.Pagination.TotalPages)
		require.Len(t, out.Data, 3)

		assert.Equal(t, newest.ID, out.Data[0].Id)
		assert.Equal(t, oldest.ID, out.Data[2].Id)
	})

	t.Run("honours page and perPage", func(t *testing.T) {
		out := list(api.ListEventsParams{Page: new(int32(2)), PerPage: new(int32(1))})

		assert.Equal(t, int32(2), out.Pagination.Page)
		assert.Equal(t, int32(1), out.Pagination.PerPage)
		assert.Equal(t, int64(3), out.Pagination.Total)
		assert.Equal(t, int32(3), out.Pagination.TotalPages)

		require.Len(t, out.Data, 1)
		assert.Equal(t, middle.ID, out.Data[0].Id)
	})

	t.Run("non-positive pagination params fall back to the defaults", func(t *testing.T) {
		out := list(api.ListEventsParams{Page: new(int32(0)), PerPage: new(int32(-5))})

		assert.Equal(t, int32(1), out.Pagination.Page)
		assert.Equal(t, int32(config.DefaultPageSize), out.Pagination.PerPage)
		assert.Len(t, out.Data, 3)
	})

	t.Run("perPage is clamped to the documented maximum", func(t *testing.T) {
		out := list(api.ListEventsParams{PerPage: new(int32(5000))})

		assert.Equal(t, int32(config.MaxPageSize), out.Pagination.PerPage)
	})

	t.Run("a page past the end is empty rather than an error", func(t *testing.T) {
		out := list(api.ListEventsParams{Page: new(int32(99))})

		assert.Empty(t, out.Data)
		assert.Equal(t, int64(3), out.Pagination.Total)
	})

	t.Run("maps the summary fields, including publication state", func(t *testing.T) {
		unpublished, err := mocks.Event(db, ctx, func(p *queries.UpsertEventParams) {
			p.Published = false
			p.StartsAt = base.Add(-time.Hour)
			p.EndsAt = base
			p.SubtitleEn = new("Subtitle")
			p.SubtitlePl = new("Podtytuł")
		})
		require.NoError(t, err)

		out := list(api.ListEventsParams{})
		require.Len(t, out.Data, 4)

		// starts_at desc puts the backdated event last.
		got := out.Data[3]
		assert.Equal(t, unpublished.ID, got.Id)
		assert.Equal(t, unpublished.Slug, got.Slug)
		assert.Equal(t, "Some event", got.TitleEn)
		assert.Equal(t, "seminar", got.EventType)
		assert.True(t, got.IsVirtual)
		assert.Equal(t, "Subtitle", *got.SubtitleEn)
		assert.Equal(t, "Podtytuł", *got.SubtitlePl)
		assert.Nil(t, got.PublishedAt)

		assert.NotNil(t, out.Data[0].PublishedAt)
	})
}

func TestEventServer_GetEvent(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	srv := api.NewEventServer(db)

	get := func(id uuid.UUID) api.GetEventResponseObject {
		resp, err := srv.GetEvent(ctx, api.GetEventRequestObject{Id: id})
		require.NoError(t, err)
		return resp
	}

	t.Run("returns the full details of a free event", func(t *testing.T) {
		event, err := mocks.Event(db, ctx, func(p *queries.UpsertEventParams) {
			p.IsVirtual = false
			p.VenueNameEn = new("Main Hall")
			p.VenueNamePl = new("Sala Główna")
			p.VenueStreet = new("ul. Testowa 1")
			p.VenueCityEn = new("Warsaw")
			p.VenueCityPl = new("Warszawa")
			p.VenuePostalCode = new("00-001")
			p.VenueCountryCode = new("PL")
		})
		require.NoError(t, err)

		out, ok := get(event.ID).(api.GetEvent200JSONResponse)
		require.True(t, ok)

		assert.Equal(t, event.ID, out.Id)
		assert.Equal(t, event.Slug, out.Slug)
		assert.Equal(t, event.TitleEn, out.TitleEn)
		assert.Equal(t, event.TitlePl, out.TitlePl)
		assert.Equal(t, "seminar", out.EventType)
		assert.Equal(t, "Some description", *out.DescriptionEn)
		assert.Equal(t, "Opis", *out.DescriptionPl)
		assert.NotNil(t, out.PublishedAt)
		assert.True(t, out.IsFree)
		assert.Nil(t, out.Price)
		assert.Nil(t, out.Currency)
		assert.False(t, out.IsVirtual)
		assert.NotNil(t, out.PublishedAt)

		assert.Equal(t, "Main Hall", *out.VenueNameEn)
		assert.Equal(t, "Sala Główna", *out.VenueNamePl)
		assert.Equal(t, "ul. Testowa 1", *out.VenueStreet)
		assert.Equal(t, "Warsaw", *out.VenueCityEn)
		assert.Equal(t, "Warszawa", *out.VenueCityPl)
		assert.Equal(t, "00-001", *out.VenuePostalCode)
		assert.Equal(t, "PL", *out.VenueCountryCode)

		assert.Empty(t, out.Hosts)
	})

	t.Run("exposes the price of a paid event as a fixed-point string", func(t *testing.T) {
		event, err := mocks.PaidEvent(db, ctx)
		require.NoError(t, err)

		out, ok := get(event.ID).(api.GetEvent200JSONResponse)
		require.True(t, ok)

		assert.False(t, out.IsFree)
		require.NotNil(t, out.Price)
		assert.Equal(t, "560.00", *out.Price)
		require.NotNil(t, out.Currency)
		assert.Equal(t, "EUR", *out.Currency)
	})

	t.Run("returns the hosts in their stored order", func(t *testing.T) {
		event, err := mocks.Event(db, ctx)
		require.NoError(t, err)

		first, err := mocks.Host(db, ctx, func(p *queries.UpsertHostParams) {
			p.GivenName = "Anna"
			p.Salutation = new("common.hosts.salutation.dr")
			p.Country = new("PL")
		})
		require.NoError(t, err)
		second, err := mocks.Host(db, ctx, func(p *queries.UpsertHostParams) {
			p.GivenName = "Bob"
		})
		require.NoError(t, err)

		for i, host := range []*queries.Host{second, first} {
			_, err := queries.New(db).InsertEventHost(ctx, &queries.InsertEventHostParams{
				EventID:  event.ID,
				HostID:   host.ID,
				Position: int32(i + 1),
			})
			require.NoError(t, err)
		}

		out, ok := get(event.ID).(api.GetEvent200JSONResponse)
		require.True(t, ok)

		require.Len(t, out.Hosts, 2)
		assert.Equal(t, second.ID, out.Hosts[0].Id)
		assert.Equal(t, first.ID, out.Hosts[1].Id)
		assert.Equal(t, "Anna", out.Hosts[1].GivenName)
		assert.Equal(t, "Doe", out.Hosts[1].FamilyName)
		assert.Equal(t, "common.hosts.salutation.dr", *out.Hosts[1].Salutation)
		assert.Equal(t, "PL", *out.Hosts[1].Country)
	})

	t.Run("returns an unpublished event", func(t *testing.T) {
		event, err := mocks.Event(db, ctx, func(p *queries.UpsertEventParams) {
			p.Published = false
		})
		require.NoError(t, err)

		out, ok := get(event.ID).(api.GetEvent200JSONResponse)
		require.True(t, ok)
		assert.Nil(t, out.PublishedAt)
	})

	t.Run("returns 404 for an unknown id", func(t *testing.T) {
		assert.IsType(t, api.GetEvent404Response{}, get(uuid.Must(uuid.NewV7())))
	})
}

func TestEventServer_CreateEvent(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	srv := api.NewEventServer(db)

	create := func(t *testing.T, body *api.EventInput) api.CreateEventResponseObject {
		t.Helper()
		resp, err := srv.CreateEvent(ctx, api.CreateEventRequestObject{Body: body})
		require.NoError(t, err)
		return resp
	}

	t.Run("creates a paid event and points Location at it", func(t *testing.T) {
		body := validEventInput("api-created-paid-event")

		out, ok := create(t, body).(api.CreateEvent201JSONResponse)
		require.True(t, ok)

		assert.Equal(t, body.Slug, out.Body.Slug)
		assert.Equal(t, body.TitleEn, out.Body.TitleEn)
		assert.Equal(t, body.TitlePl, out.Body.TitlePl)
		assert.Equal(t, "webinar", out.Body.EventType)
		assert.True(t, out.Body.IsVirtual)
		assert.False(t, out.Body.IsFree)
		assert.Equal(t, "250.00", *out.Body.Price)
		assert.Equal(t, "PLN", *out.Body.Currency)
		assert.Nil(t, out.Body.PublishedAt)

		require.NotNil(t, out.Headers.Location)
		assert.Equal(t, "/events/"+out.Body.Id.String(), *out.Headers.Location)

		// The row is readable through GetEvent afterwards.
		fetched, err := srv.GetEvent(ctx, api.GetEventRequestObject{Id: out.Body.Id})
		require.NoError(t, err)
		assert.IsType(t, api.GetEvent200JSONResponse{}, fetched)
	})

	t.Run("creates a free event without a product", func(t *testing.T) {
		body := validEventInput("api-created-free-event")
		body.Price = nil
		body.Currency = nil

		out, ok := create(t, body).(api.CreateEvent201JSONResponse)
		require.True(t, ok)

		assert.True(t, out.Body.IsFree)
		assert.Nil(t, out.Body.Price)
		assert.Nil(t, out.Body.Currency)

		event, err := queries.New(db).GetEventById(ctx, out.Body.Id)
		require.NoError(t, err)
		assert.Nil(t, event.ProductID)
	})

	t.Run("treats a zero price as free", func(t *testing.T) {
		body := validEventInput("api-created-zero-price-event")
		body.Price = new("0")

		out, ok := create(t, body).(api.CreateEvent201JSONResponse)
		require.True(t, ok)

		assert.True(t, out.Body.IsFree)

		event, err := queries.New(db).GetEventById(ctx, out.Body.Id)
		require.NoError(t, err)
		assert.Nil(t, event.ProductID)
	})

	t.Run("persists venue details for a physical event", func(t *testing.T) {
		body := validEventInput("api-created-physical-event")
		body.EventType = "seminar"
		body.IsVirtual = false
		body.VenueNameEn = new("Main Hall")
		body.VenueNamePl = new("Sala Główna")
		body.VenueStreet = new("ul. Testowa 1")
		body.VenueCityEn = new("Warsaw")
		body.VenueCityPl = new("Warszawa")
		body.VenuePostalCode = new("00-001")
		body.VenueCountryCode = new("PL")

		out, ok := create(t, body).(api.CreateEvent201JSONResponse)
		require.True(t, ok)

		fetched, err := srv.GetEvent(ctx, api.GetEventRequestObject{Id: out.Body.Id})
		require.NoError(t, err)
		details, ok := fetched.(api.GetEvent200JSONResponse)
		require.True(t, ok)

		assert.False(t, details.IsVirtual)
		assert.Equal(t, "Main Hall", *details.VenueNameEn)
		assert.Equal(t, "Warszawa", *details.VenueCityPl)
		assert.Equal(t, "PL", *details.VenueCountryCode)
	})

	t.Run("associates hosts in payload order", func(t *testing.T) {
		first, err := mocks.Host(db, ctx)
		require.NoError(t, err)
		second, err := mocks.Host(db, ctx)
		require.NoError(t, err)

		body := validEventInput("api-created-event-with-hosts")
		body.HostIds = []uuid.UUID{second.ID, first.ID}

		out, ok := create(t, body).(api.CreateEvent201JSONResponse)
		require.True(t, ok)

		require.Len(t, out.Body.Hosts, 2)
		assert.Equal(t, second.ID, out.Body.Hosts[0].Id)
		assert.Equal(t, first.ID, out.Body.Hosts[1].Id)

		fetched, err := srv.GetEvent(ctx, api.GetEventRequestObject{Id: out.Body.Id})
		require.NoError(t, err)
		details, ok := fetched.(api.GetEvent200JSONResponse)
		require.True(t, ok)

		require.Len(t, details.Hosts, 2)
		assert.Equal(t, second.ID, details.Hosts[0].Id)
		assert.Equal(t, first.ID, details.Hosts[1].Id)
	})

	t.Run("rejects invalid payloads with 422", func(t *testing.T) {
		examples := map[string]struct {
			mutate func(body *api.EventInput)
			field  string
		}{
			"missing English title": {func(b *api.EventInput) { b.TitleEn = "" }, "titleEn"},
			"missing Polish title":  {func(b *api.EventInput) { b.TitlePl = "" }, "titlePl"},
			"malformed slug":        {func(b *api.EventInput) { b.Slug = "Not A Slug" }, "slug"},
			"unknown event type":    {func(b *api.EventInput) { b.EventType = "workshop" }, "eventType"},
			"price without currency": {func(b *api.EventInput) {
				b.Currency = nil
			}, "currency"},
			"unsupported currency": {func(b *api.EventInput) { b.Currency = new("USD") }, "currency"},
			"end before start": {func(b *api.EventInput) {
				b.EndsAt = b.StartsAt.Add(-time.Hour)
			}, "endsAt"},
			"physical event without a venue": {func(b *api.EventInput) { b.IsVirtual = false }, "venueNameEn"},
		}

		i := 0
		for name, example := range examples {
			i++
			t.Run(name, func(t *testing.T) {
				body := validEventInput(fmt.Sprintf("api-invalid-payload-%d", i))
				example.mutate(body)

				out, ok := create(t, body).(api.CreateEvent422JSONResponse)
				require.True(t, ok, "expected 422")
				assert.Contains(t, out.Errors, example.field)
			})
		}
	})

	t.Run("rejects a duplicate slug with 422", func(t *testing.T) {
		body := validEventInput("api-duplicate-slug-event")
		require.IsType(t, api.CreateEvent201JSONResponse{}, create(t, body))

		out, ok := create(t, body).(api.CreateEvent422JSONResponse)
		require.True(t, ok, "expected 422")
		assert.Contains(t, out.Errors, "slug")
	})

	t.Run("rejects a malformed price string with 422", func(t *testing.T) {
		body := validEventInput("api-malformed-price-event")
		body.Price = new("not a number")

		out, ok := create(t, body).(api.CreateEvent422JSONResponse)
		require.True(t, ok, "expected 422")
		assert.Contains(t, out.Errors, "price")
	})
}

func TestEventServer_UpdateEvent(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	srv := api.NewEventServer(db)

	update := func(t *testing.T, id uuid.UUID, body *types.PatchEventInput) api.UpdateEventResponseObject {
		t.Helper()
		resp, err := srv.UpdateEvent(ctx, api.UpdateEventRequestObject{Id: id, Body: body})
		require.NoError(t, err)
		return resp
	}

	t.Run("applies the fields present in the payload", func(t *testing.T) {
		event, err := mocks.Event(db, ctx)
		require.NoError(t, err)

		out, ok := update(t, event.ID, &types.PatchEventInput{
			TitleEn:       types.Some("Updated title"),
			TitlePl:       types.Some("Zaktualizowany tytuł"),
			SubtitleEn:    types.Some("Updated subtitle"),
			DescriptionEn: types.Some("Updated description"),
			Slug:          types.Some("api-updated-slug"),
		}).(api.UpdateEvent200JSONResponse)
		require.True(t, ok)

		assert.Equal(t, event.ID, out.Id)
		assert.Equal(t, "Updated title", out.TitleEn)
		assert.Equal(t, "Zaktualizowany tytuł", out.TitlePl)
		assert.Equal(t, "Updated subtitle", *out.SubtitleEn)
		assert.Equal(t, "Updated description", *out.DescriptionEn)
		assert.Equal(t, "api-updated-slug", out.Slug)
		assert.Equal(t, "Opis", *out.DescriptionPl)
		assert.NotNil(t, out.PublishedAt)
		// updated_at is a timestamp(0), so a same-second update reads back equal.
		assert.False(t, out.UpdatedAt.Before(event.UpdatedAt))
	})

	t.Run("an empty payload leaves the event untouched", func(t *testing.T) {
		event, err := mocks.Event(db, ctx)
		require.NoError(t, err)

		out, ok := update(t, event.ID, &types.PatchEventInput{}).(api.UpdateEvent200JSONResponse)
		require.True(t, ok)

		assert.Equal(t, event.TitleEn, out.TitleEn)
		assert.Equal(t, event.UpdatedAt.UTC(), out.UpdatedAt.UTC())
	})

	t.Run("an explicit null clears a nullable column", func(t *testing.T) {
		event, err := mocks.Event(db, ctx, func(p *queries.UpsertEventParams) {
			p.SubtitleEn = new("Subtitle")
		})
		require.NoError(t, err)

		out, ok := update(t, event.ID, &types.PatchEventInput{
			SubtitleEn: types.Null[string](),
		}).(api.UpdateEvent200JSONResponse)
		require.True(t, ok)

		assert.Nil(t, out.SubtitleEn)
	})

	t.Run("pricing a free event reports the new price", func(t *testing.T) {
		event, err := mocks.Event(db, ctx)
		require.NoError(t, err)

		out, ok := update(t, event.ID, &types.PatchEventInput{
			Price:    types.Some(decimal.NewFromInt(250)),
			Currency: types.Some("PLN"),
		}).(api.UpdateEvent200JSONResponse)
		require.True(t, ok)

		assert.False(t, out.IsFree)
		assert.Equal(t, "250.00", *out.Price)
		assert.Equal(t, "PLN", *out.Currency)
	})

	t.Run("zeroing a paid event reports it as free", func(t *testing.T) {
		event, err := mocks.PaidEvent(db, ctx)
		require.NoError(t, err)

		out, ok := update(t, event.ID, &types.PatchEventInput{
			Price: types.Null[decimal.Decimal](),
		}).(api.UpdateEvent200JSONResponse)
		require.True(t, ok)

		assert.True(t, out.IsFree)
		assert.Nil(t, out.Price)
		assert.Nil(t, out.Currency)
	})

	t.Run("replacing the host list", func(t *testing.T) {
		event, err := mocks.Event(db, ctx)
		require.NoError(t, err)
		host, err := mocks.Host(db, ctx)
		require.NoError(t, err)

		out, ok := update(t, event.ID, &types.PatchEventInput{
			HostIds: types.Some([]uuid.UUID{host.ID}),
		}).(api.UpdateEvent200JSONResponse)
		require.True(t, ok)

		require.Len(t, out.Hosts, 1)
		assert.Equal(t, host.ID, out.Hosts[0].Id)

		fetched, err := srv.GetEvent(ctx, api.GetEventRequestObject{Id: event.ID})
		require.NoError(t, err)
		details, ok := fetched.(api.GetEvent200JSONResponse)
		require.True(t, ok)
		require.Len(t, details.Hosts, 1)
		assert.Equal(t, host.ID, details.Hosts[0].Id)
	})

	t.Run("rejects invalid payloads with 422", func(t *testing.T) {
		event, err := mocks.Event(db, ctx)
		require.NoError(t, err)

		other, err := mocks.Event(db, ctx)
		require.NoError(t, err)

		examples := map[string]struct {
			body  types.PatchEventInput
			field string
		}{
			"malformed slug":     {types.PatchEventInput{Slug: types.Some("Not A Slug")}, "slug"},
			"duplicate slug":     {types.PatchEventInput{Slug: types.Some(other.Slug)}, "slug"},
			"blank title":        {types.PatchEventInput{TitleEn: types.Some("   ")}, "titleEn"},
			"cleared title":      {types.PatchEventInput{TitlePl: types.Null[string]()}, "titlePl"},
			"negative price":     {types.PatchEventInput{Price: types.Some(decimal.NewFromInt(-1))}, "price"},
			"bad currency":       {types.PatchEventInput{Currency: types.Some("USD")}, "currency"},
			"duplicate hostIds":  {types.PatchEventInput{HostIds: types.Some([]uuid.UUID{other.ID, other.ID})}, "hostIds"},
			"unknown host":       {types.PatchEventInput{HostIds: types.Some([]uuid.UUID{uuid.Must(uuid.NewV7())})}, "hostIds"},
			"price without cur.": {types.PatchEventInput{Price: types.Some(decimal.NewFromInt(250))}, "currency"},
		}

		for name, example := range examples {
			t.Run(name, func(t *testing.T) {
				out, ok := update(t, event.ID, &example.body).(api.UpdateEvent422JSONResponse)
				require.True(t, ok, "expected 422")
				assert.Contains(t, out.Errors, example.field)
			})
		}
	})

	t.Run("returns 404 for an unknown id", func(t *testing.T) {
		resp := update(t, uuid.Must(uuid.NewV7()), &types.PatchEventInput{
			TitleEn: types.Some("Does not matter"),
		})
		assert.IsType(t, api.UpdateEvent404Response{}, resp)
	})
}

func TestEventServer_PublishEvent(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	srv := api.NewEventServer(db)

	publish := func(t *testing.T, id uuid.UUID) api.PublishEventResponseObject {
		t.Helper()
		resp, err := srv.PublishEvent(ctx, api.PublishEventRequestObject{Id: id})
		require.NoError(t, err)
		return resp
	}

	t.Run("publishes a complete draft", func(t *testing.T) {
		event, err := mocks.Event(db, ctx, func(p *queries.UpsertEventParams) {
			p.Published = false
		})
		require.NoError(t, err)

		assert.IsType(t, api.PublishEvent204Response{}, publish(t, event.ID))

		stored, err := queries.New(db).GetEventById(ctx, event.ID)
		require.NoError(t, err)
		assert.NotNil(t, stored.PublishedAt)
	})

	t.Run("rejects a draft with missing descriptions", func(t *testing.T) {
		for field, override := range map[string]func(p *queries.UpsertEventParams){
			"descriptionEn": func(p *queries.UpsertEventParams) { p.DescriptionEn = nil },
			"descriptionPl": func(p *queries.UpsertEventParams) { p.DescriptionPl = nil },
		} {
			t.Run(field, func(t *testing.T) {
				event, err := mocks.Event(db, ctx, func(p *queries.UpsertEventParams) {
					p.Published = false
					override(p)
				})
				require.NoError(t, err)

				out, ok := publish(t, event.ID).(api.PublishEvent422JSONResponse)
				require.True(t, ok, "expected 422")
				assert.Contains(t, out.Errors, field)

				stored, err := queries.New(db).GetEventById(ctx, event.ID)
				require.NoError(t, err)
				assert.Nil(t, stored.PublishedAt)
			})
		}
	})

	t.Run("rejects an already published event", func(t *testing.T) {
		event, err := mocks.Event(db, ctx)
		require.NoError(t, err)

		out, ok := publish(t, event.ID).(api.PublishEvent422JSONResponse)
		require.True(t, ok, "expected 422")
		assert.Contains(t, out.Errors, "publishedAt")
	})

	t.Run("returns 404 for an unknown id", func(t *testing.T) {
		assert.IsType(t, api.PublishEvent404Response{}, publish(t, uuid.Must(uuid.NewV7())))
	})
}

// TestEventServer_HTTP exercises the same operations through the generated
// router, which is where response objects turn into status codes and JSON.
func TestEventServer_HTTP(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	server := httptest.NewServer(api.NewServer(db).Handler(""))
	defer server.Close()

	decode := func(t *testing.T, resp *http.Response) map[string]any {
		t.Helper()
		defer resp.Body.Close()
		var body map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		return body
	}

	t.Run("GET /events", func(t *testing.T) {
		_, err := mocks.Event(db, ctx)
		require.NoError(t, err)

		resp, err := http.Get(server.URL + "/events?page=1&perPage=2")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body := decode(t, resp)
		assert.Len(t, body["data"], 2)
		assert.Equal(t, float64(2), body["pagination"].(map[string]any)["perPage"])
	})

	t.Run("GET /events/{id}", func(t *testing.T) {
		event, err := mocks.Event(db, ctx)
		require.NoError(t, err)

		resp, err := http.Get(fmt.Sprintf("%s/events/%s", server.URL, event.ID))
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body := decode(t, resp)
		assert.Equal(t, event.ID.String(), body["id"])
		assert.Equal(t, true, body["isFree"])
		assert.Nil(t, body["price"])
		assert.Nil(t, body["currency"])
	})

	t.Run("GET /events/{id} with an unknown id", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/events/%s", server.URL, uuid.Must(uuid.NewV7())))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("GET /events/{id} with a malformed id", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/events/not-a-uuid")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("POST /events", func(t *testing.T) {
		body := validEventInput("api-http-created-event")
		payload, err := json.Marshal(body)
		require.NoError(t, err)

		resp, err := http.Post(server.URL+"/events", "application/json", strings.NewReader(string(payload)))
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		created := decode(t, resp)
		assert.Equal(t, "/events/"+created["id"].(string), resp.Header.Get("Location"))
		assert.Equal(t, "250.00", created["price"])
	})

	t.Run("POST /events with an invalid payload", func(t *testing.T) {
		body := validEventInput("Not A Slug")
		payload, err := json.Marshal(body)
		require.NoError(t, err)

		resp, err := http.Post(server.URL+"/events", "application/json", strings.NewReader(string(payload)))
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

		errs := decode(t, resp)["errors"].(map[string]any)
		assert.Contains(t, errs, "slug")
	})

	t.Run("PATCH /events/{id}", func(t *testing.T) {
		event, err := mocks.Event(db, ctx)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPatch,
			fmt.Sprintf("%s/events/%s", server.URL, event.ID),
			strings.NewReader(`{"titleEn":"Patched over HTTP"}`))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		assert.Equal(t, "Patched over HTTP", decode(t, resp)["titleEn"])
	})

	t.Run("POST /events/{id}/publish", func(t *testing.T) {
		event, err := mocks.Event(db, ctx, func(p *queries.UpsertEventParams) {
			p.Published = false
		})
		require.NoError(t, err)

		resp, err := http.Post(fmt.Sprintf("%s/events/%s/publish", server.URL, event.ID), "application/json", nil)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}
