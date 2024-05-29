package server

import (
	"context"
	"errors"
	"fileUploader/api"
	mq "fileUploader/infra/db/mysql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	addr = ":8888"
)

// サーバー起動
func NewServer() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	//TODO Logハンドラ設定

	db, err := mq.Connect()
	if err != nil {
		// LOG ハンドラ設定
		log.Fatalf("%v", err)
	}
	r := api.NewRouter(db)

	// サーバー作成
	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	//TODO サーバーログ作成
	fmt.Println("[info] Server Start")
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			//TODO サーバーログを保存
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
