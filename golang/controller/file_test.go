package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
		Body        []*model.File
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
				m.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(dummy.GetListRequest, nil)
			},
			request: httptest.NewRequest("GET", "/api/files", nil),
			exptected: exptected{
				StatusCode:  200,
				ContentType: "application/json",
				Body:        dummy.GetListRequest,
			},
		},
		// クエリパラメータの処理方法がわからんのでそのまま
		{
			name: "リクエストを送るパターン",
			mockBehavior: func(m *mock_repository.MockFileRepository) {
				m.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(dummy.GetListRequest, nil).AnyTimes()
			},
			request: httptest.NewRequest("GET", "/api/files?sort=filename&type=jpg", nil),
			exptected: exptected{
				StatusCode:  200,
				ContentType: "application/json",
				Body:        dummy.GetListRequest,
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

		})
	}
}
