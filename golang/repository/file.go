package repository

import (
	"context"
	"fileUploader/model"
)

// このファイルはデータベースのアクセスの抽象化を担当する

type FileRepository interface {
	GetAll(ctx context.Context, params *model.GetQueryParam) ([]*model.File, error)
	Get(ctx context.Context, id model.FileID) (*model.File, error)
	GetData(ctx context.Context, id model.FileID) (*model.FileBlob, error)
	Add(ctx context.Context, file *model.File, fileData *model.FileBlob) error
	Delete(ctx context.Context, id model.FileID) error
}
