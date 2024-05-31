package server

import (
	"context"
	"errors"
	"fileUploader/api"
	mq "fileUploader/infra/db/mysql"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
)

// サーバー起動
func NewServer() {

	// .env読み込み
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	//　ログファイル設定
	serverLog, err := os.OpenFile("../log/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer serverLog.Close()
	serverWriter := io.MultiWriter(os.Stdout, serverLog)

	httpLog, err := os.OpenFile("../log/http.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer httpLog.Close()
	httpWriter := io.MultiWriter(os.Stdout, httpLog)

	var programLevel = new(slog.LevelVar) // Info by default
	serverLogger := slog.New(slog.NewJSONHandler(serverWriter, &slog.HandlerOptions{Level: programLevel}))
	httpLogger := slog.New(slog.NewJSONHandler(httpWriter, &slog.HandlerOptions{Level: programLevel}))

	logLevel := os.Getenv("LOG_LEVEL")

	if logLevel == "ERROR" {
		programLevel.Set(slog.LevelError)
	}
	if logLevel == "INFO" {
		programLevel.Set(slog.LevelInfo)
	}

	// サーバー起動
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// DB接続
	db, err := mq.Connect()
	if err != nil {
		serverLogger.Error(err.Error(), slog.String("func", "mq.Connect()"))
		return
	}

	r := api.NewRouter(db, httpLogger)

	// サーバー作成
	server := http.Server{
		Addr:    os.Getenv("ADDRESS"),
		Handler: r,
	}

	serverLogger.Info("Server start")

	// サーバー起動
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverLogger.Error("Server unexpected closed", slog.String("error", err.Error()))
		}
	}()

	// 優雅なシャットダウン
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	serverLogger.Info("Server closed")
}
