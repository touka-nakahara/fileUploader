package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fileUploader/model"
	"fileUploader/service"
	"io"
	"net/http"
	"strconv"

	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

type fileController struct {
	service service.FileService
	config  *FileControllerConfig
}

type FileControllerConfig struct {
	MaxUploadSize int64
}

func NewFileController(service service.FileService, config *FileControllerConfig) *fileController {
	return &fileController{
		service: service,
		config:  config,
	}
}

// GET /files, GET/files?
func (c *fileController) GetFileListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("GetFileListHandler").Start(ctx, "GetFileListHandler")
	defer span.End()

	queryParams := r.URL.Query()

	files, err := c.service.GetFileListService(ctx, queryParams)

	if err != nil {
		errorHandler(w, r, 500, service.ErrServerIntarnal.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	var res Response = Response{
		Message: "OK",
		Data:    files,
	}
	json.NewEncoder(w).Encode(res)
}

// GET /files/id
func (c *fileController) GetFileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("GetFileHandler").Start(ctx, "GetFileHandler")
	defer span.End()

	pathParam := r.PathValue("id")
	fileID, err := strconv.Atoi(pathParam)
	if err != nil || fileID == 0 {
		errorHandler(w, r, 400, service.ErrInvalidRequest.Error())
		return
	}

	file, err := c.service.GetFileService(ctx, model.FileID(fileID))
	if err != nil {

		// 有効期限
		if errors.Is(err, service.ErrFileNotFound) {
			errorHandler(w, r, 404, err.Error())
			return
		}

		// 存在しない
		if errors.Is(err, sql.ErrNoRows) {
			errorHandler(w, r, 404, "ファイルが見つかりません")
			return
		}

		// それ以外
		errorHandler(w, r, 500, err.Error())
		return
	}

	// HTTPヘッダーの設定
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	var res Response = Response{
		Message: "OK",
		Data:    file,
	}
	json.NewEncoder(w).Encode(res)

}

func (c *fileController) PostFileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("PostFileHandler").Start(ctx, "PostFileHandler")
	defer span.End()

	type Message struct {
		File model.File     `json:"file"`
		Data model.FileBlob `json:"data"`
	}
	var m Message

	MaxUploadSize := c.config.MaxUploadSize

	// ファイルサイズを制限
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		errorHandler(w, r, 400, "500MB以下のファイルを選択してください。")
		return
	}

	// ファイルメタデータの読み取り
	//RV nakaharaY file.File -> file.Metadata ?
	file := r.FormValue("file")
	err := json.Unmarshal([]byte(file), &m.File)
	if err != nil {
		//RV nakaharaY ここでこけたらサーバーで直さないといけないからログとらないことまりそう
		errorHandler(w, r, 500, err.Error())
	}

	// アップロードされたファイルの読み取り
	data, _, err := r.FormFile("data")
	if err != nil {
		errorHandler(w, r, 400, "不明なファイルが送信されました")
		return
	}
	defer data.Close()

	fileData, err := io.ReadAll(data)
	if err != nil {
		//RV nakaharaY こういうところのエラーログ残ってて欲しいんだけどそうしたら中身をチェックするってことだからな〜どうしたらいいんだろう？
		errorHandler(w, r, 400, "不明なファイルが送信されました")
		return
	}

	m.Data.Data = fileData

	if err := c.service.PostFileService(ctx, &m.File, &m.Data); err != nil {
		errorHandler(w, r, 500, err.Error())
		return
	}

	w.WriteHeader(204)
}

func (c *fileController) DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("DeleteFileHandler").Start(ctx, "DeleteFileHandler")
	defer span.End()

	// どこのファイルを削除するか確認
	pathParam := r.PathValue("id")
	fileID, err := strconv.Atoi(pathParam)
	if err != nil {
		errorHandler(w, r, 400, service.ErrInvalidRequest.Error())
		return
	}

	// passwordを取得
	type Message struct {
		Password string `json:"password"`
	}
	var m Message
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		errorHandler(w, r, 500, err.Error())
		return
	}

	// ファイルを削除
	if err := c.service.DeleteFileService(ctx, model.FileID(fileID), m.Password); err != nil {

		// パスワード
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			errorHandler(w, r, 400, service.ErrUnmatchPassword.Error())
			return
		}

		// 有効期限
		if errors.Is(err, service.ErrFileNotFound) {
			errorHandler(w, r, 400, err.Error())
			return
		}

		errorHandler(w, r, 400, err.Error())
		return
	}

	w.WriteHeader(204)
}

// GET files/id/download
func (c *fileController) GetFileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("GetFileDownloadHandler").Start(ctx, "GetFileDownloadHandler")
	defer span.End()

	// どこのファイルをダウンロードするか確認
	pathParam := r.PathValue("id")
	fileID, err := strconv.Atoi(pathParam)
	if err != nil || fileID == 0 {
		errorHandler(w, r, 400, "不明なリクエストです")
		return
	}

	// passwordを取得
	type Message struct {
		Password string `json:"password"`
	}
	var m Message
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		errorHandler(w, r, 400, err.Error())
		return
	}

	// ファイルをダウンロードする
	fileData, err := c.service.GetFileDownloadService(ctx, model.FileID(fileID), m.Password)

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			errorHandler(w, r, 400, "パスワードが違います")
			return
		}

		if errors.Is(err, service.ErrFileNotFound) {
			errorHandler(w, r, 400, err.Error())
			return
		}

		errorHandler(w, r, 500, err.Error())
		return
	}

	// すべて成功した場合
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	var res Response = Response{
		Data: fileData,
	}
	json.NewEncoder(w).Encode(res)
}
