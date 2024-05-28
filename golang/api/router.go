package api

import (
	"fileUploader/controller"
	"fileUploader/service"
	"net/http"
)

func NewRouter() *http.ServeMux {
	// DB接続
	// コントローラー作成
	s := service.NewFileService(nil)
	fileController := controller.NewFileController(s)

	r := http.NewServeMux()

	// 静的ファイル

	// GET /
	r.Handle("GET /", http.FileServer(http.Dir("static/root")))
	// GET /files/id
	r.Handle("GET /{id}", http.FileServer(http.Dir("static/files/id")))
	// POST /files/new
	r.Handle("GET /{id}", http.FileServer(http.Dir("static/files/new")))

	// API

	// GET /files?
	r.HandleFunc("GET /api/files", fileController.GetFileListHandler)
	// GET /files/id
	r.HandleFunc("GET /api/files/id", fileController.GetFileHandler)
	// GET files/id/download
	r.HandleFunc("GET /api/files/id/download", fileController.GetFileDownloadHandler)
	// POST /files
	r.HandleFunc("POST /api/files", fileController.PostFileHandler)
	// PUT /files/id
	r.HandleFunc("PUT /api/files", fileController.PutFileHandler)
	// DELETE /files/id
	r.HandleFunc("DELETE /api/files/id", fileController.DeleteFileHandler)

	return r
}
