package controller

import (
	"encoding/json"
	errorHandle "fileUploader/controller/error"
	"fileUploader/model"
	"fileUploader/service"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
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
	ctx, span := otel.Tracer("Test").Start(ctx, "GetFileListHandler")
	defer span.End()

	queryParams := r.URL.Query()

	// クエリパラメータ構造体に変換
	getQueryParam := &model.GetQueryParam{}

	// ファイルタイプ
	if extension := queryParams.Get("type"); extension != "" {
		getQueryParam.Extension = extension
	}

	// deleteParam
	if isAvailable := queryParams.Get("is_available"); isAvailable != "" {
		getQueryParam.Is_available = isAvailable
	}

	// searchParam
	if searchParam := queryParams.Get("search"); searchParam != "" {
		getQueryParam.Search = searchParam
	}

	// sort
	if sort_name := queryParams.Get("sort"); sort_name != "" {
		getQueryParam.Sort = sort_name
	}

	// order
	if direction := queryParams.Get("ordered"); direction != "" {
		getQueryParam.Ordered = direction
	}

	// page
	if page := queryParams.Get("page"); page != "" {
		pageInt, err := strconv.Atoi(page)
		if err == nil {
			getQueryParam.Page = pageInt
		}
	}

	files, err := c.service.GetFileListService(ctx, getQueryParam)

	if err != nil {
		errorHandle.ErrorHandler(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	var res model.Response = model.Response{
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
		errorHandle.ErrorHandler(w, r, errorHandle.ErrInvalidRequest)
		return
	}

	file, err := c.service.GetFileService(ctx, model.FileID(fileID))
	if err != nil {
		errorHandle.ErrorHandler(w, r, err)
		return
	}

	// HTTPヘッダーの設定
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	var res model.Response = model.Response{
		Message: "OK",
		Data:    file,
	}
	json.NewEncoder(w).Encode(res)

}

func (c *fileController) PostFileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("PostFileHandler").Start(ctx, "PostFileHandler")
	defer span.End()
	span.AddEvent("start handle")

	type Message struct {
		File model.File     `json:"file"`
		Data model.FileBlob `json:"data"`
	}
	var m Message

	// ファイルサイズを制限
	span.AddEvent("Parse MultipartForm")
	// MaxUploadSize := c.config.MaxUploadSize
	// r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)

	// io.Pipeにファイルを入れておきたい
	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type")) // params["boundary"] に boundaryが入ってる
	if err != nil {
		log.Fatal(err)
	}

	rp, wp := io.Pipe()

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(r.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			if p.FormName() == "file" {
				err := json.NewDecoder(p).Decode(&m.File)
				if err != nil {
					errorHandle.ErrorHandler(w, r, err)
					return
				}
			}
			if p.FormName() == "data" {
				go func() {
					_, err := io.Copy(wp, p)
					if err != nil {
						fmt.Println("copy error")
						fmt.Println(err)
					}
				}()
			}
		}
	}

	span.AddEvent("End parse")

	fileID, err := c.service.PostFileService(ctx, &m.File, rp)
	if err != nil {
		errorHandle.ErrorHandler(w, r, err)
		return
	}

	// メタデータを返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	upResp := &model.UploadResponse{
		ID: *fileID,
	}

	resp := &model.Response{
		Message: "ok",
		Data:    upResp,
	}

	json.NewEncoder(w).Encode(resp)
}

func (c *fileController) DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("DeleteFileHandler").Start(ctx, "DeleteFileHandler")
	defer span.End()

	// どこのファイルを削除するか確認
	pathParam := r.PathValue("id")
	fileID, err := strconv.Atoi(pathParam)
	if err != nil {
		errorHandle.ErrorHandler(w, r, errorHandle.ErrInvalidRequest)
		return
	}

	// passwordを取得
	type Message struct {
		Password string `json:"password"`
	}
	var m Message
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		errorHandle.ErrorHandler(w, r, err)
		return
	}

	// ファイルを削除
	if err := c.service.DeleteFileService(ctx, model.FileID(fileID), m.Password); err != nil {
		errorHandle.ErrorHandler(w, r, err)
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
		errorHandle.ErrorHandler(w, r, errorHandle.ErrInvalidRequest)
		return
	}

	// passwordを取得
	type Message struct {
		Password string `json:"password"`
	}
	var m Message
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		errorHandle.ErrorHandler(w, r, err)
		return
	}

	// ファイルをダウンロードする
	fileData, err := c.service.GetFileDownloadService(ctx, model.FileID(fileID), m.Password)

	if err != nil {
		errorHandle.ErrorHandler(w, r, err)
		return
	}

	// すべて成功した場合
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	var res model.Response = model.Response{
		Data: fileData,
	}
	span.AddEvent("encode Json")
	json.NewEncoder(w).Encode(res)
	span.AddEvent("end encoding")
}
