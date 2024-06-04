package service

// このファイルはDBとAPIの境界線の接続を担当する
import (
	"context"
	"errors"
	"fileUploader/model"
	"fileUploader/repository"
	"net/url"
	"time"

	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

// RV nakaharaY ほんとにここか？ どこならいいのかわからないやつ
var ErrFileNotFound = errors.New("ファイルが見つかりません")
var ErrInvalidRequest = errors.New("不明なリクエストです")
var ErrUnmatchPassword = errors.New("パスワードが違います")
var ErrServerIntarnal = errors.New("サーバー内部エラーです")

type FileService interface {
	GetFileListService(ctx context.Context, queryParams url.Values) ([]*model.File, error)
	GetFileService(ctx context.Context, fileID model.FileID) (*model.File, error)
	GetFileDownloadService(ctx context.Context, fileID model.FileID, password string) (*model.FileBlob, error)
	PostFileService(ctx context.Context, file *model.File, fileData *model.FileBlob) error
	DeleteFileService(ctx context.Context, fileID model.FileID, password string) error
}

type fileService struct {
	// データベースアクセスを外部から注入
	fileRepository repository.FileRepository
}

var _ FileService = (*fileService)(nil)

func NewFileService(fileRepo repository.FileRepository) *fileService {
	// 注入
	return &fileService{
		fileRepository: fileRepo,
	}
}

func (s *fileService) GetFileListService(ctx context.Context, queryPrams url.Values) ([]*model.File, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("GetFileListService").Start(ctx, "GetFileListService")
	defer span.End()

	files, err := s.fileRepository.GetAll(ctx, queryPrams)

	if err != nil {
		//RV ログを日本語にするか問題
		return nil, err
	}

	// パスワードを削除
	for _, file := range files {
		if file.Password != "" {
			file.HasPassword = true
			file.Password = ""
		}
	}

	return files, nil
}

func (s *fileService) GetFileService(ctx context.Context, fileID model.FileID) (*model.File, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("GetFileService").Start(ctx, "GetFileService")
	defer span.End()

	file, err := s.fileRepository.Get(ctx, fileID)

	if err != nil {
		return nil, err
	}

	// 有効期限チェック
	//RV 　DB側では有効期限チェックをしていないのはいいのだろうか？
	if file.IsAvailable.Before(time.Now()) {
		return nil, ErrFileNotFound
	}

	// パスワードを削除
	if file.Password != "" {
		file.HasPassword = true
		file.Password = ""
	}

	return file, nil
}

func (s *fileService) PostFileService(ctx context.Context, file *model.File, fileData *model.FileBlob) error {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("PostFileService").Start(ctx, "PostFileService")
	defer span.End()

	if err := s.fileRepository.Add(ctx, file, fileData); err != nil {
		return err
	}
	return nil
}

func (s *fileService) DeleteFileService(ctx context.Context, fileID model.FileID, password string) error {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("DeleteFileService").Start(ctx, "DeleteFileService")
	defer span.End()

	//　パスワード認証
	file, err := s.fileRepository.Get(ctx, fileID)

	if err != nil {
		return err
	}

	if file.Password != "" {
		result := bcrypt.CompareHashAndPassword([]byte(file.Password), []byte(password))
		if result != nil {
			return result
		}
	}

	// 有効期限チェック
	if file.IsAvailable.Before(time.Now()) {
		return ErrFileNotFound
	}

	if err := s.fileRepository.Delete(ctx, fileID); err != nil {
		return err
	}
	return nil
}

func (s *fileService) GetFileDownloadService(ctx context.Context, fileID model.FileID, password string) (*model.FileBlob, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("GetFileDownloadService").Start(ctx, "GetFileDownloadService")
	defer span.End()

	// パスワードチェック
	file, err := s.fileRepository.Get(ctx, fileID)

	if err != nil {
		return nil, err
	}

	// 有効期限チェック
	if file.IsAvailable.Before(time.Now()) {
		return nil, ErrFileNotFound
	}

	if file.Password != "" {
		result := bcrypt.CompareHashAndPassword([]byte(file.Password), []byte(password))
		if result != nil {
			return nil, result
		}
	}

	fileData, err := s.fileRepository.GetData(ctx, fileID)

	if err != nil {
		return nil, err
	}

	return fileData, nil
}
