package controller

import (
	"encoding/json"
	"fileUploader/model"
	"fileUploader/service"
	"net/http"
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
		errorHandler(w, r, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	var res Response = Response{
		Message: "OK",
		Data:    files,
	}
	json.NewEncoder(w).Encode(res)

	//TODO LOG
}

// GET /files/id
func (c *fileController) GetFileHandler(w http.ResponseWriter, r *http.Request) {
	fileID := r.PathValue("id")

	files, err := c.service.GetFileService(r.Context(), model.FileID(fileID))
	if err != nil {
		errorHandler(w, r, 500, err.Error())
		return
	}

	// HTTPヘッダーの設定
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	var res Response = Response{
		Message: "OK",
		Data:    files,
	}
	json.NewEncoder(w).Encode(res)

	//TODO LOG
}

func (c *fileController) PostFileHandler(w http.ResponseWriter, r *http.Request) {
	type Message struct {
		File model.File     `json:"file"`
		Data model.FileBlob `json:"data"`
	}
	var m Message

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		errorHandler(w, r, 400, err.Error())
		return
	}

	//TODO 2つのテーブル間のinsertどちらも成功しないといけない
	if err := c.service.PostFileService(r.Context(), &m.File, &m.Data); err != nil {
		errorHandler(w, r, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	var res = map[string]string{"message": "Created"}
	json.NewEncoder(w).Encode(res)
}

// POST /files
func (c *fileController) PostFilesHandler(w http.ResponseWriter, r *http.Request) {
	// JSONリクエストからmodel.fileを構成
	var request FilesUploadRequest
	var response FilesUplaodResponse
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		//TODO どうしてこうなった...
		var message Message
		message.Message = err.Error()
		message.ID = ""
		message.StatusCode = 400
		response.Messages = append(response.Messages, message)
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(request.Data) == 0 {
		var message Message
		message.Message = "リクエストがありません"
		message.ID = ""
		message.StatusCode = 400
		response.Messages = append(response.Messages, message)
		json.NewEncoder(w).Encode(response)
		return
	}

	var messages []Message

	//TODO サービスでLoopしてたら結局DBのアクセスが1回にならないので意味ない
	c.service.PostFilesService()

	// すべて成功した場合
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	var res FilesUplaodResponse = FilesUplaodResponse{
		Messages: messages,
	}
	json.NewEncoder(w).Encode(res)
}

// GET files/id/download
func (c *fileController) GetFileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	// どこのファイルをダウンロードするか確認
	fileID := r.PathValue("id")
	// ファイルを取得
	fileData, err := c.service.GetFileDownloadService(r.Context(), model.FileID(fileID))

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		var res Response = Response{
			Message: err.Error(),
			Data:    nil,
		}
		json.NewEncoder(w).Encode(res)
		return
	}

	// すべて成功した場合
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	var res Response = Response{
		Message: "OK",
		Data:    fileData,
	}
	json.NewEncoder(w).Encode(res)
}

func (c *fileController) DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	// どこのファイルを削除するか確認
	// ファイルを削除
	// ファイルの取得でエラーが発生した時 (ファイルがない, パスワードが違う)
	// jsonに詰めて送信
}

// GET files/download?
func (c *fileController) GetFilesDownloadHandler(w http.ResponseWriter, r *http.Request) {
	// どこのファイルをダウンロードするか確認
	// ファイルを取得
	// ファイルの取得でエラーが発生した時 (ファイルがない, パスワードが違う)
	// jsonに詰めて送信
}

func (c *fileController) DeleteFilesHandler(w http.ResponseWriter, r *http.Request) {
	// どこのファイルを削除するか確認
	// ファイルを削除
	// ファイルの取得でエラーが発生した時 (ファイルがない, パスワードが違う)
	// jsonに詰めて送信
}

// DELETE /files/id

// POST /tree
// PUT /tree
// POST /signing
// POST /signup
