package repository

import (
	"context"
	"fileUploader/model"
)

type FileBlobRepository interface {
	Get(ctx context.Context, id model.FileID) (*model.FileBlob, error)
	Add(ctx context.Context, fileBlob *model.FileBlob) error
}
