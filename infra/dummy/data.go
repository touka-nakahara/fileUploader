package dummy

import (
	"fileUploader/model"
	"time"
)

var DummyFiles = []*model.File{
	{
		ID:          "1",
		Name:        "example1.txt",
		Size:        1234,
		Extension:   "txt",
		Description: "This is an example text file",
		Password:    "password123",
		UUID:        "550e8400-e29b-41d4-a716-446655440000",
		FileID:      "1",
		Thumbnail:   []byte{0x89, 0x50, 0x4e, 0x47},
		Data:        []byte("Example file content 1"),
		IsAvailable: time.Now(),
		UpdateDate:  time.Now().Add(-24 * time.Hour),
		UploadDate:  time.Now().Add(-48 * time.Hour),
	},
	{
		ID:          "2",
		Name:        "example2.jpg",
		Size:        5678,
		Extension:   "jpg",
		Description: "This is an example JPEG image",
		Password:    "password456",
		UUID:        "550e8400-e29b-41d4-a716-446655440001",
		FileID:      "2",
		Thumbnail:   []byte{0xff, 0xd8, 0xff, 0xe0},
		Data:        []byte("Example file content 2"),
		IsAvailable: time.Now(),
		UpdateDate:  time.Now().Add(-24 * time.Hour),
		UploadDate:  time.Now().Add(-48 * time.Hour),
	},
	{
		ID:          "3",
		Name:        "example3.pdf",
		Size:        91011,
		Extension:   "pdf",
		Description: "This is an example PDF document",
		Password:    "password789",
		UUID:        "550e8400-e29b-41d4-a716-446655440002",
		FileID:      "3",
		Thumbnail:   []byte{0x25, 0x50, 0x44, 0x46},
		Data:        []byte("Example file content 3"),
		IsAvailable: time.Now(),
		UpdateDate:  time.Now().Add(-24 * time.Hour),
		UploadDate:  time.Now().Add(-48 * time.Hour),
	},
}
