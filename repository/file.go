package repository

import (
	"context"
	"fileUploader/model"
	"net/url"
)

// このファイルはデータベースのアクセスの抽象化を担当する

type FileRepository interface {
	GetAll(ctx context.Context, params url.Values) ([]*model.File, error)
	Get(ctx context.Context, id model.FileID) (*model.File, error)
	Add(ctx context.Context, file *model.File) error
	Delete(ctx context.Context, id model.FileID) error
}

