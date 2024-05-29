package controller

import (
	"encoding/json"
	"fileUploader/model"
	"fileUploader/service"
	"fmt"
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

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

	file, err := c.service.GetFileService(r.Context(), model.FileID(fileID))
	if err != nil {
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

	//TODO LOG
}

func (c *fileController) PostFileHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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
	fmt.Println(file)
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

// GET files/id/download
func (c *fileController) GetFileDownloadHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// どこのファイルをダウンロードするか確認
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
		Data: fileData,
	}
	json.NewEncoder(w).Encode(res)
}

func (c *fileController) DeleteFileHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

	// ファイルを削除
	if err := c.service.DeleteFileService(r.Context(), model.FileID(fileID)); err != nil {
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

// func (c *fileController) PostFilesHandler(w http.ResponseWriter, r *http.Request) {
// 	// データをパース
// 	type FilesUploadRequest struct {
// 		File     *model.File       `json:"file"`
// 		FileBlob []*model.FileBlob `json:"data"`
// 	}

// 	type Message struct {
// 		StatusCode int    `json:"statuscode"`
// 		Message    string `json:"message"`
// 		ID         string `json:"ID"`
// 	}

// 	type FilesUplaodResponse struct {
// 		Messages []Message `json:"messages"`
// 	}

// 	var request FilesUploadRequest
// 	err := json.NewDecoder(r.Body).Decode(&request)

// 	var response FilesUplaodResponse

// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(400)
// 		var message Message
// 		message.Message = err.Error()
// 		message.ID = ""
// 		message.StatusCode = 400
// 		response.Messages = append(response.Messages, message)
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	if len(request.FileBlob) == 0 {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(400)
// 		var message Message
// 		message.Message = "リクエストがありません"
// 		message.ID = ""
// 		message.StatusCode = 400
// 		response.Messages = append(response.Messages, message)
// 		json.NewEncoder(w).Encode(response)
// 		return
// 	}

// 	var messages []Message
// 	for _, fileData := range request.FileBlob {
// 		err := c.service.PostFileService(r.Context(), request.File, fileData)

// 		//TODO UUIDの作成
// 		var message Message
// 		if err != nil {
// 			message.Message = err.Error()
// 			message.ID = ""
// 			message.StatusCode = 400
// 		} else {
// 			message.Message = "OK"
// 			message.ID = ""
// 			message.StatusCode = 200
// 		}
// 		messages = append(messages, message)
// 	}
// 	// 一部分 or 全て成功した場合 or 全て失敗した場合
// 	//TODO いい書き方が思いつかないので200で返す
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(200)

// 	response.Messages = messages
// 	json.NewEncoder(w).Encode(response)
// }

// GET files/download?
// func (c *fileController) GetFilesDownloadHandler(w http.ResponseWriter, r *http.Request) {
// 	// どこのファイルをダウンロードするか確認
// 	// ファイルを取得
// 	// ファイルの取得でエラーが発生した時 (ファイルがない, パスワードが違う)
// 	// jsonに詰めて送信
// }

// func (c *fileController) DeleteFilesHandler(w http.ResponseWriter, r *http.Request) {
// 	// どこのファイルを削除するか確認
// 	// ファイルを削除
// 	// ファイルの取得でエラーが発生した時 (ファイルがない, パスワードが違う)
// 	// jsonに詰めて送信
// }

// DELETE /files/id

// POST /tree
// PUT /tree
// POST /signing
// POST /signup
