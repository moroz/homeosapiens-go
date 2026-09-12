package types

import "io"

type UploadAssetInput struct {
	Filename    string
	ContentType string
	File        io.ReadCloser
}
