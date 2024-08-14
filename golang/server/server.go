package server

import (
	"context"
	"errors"
	"fileUploader/api"
	mq "fileUploader/infra/db/mysql"
	myotel "fileUploader/otel"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// サーバー起動
func NewServer() {

	// .env読み込み
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	//　ログファイル設定
	applicationLog, err := os.OpenFile("log/application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer applicationLog.Close()
	applicationLogWriter := io.MultiWriter(os.Stdout, applicationLog)

	var programLevel = new(slog.LevelVar) // Info by default
	applicationLogger := slog.New(slog.NewJSONHandler(applicationLogWriter, &slog.HandlerOptions{Level: programLevel}))

	logLevel := os.Getenv("LOG_LEVEL")

	switch logLevel {
	case "INFO":
		programLevel.Set(slog.LevelInfo)
	case "ERROR":
		programLevel.Set(slog.LevelError)
	case "DEBUG":
		programLevel.Set(slog.LevelDebug)
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
	otelShutdown, err := myotel.SetupOTelSDK(ctx)
	if err != nil {
		slog.Error(err.Error(), slog.String("func", "otel.SetupOtelSDK()")) //TODO これスタックトレースとかからうまくできないかなできないかな https://zenn.dev/ryo_yamaoka/articles/858fa01e9e11d0 ?
		return
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	var meter = otel.Meter("")
	_, err = meter.Int64UpDownCounter(
		"items.counter",
		metric.WithDescription("Number of items."),
		metric.WithUnit("{item}"),
	)

	if err != nil {
		slog.Error(err.Error(), slog.String("func", "otel.Meter().Int64UpDownCounter"))
		return
	}

	// サーバー作成
	server := http.Server{
		Addr:         os.Getenv("ADDRESS") + ":" + os.Getenv("PORT"),
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		Handler:      r,
	}

	slog.Debug("Server setting:",
		slog.String("LOG_LEVEL", os.Getenv("LOG_LEVEL")),
		slog.String("ADDRESS", os.Getenv("ADDRESS")),
		slog.String("PORT", os.Getenv("PORT")),
		slog.String("DB_NAME", os.Getenv("DB_NAME")),
		slog.String("DB_USER", os.Getenv("DB_USER")),
		slog.String("DB_ADDRESS", os.Getenv("DB_ADDRESS")),
		slog.String("DB_PORT", os.Getenv("DB_PORT")),
		slog.String("MAX_UPLOAD_SIZE", os.Getenv("MAX_UPLOAD_SIZE")),
	)

	slog.Info("Server start",
		slog.String("ADDRESS", os.Getenv("ADDRESS")),
		slog.String("PORT", os.Getenv("PORT")),
	)

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
