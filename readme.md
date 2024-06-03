# ファイルアップローダー

### エンドポイント

GET / : index.html

### API

- GET /api/files : ファイルの一覧取得
- GET /api/files/{id} : ファイル取得
- POST /api/files/{id}/download : ファイルのダウンロード
- POST /api/files : ファイルのアップロード
- POST /api/files/{id} : ファイルの削除

### 起動方法
```
docker compose run -d

cd fileUploader/react
npm install

cd fileUploader/golang/cmd
go run main.go
```
http://localhost:3000 へアクセス
