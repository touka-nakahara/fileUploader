package controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"fileUploader/controller"
	"fileUploader/infra/dummy"
	mock_repository "fileUploader/infra/mock"
	"fileUploader/model"
	"fileUploader/service"

	"go.uber.org/mock/gomock"
)

func Test_fileController_GetFileListHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := mock_repository.NewMockFileRepository(ctrl)
	s := service.NewFileService(m)
	c := controller.NewFileController(s)
	type exptected struct {
		StatusCode  int
		ContentType string
		Body        controller.Response
		Query       map[string]string
	}
	tests := []struct {
		name         string
		mockBehavior func(m *mock_repository.MockFileRepository)
		exptected    exptected
		request      *http.Request
	}{
		{
			name: "成功パターン",

			mockBehavior: func(m *mock_repository.MockFileRepository) {
				m.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(dummy.DummyFiles, nil)
			},
			request: httptest.NewRequest("GET", "/api/files", nil),
			exptected: exptected{
				StatusCode:  200,
				ContentType: "application/json",
				Body: controller.Response{
					Message: "OK",
					Data:    dummy.DummyFiles,
				},
				Query: nil,
			},
		},
		{
			name: "サーバーエラーパターン",
			mockBehavior: func(m *mock_repository.MockFileRepository) {
				m.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(nil, errors.New("データベース内部エラーです")).AnyTimes()
			},
			request: httptest.NewRequest("GET", "/api/files", nil),
			exptected: exptected{
				StatusCode:  500,
				ContentType: "application/json",
				Body: controller.Response{
					Message: "サーバー内部エラーです",
					Data:    nil,
				},
				Query: nil,
			},
		},
		{
			name: "リクエストを送るパターン",
			mockBehavior: func(m *mock_repository.MockFileRepository) {
				m.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(nil, errors.New("データベース内部エラーです")).AnyTimes()
			},
			request: httptest.NewRequest("GET", "/api/files?sort=filename", nil),
			exptected: exptected{
				StatusCode:  200,
				ContentType: "application/json",
				Body: controller.Response{
					Message: "サーバー内部エラーです",
					Data:    nil,
				},
				Query: map[string]string{
					"sort": "filename",
				},
			},
		},
		{
			name: "複数のリクエストを同時に送るパターン",
			mockBehavior: func(m *mock_repository.MockFileRepository) {
				m.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(nil, errors.New("データベース内部エラーです")).AnyTimes()
			},
			request: httptest.NewRequest("GET", "/api/files?sort=filename&type=jpg", nil),
			exptected: exptected{
				StatusCode:  500,
				ContentType: "application/json",
				Body: controller.Response{
					Message: "サーバー内部エラーです",
					Data:    nil,
				},
				Query: map[string]string{
					"sort": "filename",
					"type": "jpg",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(m)
			roc := httptest.NewRecorder()

			c.GetFileListHandler(roc, tt.request)

			// fmt.Println(res.Result().StatusCode)
			if roc.Result().StatusCode != tt.exptected.StatusCode {
				t.Errorf("got = [%v], want [%v]", roc.Result().StatusCode, tt.exptected.StatusCode)
			}

			// fmt.Println(res.Header().Get("Content-Type"))
			if roc.Header().Get("Content-Type") != tt.exptected.ContentType {
				t.Errorf("got = [%v], want [%v]", roc.Header().Get("Content-Type"), tt.exptected.ContentType)
			}

			// fmt.Println(roc.Body)
			if roc.Body == nil {
				t.Errorf("Request Body is nil!")
			}

			var res controller.Response
			json.Unmarshal(roc.Body.Bytes(), &res)

			if res.Message != tt.exptected.Body.Message {
				t.Errorf("got = [%v], want [%v]", res.Message, tt.exptected.Body.Message)
			}

			//TODD データの完全一致を見る

		})
	}
}

func Test_fileController_GetFileHandler(t *testing.T) {
	// モックの作成
	ctrl := gomock.NewController(t)
	m := mock_repository.NewMockFileRepository(ctrl)
	s := service.NewFileService(m)
	c := controller.NewFileController(s)

	// 期待値
	type exptected struct {
		StatusCode  int
		ContentType string
		Body        controller.Response
		Query       map[string]string
	}

	tests := []struct {
		name         string
		mockBehavior func(m *mock_repository.MockFileRepository)
		exptected    exptected
		request      *http.Request
	}{
		{
			name: "成功パターン",

			mockBehavior: func(m *mock_repository.MockFileRepository) {
				m.EXPECT().Get(gomock.Any(), gomock.Any()).Return(dummy.DummyFiles[0], nil)
			},

			request: httptest.NewRequest("GET", "/api/files/", nil),
			exptected: exptected{
				StatusCode:  200,
				ContentType: "application/json",
				Body: controller.Response{
					Message: "OK",
					Data:    dummy.DummyFiles,
				},
				Query: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(m)
			roc := httptest.NewRecorder()

			tt.request.SetPathValue("id", "550e8400-e29b-41d4-a716-446655440000")
			c.GetFileHandler(roc, tt.request)

			// fmt.Println(res.Result().StatusCode)
			if roc.Result().StatusCode != tt.exptected.StatusCode {
				t.Errorf("got = [%v], want [%v]", roc.Result().StatusCode, tt.exptected.StatusCode)
			}

			// fmt.Println(res.Header().Get("Content-Type"))
			if roc.Header().Get("Content-Type") != tt.exptected.ContentType {
				t.Errorf("got = [%v], want [%v]", roc.Header().Get("Content-Type"), tt.exptected.ContentType)
			}

			// fmt.Println(roc.Body)
			if roc.Body == nil {
				t.Errorf("Request Body is nil!")
			}

			var res controller.Response
			json.Unmarshal(roc.Body.Bytes(), &res)

			if res.Message != tt.exptected.Body.Message {
				t.Errorf("got = [%v], want [%v]", res.Message, tt.exptected.Body.Message)
			}

			//TODD データの完全一致を見る

		})
	}
}

