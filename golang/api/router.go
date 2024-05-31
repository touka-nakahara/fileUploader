package api

import (
	"database/sql"
	"fileUploader/api/middleware"
	"fileUploader/controller"
	mq "fileUploader/infra/db/mysql"
	"fileUploader/service"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

func NewRouter(db *sql.DB, httpLogger *slog.Logger) http.Handler {

	connection := mq.NewFileDB(db)
	s := service.NewFileService(connection)

	//RV nakaharaY こんなところにいるはずもないのに
	maxUploadSize, err := strconv.ParseInt(os.Getenv("MAX_UPLOAD_SIZE"), 10, 64)
	if err != nil {
		log.Fatal(err)

	}
	fileConfig := controller.FileControllerConfig{MaxUploadSize: maxUploadSize}
	fileController := controller.NewFileController(s, &fileConfig)

	r := http.NewServeMux()

	// GET /
	//RV 静的ファイルとしてbuildできない
	// r.Handle("GET /", http.FileServer(http.Dir("static/root")))

	// API
	//RV nakaharaY ハンドラ名メソッドに紐づけるとズレた時めんどいのでやめるべき

	// GET /files? (get list)
	r.HandleFunc("GET /api/files", fileController.GetFileListHandler)
	// GET /files/id (get detail)
	r.HandleFunc("GET /api/files/{id}", fileController.GetFileHandler)
	// POST files/id/download (download)
	r.HandleFunc("POST /api/files/{id}/download", fileController.GetFileDownloadHandler)
	// POST /files (upload)
	r.HandleFunc("POST /api/files", fileController.PostFileHandler)
	// DELETE /files/id (delte)
	r.HandleFunc("POST /api/files/{id}", fileController.DeleteFileHandler)

	// PUT /files/id (change)
	// r.HandleFunc("PUT /api/files/{id}", fileController.PutFileHandler)

	//RV middlewareの適用順番などを自由にできるパターンが欲しい
	h := middleware.CORSMiddleware(r)
	h = middleware.LoggingMiddleware(h, httpLogger)

	return h
}
