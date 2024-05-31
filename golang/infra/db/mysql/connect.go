package mq

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {

	c := mysql.Config{
		DBName:               os.Getenv("DB_NAME"),
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_ADDRESS"),
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", c.FormatDSN())
	if err != nil {
		return nil, err
	}

	//　接続
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	//　接続確認
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