func Test_fileController_PostFileHandler(t *testing.T) {
	// モックの作成
	ctrl := gomock.NewController(t)
	m := mock_repository.NewMockFileRepository(ctrl)
	s := service.NewFileService(m)
	c := controller.NewFileController(s)

	// 期待値
	type exptected struct {
		StatusCode   int
		ContentType  string
		ResponseBody controller.FilesUplaodResponse
		RequestBody  *model.File
	}

	// body
	tests := []struct {
		name         string
		mockBehavior func(m *mock_repository.MockFileRepository)
		exptected    exptected
		request      controller.FilesUploadRequest
	}{
		{
			name: "成功パターン",

			mockBehavior: func(m *mock_repository.MockFileRepository) {
				m.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},

			request: controller.FilesUploadRequest{
				Data: []model.File{
					*dummy.DummyFiles[0],
				},
			},
			exptected: exptected{
				StatusCode:  200,
				ContentType: "application/json",
				ResponseBody: controller.FilesUplaodResponse{
					Messages: []controller.Message{
						{Message: "OK", StatusCode: 200, ID: "550e8400-e29b-41d4-a716-446655440000"},
					},
				},
			},
		},
		{
			name: "失敗パターン",

			mockBehavior: func(m *mock_repository.MockFileRepository) {
				m.EXPECT().Add(gomock.Any(), gomock.Any()).Return(errors.New("データベース内部エラーです")).Times(1)
			},

			request: controller.FilesUploadRequest{
				Data: []model.File{
					*dummy.DummyFiles[0],
				},
			},
			exptected: exptected{
				StatusCode:  200,
				ContentType: "application/json",
				ResponseBody: controller.FilesUplaodResponse{
					Messages: []controller.Message{
						{Message: "データベース内部エラーです", StatusCode: 500, ID: "550e8400-e29b-41d4-a716-446655440000"},
					},
				},
			},
		},
		{
			name: "失敗パターン2",

			mockBehavior: func(m *mock_repository.MockFileRepository) {},
			request: controller.FilesUploadRequest{
				Data: []model.File{},
			},
			exptected: exptected{
				StatusCode:  400,
				ContentType: "application/json",
				ResponseBody: controller.FilesUplaodResponse{
					Messages: []controller.Message{
						{Message: "リクエストがありません", StatusCode: 500, ID: "550e8400-e29b-41d4-a716-446655440000"},
					},
				},
			},
		},
		// {
		// 	name: "複数成功パターン",

		// 	mockBehavior: func(m *mock_repository.MockFileRepository) {
		// 		m.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil).Times(2)
		// 	},

		// 	request: controller.FilesUploadRequest{
		// 		Data: []model.File{
		// 			*dummy.DummyFiles[0],
		// 			*dummy.DummyFiles[1],
		// 		},
		// 	},
		// 	exptected: exptected{
		// 		StatusCode:  201,
		// 		ContentType: "application/json",
		// 		ResponseBody: controller.Response{
		// 			Message: "OK",
		// 		},
		// 	},
		// },
		// {
		// 	name: "一部失敗パターン",

		// 	mockBehavior: func(m *mock_repository.MockFileRepository) {
		// 		m.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		// 	},

		// 	request: controller.FilesUploadRequest{
		// 		Data: []model.File{
		// 			*dummy.DummyFiles[0],
		// 		},
		// 	},
		// 	exptected: exptected{
		// 		StatusCode:  201,
		// 		ContentType: "application/json",
		// 		ResponseBody: controller.Response{
		// 			Message: "OK",
		// 		},
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(m)
			roc := httptest.NewRecorder()

			body, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("unexpected Error !: %v", err)
			}
			request := httptest.NewRequest("POST", "/api/files", bytes.NewReader(body))

			c.PostFileHandler(roc, request)

			// fmt.Println(res.Result().StatusCode)
			if roc.Result().StatusCode != tt.exptected.StatusCode {
				t.Errorf("got = [%v], want [%v]", roc.Result().StatusCode, tt.exptected.StatusCode)
			}

			// fmt.Println(res.Header().Get("Content-Type"))
			if roc.Header().Get("Content-Type") != tt.exptected.ContentType {
				t.Errorf("got = [%v], want [%v]", roc.Header().Get("Content-Type"), tt.exptected.ContentType)
			}

			// fmt.Println(roc.Body)
			if roc.Body == nil {
				t.Errorf("Request Body is nil!")
			}

			var res controller.FilesUplaodResponse
			json.Unmarshal(roc.Body.Bytes(), &res)

			if !reflect.DeepEqual(tt.exptected.ResponseBody, res) {
				t.Errorf("expected and actual structs are not equal. Expected: %v, Actual: %v", tt.exptected.ResponseBody, res.Messages)
			}

			//TODD データの完全一致を見る

		})
	}
}
