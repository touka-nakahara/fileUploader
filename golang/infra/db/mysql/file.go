package mq

import (
	"context"
	"database/sql"
	"errors"
	"fileUploader/model"
	"fileUploader/repository"
	"fmt"
	"strings"
)

type fileDB struct {
	connection *sql.DB
}

var _ repository.FileRepository = (*fileDB)(nil)

func NewFileDB(db *sql.DB) *fileDB {
	return &fileDB{
		connection: db,
	}
}

// 　全てのファイルを取得する
func (d *fileDB) GetAll(ctx context.Context, params *model.GetQueryParam) ([]*model.File, error) {
	// クエリパラメータの処理
	query := "SELECT id, name, size, extension, description, password, thumbnail, is_available, update_date, upload_date FROM file.File"
	var conditions []string
	var args []interface{}

	// ファイルタイプ
	if extension := params.Extension; extension != "" {
		conditions = append(conditions, "c = ?")
		args = append(args, extension)
	}

	// 1時間以内に削除されたものを見る
	if isAvailable := params.Is_available; isAvailable != "" {
		if isAvailable == "false" {
			conditions = append(conditions, "is_available >= DATE_SUB(NOW(), INTERVAL 1 HOUR) AND is_available <= NOW()")
		}
	} else {
		conditions = append(conditions, "is_available > NOW()")
	}

	// 検索
	if searchParam := params.Search; searchParam != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+searchParam+"%")
	}

	// WHEREクエリの結合
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	//　ソート
	var sort_query string
	if sort_name := params.Sort; sort_name != "" {
		if sort_name == "name" || sort_name == "update_date" || sort_name == "size" {
			sort_query += "order by" + " " + sort_name
		}
	}

	// オーダー
	if direction := params.Ordered; direction != "" {
		if direction == "asc" || direction == "desc" {
			if sort_query != "" {
				sort_query += " " + direction
			} else {
				// 指定していない場合は名前でソート
				sort_query += "order by name" + " " + direction
			}
		}
	}

	// クエリ結合
	if sort_query != "" {
		query += " " + sort_query
	}

	// 最大数制限
	query += " " + "limit 20"

	//　ページ
	if page := params.Page; page != 0 {
		offset := (page - 1) * 20

		query += fmt.Sprintf(" offset %d", offset)
	}

	rows, err := d.connection.QueryContext(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var files []*model.File
	for rows.Next() {
		file := new(model.File)
		err := rows.Scan(
			&file.ID,
			&file.Name,
			&file.Size,
			&file.Extension,
			&file.Description,
			&file.Password,
			&file.Thumbnail,
			&file.IsAvailable,
			&file.UpdateDate,
			&file.UploadDate,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	// エラーを起こす
	return nil, errors.New("Test Error")
	return files, nil
}

func (d *fileDB) Get(ctx context.Context, id model.FileID) (*model.File, error) {
	query := "SELECT id, name, size, extension, description, password,  thumbnail, is_available, update_date, upload_date FROM file.File WHERE id = ?"
	row := d.connection.QueryRowContext(ctx, query, id)

	file := new(model.File)

	err := row.Scan(
		&file.ID,
		&file.Name,
		&file.Size,
		&file.Extension,
		&file.Description,
		&file.Password,
		&file.Thumbnail,
		&file.IsAvailable,
		&file.UpdateDate,
		&file.UploadDate,
	)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (d *fileDB) GetData(ctx context.Context, id model.FileID) (*model.FileBlob, error) {
	query := "SELECT id, data FROM file.Data WHERE file_id = ?"

	row := d.connection.QueryRowContext(ctx, query, id)

	file := new(model.FileBlob)
	err := row.Scan(
		&file.ID,
		&file.Data,
	)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (d *fileDB) Add(ctx context.Context, file *model.File, fileData *model.FileBlob) error {

	// データの実態とメタデータの保存をトランザクションで行う
	tx, err := d.connection.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	if err != nil {
		return err
	}

	query := `
		INSERT INTO file.File (
			name, size, extension, description, password, thumbnail
		) VALUES (?, ?, ?, ?, ?, ?)`

	result, execErr := tx.ExecContext(ctx, query,
		file.Name, file.Size, file.Extension, file.Description,
		file.Password, file.Thumbnail)

	if execErr != nil {
		tx.Rollback()
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		tx.Rollback()
		return err
	}

	query = `
		INSERT INTO file.Data (
			file_id, data
		) VALUES (?, ?)`

	_, execErr = tx.ExecContext(ctx, query, id, fileData.Data)

	if execErr != nil {
		tx.Rollback()
		return err
	}

	//　トランザクション実行
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (d *fileDB) Delete(ctx context.Context, id model.FileID) error {
	query := "UPDATE file.File SET is_available=NOW() WHERE id=?;"

	_, err := d.connection.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

// RV nakaharaY PUTにするにしてはmodel.Fileはオーバースペックのような気がする
// TODO 実装する
func (d *fileDB) Put(ctx context.Context, id model.FileID, file *model.File) error {
	query := "UPDATE file.File"
	var conditions []string
	var args []interface{}
	//　クエリの作成
	if file.Name != "" {
		conditions = append(conditions, "name = ?")
		args = append(args, file.Name)
	}
	if file.Password != "" {
		conditions = append(conditions, "password = ?")
		args = append(args, file.Password)
	}
	if file.Description != "" {
		conditions = append(conditions, "description = ?")
		args = append(args, file.Description)
	}

	if len(conditions) > 0 {
		query += " SET " + strings.Join(conditions, " , ")
	}
	//TODO UPDATEだけど何も変化がない場合

	_, err := d.connection.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}
	return nil
}
