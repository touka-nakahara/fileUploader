package mq

import (
	"context"
	"database/sql"
	"fileUploader/model"
	"fileUploader/repository"
	"fmt"
	"os"
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
	var sortQuery string
	if sort_name := params.Sort; sort_name != "" {
		if sort_name == "name" || sort_name == "update_date" || sort_name == "size" {
			sortQuery += "order by" + " " + sort_name
		}
	} else { // 指定しない場合はupdate_date
		sortQuery += "order by update_date"
	}

	// オーダー
	if direction := params.Ordered; direction != "" {
		if direction == "asc" || direction == "desc" {
			if sortQuery != "" {
				sortQuery += " " + direction
			} else {
				// 指定していない場合は名前でソート
				sortQuery += "order by name" + " " + direction
			}
		}
	} else { // 指定していない場合は降順
		sortQuery += " " + "desc"
	}

	// クエリ結合
	if sortQuery != "" {
		query += " " + sortQuery
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

func (d *fileDB) GetData(ctx context.Context, id model.FileID, uuid string) (*model.FileBlob, error) {

	localPath := os.Getenv("FILE_DIR")
	data, err := os.ReadFile(localPath + "/" + uuid)
	if err != nil {
		return nil, err
	}

	file := new(model.FileBlob)
	file.ID = id
	file.Data = data

	return file, nil
}

func (d *fileDB) Add(ctx context.Context, file *model.File, fileData *model.FileBlob) (*model.FileID, error) {
	// データの実態とメタデータの保存をトランザクションで行う

	tx, err := d.connection.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	if err != nil {
		return nil, err
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
		return nil, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//TODO これはDBじゃないからなんかいい感じにしたいけど...
	localPath := os.Getenv("FILE_DIR")

	err = os.WriteFile(localPath+"/"+file.Uuid, fileData.Data, 0666)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	//　トランザクション実行
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	fileID := model.FileID(id)

	return &fileID, nil
}

func (d *fileDB) Delete(ctx context.Context, id model.FileID) error {
	query := "UPDATE file.File SET is_available=NOW() WHERE id=?;"

	_, err := d.connection.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}
