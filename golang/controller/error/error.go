package error

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

var ErrFileNotFound = errors.New("ファイルが見つかりません")
var ErrInvalidRequest = errors.New("不明なリクエストです")
var ErrUnmatchPassword = errors.New("パスワードが違います")
var ErrServerIntarnal = errors.New("サーバー内部エラーです")
var ErrTooLarge = errors.New("ファイルの容量が大きすぎます") //TODO 許容ファイルサイズを追加？

// エラーが発生したときのレスポンス処理をここで行う
func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {

	// エラーハンドリング
	var statusCode int

	// ハンドリングするものでなかったらServerInternalを返す
	//TODO defaultの処理だけ別でログを残す
	switch {
	case errors.Is(err, sql.ErrNoRows):
		statusCode = 404
		err = ErrFileNotFound
	case errors.Is(err, ErrInvalidRequest):
		statusCode = 400
	case errors.Is(err, ErrFileNotFound):
		statusCode = 404
	case errors.Is(err, ErrTooLarge):
		statusCode = 400
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		statusCode = 400
		err = ErrUnmatchPassword
	case errors.Is(err, ErrUnmatchPassword):
		statusCode = 400
	default:
		spanCtx := trace.SpanContextFromContext(r.Context())
		slog.Error("error check", slog.String("error message", err.Error()), slog.String("trace-id", spanCtx.TraceID().String()))
		statusCode = 500
		err = ErrServerIntarnal
	}

	type ErrorMessage struct {
		Message string `json:"message"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&ErrorMessage{Message: err.Error()})
}
