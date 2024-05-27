package service

// このファイルはDBとAPIの境界線の接続を担当する
import (
	"context"
	"fileUploader/model"
	"fileUploader/repository"
)

type FileService interface {
	GetFileListService(ctx context.Context) ([]*model.File, error)
	GetFileService(ctx context.Context, fileID model.FileID) (*model.File, error)
	PostFileService(ctx context.Context, file *model.File) error
	DeleteFileService(ctx context.Context, fileID model.FileID) error
}

type fileService struct {
	// データベースアクセスを外部から注入
	fileRepository repository.FileRepository
}

var _ FileService = (*fileService)(nil)

func NewFileService(repo repository.FileRepository) *fileService {
	// 注入
	return &fileService{fileRepository: repo}
}

func (s *fileService) GetFileListService(ctx context.Context) ([]*model.File, error) {
	files, err := s.fileRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (s *fileService) GetFileService(ctx context.Context, fileID model.FileID) (*model.File, error) {
	file, err := s.fileRepository.Get(ctx, fileID)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *fileService) PostFileService(ctx context.Context, file *model.File) error {
	if err := s.fileRepository.Add(ctx, file); err != nil {
		return err
	}
	return nil
}

func (s *fileService) DeleteFileService(ctx context.Context, fileID model.FileID) error {
	if err := s.fileRepository.Delete(ctx, fileID); err != nil {
		return err
	}
	return nil
}
