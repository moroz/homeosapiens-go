package api

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moroz/homeosapiens-go/services"
	"github.com/moroz/homeosapiens-go/types"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type assetServer struct {
	db *pgxpool.Pool
}

func NewAssetServer(db *pgxpool.Pool) *assetServer {
	return &assetServer{db}
}

func guessContentType(r multipart.File) (string, error) {
	buf := make([]byte, 512)
	if _, err := r.Read(buf); err != nil {
		return "", err
	}
	if _, err := r.Seek(0, 0); err != nil {
		return "", err
	}

	return http.DetectContentType(buf), nil
}

func getFileFromMultipartForm(form *multipart.Form) (*types.UploadAssetInput, error) {
	for _, fheaders := range form.File {
		for _, header := range fheaders {
			file, err := header.Open()
			if err != nil {
				return nil, err
			}

			contentType, err := guessContentType(file)
			if err != nil {
				return nil, err
			}

			return &types.UploadAssetInput{
				Filename:    header.Filename,
				ContentType: contentType,
				File:        file,
			}, nil
		}
	}

	return nil, errors.New("file not found in multipart form")
}

func (s *assetServer) UploadAsset(ctx context.Context, request UploadAssetRequestObject) (UploadAssetResponseObject, error) {
	form, err := request.Body.ReadForm(16 << 10) // 16 MB
	if err != nil {
		return nil, err
	}

	srv, err := services.NewAssetService(s.db)
	if err != nil {
		return nil, err
	}

	input, err := getFileFromMultipartForm(form)
	if err != nil {
		return nil, err
	}

	asset, err := srv.UploadAsset(ctx, input)
	if err != nil {
		return nil, err
	}

	return UploadAsset201JSONResponse{
		Id:               openapi_types.UUID(asset.ID),
		InsertedAt:       asset.InsertedAt,
		ObjectKey:        *asset.ObjectKey,
		OriginalFilename: asset.OriginalFilename,
	}, nil
}
