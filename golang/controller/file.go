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
)

type fileController struct {
	service service.FileService
}

func NewFileController(service service.FileService) *fileController {
	return &fileController{service: service}
}

// GET /files, GET/files?
func (c *fileController) GetFileListHandler(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()

	files, err := c.service.GetFileListService(r.Context(), queryParams)

	if err != nil {
		errorHandler(w, r, 500, "サーバー内部エラーです")
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

	pathParam := r.PathValue("id")
	fileID, err := strconv.Atoi(pathParam)
	if err != nil || fileID == 0 {
		errorHandler(w, r, 400, "不明なリクエストです")
		return
	}

	file, err := c.service.GetFileService(r.Context(), model.FileID(fileID))
	if err != nil {
		if err.Error() == "ファイルが見つかりません" {
			errorHandler(w, r, 404, err.Error())
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			errorHandler(w, r, 404, "ファイルが見つかりません")
			return
		}
		errorHandler(w, r, 500, err.Error())
		return
	}

	// HTTPヘッダーの設定
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	var res Response = Response{
		//RV nakaharaY これいらない？
		Message: "OK",
		Data:    file,
	}
	json.NewEncoder(w).Encode(res)

}

// GET files/id/download
func (c *fileController) GetFileDownloadHandler(w http.ResponseWriter, r *http.Request) {

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
	fileData, err := c.service.GetFileDownloadService(r.Context(), model.FileID(fileID), m.Password)

	if err != nil || fileID == 0 {
		errorHandler(w, r, 400, err.Error())
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

func (c *fileController) PostFileHandler(w http.ResponseWriter, r *http.Request) {

	type Message struct {
		File model.File     `json:"file"`
		Data model.FileBlob `json:"data"`
	}
	var m Message

	//TODO マジックナンバー
	const MaxUploadSize = 1024 * 1024 * 500

	// ファイルサイズを制限
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		http.Error(w, "500MB以下のファイルを選択してください。", http.StatusBadRequest)
	}

	file := r.FormValue("file")
	err := json.Unmarshal([]byte(file), &m.File)
	if err != nil {
		errorHandler(w, r, 500, err.Error())
	}

	data, _, err := r.FormFile("data")
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer data.Close()

	fileData, err := io.ReadAll(data)
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}

	m.Data.Data = fileData

	//TODO 2つのテーブル間のinsertどちらも成功しないといけない
	if err := c.service.PostFileService(r.Context(), &m.File, &m.Data); err != nil {
		errorHandler(w, r, 500, err.Error())
		return
	}

	w.WriteHeader(204)
}

func (c *fileController) DeleteFileHandler(w http.ResponseWriter, r *http.Request) {

	// どこのファイルを削除するか確認
	pathParam := r.PathValue("id")
	fileID, err := strconv.Atoi(pathParam)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		var res Response = Response{
			Message: err.Error(), //　不明なアドレスです
			Data:    nil,
		}
		json.NewEncoder(w).Encode(res)
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

	// ファイルを削除
	if err := c.service.DeleteFileService(r.Context(), model.FileID(fileID), m.Password); err != nil {
		errorHandler(w, r, 400, err.Error())
		return
	}

	w.WriteHeader(204)
}

// PUT /files/id
func (c *fileController) PutFileHandler(w http.ResponseWriter, r *http.Request) {
	pathParam := r.PathValue("id")
	fileID, err := strconv.Atoi(pathParam)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		var res Response = Response{
			Message: err.Error(), //　不明なアドレスです
			Data:    nil,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	// パラメータの取得
	type Message struct {
		File model.File `json:"file"`
	}
	var m Message
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		errorHandler(w, r, 400, err.Error())
		return
	}

	if err := c.service.PutFileService(r.Context(), model.FileID(fileID), &m.File); err != nil {
		errorHandler(w, r, 400, err.Error())
		return
	}

	w.WriteHeader(204)
}
