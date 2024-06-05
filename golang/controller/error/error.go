package error

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var FileNotFound = errors.New("ファイルが見つかりません")
var InvalidRequest = errors.New("不明なリクエストです")
var UnmatchPassword = errors.New("パスワードが違います")
var ServerIntarnal = errors.New("サーバー内部エラーです")
var TooLarge = errors.New("ファイルの容量が大きすぎます") //TODO 許容ファイルサイズを追加？

// エラーが発生したときのレスポンス処理をここで行う
func ErrorHandler(w http.ResponseWriter, _ *http.Request, err error) {

	// エラーハンドリング
	var statusCode int

	// ハンドリングするものでなかったらServerInternalを返す
	//TODO defaultの処理だけ別でログを残す
	switch {
	case errors.Is(err, sql.ErrNoRows):
		statusCode = 404
		err = FileNotFound
	case errors.Is(err, InvalidRequest):
		statusCode = 400
	case errors.Is(err, FileNotFound):
		statusCode = 404
	case errors.Is(err, TooLarge):
		statusCode = 400
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		statusCode = 400
		err = UnmatchPassword
	default:
		statusCode = 500
		err = ServerIntarnal
	}

	type ErrorMessage struct {
		Message string `json:"message"`
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&ErrorMessage{Message: err.Error()})
}
