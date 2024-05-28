package mq

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	//TODO マジックナンバーをconfig化
	db, err := sql.Open("mysql", "root:root@tcp(:3307)/mysql?parseTime=true")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// 接続保証
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
