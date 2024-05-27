package controller_test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"fileUploader/controller"
	"fileUploader/infra/dummy"
	mock_repository "fileUploader/infra/mock"
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
	}
	tests := []struct {
		name         string
		mockBehavior func(m *mock_repository.MockFileRepository)
		exptected    exptected
	}{
		{
			name: "成功パターン",
			mockBehavior: func(m *mock_repository.MockFileRepository) {
				m.EXPECT().GetAll(gomock.Any()).Return(dummy.DummyFiles, nil)
			},
			exptected: exptected{
				StatusCode:  200,
				ContentType: "application/json",
				Body: controller.Response{
					Message: "OK",
					Data:    dummy.DummyFiles,
				},
			},
		},
		{
			name: "サーバーエラーパターン",
			mockBehavior: func(m *mock_repository.MockFileRepository) {
				m.EXPECT().GetAll(gomock.Any()).Return(nil, errors.New("データベース内部エラーです")).AnyTimes()
			},
			exptected: exptected{
				StatusCode:  500,
				ContentType: "application/json",
				Body: controller.Response{
					Message: "サーバー内部エラーです",
					Data:    nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(m)
			req := httptest.NewRequest("GET", "/api/files", nil)
			roc := httptest.NewRecorder()

			c.GetFileListHandler(roc, req)

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
