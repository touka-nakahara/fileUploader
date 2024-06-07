package service

// このファイルはDBとAPIの境界線の接続を担当する
import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	errorHandle "fileUploader/controller/error"
	"fileUploader/model"
	"fileUploader/repository"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/argon2"
)

type FileService interface {
	GetFileListService(ctx context.Context, queryParams *model.GetQueryParam) ([]*model.File, error)
	GetFileService(ctx context.Context, fileID model.FileID) (*model.File, error)
	GetFileDownloadService(ctx context.Context, fileID model.FileID, password string) (*model.FileBlob, error)
	PostFileService(ctx context.Context, file *model.File, fileData *model.FileBlob) (*model.FileID, error)
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

func (s *fileService) GetFileListService(ctx context.Context, queryPrams *model.GetQueryParam) ([]*model.File, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("GetFileListService").Start(ctx, "GetFileListService")
	defer span.End()

	files, err := s.fileRepository.GetAll(ctx, queryPrams)

	if err != nil {
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
	if file.IsAvailable.Before(time.Now()) {
		return nil, errorHandle.ErrFileNotFound
	}

	// パスワードを削除
	if file.Password != "" {
		file.HasPassword = true
		file.Password = ""
	}

	return file, nil
}

type Argon2Config struct {
	Memory      int
	Inerations  int
	Parallelism int
	KeyLength   int
}

func (s *fileService) PostFileService(ctx context.Context, file *model.File, fileData *model.FileBlob) (*model.FileID, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("PostFileService").Start(ctx, "PostFileService")
	defer span.End()

	// 暗号化処理を追加
	if file.Password != "" {
		span.AddEvent("Start crypt")

		hash, err := cryptPassword(file.Password)
		if err != nil {
			return nil, err
		}
		file.Password = *hash

		span.AddEvent("End crypt")
	}

	// uuidの作成
	_uuid := uuid.New()
	file.Uuid = _uuid.String()

	fileID, err := s.fileRepository.Add(ctx, file, fileData)

	if err != nil {
		return nil, err
	}
	return fileID, nil
}

func (s *fileService) DeleteFileService(ctx context.Context, fileID model.FileID, password string) error {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("DeleteFileService").Start(ctx, "DeleteFileService")
	defer span.End()

	//　パスワード認証
	file, err := s.fileRepository.Get(ctx, fileID)

	if err != nil {
		return err
	}

	//TODO 関数化
	if file.Password != "" {
		span.AddEvent("Start crypt")
		vals := strings.Split(file.Password, "$")
		if len(vals) != 6 {
			return errorHandle.ErrServerIntarnal
		}

		var version, memory, iterations, parallelism int
		_, err = fmt.Sscanf(vals[2], "v=%d", &version)
		if err != nil {
			return errorHandle.ErrServerIntarnal
		}

		//TODO vals[3]はENVにあれば良い埋め込む必要はない？
		_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
		if err != nil {
			return errorHandle.ErrServerIntarnal
		}
		salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
		if err != nil {
			return errorHandle.ErrServerIntarnal
		}
		hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
		if err != nil {
			return errorHandle.ErrServerIntarnal
		}

		otherHash := argon2.IDKey([]byte(password), salt, uint32(iterations), uint32(memory), uint8(parallelism), uint32(len(hash)))

		if subtle.ConstantTimeCompare(hash, otherHash) != 1 {
			return errorHandle.ErrUnmatchPassword
		}
		span.AddEvent("End crypt")
	}

	// 有効期限チェック
	if file.IsAvailable.Before(time.Now()) {
		return errorHandle.ErrFileNotFound
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
		return nil, errorHandle.ErrFileNotFound
	}

	//TODO 関数化 deleteと同じに
	if file.Password != "" {
		vals := strings.Split(file.Password, "$")
		if len(vals) != 6 {
			return nil, errorHandle.ErrServerIntarnal
		}

		var version, memory, iterations, parallelism int
		_, err = fmt.Sscanf(vals[2], "v=%d", &version)
		if err != nil {
			return nil, errorHandle.ErrServerIntarnal
		}

		//TODO vals[3]はENVにあれば良い埋め込む必要はない？
		_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
		if err != nil {
			return nil, errorHandle.ErrServerIntarnal
		}
		salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
		if err != nil {
			return nil, errorHandle.ErrServerIntarnal
		}
		hash, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
		if err != nil {
			return nil, errorHandle.ErrServerIntarnal
		}

		otherHash := argon2.IDKey([]byte(password), salt, uint32(iterations), uint32(memory), uint8(parallelism), uint32(len(hash)))

		if subtle.ConstantTimeCompare(hash, otherHash) != 1 {
			return nil, errorHandle.ErrUnmatchPassword
		}
	}

	fileData, err := s.fileRepository.GetData(ctx, fileID)

	if err != nil {
		return nil, err
	}

	return fileData, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func cryptPassword(password string) (*string, error) {
	salt, err := generateRandomBytes(16)
	if err != nil {
		return nil, err
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, 64*1024, 1, 4, b64Salt, b64Hash)

	return &encodedHash, nil
}
