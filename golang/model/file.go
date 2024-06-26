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

	Thumbnail []byte `json:"thumnbnail"`

	IsAvailable time.Time `json:"is_available"`
	UpdateDate  time.Time `json:"update_date"`
	UploadDate  time.Time `json:"upload_date"`

	HasPassword bool `json:"has_password"`
}

type FileBlob struct {
	ID   FileID `json:"id"`
	Data []byte `json:"data"`
}
