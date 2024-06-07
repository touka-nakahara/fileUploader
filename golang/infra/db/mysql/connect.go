package mq

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

func Connect() (*sql.DB, error) {

	c := mysql.Config{
		DBName:               os.Getenv("DB_NAME"),
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_ADDRESS") + ":" + os.Getenv("DB_PORT"),
		ParseTime:            true,
		AllowNativePasswords: true,
		InterpolateParams:    true,
	}

	db, err := otelsql.Open("mysql", c.FormatDSN(), otelsql.WithAttributes(semconv.DBSystemMySQL), otelsql.WithSQLCommenter(true))
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

	// meter
	err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(semconv.DBSystemMySQL))
	if err != nil {
		return nil, err
	}

	return db, nil
}
