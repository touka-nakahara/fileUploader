package controller

import (
	"encoding/json"
	"fileUploader/service"
	"net/http"
)

type fileController struct {
	service service.FileService
}

func NewFileController(service service.FileService) *fileController {
	return &fileController{service: service}
}

// GET /, GET/?
func (c *fileController) GetFileListHandler(w http.ResponseWriter, r *http.Request) {
	files, err := c.service.GetFileListService(r.Context())
	if err != nil {
		// GETALLが失敗 = サーバーエラー
		//TODO ログにはSQLのエラー内容を含める
		errorHandler(w, r, 500, "サーバー内部エラーです")
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
	// HTTPレスポンスの設定
	//TODO LOG
}

// GET /files/id
// GET /files/new
// GET files/id/download
// GET files/download?
// DELETE /files/id

// POST /tree
// PUT /tree
// POST /signing
// POST /signup
