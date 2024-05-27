package model

import "time"

type FileID int

type File struct {
	ID          FileID `json:"id"`
	Name        string `json:"name"`
	Size        int    `json:"size"`
	Extension   string `json:"extension"`
	Description string `json:"description"`

	Password string `json:"password"`

	UUID string `json:"uuid"`

	FileID    FileID `json:"file_id"`
	Thumbnail []byte `json:"thumnbnail"`
	Data      []byte `json:"file_data"`

	IsAvailable time.Time `json:"is_available"`
	UpdateDate  time.Time `json:"update_date"`
	UploadDate  time.Time `json:"upload_date"`
}
