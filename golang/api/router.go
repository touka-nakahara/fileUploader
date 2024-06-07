package api

import (
	"database/sql"
	"fileUploader/api/middleware"
	"fileUploader/controller"
	mq "fileUploader/infra/db/mysql"
	"fileUploader/service"
	"log"
	"net/http"
	"os"
	"strconv"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"

	"net/http/pprof"
)

func NewRouter(db *sql.DB) http.Handler {

	connection := mq.NewFileDB(db)
	s := service.NewFileService(connection)

	maxUploadSize, err := strconv.ParseInt(os.Getenv("MAX_UPLOAD_SIZE"), 10, 64)
	if err != nil {
		log.Fatal(err)

	}

	fileConfig := controller.FileControllerConfig{MaxUploadSize: maxUploadSize}
	fileController := controller.NewFileController(s, &fileConfig)

	mux := http.NewServeMux()

	// 経路を保存するためにotelのhuncに置き換える
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		handler := otelhttp.WithRouteTag(pattern, http.HandlerFunc(handlerFunc))
		mux.Handle(pattern, handler)
	}

	// GET /files
	handleFunc("GET /api/files", fileController.GetFileListHandler)
	// GET /files/id (get detail)
	handleFunc("GET /api/files/{id}", fileController.GetFileHandler)
	// POST files/id/download (download)
	handleFunc("POST /api/files/{id}/download", fileController.GetFileDownloadHandler)
	// POST /files (upload)
	handleFunc("POST /api/files", fileController.PostFileHandler)
	// DELETE /files/id (delte)
	handleFunc("POST /api/files/{id}", fileController.DeleteFileHandler)

	// pprof
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// // ミドルウェアの適用
	h := middleware.CORSMiddleware(mux)
	h = middleware.LoggingMiddleware(h)
	meterHanlder := middleware.NewMetricMiddleware(otel.Meter("Test1_metrics"))
	h = meterHanlder(h)

	handler := otelhttp.NewHandler(h, "/", otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents))

	return handler
}
