package services

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"slices"
	"strings"
	"uuid"

	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/config"
	"github.com/moroz/homeosapiens-go/db/queries"
	"github.com/moroz/homeosapiens-go/types"
)

type AssetService struct {
	db *pgxpool.Pool
	s3 *s3.Client
}

func NewAssetService(db *pgxpool.Pool) (*AssetService, error) {
	cfg, err := aws_config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &AssetService{
		db: db,
		s3: client,
	}, nil
}

var AssetContentTypeWhitelist = []string{
	"image/jpeg",
	"image/png",
}

var ErrDisallowedAssetContentType = errors.New("disallowed Content-Type")

func (s *AssetService) UploadAsset(ctx context.Context, input *types.UploadAssetInput) (*queries.Asset, error) {
	if !slices.Contains(AssetContentTypeWhitelist, input.ContentType) {
		return nil, fmt.Errorf("%w: must be one of: %s, got: %v", ErrDisallowedAssetContentType, strings.Join(AssetContentTypeWhitelist, ", "), input.ContentType)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("UploadAsset: %w", err)
	}
	defer tx.Rollback(ctx)

	id := uuid.NewV7()
	extension, err := mime.ExtensionsByType(input.ContentType)
	if err != nil || len(extension) == 0 {
		return nil, fmt.Errorf("UploadAsset: could not guess extension for MIME type %v: %w", input.ContentType, err)
	}
	objectKey := fmt.Sprintf("assets/%s/original%s", id, extension[0])

	asset, err := queries.New(s.db).InsertAsset(ctx, &queries.InsertAssetParams{
		ID:               id,
		ObjectKey:        &objectKey,
		OriginalFilename: &input.Filename,
	})
	if err != nil {
		return nil, fmt.Errorf("UploadAsset: %w", err)
	}

	_, err = s.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:       &config.AssetBucketName,
		Key:          &objectKey,
		Body:         input.File,
		CacheControl: new("public, max-age=31536000, immutable"),
	})
	if err != nil {
		return nil, fmt.Errorf("UploadAsset: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("UploadAsset: %w", err)
	}

	return asset, nil
}
