package server

import (
	"context"
	"errors"
	"fileUploader/api"
	mq "fileUploader/infra/db/mysql"
	"fileUploader/otel"
	"io"
	"log"
	"log/slog"
	"net"
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
	applicationLog, err := os.OpenFile("../log/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer applicationLog.Close()
	applicationLogWriter := io.MultiWriter(os.Stdout, applicationLog)

	var programLevel = new(slog.LevelVar) // Info by default
	applicationLogger := slog.New(slog.NewJSONHandler(applicationLogWriter, &slog.HandlerOptions{Level: programLevel}))

	logLevel := os.Getenv("LOG_LEVEL")

	if logLevel == "ERROR" {
		programLevel.Set(slog.LevelError)
	}
	if logLevel == "INFO" {
		programLevel.Set(slog.LevelInfo)
	}

	slog.SetDefault(applicationLogger)

	// DB接続
	db, err := mq.Connect()
	if err != nil {
		slog.Error(err.Error(), slog.String("func", "mq.Connect()"))
		return
	}

	r := api.NewRouter(db)

	// Handle SIGINT (CTRL + C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Set up Otel
	otelShutdown, err := otel.SetupOTelSDK(ctx)
	if err != nil {
		slog.Error(err.Error(), slog.String("func", "otel.SetupOtelSDK()")) //TODO これスタックトレースとかからうまくできないかなできないかな https://zenn.dev/ryo_yamaoka/articles/858fa01e9e11d0 ?
		return
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// サーバー作成
	server := http.Server{
		Addr:         os.Getenv("ADDRESS"),
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}

	slog.Info("Server start")

	srvErr := make(chan error, 1)
	// サーバー起動
	go func() {
		srvErr <- server.ListenAndServe()
	}()

	// 優雅なシャットダウン
	select {
	case err = <-srvErr:
		if errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server unexpected closed", slog.String("error", err.Error()))
			return
		}
	case <-ctx.Done():
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
		slog.Info("Server closed gracefully")
	}
}
