package server

import (
	"context"
	"errors"
	"fmt"
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
	//TODO ルーティング設定
	r := http.NewServeMux()

	// 静的ファイル

	// GET /
	r.Handle("GET /", http.FileServer(http.Dir("static/root")))
	// GET /files/id
	// POST /files
	// PUT /files/id

	// GET /signin
	// GET /signup

	// API

	// GET /?

	// GET /files/id
	// GET /files/new
	// GET files/id/download
	// GET files/download?
	// DELETE /files/id

	// POST /tree
	// PUT /tree
	// POST /signing
	// POST /signup

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
