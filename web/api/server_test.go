package api_test

import (
	"testing"

	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/web/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_ListHosts(t *testing.T) {
	ctx := t.Context()
	db, err := initDB(ctx)
	require.NoError(t, err)
	defer db.Close()

	_, err = mocks.Host(db, ctx)
	require.NoError(t, err)

	srv := api.NewServer(db)

	list := func(params api.ListHostsParams) api.ListHosts200JSONResponse {
		resp, err := srv.ListHosts(ctx, api.ListHostsRequestObject{Params: params})
		require.NoError(t, err)
		out, ok := resp.(api.ListHosts200JSONResponse)
		require.True(t, ok, "expected 200, got %T", resp)
		return out
	}

	t.Run("omitted pagination params fall back to the defaults", func(t *testing.T) {
		out := list(api.ListHostsParams{})

		assert.Equal(t, int32(1), out.Pagination.Page)
		assert.Equal(t, int32(config.DefaultPageSize), out.Pagination.PerPage)
		assert.NotEmpty(t, out.Data)
	})

	t.Run("reports the page that was served", func(t *testing.T) {
		out := list(api.ListHostsParams{Page: new(int32(2)), PerPage: new(int32(1))})

		assert.Equal(t, int32(2), out.Pagination.Page)
		assert.Equal(t, int32(1), out.Pagination.PerPage)
	})
}
