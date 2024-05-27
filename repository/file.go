package repository

import (
	"context"
	"fileUploader/model"
)

// このファイルはデータベースのアクセスの抽象化を担当する

type FileRepository interface {
	GetAll(ctx context.Context) ([]*model.File, error)
	Get(ctx context.Context, id model.FileID) (*model.File, error)
	Add(ctx context.Context, singer *model.File) error
	Delete(ctx context.Context, id model.FileID) error
}
