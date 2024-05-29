package mq

import (
	"context"
	"database/sql"
	"fileUploader/model"
	"fileUploader/repository"
	"net/url"
	"strings"
)

type fileDB struct{ connection *sql.DB }

var _ repository.FileRepository = (*fileDB)(nil)

func NewFileDB(db *sql.DB) *fileDB {
	return &fileDB{
		connection: db,
	}
}

// 　全てのファイルを取得する
func (d *fileDB) GetAll(ctx context.Context, params url.Values) ([]*model.File, error) {
	// クエリパラメータの処理
	query := "SELECT id, name, size, extension, description, password, UUID, thumbnail, is_available, update_date, upload_date FROM file.File"
	var conditions []string
	var args []interface{}

	//TODO クエリパラメータが複数 (jpg, pdf)の時対応できてない
	if extension := params.Get("type"); extension != "" {
		//TODO サニタイズ
		conditions = append(conditions, "extension = ?")
		args = append(args, extension)
	}

	if isAvailable := params.Get("is_available"); isAvailable != "" {
		if isAvailable == "false" {
			conditions = append(conditions, "is_available <= DATE_SUB(NOW(), INTERVAL 1 HOUR")
		}
	} else {
		conditions = append(conditions, "is_available > NOW()")
	}

	if searchParam := params.Get("search"); searchParam != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, searchParam)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	//　ソート
	var sort_query string
	if sort_name := params.Get("sort"); sort_name != "" {
		if sort_name == "name" || sort_name == "update_date" || sort_name == "size" {
			sort_query += "order by" + " " + sort_name
		}
	}

	if direction := params.Get("ordered"); direction != "" {
		// 昇順 or 降順
		if direction == "asce" || direction == "desc" {
			if sort_query != "" {
				sort_query += " " + direction
			} else {
				sort_query += "order by name" + " " + direction
			}
		}
	}

	if sort_query != "" {
		query += " " + sort_query
	}

	//TODO ページング
	query += " " + "limit 20"

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
			&file.UUID,
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

	return files, nil
}

func (d *fileDB) Get(ctx context.Context, id model.FileID) (*model.File, error) {
	query := "SELECT id, name, size, extension, description, password, UUID, thumbnail, is_available, update_date, upload_date FROM file.File WHERE id = ?"
	row := d.connection.QueryRowContext(ctx, query, id)

	file := new(model.File)
	err := row.Scan(
		&file.ID,
		&file.Name,
		&file.Size,
		&file.Extension,
		&file.Description,
		&file.Password,
		&file.UUID,
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
	// トランザクションの開始
	tx, err := d.connection.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return err
	}

	// file.Fileへデータ保存
	query := `
		INSERT INTO file.File (
			name, size, extension, description, password, UUID, thumbnail
		) VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, execErr := tx.ExecContext(ctx, query,
		file.Name, file.Size, file.Extension, file.Description,
		file.Password, file.UUID, file.Thumbnail)

	if execErr != nil {
		tx.Rollback()
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	// file.Dataへデータ保存
	query = `
		INSERT INTO file.Data (
			file_id, data
		) VALUES (?, ?)`

	_, execErr = tx.ExecContext(ctx, query,
		id, fileData.Data)

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

// TODO PUTにするにしては構造体として大きすぎるよな〜
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

// TODO 外部キー制約があるので消せないわ...
func (d *fileDB) Delete(ctx context.Context, id model.FileID) error {
	query := "UPDATE file.File SET is_available=NOW() WHERE id=?;"
	_, err := d.connection.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
