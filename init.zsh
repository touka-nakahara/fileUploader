# アップロードされたファイルを物理削除する
rm -f ./golang/infra/db/local/* 
rm -f ./sql/log/slow.log
cp /dev/null ./sql/log/slow.log
docker compose down --volumes
docker compose up -d
# 取扱注意！！！！