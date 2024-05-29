package api

import (
	"database/sql"
	"fileUploader/controller"
	mq "fileUploader/infra/db/mysql"
	"fileUploader/service"
	"net/http"
)

func NewRouter(db *sql.DB) *http.ServeMux {
	// DB接続
	// コントローラー作成
	connection := mq.NewFileDB(db)
	s := service.NewFileService(connection)
	fileController := controller.NewFileController(s)

	r := http.NewServeMux()

	// 静的ファイル

	// GET /
	r.Handle("GET /", http.FileServer(http.Dir("static/root")))
	// TODO 何がきてもindexへredirectする (api以外は)

	// API

	// GET /files? (get list)
	r.HandleFunc("GET /api/files", fileController.GetFileListHandler)
	// GET /files/id (get detail)
	r.HandleFunc("GET /api/files/{id}", fileController.GetFileHandler)
	// POST files/id/download (download)
	r.HandleFunc("POST /api/files/{id}/download", fileController.GetFileDownloadHandler)
	// POST /files (upload)
	r.HandleFunc("POST /api/files", fileController.PostFileHandler)
	// PUT /files/id (change)
	r.HandleFunc("PUT /api/files", fileController.PutFileHandler)
	// DELETE /files/id (delte)
	r.HandleFunc("POST /api/files/{id}", fileController.DeleteFileHandler)

	return r
}
