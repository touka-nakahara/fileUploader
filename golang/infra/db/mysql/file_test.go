package mq_test

import (
	"context"
	mq "fileUploader/infra/db/mysql"
	"fileUploader/model"
	"fmt"
	"net/url"
	"strconv"
	"testing"
)

func Test_fileDB_GetAll(t *testing.T) {
	db, err := mq.Connect()
	if err != nil {
		t.Fatalf("Unexpected Error %v", err)
	}
	fileDB := mq.NewFileDB(db)
	tests := []struct {
		name   string
		params url.Values
	}{
		{name: "test1", params: nil},
		{name: "test2", params: url.Values{"type": {"jpg"}}},
		{name: "test3", params: url.Values{"type": {"jpg", "pdf"}}},
		{name: "test4", params: url.Values{"sort": {"update_date"}, "ordered": {"desc"}}},
		{name: "test4", params: url.Values{"sort": {"size"}}},
		{name: "test4", params: url.Values{"sort": {"name"}}},
		{name: "test4", params: url.Values{"ordered": {"desc"}}},
		//TODO 検索テスト
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gots, err := fileDB.GetAll(context.Background(), tt.params)
			if err != nil {
				t.Errorf("fileDB.GetAll() error = %v", err)
				return
			}
			for _, got := range gots {
				fmt.Printf("%#v\n", got)
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("fileDB.GetAll() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func Test_fileDB_Get(t *testing.T) {
	db, err := mq.Connect()
	if err != nil {
		t.Fatalf("Unexpected Error %v", err)
	}
	fileDB := mq.NewFileDB(db)
	tests := []struct {
		name   string
		params url.Values
	}{
		{name: "test1", params: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fileDB.Get(context.Background(), 5)
			if err != nil {
				t.Errorf("fileDB.GetAll() error = %v", err)
				return
			}
			fmt.Printf("%#v\n", got)
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("fileDB.GetAll() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func Test_PlayGround(t *testing.T) {
	s, e := strconv.Atoi("")
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(s)
}

func Test_fileDB_Put(t *testing.T) {
	db, err := mq.Connect()
	if err != nil {
		t.Fatalf("Unexpected Error %v", err)
	}
	fileDB := mq.NewFileDB(db)

	tests := []struct {
		name   string
		params *model.File
	}{
		{name: "test1", params: &model.File{Name: "hogehoge"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fileDB.Put(context.Background(), 2, tt.params)
			if err != nil {
				t.Errorf("fileDB.GetAll() error = %v", err)
				return
			}
			got, err := fileDB.Get(context.Background(), 2)
			if err != nil {
				t.Errorf("fileDB.GetAll() error = %v", err)
				return
			}
			fmt.Printf("%#v", got)
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("fileDB.GetAll() = %v, want %v", got, tt.want)
			// }
		})
	}
}
