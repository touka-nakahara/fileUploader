package mock_repository

import (
	"context"
	"fmt"
	"testing"

	"fileUploader/infra/dummy"

	gomock "go.uber.org/mock/gomock"
)

func TestNewServer(t *testing.T) {
	ctrl := gomock.NewController(t)

	m := NewMockFileRepository(ctrl)

	m.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(dummy.DummyFiles, nil)
	tests := []struct {
		name string
	}{
		{name: "Just Run"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			files, _ := m.GetAll(ctx, nil)
			for _, file := range files {
				fmt.Printf("File ID: %d, Name: %s, Size: %d, Extension: %s, Description: %s\n",
					file.ID, file.Name, file.Size, file.Extension, file.Description)
			}
		})
	}
}
