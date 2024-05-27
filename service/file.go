package service

// このファイルはDBとAPIの境界線の接続を担当する
import (
	"context"
	"fileUploader/model"
	"fileUploader/repository"
	"net/url"
)

type FileService interface {
	GetFileListService(ctx context.Context, queryParams url.Values) ([]*model.File, error)
	GetFileService(ctx context.Context, fileID model.FileID) (*model.File, error)
	DeleteFileService(ctx context.Context, fileID model.FileID) error

	GetFileDownloadService(ctx context.Context, fileID model.FileID) (*model.FileBlob, error)
	GetFilesDownloadService(ctx context.Context, fileID model.FileID) ([]*model.FileBlob, error)

	PostFilesService(ctx context.Context, file []*model.File, fileBlob []*model.FileBlob) error
	PostFileService(ctx context.Context, file *model.File, fileBlob *model.FileBlob) error
}

type fileService struct {
	// データベースアクセスを外部から注入
	fileRepository     repository.FileRepository
	fileBlobRepository repository.FileBlobRepository
}

var _ FileService = (*fileService)(nil)

func NewFileService(fileRepo repository.FileRepository, fileBlobRepo repository.FileBlobRepository) *fileService {
	// 注入
	return &fileService{
		fileRepository:     fileRepo,
		fileBlobRepository: fileBlobRepo,
	}
}

func (s *fileService) GetFileListService(ctx context.Context, queryPrams url.Values) ([]*model.File, error) {
	files, err := s.fileRepository.GetAll(ctx, queryPrams)
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

func (s *fileService) PostFileService(ctx context.Context, file *model.File, fileBlob *model.FileBlob) error {
	if err := s.fileRepository.Add(ctx, file); err != nil {
		return err
	}
	if err := s.fileBlobRepository.Add(ctx, fileBlob); err != nil {
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
