package handlers_test

import (
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/web/router"
	"github.com/moroz/homeosapiens-go/web/sessions"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Knowing the URL of a video in a paid series must not be enough to watch it:
// the page has to withhold the sources, which are the only thing pointing at the
// asset on the CDN.
func TestVideoController_Show_Paywall(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	store, err := sessions.NewStore(config.SessionKey)
	require.NoError(t, err)

	server := httptest.NewServer(router.Router(db, store, mocks.NewMockStripeService(t)))
	defer server.Close()

	product, err := mocks.Product(db, ctx, func(p *queries.InsertProductParams) {
		p.BasePriceAmount = decimal.NewFromInt(199)
		p.BasePriceCurrency = "PLN"
	})
	require.NoError(t, err)

	group, err := mocks.VideoGroup(db, ctx, func(p *queries.InsertVideoGroupParams) {
		p.ProductID = &product.ID
	})
	require.NoError(t, err)

	video, err := mocks.Video(db, ctx)
	require.NoError(t, err)

	_, err = queries.New(db).AddVideoToVideoGroup(ctx, &queries.AddVideoToVideoGroupParams{
		VideoID:      video.ID,
		VideoGroupID: group.ID,
	})
	require.NoError(t, err)

	const objectKey = "videos/secret-master.m3u8"
	_, err = queries.New(db).UpsertVideoSource(ctx, &queries.UpsertVideoSourceParams{
		ID:          uuid.Must(uuid.NewV7()),
		VideoID:     video.ID,
		ContentType: "application/x-mpegURL",
		ObjectKey:   objectKey,
	})
	require.NoError(t, err)

	user, err := mocks.User(db, ctx)
	require.NoError(t, err)

	client, err := mocks.ClientWithUser(&mocks.ClientWithUserInput{
		Store:   store,
		Server:  server,
		User:    user,
		DB:      db,
		Context: ctx,
	})
	require.NoError(t, err)

	videoURL := fmt.Sprintf("%s/videos/%s/%s", server.URL, group.Slug, video.Slug)

	get := func() string {
		resp, err := client.Get(videoURL)
		require.NoError(t, err)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		return string(body)
	}

	body := get()
	assert.NotContains(t, body, objectKey)
	assert.NotContains(t, body, "<source")
	assert.Contains(t, body, "Paid content")

	// Buying the product is what unlocks the series, and with it the sources.
	_, err = db.Exec(ctx, "insert into user_product_access (user_id, product_id) values ($1, $2)", user.ID, product.ID)
	require.NoError(t, err)

	body = get()
	assert.Contains(t, body, objectKey)
}
