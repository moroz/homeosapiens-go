package services_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/services/mocks"
	"github.com/moroz/homeosapiens-go/types"
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

		actual, err := services.NewVideoService(db).ListVideoGroupsForUser(ctx, user.ID)
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

		actual, err := services.NewVideoService(db).ListVideoGroupsForUser(ctx, user.ID)
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

		actual, err = services.NewVideoService(db).ListVideoGroupsForUser(ctx, admin.ID)
		assert.NoError(t, err)

		ids = nil
		for _, v := range actual {
			ids = append(ids, v.VideoGroup.ID)
			assert.True(t, v.HasAccess)
		}
		assert.Equal(t, []uuid.UUID{paid.ID, free.ID}, ids)
	})
}
