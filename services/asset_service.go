package services

import (
	"context"
	"uuid"

	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/types"
)

type AssetService struct {
	db queries.DBTX
}

func (s *AssetService) UploadAsset(ctx context.Context, input *types.UploadAssetInput) (*queries.Asset, error) {
	id := uuid.NewV7()
}
